// Copyright 16-May-2020 ºDeme
// GNU General Public License - V3 <http://www.gnu.org/licenses/>

package primitives

import (
	"github.com/dedeme/dmstack/imports"
	"github.com/dedeme/dmstack/machine"
	"github.com/dedeme/dmstack/reader"
	"github.com/dedeme/dmstack/stack"
	"github.com/dedeme/dmstack/symbol"
	"github.com/dedeme/dmstack/token"
	"io/ioutil"
	"os"
	"path"
	"sync"
)

// Execute a procedure saved in stack.
//    m: Virtual machine.
//    run: Function which running a machine.
func prRun(m *machine.T, run func(m *machine.T)) {
	tk := m.PopT(token.Procedure)
	proc, _ := tk.P()
	run(machine.New(m.SourceDir, m.Pmachines, proc))
}

// Execute a procedure in an isolated machine. The stack of the isolated
// machine is put en the current stack as a list.
//    m: Virtual machine
//    run: Function which running a machine.
func prData(m *machine.T, run func(m *machine.T)) {
	tk := m.PopT(token.Procedure)
	proc, _ := tk.P()
	m2 := machine.NewIsolate(m.SourceDir, m.Pmachines, proc)
	run(m2)
	m.Push(token.NewL(*m2.Stack, m.MkPos()))
}

// Process an 'import'.
//    m: Virtual machine.
//    run: Function which running a machine.
func prImport(m *machine.T, run func(m *machine.T)) {
	tk := m.Pop()
	sourceKv, err := imports.ReadSymbol(tk)
	if err != nil {
		m.Fail("Import error", err.Error())
	}
	source := sourceKv.Value
	sourcef := path.Clean(path.Join(m.SourceDir, source.String()))
	sourcefSym := symbol.New(sourcef)

	onWay := imports.IsOnWay(sourcefSym)
	_, imported := imports.Get(sourcefSym)
	if !onWay && !imported {
		f := sourcef + ".dms"
		if _, err := os.Stat(f); err != nil {
			m.Fail("Import error", "File '%v' not found.", f)
			return
		}

		imports.PutOnWay(sourcefSym)
		bs, err := ioutil.ReadFile(f)
		if err != nil {
			panic(err)
		}
		rd := reader.New(sourcef, string(bs))

		prg, ok := rd.Process().P()
		if !ok {
			m.Fail("Import error", "Reader process does not return a Program.")
		}

		m2 := machine.NewIsolate(path.Dir(f), m.Pmachines, prg)
		run(m2)

		imports.Add(sourcefSym, m2.Heap)
		imports.QuitOnWay(sourcefSym)
	} else if onWay {
		m.Fail("Import error", "Cyclic import of '%v'.", sourcefSym)
	}
	m.ImportsAdd(sourcefSym)
}

// Process an 'if'.
//    m: Virtual machine.
//    run: Function which running a machine.
func prIf(m *machine.T, run func(m *machine.T)) {
	tk := m.Pop()
	p, ok := tk.P()
	if ok {
		tk2 := m.Pop()
		b, ok := tk2.B()
		if ok {
			if b {
				run(machine.New(m.SourceDir, m.Pmachines, p))
			}
			return
		}
		sy, ok := tk2.Sy()
		if !ok || sy != symbol.Else {
			m.Failt(
				"\n  Expected: Bool or 'else'.\n  Actual  : %s.",
				tk2,
			)
		}
	} else {
		sy, ok := tk.Sy()
		if !ok || sy != symbol.Else {
			m.Failt(
				"\n  Expected: Procedure or 'else'.\n  Actual  : %s.",
				tk,
			)
		}
		p = nil
	}

	var a []*token.T
	for {
		a = append(a, m.PopT(token.Procedure))
		a = append(a, m.PopT(token.Procedure))

		tk2, ok := stack.Peek(*m.Stack)
		if !ok {
			break
		}
		sy, ok := tk2.Sy()
		if !ok || sy != symbol.Else {
			break
		}
		m.Pop()
	}

	for {
		if len(a) == 0 {
			break
		}
		a2, tk1, _ := stack.Pop(a)
		a3, tk2, _ := stack.Pop(a2)
		a = a3
		p1, _ := tk1.P()
		m2 := machine.NewIsolate(m.SourceDir, m.Pmachines, p1)
		run(m2)
		if len(*m2.Stack) != 1 {
			m.Failt(
				"\n  Expected: Only one value returned of type Bool."+
					"\n  Actual  : %v values returned."+
					"\nIn '%v'.",
				len(*m2.Stack), tk1,
			)
		}
		tk11 := m2.Pop()
		b, ok := tk11.B()
		if !ok {
			m.Failt(
				"\n  Expected: Value of type Bool."+
					"\n  Actual  : %v."+
					"\nIn '%v'.",
				tk11, tk1,
			)
		}
		if b {
			p2, _ := tk2.P()
			run(machine.New(m.SourceDir, m.Pmachines, p2))
			return
		}
	}
	if p != nil {
		run(machine.New(m.SourceDir, m.Pmachines, p))
	}
}

// Process an 'elif'.
//    m: Virtual machine.
//    run: Function which running a machine.
func prElif(m *machine.T, run func(m *machine.T)) {
	tk1 := m.PopT(token.Procedure)
	tk2 := m.PopT(token.Procedure)
	tk3 := m.PopT(token.Bool)
	b, _ := tk3.B()
	if b {
		p, _ := tk2.P()
		run(machine.New(m.SourceDir, m.Pmachines, p))
		return
	}
	p, _ := tk1.P()
	run(machine.New(m.SourceDir, m.Pmachines, p))
}

// Process a loop.
//    m: Virtual machine.
//    run: Function which running a machine.
func prLoop(m *machine.T, run func(m *machine.T)) {
	tk := m.PopT(token.Procedure)
	p, _ := tk.P()
	for {
		run(machine.New(m.SourceDir, m.Pmachines, p))
		tk2, ok := stack.Peek(*m.Stack)
		if ok {
			sy, ok2 := tk2.Sy()
			if ok2 && sy == symbol.Break {
				m.Pop()
				break
			}
		}
	}
}

// Process a while.
//    m: Virtual machine.
//    run: Function which running a machine.
func prWhile(m *machine.T, run func(m *machine.T)) {
	tk1 := m.PopT(token.Procedure)
	tk2 := m.PopT(token.Procedure)
	p, _ := tk1.P()
	pb, _ := tk2.P()
	for {
		m2 := machine.NewIsolate(m.SourceDir, m.Pmachines, pb)
		run(m2)

		if len(*m2.Stack) != 1 {
			m.Failt(
				"\n  Expected: Only one value returned of type Bool."+
					"\n  Actual  : %v values returned."+
					"\nIn '%v'.",
				len(*m2.Stack), tk2,
			)
		}
		tk21 := m2.Pop()
		b, ok := tk21.B()
		if !ok {
			m.Failt(
				"\n  Expected: Value of type Bool."+
					"\n  Actual  : %v."+
					"\nIn '%v'.",
				tk21, tk2,
			)
		}
		if !b {
			break
		}

		run(machine.New(m.SourceDir, m.Pmachines, p))
		tk12, ok := stack.Peek(*m.Stack)
		if ok {
			sy, ok2 := tk12.Sy()
			if ok2 && sy == symbol.Break {
				m.Pop()
				break
			}
		}
	}
}

// Process a for.
//    m: Virtual machine.
//    run: Function which running a machine.
func prFor(m *machine.T, run func(m *machine.T)) {
	tk1 := m.PopT(token.Procedure)
	p, _ := tk1.P()
	tk2 := m.Pop()
	step := int64(1)
	begin := int64(0)
	end, ok := tk2.I()
	if !ok {
		l, ok := tk2.L()
		if !ok {
			m.Failt(
				"\n  Expected: Value of type Int or List."+
					"\n  Actual  : %v.",
				tk2.StringDraft(),
			)
		}
		ok = len(l) == 2 || len(l) == 3
		var ok1, ok2, ok3 bool
		if ok {
			begin, ok1 = l[0].I()
			end, ok2 = l[1].I()
			ok3 = true
			if len(l) == 3 {
				step, ok3 = l[2].I()
			}
		}
		if !(ok && ok1 && ok2 && ok3) {
			m.Fail(
				"\n  Expected: [Int, Int] or [Int, Int, Int]."+
					"\n  Actual  : %v.",
				tk2.StringDraft(),
			)
		}
		if step == 0 {
			m.Fail("For error", "Step can not be 0")
		}
	}

	pos := m.MkPos()
	if step > 0 {
		for i := begin; i < end; i += step {
			m2 := machine.New(m.SourceDir, m.Pmachines, p)
			m2.Push(token.NewI(i, pos))
			run(m2)
		}
		return
	}

	for i := begin; i > end; i += step {
		m2 := machine.New(m.SourceDir, m.Pmachines, p)
		m2.Push(token.NewI(i, pos))
		run(m2)
	}
}

var mutex sync.Mutex

// Synchronization of procedures.
//    m: Virtual machine.
//    run: Function which running a machine.
func prSync(m *machine.T, run func(m *machine.T)) {
	tk1 := m.PopT(token.Procedure)
	p, _ := tk1.P()

	mutex.Lock()
	run(machine.New(m.SourceDir, m.Pmachines, p))
	mutex.Unlock()
}

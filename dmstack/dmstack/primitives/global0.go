// Copyright 06-Jan-2021 ºDeme
// GNU General Public License - V3 <http://www.gnu.org/licenses/>

package primitives

import (
	"fmt"
	"github.com/dedeme/dmstack/imports"
	"github.com/dedeme/dmstack/machine"
	"github.com/dedeme/dmstack/reader"
	"github.com/dedeme/dmstack/stack"
	"github.com/dedeme/dmstack/symbol"
	"github.com/dedeme/dmstack/token"
	"io/ioutil"
	"sync"
)

// Execute a procedure saved in stack.
//    m: Virtual machine.
//    run: Function which running a machine.
func prRun(m *machine.T, run func(m *machine.T)) {
	tk := m.PopT(token.Procedure)
	run(machine.New(m.Source, m.Pmachines, tk))
}

// Process an 'import'.
//    m: Virtual machine.
//    run: Function which running a machine.
func prImport(m *machine.T, run func(m *machine.T)) {
	tk := m.PopT(token.Symbol)
	sym, _ := tk.Sy()
	if heap, ok := imports.Get(sym); ok {
		if heap == nil {
			bs, err := ioutil.ReadFile(sym.String() + ".dms")
			if err != nil {
				m.Fail(
					machine.EImport(),
					"File '%v' can not be read", sym.String()+".dms",
				)
			}

			rd := reader.New(sym, string(bs))
			tkP := rd.Process()
			_, ok := tkP.P()
			if !ok {
				m.Fail(
					machine.EImport(),
					"File '%v' is not a valid program", sym.String()+".dms",
				)
			}
			if rd.LastChar() != "" {
				m.Fail(
					machine.EImport(),
					"Unexpected end of program (%v) in %v",
					rd.LastChar(),
					sym.String()+".dms",
				)
			}

			imports.Initialize(sym)
			m2 := machine.NewThread(sym, []*machine.T{}, tkP)
			run(m2)
			for k, v := range m2.Heap {
				imports.AddKey(sym, k, v)
			}
		}
		return
	}
	// sym should be registred by the reader.
	panic(fmt.Sprintf("Symbol '%v' not found", sym))
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
				run(machine.New(m.Source, m.Pmachines, tk))
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
		m2 := machine.New(m.Source, m.Pmachines, tk1)
		run(m2)
		l := len(*m.Stack) - 1
		if l < 0 || (*m.Stack)[l].Type() != token.Bool {
			m.Failt(
				"\n  Expected: Procedure with a Bool return."+
					"\n  Actual  : '%v'.",
				tk1.StringDraft(),
			)
		}

		tk11 := m.Pop()
		b, _ := tk11.B()
		if b {
			run(machine.New(m.Source, m.Pmachines, tk2))
			return
		}
	}
	if p != nil {
		run(machine.New(m.Source, m.Pmachines, tk))
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
		run(machine.New(m.Source, m.Pmachines, tk2))
		return
	}
	run(machine.New(m.Source, m.Pmachines, tk1))
}

// Process a loop.
//    m: Virtual machine.
//    run: Function which running a machine.
func prLoop(m *machine.T, run func(m *machine.T)) {
	tk := m.PopT(token.Procedure)
	for {
		run(machine.New(m.Source, m.Pmachines, tk))
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
	for {
		m2 := machine.New(m.Source, m.Pmachines, tk2)
		run(m2)
		l := len(*m.Stack) - 1
		if l < 0 || (*m.Stack)[l].Type() != token.Bool {
			m.Failt(
				"\n  Expected: Procedure with a Bool return."+
					"\n  Actual  : '%v'.",
				tk2.StringDraft(),
			)
		}

		tk21 := m.Pop()
		b, _ := tk21.B()
		if !b {
			break
		}

		run(machine.New(m.Source, m.Pmachines, tk1))
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
	tk2 := m.Pop()
	step := int64(1)
	begin := int64(0)
	end, ok := tk2.I()
	if !ok {
		a, ok := tk2.A()
		if !ok {
			m.Failt(
				"\n  Expected: Value of type Int or List."+
					"\n  Actual  : %v.",
				tk2.StringDraft(),
			)
		}
		ok = len(a) == 2 || len(a) == 3
		var ok1, ok2, ok3 bool
		if ok {
			begin, ok1 = a[0].I()
			end, ok2 = a[1].I()
			ok3 = true
			if len(a) == 3 {
				step, ok3 = a[2].I()
			}
		}
		if !(ok && ok1 && ok2 && ok3) {
			m.Failt(
				"\n  Expected: [Int, Int] or [Int, Int, Int]."+
					"\n  Actual  : %v.",
				tk2.StringDraft(),
			)
		}
		if step == 0 {
			m.Fail(machine.EMachine(), "For step can not be 0")
		}
	}

	pos := m.MkPos()
	if step > 0 {
		for i := begin; i < end; i += step {
			m2 := machine.New(m.Source, m.Pmachines, tk1)
			m2.Push(token.NewI(i, pos))
			run(m2)
		}
		return
	}

	for i := begin; i > end; i += step {
		m2 := machine.New(m.Source, m.Pmachines, tk1)
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
	mutex.Lock()
	run(machine.New(m.Source, m.Pmachines, tk1))
	mutex.Unlock()
}

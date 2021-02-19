// Copyright 07-Jan-2021 ºDeme
// GNU General Public License - V3 <http://www.gnu.org/licenses/>

package primitives

import (
	"github.com/dedeme/dmstack/machine"
	"github.com/dedeme/dmstack/operator"
	"github.com/dedeme/dmstack/stack"
	"github.com/dedeme/dmstack/symbol"
	"github.com/dedeme/dmstack/token"
)

func prAssign(m *machine.T) {
	var syms []symbol.T
	for {
		tk, ok := m.PrgNext()
		if !ok {
			m.Fail(
        machine.EMachine(),
        "Bad assignation variable or ',' missing in literal map",
      )
		}
		if sym, ok := tk.Sy(); ok {
			syms = append(syms, sym)
			continue
		}
		op, ok := tk.O()
		if ok && op == operator.Equals {
			if len(syms) == 0 {
				m.Fail(machine.EMachine(), "There are no declared variables")
			}

			for i, j := 0, len(syms)-1; i < j; i, j = i+1, j-1 {
				syms[i], syms[j] = syms[j], syms[i]
			}

			for _, k := range syms {
				if k.Reserved() {
					m.Fail(machine.EMachine(), "Symbol '%v' is reserved", k)
				}
				if _, ok := m.Heap[k]; ok {
					m.Fail(machine.EMachine(), "Symbol '%v' is duplicated", k)
				}
				m.Heap[k] = m.Pop()
			}
			break
		}
		m.Failt(
			"\n  Expected: token of type 'Symbol' or '='.\n  Actual  : '%v'.",
			tk.StringDraft(),
		)
	}
}

func prProcHeap(m *machine.T) {
	tk := m.PopT(token.Procedure)
	tk.SetHeap(m.Heap)
	m.Push(tk)
}

func prStackCheck(m *machine.T) {
	tk := m.PopT(token.String)
	s, _ := tk.S()
	_, err := stack.TypesOk(*m.Stack, s)
	if err == nil {
		m.Push(token.NewB(true, m.MkPos()))
		return
	}
	m.Push(token.NewB(false, m.MkPos()))
}

func prStackOpen(m *machine.T) {
	tk := m.PopT(token.String)
	s, _ := tk.S()
	n, err := stack.TypesOk(*m.Stack, s)
	if err != nil {
		m.Fail(machine.EStack(), err.Error())
	}
	var st2 []*token.T
	for i := 0; i < n; i++ {
		st2 = append(st2, m.Pop())
	}
	m.Push(token.NewO(operator.StackStop, m.MkPos()))
	for i := n - 1; i >= 0; i-- {
		m.Push(st2[i])
	}
}

func prStackClose(m *machine.T) {
	tk := m.PopT(token.String)
	s, _ := tk.S()
	n, err := stack.StopTypesOk(*m.Stack, s)
	if err != nil {
		m.Fail(machine.EStack(), err.Error())
	}
	var st2 []*token.T
	for i := 0; i < n; i++ {
		st2 = append(st2, m.Pop())
	}
	m.Pop()
	for i := n - 1; i >= 0; i-- {
		m.Push(st2[i])
	}
}

func prStack(m *machine.T) {
	tk := m.PopT(token.String)
	s, _ := tk.S()
	_, err := stack.TypesOk(*m.Stack, s)
	if err != nil {
		m.Fail(machine.EStack(), err.Error())
	}
}

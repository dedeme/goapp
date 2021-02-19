// Copyright 07-Jan-2021 ºDeme
// GNU General Public License - V3 <http://www.gnu.org/licenses/>

package primitives

import (
	"github.com/dedeme/dmstack/machine"
	"github.com/dedeme/dmstack/token"
)

// Returns the element of list at postion ix.
//    m: Virtual machine.
func prRefGet(m *machine.T) {
	tk := m.PopT(token.Array)
	a, _ := tk.A()
	if len(a) != 1 {
		m.Failt(
			"\n  Expected: Reference with one element.\n  Actual  : '%v'.",
			tk.StringDraft(),
		)
	}
	m.Push(a[0])
}

// Sets the element of list at postion ix.
//    m: Virtual machine.
func prRefSet(m *machine.T) {
	tkV := m.Pop()
	tk := m.PopT(token.Array)
	a, _ := tk.A()
	if len(a) != 1 {
		m.Failt(
			"\n  Expected: Reference with one element.\n  Actual  : '%v'.",
			tk.StringDraft(),
		)
	}
	a[0] = tkV
}

//  Updates the element of list at position ix.
//    m: Virtual machine.
//    run: Function which running a machine.
func prRefUp(m *machine.T, run func(m *machine.T)) {
	tkP := m.PopT(token.Procedure)
	tk := m.PopT(token.Array)
	a, _ := tk.A()
	if len(a) != 1 {
		m.Failt(
			"\n  Expected: Reference with one element.\n  Actual  : '%v'.",
			tk.StringDraft(),
		)
	}
	m2 := machine.New(m.Source, m.Pmachines, tkP)
	m2.Push(a[0])
	run(m2)
	if len(*m2.Stack) == 0 {
		m.Failt(
			"\n  Expected: Prodedure which returns one element.\n  Actual  : '%v'.",
			tkP.StringDraft(),
		)
	}
	a[0] = m2.Pop()
}

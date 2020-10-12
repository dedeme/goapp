// Copyright 16-May-2020 ºDeme
// GNU General Public License - V3 <http://www.gnu.org/licenses/>

package primitives

import (
	"github.com/dedeme/dmstack/machine"
	"github.com/dedeme/dmstack/token"
)

// Returns the element of list at postion ix.
//    m: Virtual machine.
func prRefGet(m *machine.T) {
	tk := m.PopT(token.List)
	l, _ := tk.L()
	if len(l) != 1 {
		m.Failt(
			"\n  Expected: Reference with one element.\n  Actual  : '%v'.",
			tk.StringDraft(),
		)
	}
	m.Push(l[0])
}

// Sets the element of list at postion ix.
//    m: Virtual machine.
func prRefSet(m *machine.T) {
	tkV := m.Pop()
	tk := m.PopT(token.List)
	l, _ := tk.L()
	if len(l) != 1 {
		m.Failt(
			"\n  Expected: Reference with one element.\n  Actual  : '%v'.",
			tk.StringDraft(),
		)
	}
	l[0] = tkV
}

//  Updates the element of list at position ix.
//    m: Virtual machine.
//    run: Function which running a machine.
func prRefUp(m *machine.T, run func(m *machine.T)) {
	tkP := m.PopT(token.Procedure)
	p, _ := tkP.P()
	tk := m.PopT(token.List)
	l, _ := tk.L()
	if len(l) != 1 {
		m.Failt(
			"\n  Expected: Reference with one element.\n  Actual  : '%v'.",
			tk.StringDraft(),
		)
	}
	m2 := machine.NewIsolate(m.SourceDir, m.Pmachines, p)
	m2.Push(l[0])
	run(m2)
	l[0] = m2.Pop()
}

// Copyright 22-Jul-2020 ºDeme
// GNU General Public License - V3 <http://www.gnu.org/licenses/>

package list

import (
	"github.com/dedeme/dmstack/machine"
	"github.com/dedeme/dmstack/token"
)

// Creates a reference.
func prRef(m *machine.T) {
	tk := m.Pop()
	m.Push(token.NewL([]*token.T{tk}, m.MkPos()))
}

// Creates a tuple of 2 elements.
func prTp(m *machine.T) {
	tk1 := m.Pop()
	tk2 := m.Pop()
	m.Push(token.NewL([]*token.T{tk2, tk1}, m.MkPos()))
}

// Creates a tuple of 3 elements.
func prTp3(m *machine.T) {
	tk1 := m.Pop()
	tk2 := m.Pop()
	tk3 := m.Pop()
	m.Push(token.NewL([]*token.T{tk3, tk2, tk1}, m.MkPos()))
}

// Set the first element of a tuple.
func prE1(m *machine.T) {
	tk0 := m.Pop()
	tk := m.PopT(token.List)
	l, _ := tk.L()
	if len(l) < 1 {
		m.Fail("Index out of range error", "0")
	}
	l[0] = tk0
}

// Set the second element of a tuple.
func prE2(m *machine.T) {
	tk1 := m.Pop()
	tk := m.PopT(token.List)
	l, _ := tk.L()
	if len(l) < 2 {
		m.Fail("Index out of range error", "1")
	}
	l[1] = tk1
}

// Set the third element of a tuple.
func prE3(m *machine.T) {
	tk2 := m.Pop()
	tk := m.PopT(token.List)
	l, _ := tk.L()
	if len(l) < 3 {
		m.Fail("Index out of range error", "2")
	}
	l[2] = tk2
}

// Creates a Either-Left
func prLeft(m *machine.T) {
	tk := m.Pop()
	m.Push(token.NewL([]*token.T{token.NewB(false, m.MkPos()), tk}, m.MkPos()))
}

// Creates a Either-Right
func prRight(m *machine.T) {
	tk := m.Pop()
	m.Push(token.NewL([]*token.T{token.NewB(true, m.MkPos()), tk}, m.MkPos()))
}

// Creates a Result-Error
func prError(m *machine.T) {
	tk := m.PopT(token.String)
	m.Push(token.NewL([]*token.T{token.NewB(false, m.MkPos()), tk}, m.MkPos()))
}

// Creates a Result-Ok
func prOk(m *machine.T) {
	tk := m.Pop()
	m.Push(token.NewL([]*token.T{token.NewB(true, m.MkPos()), tk}, m.MkPos()))
}

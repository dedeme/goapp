// Copyright 07-Jan-2021 ºDeme
// GNU General Public License - V3 <http://www.gnu.org/licenses/>

package primitives

import (
	"fmt"
	"github.com/dedeme/dmstack/machine"
	"github.com/dedeme/dmstack/token"
)

// Shows and consume a draft of the last element of stack.
//    m: Virtual machine.
func prPuts(m *machine.T) {
	tk := m.Pop()
	fmt.Println(tk.StringDraft())
}

// Push the value in stack converted to string.
//    m: Virtual machine.
func prToString(m *machine.T) {
	m.Push(token.NewS(m.Pop().String(), m.MkPos()))
}

// Push a new Token which is a deep copy of the last value of stack.
//    If token is of type Native, the native object is sallowly copied.
func prClone(m *machine.T) {
	m.Push(m.Pop().Clone())
}

// Raise an "Assert error" if the last element of stack is not a Bool or it
// is 'false'.
//    m: Virtual machine.
func prAssert(m *machine.T) {
	okV, ok := m.PopT(token.Bool).B()
	if ok {
		if okV {
			return
		}
		m.Fail(machine.EAssert(), "Assert failed")
	}
	m.Fail(machine.EAssert(), "Assert value is not a Bool value")
}

// Raise a "Expect error" if the two last elements of stack are not equals.
//    m: Virtual machine.
func prExpect(m *machine.T) {
	expected := m.Pop()
	actual := m.Pop()
	if expected.Eq(actual) {
		return
	}
	m.Fail(
		machine.EExpect(), "\n  Expected: '%s'.\n  Actual  : '%s'.",
		expected.StringDraft(), actual.StringDraft(),
	)
}

// Generate a generic panic response.
//    m: Virtual machine.
func prFail(m *machine.T) {
	tk := m.PopT(token.String)
	s, _ := tk.S()
	panic(&machine.Error{m, machine.ERuntime(), s})
}

// Generate a panic response.
//    m: Virtual machine.
func prThrow(m *machine.T) {
	tk2 := m.PopT(token.String)
	s2, _ := tk2.S()
	tk1 := m.PopT(token.String)
	s1, _ := tk1.S()
	panic(&machine.Error{m, machine.ECustom(s1), s2})
}

// Run a procedure which can fail, and recover it in such case.
//    m: Virtual machine.
//    run: Function which running a machine.
func prTry(m *machine.T, run func(m *machine.T)) {
	catch := m.PopT(token.Procedure)
	tk2 := m.PopT(token.Procedure)
	stack := append([]*token.T{}, *m.Stack...)

	defer func() {
		if err := recover(); err != nil {
			m.Stack = &stack
			e := err.(*machine.Error)
			m2 := machine.New(m.Source, m.Pmachines, catch)
			m2.Push(token.NewS(e.Type.String()+":"+e.Message, m2.MkPos()))
			run(m2)
		}
	}()

	m2 := machine.New(m.Source, m.Pmachines, tk2)
	run(m2)
}

// Interchanges the two last values of stack.
//    m: Virtual machine.
func prSwap(m *machine.T) {
	tk1 := m.Pop()
	tk2 := m.Pop()
	m.Push(tk1)
	m.Push(tk2)
}

// Removes the last element of stack.
//    m: Virtual machine.
func prPop(m *machine.T) {
	m.Pop()
}

// Duplicates (not clone, only reference) the last element of stack.
//    m: Virtual machine.
func prDup(m *machine.T) {
	tk := m.Pop()
	m.Push(tk)
	m.Push(tk)
}

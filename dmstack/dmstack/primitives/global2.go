// Copyright 07-Jan-2021 ºDeme
// GNU General Public License - V3 <http://www.gnu.org/licenses/>

package primitives

import (
	"github.com/dedeme/dmstack/machine"
	"github.com/dedeme/dmstack/stack"
	"github.com/dedeme/dmstack/token"
)

// Function applyable to
//    Int Int -> Int
//    Int Float -> Float
//    Float Int -> Float
//    Float Float -> Float
//    String String -> String
// Parameter:
//    m : Virtual machine.
func prPlus(m *machine.T) {
	tk1 := m.Pop()
	tk2 := m.Pop()
	i1, ok := tk1.I()
	if ok {
		i2, ok := tk2.I()
		if ok {
			m.Push(token.NewI(i2+i1, m.MkPos()))
			return
		}
		f2, ok := tk2.F()
		if ok {
			m.Push(token.NewF(f2+float64(i1), m.MkPos()))
			return
		}
		m.Failt(
			"\n  Expected: token of type 'Int' or 'Float'\n  Actual  : '%v'.",
			tk2.StringDraft(),
		)
	}
	f1, ok := tk1.F()
	if ok {
		i2, ok := tk2.I()
		if ok {
			m.Push(token.NewF(float64(i2)+f1, m.MkPos()))
			return
		}
		f2, ok := tk2.F()
		if ok {
			m.Push(token.NewF(f2+f1, m.MkPos()))
			return
		}
		m.Failt(
			"\n  Expected: token of type 'Int' or 'Float'.\n  Actual  : '%v'.",
			tk2.StringDraft(),
		)
	}
	s1, ok := tk1.S()
	if ok {
		s2, ok := tk2.S()
		if ok {
			m.Push(token.NewS(s2+s1, m.MkPos()))
			return
		}
		m.Failt(
			"\n  Expected: token of type 'String'.\n  Actual  : '%v'.",
			tk2.StringDraft(),
		)
	}

	a1, ok := tk1.A()
	if ok {
		a2, ok := tk2.A()
		if ok {
			for _, e := range a1 {
				a2 = append(a2, e)
			}
			m.Push(token.NewA(a2, m.MkPos()))
			return
		}
		m.Failt(
			"\n  Expected: token of type 'List'.\n  Actual  : '%v'.",
			tk2.StringDraft(),
		)
	}

	m.Failt(
		"\n  Expected: token of type 'Int', 'Float', 'String' or 'List'.\n"+
			"  Actual  : '%v'.",
		tk1.StringDraft(),
	)
}

// Function applyable to
//    Int Int -> Int
//    Int Float -> Float
//    Float Int -> Float
//    Float Float -> Float
//    String String -> String
// Parameter:
//    m : Virtual machine.
func prMinus(m *machine.T) {
	tk1 := m.Pop()
	tk2 := m.Pop()
	i1, ok := tk1.I()
	if ok {
		i2, ok := tk2.I()
		if ok {
			m.Push(token.NewI(i2-i1, m.MkPos()))
			return
		}
		f2, ok := tk2.F()
		if ok {
			m.Push(token.NewF(f2-float64(i1), m.MkPos()))
			return
		}
		m.Failt(
			"\n  Expected: token of type 'Int' or 'Float'.\n  Actual  : '%v'.",
			tk2.StringDraft(),
		)
	}
	f1, ok := tk1.F()
	if ok {
		i2, ok := tk2.I()
		if ok {
			m.Push(token.NewF(float64(i2)-f1, m.MkPos()))
			return
		}
		f2, ok := tk2.F()
		if ok {
			m.Push(token.NewF(f2-f1, m.MkPos()))
			return
		}
		m.Failt(
			"\n  Expected: token of type 'Int' or 'Float'.\n  Actual  : '%v'.",
			tk2.StringDraft(),
		)
	}
	m.Failt(
		"\n  Expected: token of type 'Int' or 'Float'.\n  Actual  : '%v'.",
		tk1.StringDraft(),
	)
}

// Function applyable to
//    Int Int -> Int
//    Int Float -> Float
//    Float Int -> Float
//    Float Float -> Float
// Parameter:
//    m : Virtual machine.
func prMul(m *machine.T) {
	tk1 := m.Pop()
	tk2 := m.Pop()
	i1, ok := tk1.I()
	if ok {
		i2, ok := tk2.I()
		if ok {
			m.Push(token.NewI(i2*i1, m.MkPos()))
			return
		}
		f2, ok := tk2.F()
		if ok {
			m.Push(token.NewF(f2*float64(i1), m.MkPos()))
			return
		}
		m.Failt(
			"\n  Expected: token of type 'Int' or 'Float'.\n  Actual  : '%v'.",
			tk2.StringDraft(),
		)
	}
	f1, ok := tk1.F()
	if ok {
		i2, ok := tk2.I()
		if ok {
			m.Push(token.NewF(float64(i2)*f1, m.MkPos()))
			return
		}
		f2, ok := tk2.F()
		if ok {
			m.Push(token.NewF(f2*f1, m.MkPos()))
			return
		}
		m.Failt(
			"\n  Expected: token of type 'Int' or 'Float'.\n  Actual  : '%v'.",
			tk2.StringDraft(),
		)
	}
	m.Failt(
		"\n  Expected: token of type 'Int' or 'Float'.\n  Actual  : '%v'.",
		tk1.StringDraft(),
	)
}

// Function applyable to
//    Int Int -> Int
//    Int Float -> Float
//    Float Int -> Float
//    Float Float -> Float
// Parameter:
//    m : Virtual machine.
func prDiv(m *machine.T) {
	tk1 := m.Pop()
	tk2 := m.Pop()
	i1, ok := tk1.I()
	if ok {
		i2, ok := tk2.I()
		if ok {
			m.Push(token.NewI(i2/i1, m.MkPos()))
			return
		}
		f2, ok := tk2.F()
		if ok {
			m.Push(token.NewF(f2/float64(i1), m.MkPos()))
			return
		}
		m.Failt(
			"\n  Expected: token of type 'Int' or 'Float'.\n  Actual  : '%v'.",
			tk2.StringDraft(),
		)
	}
	f1, ok := tk1.F()
	if ok {
		i2, ok := tk2.I()
		if ok {
			m.Push(token.NewF(float64(i2)/f1, m.MkPos()))
			return
		}
		f2, ok := tk2.F()
		if ok {
			m.Push(token.NewF(f2/f1, m.MkPos()))
			return
		}
		m.Failt(
			"\n  Expected: token of type 'Int' or 'Float'.\n  Actual  : '%v'.",
			tk2.StringDraft(),
		)
	}
	m.Failt(
		"\n  Expected: token of type 'Int' or 'Float'.\n  Actual  : '%v'.",
		tk1.StringDraft(),
	)
}

// Function applyable to
//    Int Int -> Int
//    Int Float -> Float
//    Float Int -> Float
//    Float Float -> Float
//    String String -> String
// Parameter:
//    m : Virtual machine.
func prMod(m *machine.T) {
	tk1 := m.PopT(token.Int)
	i1, _ := tk1.I()
	tk2 := m.PopT(token.Int)
	i2, _ := tk2.I()
	m.Push(token.NewI(i2%i1, m.MkPos()))
}

// Function applyable to
//    Int -> Int
//    String .... -> String
//    List .... -> List
func prPlusPlus(m *machine.T) {
	tk := m.Pop()
	i, ok := tk.I()
	if ok {
		m.Push(token.NewI(i+1, m.MkPos()))
		return
	}
	s, ok := tk.S()
	if ok {
		st := m.Stack
		for {
			tk2, ok := stack.Peek(*st)
			if ok {
				s2, ok := tk2.S()
				if ok {
					s = s2 + s
					m.Pop()
					continue
				}
			}
			break
		}
		m.Push(token.NewS(s, m.MkPos()))
		return
	}
	a, ok := tk.A()
	if ok {
		st := m.Stack
		for {
			tk2, ok := stack.Peek(*st)
			if ok {
				a2, ok := tk2.A()
				if ok {
					for _, e := range a {
						a2 = append(a2, e)
					}
					a = a2
					m.Pop()
					continue
				}
			}
			break
		}
		m.Push(token.NewA(a, m.MkPos()))
		return
	}
	m.Failt(
		"\n  Expected: token of type 'Int', 'String' or 'List'.\n  Actual  : '%v'.",
		tk.StringDraft(),
	)
}

// Function applyable to
//    Int -> Int
func prMinusMinus(m *machine.T) {
	tk := m.PopT(token.Int)
	i, _ := tk.I()
	m.Push(token.NewI(i-1, m.MkPos()))
}

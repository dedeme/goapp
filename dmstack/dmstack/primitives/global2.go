// Copyright 16-May-2020 ºDeme
// GNU General Public License - V3 <http://www.gnu.org/licenses/>

package primitives

import (
	"github.com/dedeme/dmstack/machine"
	"github.com/dedeme/dmstack/token"
	"github.com/dedeme/dmstack/stack"
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
		m.Failf(
			"Stack:\n  Expected token of type 'Int' or 'Float'\n  Found %v.",
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
		m.Failf(
			"Stack:\n  Expected token of type 'Int' or 'Float'\n  Found %v.",
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
		m.Failf(
			"Stack:\n  Expected token of type 'String'\n  Found %v.",
			tk2.StringDraft(),
		)
	}

	l1, ok := tk1.L()
	if ok {
		l2, ok := tk2.L()
		if ok {
      for _, e := range l1 {
        l2 = append(l2, e)
      }
			m.Push(token.NewL(l2, m.MkPos()))
			return
		}
		m.Failf(
			"Stack:\n  Expected token of type 'List'\n  Found %v.",
			tk2.StringDraft(),
		)
	}

	m.Failf(
		"Stack:\n  Expected token of type 'Int', 'Float', 'String' or 'List'\n" +
    "  Found %v.",
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
		m.Failf(
			"Stack:\n  Expected token of type 'Int' or 'Float'\n  Found %v.",
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
		m.Failf(
			"Stack:\n  Expected token of type 'Int' or 'Float'\n  Found %v.",
			tk2.StringDraft(),
		)
	}
	m.Failf(
		"Stack:\n  Expected token of type 'Int' or 'Float'\n  Found %v.",
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
		m.Failf(
			"Stack:\n  Expected token of type 'Int' or 'Float'\n  Found %v.",
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
		m.Failf(
			"Stack:\n  Expected token of type 'Int' or 'Float'\n  Found %v.",
			tk2.StringDraft(),
		)
	}
	m.Failf(
		"Stack:\n  Expected token of type 'Int' or 'Float'\n  Found %v.",
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
		m.Failf(
			"Stack:\n  Expected token of type 'Int' or 'Float'\n  Found %v.",
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
		m.Failf(
			"Stack:\n  Expected token of type 'Int' or 'Float'\n  Found %v.",
			tk2.StringDraft(),
		)
	}
	m.Failf(
		"Stack:\n  Expected token of type 'Int' or 'Float'\n  Found %v.",
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
    m.Push(token.NewI(i + 1, m.MkPos()))
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
  l, ok := tk.L()
  if ok {
    st := m.Stack
    for {
      tk2, ok := stack.Peek(*st)
      if ok {
        l2, ok := tk2.L()
        if ok {
          for _, e := range l {
            l2 = append(l2, e)
          }
          l = l2
          m.Pop()
          continue
        }
      }
      break
    }
    m.Push(token.NewL(l, m.MkPos()))
    return
  }
	m.Failf(
		"Stack:\n  Expected token of type 'Int', 'String' or 'List'\n  Found %v.",
		tk.StringDraft(),
	)
}

// Function applyable to
//    Int -> Int
func prMinusMinus(m *machine.T) {
  tk := m.PopT(token.Int)
  i, _ := tk.I()
  m.Push(token.NewI(i - 1, m.MkPos()))
}

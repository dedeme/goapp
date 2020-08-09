// Copyright 16-May-2020 ºDeme
// GNU General Public License - V3 <http://www.gnu.org/licenses/>

package primitives

import (
	"github.com/dedeme/dmstack/machine"
	"github.com/dedeme/dmstack/token"
)

// Pushes true if tk1 is equals to tk2.
//    * * -> Bool
// Parameter:
//    m : Virtual machine.
func prEq(m *machine.T) {
	tk1 := m.Pop()
	tk2 := m.Pop()
	m.Push(token.NewB(tk1.Eq(tk2), m.MkPos()))
}

// Pushes false if tk1 is equals to tk2.
//    * * -> Bool
// Parameter:
//    m : Virtual machine.
func prNeq(m *machine.T) {
	tk1 := m.Pop()
	tk2 := m.Pop()
	m.Push(token.NewB(!tk1.Eq(tk2), m.MkPos()))
}

// Pushes true if tk2 is less than tk1.
//    Bool Bool -> Bool
//    Int Int -> Bool
//    Float Float -> Bool
//    String String -> Bool
// Parameter:
//    m : Virtual machine.
func prLess(m *machine.T) {
	tk1 := m.Pop()
	b1, ok := tk1.B()
	if ok {
		tk2 := m.PopT(token.Bool)
		b2, _ := tk2.B()
		r := false
		if b1 && !b2 {
			r = true
		}
		m.Push(token.NewB(r, m.MkPos()))
		return
	}
	i1, ok := tk1.I()
	if ok {
		tk2 := m.PopT(token.Int)
		i2, _ := tk2.I()
		r := false
		if i2 < i1 {
			r = true
		}
		m.Push(token.NewB(r, m.MkPos()))
		return
	}
	f1, ok := tk1.F()
	if ok {
		tk2 := m.PopT(token.Float)
		f2, _ := tk2.F()
		r := false
		if f2 < f1 {
			r = true
		}
		m.Push(token.NewB(r, m.MkPos()))
		return
	}
	s1, ok := tk1.S()
	if ok {
		tk2 := m.PopT(token.String)
		s2, _ := tk2.S()
		r := false
		if s2 < s1 {
			r = true
		}
		m.Push(token.NewB(r, m.MkPos()))
		return
	}
	m.Failf(
		"Stack:\n  Expected token of type 'Bool', 'Int', 'Float' or 'String'\n"+
			"  Found %v.",
		tk1.StringDraft(),
	)
}

// Pushes true if tk2 is less or equals to tk1.
//    Bool Bool -> Bool
//    Int Int -> Bool
//    Float Float -> Bool
//    String String -> Bool
// Parameter:
//    m : Virtual machine.
func prLessEq(m *machine.T) {
	tk1 := m.Pop()
	tk2 := m.Pop()
	tp := tk1.Type()
	if tp == token.Bool || tp == token.Int ||
		tp == token.Float || tp == token.String {
		if tk1.Eq(tk2) {
			m.Push(token.NewB(true, m.MkPos()))
			return
		}
		m.Push(tk2)
		m.Push(tk1)
		prLess(m)
		return
	}
	m.Failf(
		"Stack:\n  Expected token of type 'Bool', 'Int', 'Float' or 'String'\n"+
			"  Found %v.",
		tk1.StringDraft(),
	)
}

// Pushes true if tk2 is greater than tk1.
//    Bool Bool -> Bool
//    Int Int -> Bool
//    Float Float -> Bool
//    String String -> Bool
// Parameter:
//    m : Virtual machine.
func prGreater(m *machine.T) {
	tk1 := m.Pop()
	b1, ok := tk1.B()
	if ok {
		tk2 := m.PopT(token.Bool)
		b2, _ := tk2.B()
		r := false
		if b2 && !b1 {
			r = true
		}
		m.Push(token.NewB(r, m.MkPos()))
		return
	}
	i1, ok := tk1.I()
	if ok {
		tk2 := m.PopT(token.Int)
		i2, _ := tk2.I()
		r := false
		if i2 > i1 {
			r = true
		}
		m.Push(token.NewB(r, m.MkPos()))
		return
	}
	f1, ok := tk1.F()
	if ok {
		tk2 := m.PopT(token.Float)
		f2, _ := tk2.F()
		r := false
		if f2 > f1 {
			r = true
		}
		m.Push(token.NewB(r, m.MkPos()))
		return
	}
	s1, ok := tk1.S()
	if ok {
		tk2 := m.PopT(token.String)
		s2, _ := tk2.S()
		r := false
		if s2 > s1 {
			r = true
		}
		m.Push(token.NewB(r, m.MkPos()))
		return
	}
	m.Failf(
		"Stack:\n  Expected token of type 'Bool', 'Int', 'Float' or 'String'\n"+
			"  Found %v.",
		tk1.StringDraft(),
	)
}

// Pushes true if tk2 is greater or equals to tk1.
//    Bool Bool -> Bool
//    Int Int -> Bool
//    Float Float -> Bool
//    String String -> Bool
// Parameter:
//    m : Virtual machine.
func prGreaterEq(m *machine.T) {
	tk1 := m.Pop()
	tp := tk1.Type()
	if tp == token.Bool || tp == token.Int ||
		tp == token.Float || tp == token.String {
		tk2 := m.Pop()
		if tk1.Eq(tk2) {
			m.Push(token.NewB(true, m.MkPos()))
			return
		}
		m.Push(tk2)
		m.Push(tk1)
		prGreater(m)
		return
	}
	m.Failf(
		"Stack:\n  Expected token of type 'Bool', 'Int', 'Float' or 'String'\n"+
			"  Found %v.",
		tk1.StringDraft(),
	)
}

// Pushes true if tk1 is true (and/or) tk2 is true
//    Bool Bool -> Bool
//    Procedure Procedure -> Bool
// Parameter:
//    m : Virtual machine.
//    run: Function which running a machine.
//    isAnd: It is true if function to process is 'And'
func prAndOr(m *machine.T, run func(m *machine.T), isAnd bool) {
	tk1 := m.Pop()
	b1, ok := tk1.B()
	if ok {
		tk2 := m.PopT(token.Bool)
		b2, _ := tk2.B()
		if isAnd {
			m.Push(token.NewB(b1 && b2, m.MkPos()))
			return
		}
		m.Push(token.NewB(b1 || b2, m.MkPos()))
		return
	}
	p1, ok := tk1.P()
	if ok {
		tk2 := m.PopT(token.Procedure)
		p2, _ := tk2.P()

		m2 := machine.NewIsolate(m.SourceDir, m.Pmachines, p2)
		run(m2)
		okReturn := false
		if len(*m2.Stack) == 1 {
			b1, ok = (*m2.Stack)[0].B()
			if ok {
				okReturn = true
			}
		}
		if !okReturn {
			m.Failf("Expected a Bool return:\b%v", tk1.StringDraft())
		}
		if (isAnd && !b1) || (!isAnd && b1) {
			m.Push(token.NewB(b1, m.MkPos()))
			return
		}

		m1 := machine.NewIsolate(m.SourceDir, m.Pmachines, p1)
		run(m1)
		if len(*m1.Stack) == 1 {
			b2, ok := (*m1.Stack)[0].B()
			if ok {
				if isAnd {
					m.Push(token.NewB(b1 && b2, m.MkPos()))
					return
				}
				m.Push(token.NewB(b1 || b2, m.MkPos()))
				return
			}
		}

		m.Failf("Expected a Bool return:\b%v", tk2.StringDraft())
	}
	m.Failf(
		"Stack:\n  Expected token of type 'Bool' or 'Procedure'\n"+
			"  Found %v.",
		tk1.StringDraft(),
	)
}

// Pushes true if tk1 is true and tk2 is true
//    Bool Bool -> Bool
//    Procedure Procedure -> Bool
// Parameter:
//    m : Virtual machine.
//    run: Function which running a machine.
func prAnd(m *machine.T, run func(m *machine.T)) {
	prAndOr(m, run, true)
}

// Pushes true if tk1 is true or tk2 is true
//    Bool Bool -> Bool
//    Procedure Procedure -> Bool
// Parameter:
//    m : Virtual machine.
//    run: Function which running a machine.
func prOr(m *machine.T, run func(m *machine.T)) {
	prAndOr(m, run, false)
}

// Pushes true if tk1 is false and false if it is true
//    Bool -> Bool
// Parameter:
//    m : Virtual machine.
//    run: Function which running a machine.
func prNot(m *machine.T) {
	tk := m.PopT(token.Bool)
	b, _ := tk.B()
	m.Push(token.NewB(!b, m.MkPos()))
}

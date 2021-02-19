// Copyright 10-Aug-2020 ºDeme
// GNU General Public License - V3 <http://www.gnu.org/licenses/>

package it

import (
	"github.com/dedeme/dmstack/machine"
	"github.com/dedeme/dmstack/token"
	"io"
	"strings"
)

// Creates a new It from procedures 'has', 'next'.
//    m: Virtual machine.
//    run: Function which running a machine.
func prNew(m *machine.T, run func(m *machine.T)) {
	tk1 := m.PopT(token.Procedure)
	nextf := func() (tk *token.T, ok bool) {
		m2 := machine.New(m.Source, m.Pmachines, tk1)
		run(m2)
		okTk := m2.PopT(token.Bool)
		ok, _ = okTk.B()
		tk = m2.Pop()
		return
	}
	pushIt(m, New(nextf))
}

// Creates a new empty It.
func prEmpty(m *machine.T) {
	next := func() (tk *token.T, ok bool) {
		return
	}
	pushIt(m, New(next))
}

// Creates a new It with one element.
func prUnary(m *machine.T) {
	tk1 := m.Pop()
	more := true
	next := func() (tk *token.T, ok bool) {
		if more {
			tk = tk1
			ok = true
			more = false
		}
		return
	}
	pushIt(m, New(next))
}

// Creates a new It from a List.
func PrFrom(m *machine.T) {
	tk := m.PopT(token.Array)
	a, _ := tk.A()
	ln := len(a)
	i := 0
	next := func() (tk *token.T, ok bool) {
		if i < ln {
			tk = a[i]
			ok = true
			i++
		}
		return
	}
	pushIt(m, New(next))
}

// Creates a range It from begin (inclusive) to end (exclusive).
func prRange(m *machine.T) {
	tk1 := m.PopT(token.Int)
	end, _ := tk1.I()
	tk2 := m.PopT(token.Int)
	c, _ := tk2.I()
	next := func() (tk *token.T, ok bool) {
		if c < end {
			tk = token.NewI(c, m.MkPos())
			ok = true
			c++
		}
		return
	}
	pushIt(m, New(next))
}

// Creates a range It from 0 (inclusive) to end (exclusive).
func prRange0(m *machine.T) {
	tk1 := m.PopT(token.Int)
	end, _ := tk1.I()
	c := int64(0)
	next := func() (tk *token.T, ok bool) {
		if c < end {
			tk = token.NewI(c, m.MkPos())
			ok = true
			c++
		}
		return
	}
	pushIt(m, New(next))
}

// Creates an iterator over the runes of s.
func prRunes(m *machine.T) {
	tk := m.PopT(token.String)
	s, _ := tk.S()

	rd := strings.NewReader(s)
	next := func() (tk *token.T, ok bool) {
		rn, _, err := rd.ReadRune()
		if err == nil {
			tk = token.NewS(string(rn), m.MkPos())
			ok = true
			return
		}
		if err != io.EOF {
			m.Failt("Str error", "String is empty or invalid")
		}
		return
	}

	pushIt(m, New(next))
}

// Return 'true' if It has more elements.
//    m: Virtual machine.
func prHas(m *machine.T) {
	i := popIt(m)
	m.Push(token.NewB(i.ok, m.MkPos()))
}

// Return the next It value.
//    m: Virtual machine.
func prPeek(m *machine.T) {
	i := popIt(m)
	m.Push(i.e)
}

// Return the next It value.
//    m: Virtual machine.
func prNext(m *machine.T) {
	i := popIt(m)
	m.Push(i.e)
	i.e, i.ok = i.next()
}

// Copyright 10-Aug-2020 ºDeme
// GNU General Public License - V3 <http://www.gnu.org/licenses/>

package it

import (
	"github.com/dedeme/dmstack/machine"
	"github.com/dedeme/dmstack/primitives/list"
	"github.com/dedeme/dmstack/token"
	"math/rand"
)

// Returns 'true' if all elements of it returns 'true' with 'p'.
//    m: Virtual machine.
//    run: Function which running a machine.
func prAll(m *machine.T, run func(m *machine.T)) {
	tk := m.PopT(token.Procedure)
	p, _ := tk.P()
	i := popIt(m)
	r := true
	for r && i.ok {
		e := i.e
		i.e, i.ok = i.next()

		m2 := machine.NewIsolate(m.SourceDir, m.Pmachines, p)
		m2.Push(e)
		run(m2)
		tk2 := m2.PopT(token.Bool)
		r, _ = tk2.B()
	}
	m.Push(token.NewB(r, m.MkPos()))
}

// Returns 'true' if any elements of it returns 'true' with 'p'.
//    m: Virtual machine.
//    run: Function which running a machine.
func prAny(m *machine.T, run func(m *machine.T)) {
	tk := m.PopT(token.Procedure)
	p, _ := tk.P()
	i := popIt(m)
	r := false
	for !r && i.ok {
		e := i.e
		i.e, i.ok = i.next()

		m2 := machine.NewIsolate(m.SourceDir, m.Pmachines, p)
		m2.Push(e)
		run(m2)
		tk2 := m2.PopT(token.Bool)
		r, _ = tk2.B()
	}
	m.Push(token.NewB(r, m.MkPos()))
}

// Returns 'true' if 'it' contains the element 'tk'.
//    m: Virtual machine.
//    run: Function which running a machine.
func prContains(m *machine.T, run func(m *machine.T)) {
	tk := m.Pop()
	i := popIt(m)
	r := false
	for !r && i.ok {
		e := i.e
		i.e, i.ok = i.next()

		r = tk.Eq(e)
	}
	m.Push(token.NewB(r, m.MkPos()))
}

// Executes 'p' with each element of 'i'.
//    m: Virtual machine.
//    run: Function which running a machine.
func prEach(m *machine.T, run func(m *machine.T)) {
	tk := m.PopT(token.Procedure)
	p, _ := tk.P()
	i := popIt(m)
	for i.ok {
		e := i.e
		i.e, i.ok = i.next()

		m2 := machine.NewIsolate(m.SourceDir, m.Pmachines, p)
		m2.Push(e)
		run(m2)
	}
}

// Executes 'p' with each element of 'i' and its index.
//    m: Virtual machine.
//    run: Function which running a machine.
func prEachIx(m *machine.T, run func(m *machine.T)) {
	tk := m.PopT(token.Procedure)
	p, _ := tk.P()
	i := popIt(m)
	ix := int64(0)
	for i.ok {
		e := i.e
		i.e, i.ok = i.next()

		m2 := machine.NewIsolate(m.SourceDir, m.Pmachines, p)
		m2.Push(e)
		m2.Push(token.NewI(ix, m.MkPos()))
		ix++
		run(m2)
	}
}

// Returns 'true' if each element of i1 is equals to i2 with 'p'.
//    m: Virtual machine.
//    run: Function which running a machine.
func prEq(m *machine.T, run func(m *machine.T)) {
	tk := m.PopT(token.Procedure)
	p, _ := tk.P()
	i2 := popIt(m)
	i1 := popIt(m)
	r := true
	for i1.ok {
		if !i2.ok {
			r = false
			break
		}
		e1 := i1.e
		i1.e, i1.ok = i1.next()
		e2 := i2.e
		i2.e, i2.ok = i2.next()

		m2 := machine.NewIsolate(m.SourceDir, m.Pmachines, p)
		m2.Push(e1)
		m2.Push(e2)
		run(m2)
		tk2 := m2.PopT(token.Bool)
		r, _ = tk2.B()
		if !r {
			break
		}
	}
	if i2.ok {
		r = false
	}
	m.Push(token.NewB(r, m.MkPos()))
}

// Returns !prEq.
//    m: Virtual machine.
//    run: Function which running a machine.
func prNeq(m *machine.T, run func(m *machine.T)) {
	tk := m.PopT(token.Procedure)
	p, _ := tk.P()
	i2 := popIt(m)
	i1 := popIt(m)
	r := false
	for i1.ok {
		if !i2.ok {
			r = true
			break
		}
		e1 := i1.e
		i1.e, i1.ok = i1.next()
		e2 := i2.e
		i2.e, i2.ok = i2.next()

		m2 := machine.NewIsolate(m.SourceDir, m.Pmachines, p)
		m2.Push(e1)
		m2.Push(e2)
		run(m2)
		tk2 := m2.PopT(token.Bool)
		r2, _ := tk2.B()
		if !r2 {
			r = true
			break
		}
	}
	if i2.ok {
		r = true
	}
	m.Push(token.NewB(r, m.MkPos()))
}

// Returns 'true' if each element of i1 is equals to i2 with '=='.
//    m: Virtual machine.
func prEquals(m *machine.T) {
	i2 := popIt(m)
	i1 := popIt(m)

	r := true
	for i1.ok {
		if !i2.ok {
			r = false
			break
		}
		e1 := i1.e
		i1.e, i1.ok = i1.next()
		e2 := i2.e
		i2.e, i2.ok = i2.next()

		if !e1.Eq(e2) {
			r = false
			break
		}
	}
	if i2.ok {
		r = false
	}
	m.Push(token.NewB(r, m.MkPos()))
}

// Returns !prEquals.
//    m: Virtual machine.
func prNequals(m *machine.T) {
	i2 := popIt(m)
	i1 := popIt(m)

	r := false
	for i1.ok {
		if !i2.ok {
			r = true
			break
		}
		e1 := i1.e
		i1.e, i1.ok = i1.next()
		e2 := i2.e
		i2.e, i2.ok = i2.next()

		if !e1.Eq(e2) {
			r = true
			break
		}
	}
	if i2.ok {
		r = true
	}
	m.Push(token.NewB(r, m.MkPos()))
}

// Returns [e] where 'e' is the first element of 'i' which returns 'true'.
// If there is no such element, returns [].
//    m: Virtual machine.
//    run: Function which running a machine.
func prFind(m *machine.T, run func(m *machine.T)) {
	tk := m.PopT(token.Procedure)
	p, _ := tk.P()
	i := popIt(m)
	var r []*token.T
	for i.ok {
		e := i.e
		i.e, i.ok = i.next()

		m2 := machine.NewIsolate(m.SourceDir, m.Pmachines, p)
		m2.Push(e)
		run(m2)
		tk2 := m2.PopT(token.Bool)
		r2, _ := tk2.B()
		if r2 {
			r = append(r, e)
			break
		}
	}
	m.Push(token.NewL(r, m.MkPos()))
}

// Returns the index of the first element of 'i' equals to tk or -1 if such
// element does not exist.
//    m: Virtual machine.
func prIndex(m *machine.T) {
	tk := m.Pop()
	i := popIt(m)

	c := int64(0)
	for i.ok {
		e := i.e
		i.e, i.ok = i.next()

		if e.Eq(tk) {
			m.Push(token.NewI(c, m.MkPos()))
			return
		}
		c++
	}
	m.Push(token.NewI(-1, m.MkPos()))
}

// Returns the index of the first element of 'i' that returns 'true' with p
// or -1 if such element does not exist.
//    m: Virtual machine.
func prIndexF(m *machine.T, run func(m *machine.T)) {
	tk := m.PopT(token.Procedure)
	p, _ := tk.P()
	i := popIt(m)

	c := int64(0)
	for i.ok {
		e := i.e
		i.e, i.ok = i.next()

		m2 := machine.NewIsolate(m.SourceDir, m.Pmachines, p)
		m2.Push(e)
		run(m2)
		tk2 := m2.PopT(token.Bool)
		r2, _ := tk2.B()

		if r2 {
			m.Push(token.NewI(c, m.MkPos()))
			return
		}
		c++
	}
	m.Push(token.NewI(-1, m.MkPos()))
}

// Returns the index of the last element of 'i' equals to tk or -1 if such
// element does not exist.
//    m: Virtual machine.
func prLastIndex(m *machine.T) {
	tk := m.Pop()
	i := popIt(m)

	r := int64(-1)
	c := int64(0)
	for i.ok {
		e := i.e
		i.e, i.ok = i.next()

		if e.Eq(tk) {
			r = c
		}
		c++
	}
	m.Push(token.NewI(r, m.MkPos()))
}

// Returns the index of the first element of 'i' that returns 'true' with p
// or -1 if such element does not exist.
//    m: Virtual machine.
func prLastIndexF(m *machine.T, run func(m *machine.T)) {
	tk := m.PopT(token.Procedure)
	p, _ := tk.P()
	i := popIt(m)

	r := int64(-1)
	c := int64(0)
	for i.ok {
		e := i.e
		i.e, i.ok = i.next()

		m2 := machine.NewIsolate(m.SourceDir, m.Pmachines, p)
		m2.Push(e)
		run(m2)
		tk2 := m2.PopT(token.Bool)
		r2, _ := tk2.B()

		if r2 {
			r = c
		}
		c++
	}
	m.Push(token.NewI(r, m.MkPos()))
}

// Returns the result of apply 'p' over 'seed' and each element of 'i' in turn.
//    m: Virtual machine.
func prReduce(m *machine.T, run func(m *machine.T)) {
	tk := m.PopT(token.Procedure)
	p, _ := tk.P()
	seed := m.Pop()
	i := popIt(m)

	for i.ok {
		e := i.e
		i.e, i.ok = i.next()

		m2 := machine.NewIsolate(m.SourceDir, m.Pmachines, p)
		m2.Push(seed)
		m2.Push(e)
		run(m2)

		seed = m2.Pop()
	}
	m.Push(seed)
}

// Returns a list with elements of it.
//    m: Virtual machine.
func PrTo(m *machine.T) {
	i := popIt(m)
	var l []*token.T
	for {
		if i.ok {
			l = append(l, i.e)
			i.e, i.ok = i.next()
			continue
		}
		break
	}
	m.Push(token.NewL(l, m.MkPos()))
}

// Returns a new iterator with elements in reverse order.
//    m: Virtual machine.
func prReverse(m *machine.T) {
	PrTo(m)
	tk := (*m.Stack)[0]
	list.PrReverse(m)
	m.Push(tk)
	PrFrom(m)
}

// Returns a new iterator with elements randomly placed.
//    m: Virtual machine.
func prShuffle(m *machine.T) {
	PrTo(m)
	tk := (*m.Stack)[0]
	list.PrShuffle(m)
	m.Push(tk)
	PrFrom(m)
}

// Returns a new iterator with elements ordered with 'p'.
//    m: Virtual machine.
func prSort(m *machine.T, run func(m *machine.T)) {
	tk := m.Pop()
	PrTo(m)
	tk2 := (*m.Stack)[0]
	m.Push(tk)
	list.PrSort(m, run)
	m.Push(tk2)
	PrFrom(m)
}

// Returns an infinite iterator with elements of a list randomly placed in rows
// such that each row consume all its elements.
//    m: Virtual machine.
func prBox(m *machine.T) {
	tk := m.PopT(token.List)
	l, _ := tk.L()
	ln := int64(len(l))
	c := ln
	pushIt(m, New(func() (r *token.T, ok bool) {
		if c == ln {
			i := ln
			for {
				if i <= 1 {
					break
				}
				j := rand.Int63n(i)
				i--
				l[i], l[j] = l[j], l[i]
			}

			c = 0
		}
		r = l[c]
		ok = true
		c++
		return
	}))
}

// Returns a 'it.box' from a list of 'tp' which first element is the number of
// duplicates and the second is the element to put in list for the 'it.box'.
//    m: Virtual machine.
func prBox2(m *machine.T) {
	tk := m.PopT(token.List)
	l, _ := tk.L()

	var l2 []*token.T
	for _, tk2 := range l {
		pair, ok := tk2.L()
		if !ok || len(pair) != 2 {
			m.Fail("Box error", "The element (%v) is not a 'pair'", tk2.StringDraft())
		}
		n, ok := pair[0].I()
		if !ok {
			m.Fail(
				"Box error", "The first element of 'pair' (%v) is not an integer",
				tk2.StringDraft(),
			)
		}
		if n <= 0 {
			m.Fail(
				"Box error", "The first element of 'pair' (%v) is <= 0",
				tk2.StringDraft(),
			)
		}
		for i := int64(0); i < n; i++ {
			l2 = append(l2, pair[1])
		}
	}

	m.Push(token.NewL(l2, m.MkPos()))
	prBox(m)
}

// Copyright 10-Aug-2020 ºDeme
// GNU General Public License - V3 <http://www.gnu.org/licenses/>

package it

import (
	"github.com/dedeme/dmstack/machine"
	"github.com/dedeme/dmstack/symbol"
	"github.com/dedeme/dmstack/token"
)

// Join two iterators.
//    m: Virtual machine.
func prPlus(m *machine.T) {
	i2 := popIt(m)
	i1 := popIt(m)
	i := New(func() (tk *token.T, ok bool) {
		if i1.ok {
			tk = i1.e
			ok = true
			i1.e, i1.ok = i1.next()
			return
		}
		if i2.ok {
			tk = i2.e
			ok = true
			i2.e, i2.ok = i2.next()
		}
		return
	})
	pushIt(m, i)
}

// Join a list of iterators.
//    m: Virtual machine.
func prPlus2(m *machine.T) {
	tk := m.PopT(token.List)
	l, _ := tk.L()
	if len(l) == 0 {
		prEmpty(m)
		return
	}

	sym, i, ok := l[0].N()
	if !ok || sym != symbol.Iterator_ {
		m.Failt("\n Expected: Iterator object.\n  Actual  : '%v'.", sym)
	}

	pushIt(m, i.(*T))
	for ix := 1; ix < len(l); ix++ {
		sym, i, ok := l[ix].N()
		if !ok {
			m.Failt("\n  Expected: Iterator object.\n  Actual  : '%v'.", l[ix])
		}
		if sym != symbol.Iterator_ {
			m.Failt("\n  Expected: Iterator object.\n  Actual  : '%v'.", sym)
		}
		pushIt(m, i.(*T))
		prPlus(m)
	}
}

// Returns an iterator with elements of i after n elements.
//    m: Virtual machine.
func prDrop(m *machine.T) {
	tk := m.PopT(token.Int)
	n, _ := tk.I()
	i := popIt(m)
	c := int64(0)
	for {
		if i.ok && c < n {
			i.e, i.ok = i.next()
			c++
			continue
		}
		break
	}
	pushIt(m, i)
}

// Returns an iterator with elements of i after the last element which
// accomplishes a condition.
//    m: Virtual machine.
//    run: Function which running a machine.
func prDropf(m *machine.T, run func(m *machine.T)) {
	tk := m.PopT(token.Procedure)
	p, _ := tk.P()
	i := popIt(m)
	for {
		if i.ok {
			m2 := machine.NewIsolate(m.SourceDir, m.Pmachines, p)
			m2.Push(i.e)
			run(m2)
			tk2 := m2.PopT(token.Bool)
			value, _ := tk2.B()
			if value {
				i.e, i.ok = i.next()
				continue
			}
		}
		break
	}
	pushIt(m, i)
}

// Returns an iterator with elements of i witch accomplish a condition.
//    m: Virtual machine.
//    run: Function which running a machine.
func prFilter(m *machine.T, run func(m *machine.T)) {
	tk := m.PopT(token.Procedure)
	p, _ := tk.P()
	i := popIt(m)
	i2 := New(func() (tkr *token.T, okr bool) {
		for {
			if i.ok {
				e := i.e
				i.e, i.ok = i.next()

				m2 := machine.NewIsolate(m.SourceDir, m.Pmachines, p)
				m2.Push(e)
				run(m2)
				tk2 := m2.PopT(token.Bool)
				value, _ := tk2.B()
				if value {
					tkr = e
					okr = true
					return
				}
				continue
			}
			return
		}
	})
	pushIt(m, i2)
}

// Returns an iterator with the result of run 'p' with each element.
//    m: Virtual machine.
//    run: Function which running a machine.
func PrMap(m *machine.T, run func(m *machine.T)) {
	tk := m.PopT(token.Procedure)
	p, _ := tk.P()
	i := popIt(m)
	i2 := New(func() (tkr *token.T, okr bool) {
		if i.ok {
			e := i.e
			i.e, i.ok = i.next()

			m2 := machine.NewIsolate(m.SourceDir, m.Pmachines, p)
			m2.Push(e)
			run(m2)
			tkr = m2.Pop()
			okr = true
		}
		return
	})
	pushIt(m, i2)
}

// Add an element at the end of an iterator.
//    m: Virtual machine.
func prPush(m *machine.T) {
	tk := m.Pop()
	i := popIt(m)
	more := true
	i2 := New(func() (tkr *token.T, okr bool) {
		if more {
			if i.ok {
				tkr = i.e
				i.e, i.ok = i.next()
			} else {
				tkr = tk
				more = false
			}
			okr = true
		}
		return
	})
	pushIt(m, i2)
}

// Add an element at the beginning of an iterator.
//    m: Virtual machine.
func prPush0(m *machine.T) {
	tk := m.Pop()
	i := popIt(m)
	added := false
	i2 := New(func() (tkr *token.T, okr bool) {
		if added {
			if i.ok {
				tkr = i.e
				i.e, i.ok = i.next()
				okr = true
			}
			return
		}
		tkr = tk
		okr = true
		added = true
		return
	})
	pushIt(m, i2)
}

// Returns an iterator with the n first elements of i.
//    m: Virtual machine.
func prTake(m *machine.T) {
	tk := m.PopT(token.Int)
	n, _ := tk.I()
	i := popIt(m)
	c := int64(0)
	i2 := New(func() (tkr *token.T, okr bool) {
		if i.ok && c < n {
			tkr = i.e
			i.e, i.ok = i.next()
			c++
			okr = true
		}
		return
	})
	pushIt(m, i2)
}

// Returns an iterator with first elements of i which accomplishes a condition.
//    m: Virtual machine.
//    run: Function which running a machine.
func prTakef(m *machine.T, run func(m *machine.T)) {
	tk := m.PopT(token.Procedure)
	p, _ := tk.P()
	i := popIt(m)
	i2 := New(func() (tkr *token.T, okr bool) {
		if i.ok {
			m2 := machine.NewIsolate(m.SourceDir, m.Pmachines, p)
			m2.Push(i.e)
			run(m2)
			tk2 := m2.PopT(token.Bool)
			value, _ := tk2.B()
			if value {
				tkr = i.e
				i.e, i.ok = i.next()
				okr = true
			}
		}
		return
	})
	pushIt(m, i2)
}

// Returns an iterator with pairs from i1 and i2.
//    m: Virtual machine.
func prZip(m *machine.T) {
	i2 := popIt(m)
	i1 := popIt(m)
	i := New(func() (tkr *token.T, okr bool) {
		if i1.ok && i2.ok {
			l := []*token.T{i1.e, i2.e}
			i1.e, i1.ok = i1.next()
			i2.e, i2.ok = i2.next()
			tkr = token.NewL(l, m.MkPos())
			okr = true
		}
		return
	})
	pushIt(m, i)
}

// Returns an iterator with triples from i1, i2 and i3.
//    m: Virtual machine.
func prZip3(m *machine.T) {
	i3 := popIt(m)
	i2 := popIt(m)
	i1 := popIt(m)
	i := New(func() (tkr *token.T, okr bool) {
		if i1.ok && i2.ok && i3.ok {
			l := []*token.T{i1.e, i2.e, i3.e}
			i1.e, i1.ok = i1.next()
			i2.e, i2.ok = i2.next()
			i3.e, i3.ok = i3.next()
			tkr = token.NewL(l, m.MkPos())
			okr = true
		}
		return
	})
	pushIt(m, i)
}

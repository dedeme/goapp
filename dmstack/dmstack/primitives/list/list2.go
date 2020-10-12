// Copyright 22-Jul-2020 ºDeme
// GNU General Public License - V3 <http://www.gnu.org/licenses/>

package list

import (
	"github.com/dedeme/dmstack/machine"
	"github.com/dedeme/dmstack/token"
)

// Auxiliar function
func cmpf(
	m *machine.T, run func(m *machine.T), proc []*token.T, e1, e2 *token.T,
) bool {
	m2 := machine.NewIsolate(m.SourceDir, m.Pmachines, proc)
	m2.Push(e1)
	m2.Push(e2)
	run(m2)
	tk2 := m2.PopT(token.Bool)
	r, _ := tk2.B()
	return r
}

// Auxiliar function
func predf(
	m *machine.T, run func(m *machine.T), proc []*token.T, e *token.T,
) bool {
	m2 := machine.NewIsolate(m.SourceDir, m.Pmachines, proc)
	m2.Push(e)
	run(m2)
	tk2 := m2.PopT(token.Bool)
	r, _ := tk2.B()
	return r
}

// Returns a new list with elements of 'l' without duplicates.
func prRemoveDup(m *machine.T, run func(m *machine.T)) {
	tk := m.PopT(token.Procedure)
	p, _ := tk.P()
	l := popList(m)
	var r []*token.T
	for _, e := range l {
		new := true
		for _, e2 := range r {
			if cmpf(m, run, p, e, e2) {
				new = false
			}
		}
		if new {
			r = append(r, e)
		}
	}
	pushList(m, r)
}

// Returns 'true' if every element of l is 'true' with 'p'.
func prAll(m *machine.T, run func(m *machine.T)) {
	tk := m.PopT(token.Procedure)
	p, _ := tk.P()
	l := popList(m)
	ok := true
	for _, e := range l {
		if !predf(m, run, p, e) {
			ok = false
			break
		}
	}
	m.Push(token.NewB(ok, m.MkPos()))
}

// Returns 'true' if some element of l is 'true' with 'p'.
func prAny(m *machine.T, run func(m *machine.T)) {
	tk := m.PopT(token.Procedure)
	p, _ := tk.P()
	l := popList(m)
	ok := false
	for _, e := range l {
		if predf(m, run, p, e) {
			ok = true
			break
		}
	}
	m.Push(token.NewB(ok, m.MkPos()))
}

// Executes 'p' with every element of l.
func prEach(m *machine.T, run func(m *machine.T)) {
	tk := m.PopT(token.Procedure)
	p, _ := tk.P()
	l := popList(m)
	for _, e := range l {
		m2 := machine.NewIsolate(m.SourceDir, m.Pmachines, p)
		m2.Push(e)
		run(m2)
	}
}

// Executes 'p' with every element of l and its index.
func prEachIx(m *machine.T, run func(m *machine.T)) {
	tk := m.PopT(token.Procedure)
	p, _ := tk.P()
	l := popList(m)
	pos := m.MkPos()
	for i, e := range l {
		m2 := machine.NewIsolate(m.SourceDir, m.Pmachines, p)
		m2.Push(e)
		m2.Push(token.NewI(int64(i), pos))
		run(m2)
	}
}

// Auxiliar function
func prEqNeq(m *machine.T, run func(m *machine.T)) bool {
	tk := m.PopT(token.Procedure)
	p, _ := tk.P()
	l1 := popList(m)
	l2 := popList(m)
	ok := false
	if len(l1) == len(l2) {
		ok = true
		for i := 0; i < len(l1); i++ {
			if !cmpf(m, run, p, l1[i], l2[i]) {
				ok = false
				break
			}
		}
	}
	return ok
}

// Returns 'true' if l1 == l2, comparing elements with 'p'.
func prEq(m *machine.T, run func(m *machine.T)) {
	m.Push(token.NewB(prEqNeq(m, run), m.MkPos()))
}

// Returns 'false' if l1 == l2, comparing elements with 'p'.
func prNeq(m *machine.T, run func(m *machine.T)) {
	m.Push(token.NewB(!prEqNeq(m, run), m.MkPos()))
}

// Returns the index of the first element of 'ls' which yield 'true' with 'p',
// or -1 if no element match the condition.
func prIndex(m *machine.T, run func(m *machine.T)) {
	tk := m.PopT(token.Procedure)
	p, _ := tk.P()
	l := popList(m)
	r := -1
	for i, e := range l {
		if predf(m, run, p, e) {
			r = i
			break
		}
	}
	m.Push(token.NewI(int64(r), m.MkPos()))
}

// Returns a option with the first element of 'ls' which yield 'true' with 'p',
// or a option.none if no element match the condition.
func prFind(m *machine.T, run func(m *machine.T)) {
	tk := m.PopT(token.Procedure)
	p, _ := tk.P()
	l := popList(m)
	r := []*token.T{}
	for _, e := range l {
		if predf(m, run, p, e) {
			r = append(r, e)
			break
		}
	}
	pushList(m, r)
}

// Returns the index of the last element of 'ls' which yield 'true' with 'p',
// or -1 if no element match the condition.
func prLastIndex(m *machine.T, run func(m *machine.T)) {
	tk := m.PopT(token.Procedure)
	p, _ := tk.P()
	l := popList(m)
	r := -1
	for i, e := range l {
		if predf(m, run, p, e) {
			r = i
		}
	}
	m.Push(token.NewI(int64(r), m.MkPos()))
}

// Applies 'proc' to (seed, element of 'l') and the result is used as new seed
// for a new application of 'p', until every element of 'l' is processed.
func prReduce(m *machine.T, run func(m *machine.T)) {
	tk := m.PopT(token.Procedure)
	p, _ := tk.P()
	seed := m.Pop()
	l := popList(m)
	for _, e := range l {
		m2 := machine.NewIsolate(m.SourceDir, m.Pmachines, p)
		m2.Push(seed)
		m2.Push(e)
		run(m2)
		seed = m2.Pop()
	}
	m.Push(seed)
}

// Returns a shallow copy of 'l'.
func prCopy(m *machine.T, run func(m *machine.T)) {
	l := popList(m)
	pushList(m, append([]*token.T{}, l...))
}

// Returns a new list with elemts of 'l', removing its 'n' first elements.
func prDrop(m *machine.T, run func(m *machine.T)) {
	tk := m.PopT(token.Int)
	n, _ := tk.I()
	l := popList(m)
	if n < 0 {
		n = int64(len(l)) + n
	}
	r := []*token.T{}
	skip := true
	for i, e := range l {
		if skip {
			if int64(i) < n {
				continue
			}
			skip = false
		}
		r = append(r, e)
	}
	pushList(m, r)
}

// Returns a new list with elemts of 'l', removing its first elements with
// mach with 'p'.
func prDropf(m *machine.T, run func(m *machine.T)) {
	tk := m.PopT(token.Procedure)
	p, _ := tk.P()
	l := popList(m)
	r := []*token.T{}
	skip := true
	for _, e := range l {
		if skip {
			if predf(m, run, p, e) {
				continue
			}
			skip = false
		}
		r = append(r, e)
	}
	pushList(m, r)
}

// Returns a new list with elements of 'l' which match 'p'.
func prFilter(m *machine.T, run func(m *machine.T)) {
	tk := m.PopT(token.Procedure)
	p, _ := tk.P()
	l := popList(m)
	r := []*token.T{}
	for _, e := range l {
		if predf(m, run, p, e) {
			r = append(r, e)
		}
	}
	pushList(m, r)
}

// Returns a new list changing each list element by its values. Only flats
// one level.
func prFlat(m *machine.T, run func(m *machine.T)) {
	l := popList(m)
	var nl []*token.T
	for _, e := range l {
		subl, ok := e.L()
		if ok {
			nl = append(nl, subl...)
		} else {
			nl = append(nl, e)
		}
	}
	pushList(m, nl)
}

// Returns a new list with elements of 'l' transformed with 'p'.
func prMap(m *machine.T, run func(m *machine.T)) {
	tk := m.PopT(token.Procedure)
	p, _ := tk.P()
	l := popList(m)
	r := []*token.T{}
	for _, e := range l {
		m2 := machine.NewIsolate(m.SourceDir, m.Pmachines, p)
		m2.Push(e)
		run(m2)
		r = append(r, m2.Pop())
	}
	pushList(m, r)
}

// Equals to 'l end list.take begin list.drop'
func prSub(m *machine.T, run func(m *machine.T)) {
	tk1 := m.PopT(token.Int)
	end, _ := tk1.I()
	tk2 := m.PopT(token.Int)
	begin, _ := tk2.I()
	l := popList(m)
	if begin < 0 {
		begin = int64(len(l)) + begin
	}
	if end < 0 {
		end = int64(len(l)) + end
	}
	r := []*token.T{}
	skip := true
	for i, e := range l {
		if skip {
			if int64(i) < begin {
				continue
			}
			skip = false
		}
		if int64(i) >= end {
			break
		}
		r = append(r, e)
	}
	pushList(m, r)
}

// Returns a new list with the 'n' first elements of 'l'.
func prTake(m *machine.T, run func(m *machine.T)) {
	tk := m.PopT(token.Int)
	n, _ := tk.I()
	l := popList(m)
	if n < 0 {
		n = int64(len(l)) + n
	}
	r := []*token.T{}
	for i, e := range l {
		if int64(i) >= n {
			break
		}
		r = append(r, e)
	}
	pushList(m, r)
}

// Returns a new list with firsts elements of 'l' which match with 'p'.
func prTakef(m *machine.T, run func(m *machine.T)) {
	tk := m.PopT(token.Procedure)
	p, _ := tk.P()
	l := popList(m)
	r := []*token.T{}
	for _, e := range l {
		if !predf(m, run, p, e) {
			break
		}
		r = append(r, e)
	}
	pushList(m, r)
}

// Returns a new list with tuples of elements from 'l1', 'l2'.
func prZip(m *machine.T, run func(m *machine.T)) {
	l2 := popList(m)
	l1 := popList(m)
	r := []*token.T{}
	ln := len(l1)
	if ln > len(l2) {
		ln = len(l2)
	}
	for i := 0; i < ln; i++ {
		r = append(r, token.NewL([]*token.T{l1[i], l2[i]}, m.MkPos()))
	}
	pushList(m, r)
}

// Returns a new list with tuples of elements from 'l1', 'l2', 'l3'.
func prZip3(m *machine.T, run func(m *machine.T)) {
	l3 := popList(m)
	l2 := popList(m)
	l1 := popList(m)
	r := []*token.T{}
	ln := len(l1)
	if ln > len(l2) {
		ln = len(l2)
	}
	if ln > len(l3) {
		ln = len(l3)
	}
	for i := 0; i < ln; i++ {
		r = append(r, token.NewL([]*token.T{l1[i], l2[i], l3[i]}, m.MkPos()))
	}
	pushList(m, r)
}

// Returns two lists with first and second values of elements (duples) from
// 'l'.
func prUnzip(m *machine.T, run func(m *machine.T)) {
	l := popList(m)
	r1 := []*token.T{}
	r2 := []*token.T{}
	for _, e := range l {
		l, ok := e.L()
		if !ok || len(l) != 2 {
			m.Failt("List elements must be duples")
		}
		r1 = append(r1, l[0])
		r2 = append(r2, l[1])
	}
	pushList(m, r1)
	pushList(m, r2)
}

// Returns three lists with first, second and third values of elements (triples)
// from 'l'.
func prUnzip3(m *machine.T, run func(m *machine.T)) {
	l := popList(m)
	r1 := []*token.T{}
	r2 := []*token.T{}
	r3 := []*token.T{}
	for _, e := range l {
		l, ok := e.L()
		if !ok || len(l) != 3 {
			m.Failt("List elements must be duples")
		}
		r1 = append(r1, l[0])
		r2 = append(r2, l[1])
		r3 = append(r3, l[2])
	}
	pushList(m, r1)
	pushList(m, r2)
	pushList(m, r3)
}

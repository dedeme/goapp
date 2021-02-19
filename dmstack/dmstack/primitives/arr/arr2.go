// Copyright 10-Jan-2021 ºDeme
// GNU General Public License - V3 <http://www.gnu.org/licenses/>

package arr

import (
	"github.com/dedeme/dmstack/machine"
	"github.com/dedeme/dmstack/token"
)

// Auxiliar function
func cmpf(
	m *machine.T, run func(m *machine.T), proc *token.T, e1, e2 *token.T,
) bool {
	m2 := machine.New(m.Source, m.Pmachines, proc)
	m2.Push(e1)
	m2.Push(e2)
	run(m2)
	tk2 := m2.PopT(token.Bool)
	r, _ := tk2.B()
	return r
}

// Auxiliar function
func predf(
	m *machine.T, run func(m *machine.T), proc *token.T, e *token.T,
) bool {
	m2 := machine.New(m.Source, m.Pmachines, proc)
	m2.Push(e)
	run(m2)
	tk2 := m2.PopT(token.Bool)
	r, _ := tk2.B()
	return r
}

// Returns a new list with elements of 'l' without duplicates.
func prRemoveDup(m *machine.T, run func(m *machine.T)) {
	tk := m.PopT(token.Procedure)
	a := popArr(m)
	var r []*token.T
	for _, e := range a {
		new := true
		for _, e2 := range r {
			if cmpf(m, run, tk, e, e2) {
				new = false
			}
		}
		if new {
			r = append(r, e)
		}
	}
	pushArr(m, r)
}

// Returns 'true' if every element of l is 'true' with 'p'.
func prAll(m *machine.T, run func(m *machine.T)) {
	tk := m.PopT(token.Procedure)
	a := popArr(m)
	ok := true
	for _, e := range a {
		if !predf(m, run, tk, e) {
			ok = false
			break
		}
	}
	m.Push(token.NewB(ok, m.MkPos()))
}

// Returns 'true' if some element of l is 'true' with 'p'.
func prAny(m *machine.T, run func(m *machine.T)) {
	tk := m.PopT(token.Procedure)
	a := popArr(m)
	ok := false
	for _, e := range a {
		if predf(m, run, tk, e) {
			ok = true
			break
		}
	}
	m.Push(token.NewB(ok, m.MkPos()))
}

// Executes 'p' with every element of l.
func prEach(m *machine.T, run func(m *machine.T)) {
	tk := m.PopT(token.Procedure)
	a := popArr(m)
	for _, e := range a {
		m2 := machine.New(m.Source, m.Pmachines, tk)
		m2.Push(e)
		run(m2)
	}
}

// Executes 'p' with every element of l and its index.
func prEachIx(m *machine.T, run func(m *machine.T)) {
	tk := m.PopT(token.Procedure)
	a := popArr(m)
	pos := m.MkPos()
	for i, e := range a {
		m2 := machine.New(m.Source, m.Pmachines, tk)
		m2.Push(e)
		m2.Push(token.NewI(int64(i), pos))
		run(m2)
	}
}

// Auxiliar function
func prEqNeq(m *machine.T, run func(m *machine.T)) bool {
	tk := m.PopT(token.Procedure)
	a1 := popArr(m)
	a2 := popArr(m)
	ok := false
	if len(a1) == len(a2) {
		ok = true
		for i := 0; i < len(a1); i++ {
			if !cmpf(m, run, tk, a1[i], a2[i]) {
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
	a := popArr(m)
	r := -1
	for i, e := range a {
		if predf(m, run, tk, e) {
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
	a := popArr(m)
	r := []*token.T{}
	for _, e := range a {
		if predf(m, run, tk, e) {
			r = append(r, e)
			break
		}
	}
	pushArr(m, r)
}

// Returns the index of the last element of 'ls' which yield 'true' with 'p',
// or -1 if no element match the condition.
func prLastIndex(m *machine.T, run func(m *machine.T)) {
	tk := m.PopT(token.Procedure)
	a := popArr(m)
	r := -1
	for i, e := range a {
		if predf(m, run, tk, e) {
			r = i
		}
	}
	m.Push(token.NewI(int64(r), m.MkPos()))
}

// Applies 'proc' to (seed, element of 'l') and the result is used as new seed
// for a new application of 'p', until every element of 'l' is processed.
func prReduce(m *machine.T, run func(m *machine.T)) {
	tk := m.PopT(token.Procedure)
	seed := m.Pop()
	a := popArr(m)
	for _, e := range a {
		m2 := machine.New(m.Source, m.Pmachines, tk)
		m2.Push(seed)
		m2.Push(e)
		run(m2)
		seed = m2.Pop()
	}
	m.Push(seed)
}

// Returns a shallow copy of 'l'.
func prCopy(m *machine.T, run func(m *machine.T)) {
	a := popArr(m)
	pushArr(m, append([]*token.T{}, a...))
}

// Returns a new list with elemts of 'l', removing its 'n' first elements.
func prDrop(m *machine.T, run func(m *machine.T)) {
	tk := m.PopT(token.Int)
	n, _ := tk.I()
	a := popArr(m)
	if n < 0 {
		n = int64(len(a)) + n
	}
	r := []*token.T{}
	skip := true
	for i, e := range a {
		if skip {
			if int64(i) < n {
				continue
			}
			skip = false
		}
		r = append(r, e)
	}
	pushArr(m, r)
}

// Returns a new list with elemts of 'l', removing its first elements with
// mach with 'p'.
func prDropf(m *machine.T, run func(m *machine.T)) {
	tk := m.PopT(token.Procedure)
	a := popArr(m)
	r := []*token.T{}
	skip := true
	for _, e := range a {
		if skip {
			if predf(m, run, tk, e) {
				continue
			}
			skip = false
		}
		r = append(r, e)
	}
	pushArr(m, r)
}

// Returns a new list with elements of 'l' which match 'p'.
func prFilter(m *machine.T, run func(m *machine.T)) {
	tk := m.PopT(token.Procedure)
	a := popArr(m)
	r := []*token.T{}
	for _, e := range a {
		if predf(m, run, tk, e) {
			r = append(r, e)
		}
	}
	pushArr(m, r)
}

// Returns a new list changing each list element by its values. Only flats
// one level.
func prFlat(m *machine.T, run func(m *machine.T)) {
	a := popArr(m)
	var na []*token.T
	for _, e := range a {
		suba, ok := e.A()
		if ok {
			na = append(na, suba...)
		} else {
			na = append(na, e)
		}
	}
	pushArr(m, na)
}

// Returns a new list with elements of 'l' transformed with 'p'.
func prMap(m *machine.T, run func(m *machine.T)) {
	tk := m.PopT(token.Procedure)
	a := popArr(m)
	r := []*token.T{}
	for _, e := range a {
		m2 := machine.New(m.Source, m.Pmachines, tk)
		m2.Push(e)
		run(m2)
		r = append(r, m2.Pop())
	}
	pushArr(m, r)
}

// Equals to 'l end list.take begin list.drop'
func prSub(m *machine.T, run func(m *machine.T)) {
	tk1 := m.PopT(token.Int)
	end, _ := tk1.I()
	tk2 := m.PopT(token.Int)
	begin, _ := tk2.I()
	a := popArr(m)
	if begin < 0 {
		begin = int64(len(a)) + begin
	}
	if end < 0 {
		end = int64(len(a)) + end
	}
	r := []*token.T{}
	skip := true
	for i, e := range a {
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
	pushArr(m, r)
}

// Returns a new list with the 'n' first elements of 'l'.
func prTake(m *machine.T, run func(m *machine.T)) {
	tk := m.PopT(token.Int)
	n, _ := tk.I()
	a := popArr(m)
	if n < 0 {
		n = int64(len(a)) + n
	}
	r := []*token.T{}
	for i, e := range a {
		if int64(i) >= n {
			break
		}
		r = append(r, e)
	}
	pushArr(m, r)
}

// Returns a new list with firsts elements of 'l' which match with 'p'.
func prTakef(m *machine.T, run func(m *machine.T)) {
	tk := m.PopT(token.Procedure)
	a := popArr(m)
	r := []*token.T{}
	for _, e := range a {
		if !predf(m, run, tk, e) {
			break
		}
		r = append(r, e)
	}
	pushArr(m, r)
}

// Returns a new list with tuples of elements from 'l1', 'l2'.
func prZip(m *machine.T, run func(m *machine.T)) {
	a2 := popArr(m)
	a1 := popArr(m)
	r := []*token.T{}
	an := len(a1)
	if an > len(a2) {
		an = len(a2)
	}
	for i := 0; i < an; i++ {
		r = append(r, token.NewA([]*token.T{a1[i], a2[i]}, m.MkPos()))
	}
	pushArr(m, r)
}

// Returns a new list with tuples of elements from 'l1', 'l2', 'l3'.
func prZip3(m *machine.T, run func(m *machine.T)) {
	a3 := popArr(m)
	a2 := popArr(m)
	a1 := popArr(m)
	r := []*token.T{}
	an := len(a1)
	if an > len(a2) {
		an = len(a2)
	}
	if an > len(a3) {
		an = len(a3)
	}
	for i := 0; i < an; i++ {
		r = append(r, token.NewA([]*token.T{a1[i], a2[i], a3[i]}, m.MkPos()))
	}
	pushArr(m, r)
}

// Returns two lists with first and second values of elements (duples) from
// 'l'.
func prUnzip(m *machine.T, run func(m *machine.T)) {
	a := popArr(m)
	r1 := []*token.T{}
	r2 := []*token.T{}
	for _, e := range a {
		a, ok := e.A()
		if !ok || len(a) != 2 {
			m.Failt("List elements must be duples")
		}
		r1 = append(r1, a[0])
		r2 = append(r2, a[1])
	}
	pushArr(m, r1)
	pushArr(m, r2)
}

// Returns three lists with first, second and third values of elements (triples)
// from 'l'.
func prUnzip3(m *machine.T, run func(m *machine.T)) {
	a := popArr(m)
	r1 := []*token.T{}
	r2 := []*token.T{}
	r3 := []*token.T{}
	for _, e := range a {
		a, ok := e.A()
		if !ok || len(a) != 3 {
			m.Failt("List elements must be duples")
		}
		r1 = append(r1, a[0])
		r2 = append(r2, a[1])
		r3 = append(r3, a[2])
	}
	pushArr(m, r1)
	pushArr(m, r2)
	pushArr(m, r3)
}

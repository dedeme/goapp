// Copyright 08-Jan-2021 ºDeme
// GNU General Public License - V3 <http://www.gnu.org/licenses/>

// arr module.
package arr

import (
	"github.com/dedeme/dmstack/machine"
	"github.com/dedeme/dmstack/token"
	"math/rand"
	"sort"
)

// Auxiliar function
func popArr(m *machine.T) (a []*token.T) {
	tk := m.PopT(token.Array)
	a, _ = tk.A()
	return
}

// Auxiliar function
func pushArr(m *machine.T, a []*token.T) {
	m.Push(token.NewA(a, m.MkPos()))
}

// Creates a new empty arr.
//    m: Virtual machine.
func prNew(m *machine.T) {
	pushArr(m, []*token.T{})
}

// Creates an arr with n elements with the value value.
// Throws a "Range error" if n < 0.
// (Value is cloned. See token.Clone).
//    m: Virtual machine.
func prMake(m *machine.T) {
	tk := m.PopT(token.Int)
	n, _ := tk.I()
	if n < 0 {
		m.Fail(machine.ERange(), "Number of elements < 0 (%v)", n)
	}
	value := m.Pop()
	var a []*token.T
	for i := int64(0); i < n; i++ {
		a = append(a, value.Clone())
	}
	pushArr(m, a)
}

// Returns the number of elements of arr.
//    m: Virtual machine.
func prSize(m *machine.T) {
	a := popArr(m)
	m.Push(token.NewI(int64(len(a)), m.MkPos()))
}

// Returns 'true' if arr is empty.
//    m: Virtual machine.
func prEmpty(m *machine.T) {
	a := popArr(m)
	m.Push(token.NewB(len(a) == 0, m.MkPos()))
}

// Adds an element at the end of arr.
//    m: Virtual machine.
func prPush(m *machine.T) {
	tkV := m.Pop()
	tk := m.PopT(token.Array)
	a, _ := tk.A()
	tk.SetA(append(a, tkV))
}

// Adds an element at the beginning of arr.
//    m: Virtual machine.
func prPush0(m *machine.T) {
	tkV := m.Pop()
	tk := m.PopT(token.Array)
	a, _ := tk.A()
	tk.SetA(append([]*token.T{tkV}, a...))
}

// Auxiliar function
//    m: Virtual machine.
func popPeek(m *machine.T, isPop bool) {
	tk := m.PopT(token.Array)
	a, _ := tk.A()
	len1 := len(a) - 1
	if len1 < 0 {
		m.Fail(
			machine.ERange(),
			"\n  Expected: List with at least 1 element.\n  Actual  : %v",
			tk.StringDraft(),
		)
	}
	m.Push(a[len1])
	if isPop {
		tk.SetA(a[:len1])
	}
}

// Removes and returns the last element of arr.
//    m: Virtual machine.
func prPop(m *machine.T) {
	popPeek(m, true)
}

// Returns but not remove the last element of arr.
//    m: Virtual machine.
func prPeek(m *machine.T) {
	popPeek(m, false)
}

func popPeek0(m *machine.T, isPop bool) {
	tk := m.PopT(token.Array)
	a, _ := tk.A()
	length := len(a)
	if length == 0 {
		m.Fail(
			machine.ERange(),
			"\n  Expected: List with at least 1 element.\n  Actual  : %v",
			tk.StringDraft(),
		)
	}
	m.Push(a[0])
	if isPop {
		tk.SetA(a[1:])
	}
}

// Removes and returns the first element of arr.
//    m: Virtual machine.
func prPop0(m *machine.T) {
	popPeek0(m, true)
}

// Returns but not remove the first element of arr.
//    m: Virtual machine.
func prPeek0(m *machine.T) {
	popPeek0(m, false)
}

// Inserts an element in an arr.
//    m: Virtual machine.
func prInsert(m *machine.T) {
	tk1 := m.Pop()
	tk2 := m.PopT(token.Int)
	i, _ := tk2.I()
	tk3 := m.PopT(token.Array)
	a, _ := tk3.A()
	ln := int64(len(a))
	if i < 0 || i > ln {
		m.Fail(machine.ERange(), "%v [0, %v]", i, ln)
	}
	var ra []*token.T
	ra = append(ra, a[:i]...)
	ra = append(ra, tk1)
	tk3.SetA(append(ra, a[i:]...))
}

// Inserts an arr in another arr.
//    m: Virtual machine.
func prInsertList(m *machine.T) {
	subList := popArr(m)
	tk2 := m.PopT(token.Int)
	i, _ := tk2.I()
	tk3 := m.PopT(token.Array)
	a, _ := tk3.A()
	ln := int64(len(a))
	if i < 0 || i > ln {
		m.Fail(machine.ERange(), "%v [0, %v]", i, ln)
	}
	var ra []*token.T
	ra = append(ra, a[:i]...)
	ra = append(ra, subList...)
	tk3.SetA(append(ra, a[i:]...))
}

// Removes an element in an arr.
//    m: Virtual machine.
func prRemove(m *machine.T) {
	tk1 := m.PopT(token.Int)
	i, _ := tk1.I()
	tk2 := m.PopT(token.Array)
	a, _ := tk2.A()
	ln := int64(len(a))
	if i < 0 || i >= ln {
		m.Fail(machine.ERange(), "%v [0, %v)", i, ln)
	}
	var ra []*token.T
	ra = append(ra, a[:i]...)
	tk2.SetA(append(ra, a[i+1:]...))
}

// Removes a range [begin-end) of elements in an arr.
//    m: Virtual machine.
func prRemoveRange(m *machine.T) {
	tk1 := m.PopT(token.Int)
	end, _ := tk1.I()
	tk2 := m.PopT(token.Int)
	begin, _ := tk2.I()
	tk3 := m.PopT(token.Array)
	a, _ := tk3.A()
	ln := int64(len(a))
	if begin < 0 || begin >= ln {
		m.Fail(machine.ERange(), "%v [0, %v)", begin, ln)
	}
	if end < 0 || end > ln {
		m.Fail(machine.ERange(), "%v [0, %v]", end, ln)
	}
	if end > begin {
		var ra []*token.T
		ra = append(ra, a[:begin]...)
		tk3.SetA(append(ra, a[end:]...))
	}
}

// Removes every element in an arr.
//    m: Virtual machine.
func prClear(m *machine.T) {
	tk := m.PopT(token.Array)
	tk.SetA([]*token.T{})
}

// Reverse arr elements.
//    m: Virtual machine.
func PrReverse(m *machine.T) {
	a := popArr(m)
	for i, j := 0, len(a)-1; i < j; i, j = i+1, j-1 {
		a[i], a[j] = a[j], a[i]
	}
}

// Ramdon reordering of an arr.
//    m: Virtual machine.
func PrShuffle(m *machine.T) {
	a := popArr(m)
	i := int64(len(a))
	for {
		if i <= 1 {
			break
		}
		j := rand.Int63n(i)
		i--
		a[i], a[j] = a[j], a[i]
	}
}

// Sorts in place an arr using 'p'.
//    m: Virtual machine.
func PrSort(m *machine.T, run func(m *machine.T)) {
	tk := m.PopT(token.Procedure)

	a := popArr(m)
	sort.Slice(a, func(i, j int) bool {
		m2 := machine.New(m.Source, m.Pmachines, tk)
		m2.Push(a[i])
		m2.Push(a[j])
		run(m2)
		tk2 := m2.PopT(token.Bool)
		r, _ := tk2.B()
		return r
	})
}

// Returns the element of arr at postion ix.
//    m: Virtual machine.
func prGet(m *machine.T) {
	tkIx := m.PopT(token.Int)
	ix, _ := tkIx.I()
	tk := m.PopT(token.Array)
	a, _ := tk.A()
	if ix < 0 || ix >= int64(len(a)) {
		m.Fail(machine.ERange(), "%v [0, %v)", ix, len(a))
	}
	m.Push(a[ix])
}

// Sets the element of arr at postion ix.
//    m: Virtual machine.
func prSet(m *machine.T) {
	tkV := m.Pop()
	tkIx := m.PopT(token.Int)
	ix, _ := tkIx.I()
	tk := m.PopT(token.Array)
	a, _ := tk.A()
	if ix < 0 || ix >= int64(len(a)) {
		m.Fail(machine.ERange(), "%v [0, %v)", ix, len(a))
	}
	a[ix] = tkV
}

//  Updates the element of arr at position ix.
//    m: Virtual machine.
//    run: Function which running a machine.
func prUp(m *machine.T, run func(m *machine.T)) {
	tkP := m.PopT(token.Procedure)
	tkIx := m.PopT(token.Int)
	ix, _ := tkIx.I()
	tk := m.PopT(token.Array)
	a, _ := tk.A()
	if ix < 0 || ix >= int64(len(a)) {
		m.Fail(machine.ERange(), "%v [0, %v)", ix, len(a))
	}
	m2 := machine.New(m.Source, m.Pmachines, tkP)
	m2.Push(a[ix])
	run(m2)
	a[ix] = m2.Pop()
}

// Fill an arr with an element. Element is cloned. (See token.Clone).
//    m: Virtual machine.
func prFill(m *machine.T) {
	tk := m.Pop()
	a := popArr(m)
	for i := range a {
		a[i] = tk.Clone()
	}
}

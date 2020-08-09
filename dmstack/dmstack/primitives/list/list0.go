// Copyright 21-Jul-2020 ºDeme
// GNU General Public License - V3 <http://www.gnu.org/licenses/>

package list

import (
	"github.com/dedeme/dmstack/machine"
	"github.com/dedeme/dmstack/token"
	"math/rand"
	"sort"
)

// Auxiliar function
func popList(m *machine.T) (l []*token.T) {
	tk := m.PopT(token.List)
	l, _ = tk.L()
	return
}

// Auxiliar function
func pushList(m *machine.T, l []*token.T) {
	m.Push(token.NewL(l, m.MkPos()))
}

// Creates a new empty list.
//    m: Virtual machine.
func prNew(m *machine.T) {
	pushList(m, []*token.T{})
}

// Creates a list with n elements with the value value.
// (Value is cloned. See token.Clone).
//    m: Virtual machine.
func prMake(m *machine.T) {
	tk := m.PopT(token.Int)
	n, _ := tk.I()
	if n < 0 {
		m.Failf("Number of elements < 0 (%v)", n)
	}
	value := m.Pop()
	var l []*token.T
	for i := int64(0); i < n; i++ {
		l = append(l, value.Clone())
	}
	pushList(m, l)
}

// Returns the number of elements of list.
//    m: Virtual machine.
func prSize(m *machine.T) {
	l := popList(m)
	m.Push(token.NewI(int64(len(l)), m.MkPos()))
}

// Adds an element at the end of list.
//    m: Virtual machine.
func prPush(m *machine.T) {
	tkV := m.Pop()
	tk := m.PopT(token.List)
	l, _ := tk.L()
	tk.SetL(append(l, tkV))
}

// Adds an element at the beginning of list.
//    m: Virtual machine.
func prPush0(m *machine.T) {
	tkV := m.Pop()
	tk := m.PopT(token.List)
	l, _ := tk.L()
	tk.SetL(append([]*token.T{tkV}, l...))
}

// Auxiliar function
//    m: Virtual machine.
func prPopPeek(m *machine.T, isPop bool) {
	tk := m.PopT(token.List)
	l, _ := tk.L()
	len1 := len(l) - 1
	if len1 < 0 {
		m.Failf(
			"Stack:\nExpected: List with at least 1 element.\nActual  : %v",
			tk.StringDraft(),
		)
	}
	m.Push(l[len1])
	if isPop {
		tk.SetL(l[:len1])
	}
}

// Removes and returns the last element of list.
//    m: Virtual machine.
func prPop(m *machine.T) {
	prPopPeek(m, true)
}

// Returns but not remove the last element of list.
//    m: Virtual machine.
func prPeek(m *machine.T) {
	prPopPeek(m, false)
}

func prPopPeek0(m *machine.T, isPop bool) {
	tk := m.PopT(token.List)
	l, _ := tk.L()
	length := len(l)
	if length == 0 {
		m.Failf(
			"Stack:\nExpected: List with at least 1 element.\nActual  : %v",
			tk.StringDraft(),
		)
	}
	m.Push(l[0])
	if isPop {
		tk.SetL(l[1:])
	}
}

// Removes and returns the first element of list.
//    m: Virtual machine.
func prPop0(m *machine.T) {
	prPopPeek0(m, true)
}

// Returns but not remove the first element of list.
//    m: Virtual machine.
func prPeek0(m *machine.T) {
	prPopPeek0(m, false)
}

// Inserts an element in a list.
//    m: Virtual machine.
func prInsert(m *machine.T) {
	tk1 := m.Pop()
	tk2 := m.PopT(token.Int)
	i, _ := tk2.I()
	tk3 := m.PopT(token.List)
	l, _ := tk3.L()
	ln := int64(len(l))
	if i < 0 || i > ln {
		m.Failf("Index out of range (%v out of [0-%v])", i, ln)
	}
	var rl []*token.T
	rl = append(rl, l[:i]...)
	rl = append(rl, tk1)
	tk3.SetL(append(rl, l[i:]...))
}

// Inserts a list in another list.
//    m: Virtual machine.
func prInsertList(m *machine.T) {
	subList := popList(m)
	tk2 := m.PopT(token.Int)
	i, _ := tk2.I()
	tk3 := m.PopT(token.List)
	l, _ := tk3.L()
	ln := int64(len(l))
	if i < 0 || i > ln {
		m.Failf("Index out of range (%v out of [0-%v])", i, ln)
	}
	var rl []*token.T
	rl = append(rl, l[:i]...)
	rl = append(rl, subList...)
	tk3.SetL(append(rl, l[i:]...))
}

// Removes an element in a list.
//    m: Virtual machine.
func prRemove(m *machine.T) {
	tk1 := m.PopT(token.Int)
	i, _ := tk1.I()
	tk2 := m.PopT(token.List)
	l, _ := tk2.L()
	ln := int64(len(l))
	if i < 0 || i >= ln {
		m.Failf("Index out of range (%v out of [0-%v))", i, ln)
	}
	var rl []*token.T
	rl = append(rl, l[:i]...)
	tk2.SetL(append(rl, l[i+1:]...))
}

// Removes a range [begin-end) of elements in a list.
//    m: Virtual machine.
func prRemoveRange(m *machine.T) {
	tk1 := m.PopT(token.Int)
	end, _ := tk1.I()
	tk2 := m.PopT(token.Int)
	begin, _ := tk2.I()
	tk3 := m.PopT(token.List)
	l, _ := tk3.L()
	ln := int64(len(l))
	if begin < 0 || begin >= ln {
		m.Failf("Index out of range (%v out of [0-%v))", begin, ln)
	}
	if end < 0 || end > ln {
		m.Failf("Index out of range (%v out of [0-%v])", end, ln)
	}
	if end > begin {
		var rl []*token.T
		rl = append(rl, l[:begin]...)
		tk3.SetL(append(rl, l[end:]...))
	}
}

// Removes every element in a list.
//    m: Virtual machine.
func prClear(m *machine.T) {
	tk := m.PopT(token.List)
	tk.SetL([]*token.T{})
}

// Reverse list elements.
//    m: Virtual machine.
func prReverse(m *machine.T) {
	l := popList(m)
	for i, j := 0, len(l)-1; i < j; i, j = i+1, j-1 {
		l[i], l[j] = l[j], l[i]
	}
}

// Ramdon reordering of a list.
//    m: Virtual machine.
func prShuffle(m *machine.T) {
	l := popList(m)
	i := int64(len(l))
	for {
		if i <= 1 {
			break
		}
		j := rand.Int63n(i)
		i--
		l[i], l[j] = l[j], l[i]
	}
}

// Sorts ascendantly a list using 'global.<'.
//    m: Virtual machine.
func prSort(m *machine.T, run func(m *machine.T)) {
	tk := m.PopT(token.Procedure)
	p, _ := tk.P()

	l := popList(m)
	sort.Slice(l, func(i, j int) bool {
		m2 := machine.NewIsolate(m.SourceDir, m.Pmachines, p)
		m2.Push(l[i])
		m2.Push(l[j])
		run(m2)
		tk2 := m2.PopT(token.Bool)
		r, _ := tk2.B()
		return r
	})
}

// Returns the element of list at postion ix.
//    m: Virtual machine.
func prGet(m *machine.T) {
	tkIx := m.PopT(token.Int)
	ix, _ := tkIx.I()
	tk := m.PopT(token.List)
	l, _ := tk.L()
	if int64(len(l)) < ix+1 {
		m.Failf(
			"Stack:\nExpected: List with at least %v elements.\nActual  : %v",
			ix+1, tk.StringDraft(),
		)
	}
	m.Push(l[ix])
}

// Sets the element of list at postion ix.
//    m: Virtual machine.
func prSet(m *machine.T) {
	tkV := m.Pop()
	tkIx := m.PopT(token.Int)
	ix, _ := tkIx.I()
	tk := m.PopT(token.List)
	l, _ := tk.L()
	if int64(len(l)) < ix+1 {
		m.Failf(
			"Stack:\nExpected: List with at least %v elements.\nActual  : %v",
			ix+1, tk.StringDraft(),
		)
	}
	l[ix] = tkV
}

//  Updates the element of list at position ix.
//    m: Virtual machine.
//    run: Function which running a machine.
func prUp(m *machine.T, run func(m *machine.T)) {
	tkP := m.PopT(token.Procedure)
	p, _ := tkP.P()
	tkIx := m.PopT(token.Int)
	ix, _ := tkIx.I()
	tk := m.PopT(token.List)
	l, _ := tk.L()
	if int64(len(l)) < ix+1 {
		m.Failf(
			"Stack:\nExpected: List with at least %v elements.\nActual  : %v",
			ix+1, tk.StringDraft(),
		)
	}
	m2 := machine.NewIsolate(m.SourceDir, m.Pmachines, p)
	m2.Push(l[ix])
	run(m2)
	st := *m2.Stack
	if len(st) != 1 {
		m.Failf(
			"Function update:"+
				"Expected: Return of one element.\nActual  : Return of %v elements.",
			len(st),
		)
	}
	l[ix] = st[0]
}

// Fill a list with an element. Element is cloned. (See token.Clone).
//    m: Virtual machine.
func prFill(m *machine.T) {
	tk := m.Pop()
	l := popList(m)
	for i := range l {
		l[i] = tk.Clone()
	}
}

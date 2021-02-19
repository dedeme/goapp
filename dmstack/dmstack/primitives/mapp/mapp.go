// Copyright 27-Sep-2020 ºDeme
// GNU General Public License - V3 <http://www.gnu.org/licenses/>

// Map procedures.
package mapp

import (
	"github.com/dedeme/dmstack/machine"
	"github.com/dedeme/dmstack/symbol"
	"github.com/dedeme/dmstack/token"
)

// Auxiliar function
func popMap(m *machine.T) (mp map[string]*token.T) {
	tk := m.PopT(token.Map)
	mp, _ = tk.M()
	return
}

// Auxiliar function
func pushMap(m *machine.T, mp map[string]*token.T) {
	m.Push(token.NewM(mp, m.MkPos()))
}

// Creates a new empty list.
//    m: Virtual machine.
func prNew(m *machine.T) {
	pushMap(m, map[string]*token.T{})
}

// Creates a new map from a list.
//    m: Virtual machine.
func prFrom(m *machine.T) {
	tk := m.PopT(token.Array)
	a, _ := tk.A()

	r := map[string]*token.T{}
	isKey := true
	k := ""
	ok := false
	for _, e := range a {
		if isKey {
			k, ok = e.S()
			if !ok {
				m.Failt(
					"Expected: key of type string.\nActual  : %v",
					e.StringDraft(),
				)
			}
			isKey = false
			continue
		}
		r[k] = e
		isKey = true
	}
	if !isKey {
		m.Failt("Value at end of list is missing in '%v'", tk.StringDraft())
	}

	pushMap(m, r)
}

// Creates a shallow copy.
//    m: Virtual machine.
func prCopy(m *machine.T) {
	mp := popMap(m)
	mp2 := map[string]*token.T{}
	for k, v := range mp {
		mp2[k] = v
	}
	pushMap(m, mp2)
}

// Returns the size of mp.
//    m: Virtual machine.
func prSize(m *machine.T) {
	mp := popMap(m)
	m.Push(token.NewI(int64(len(mp)), m.MkPos()))
}

// Returns the value of a key. If key is not found produces a "Map error".
//    m: Virtual machine.
func prGet(m *machine.T) {
	tk := m.PopT(token.String)
	k, _ := tk.S()
	mp := popMap(m)
	v, ok := mp[k]
	if !ok {
		m.Fail(machine.ENotFound(), "Key '%v' not found", k)
	}
	m.Push(v)
}

// Returns an option (a list with only one element) with the value of a key. If
// key is not found returns an empty option (an empty list).
//    m: Virtual machine.
func prOget(m *machine.T) {
	tk := m.PopT(token.String)
	k, _ := tk.S()
	mp := popMap(m)
	v, ok := mp[k]
	r := []*token.T{}
	if ok {
		r = append(r, v)
	}
	m.Push(token.NewA(r, m.MkPos()))
}

// Returns 'true' if 'mp' has the key 'k'.
//    m: Virtual machine.
func prHas(m *machine.T) {
	tk := m.PopT(token.String)
	k, _ := tk.S()
	mp := popMap(m)
	_, ok := mp[k]
	m.Push(token.NewB(ok, m.MkPos()))
}

// Adds a new value with the key 'k'.
//    m: Virtual machine.
func prPut(m *machine.T) {
	v := m.Pop()
	tk := m.PopT(token.String)
	k, _ := tk.S()
	mp := popMap(m)
	mp[k] = v
}

// Update the value of 'k' ('v') with the value returned for 'p' applied
// to 'v'. If key is not found produces a "Map error".
//    m: Virtual machine.
func prUp(m *machine.T, run func(m *machine.T)) {
	tk1 := m.PopT(token.Procedure)
	tk2 := m.PopT(token.String)
	k, _ := tk2.S()
	mp := popMap(m)
	v, ok := mp[k]
	if !ok {
		m.Fail(machine.ENotFound(), "Key '%v' not found", k)
	}
	m2 := machine.New(m.Source, m.Pmachines, tk1)
	m2.Push(v)
	run(m2)

	mp[k] = m2.Pop()
}

// Removes a pair with key 'k'.
//    m: Virtual machine.
func prRemove(m *machine.T) {
	tk := m.PopT(token.String)
	k, _ := tk.S()
	mp := popMap(m)
	delete(mp, k)
}

func eq(
	m *machine.T, run func(m *machine.T), tk *token.T,
	m1 map[string]*token.T, m2 map[string]*token.T,
) bool {
	if len(m1) != len(m2) {
		return false
	}
	for k, v := range m1 {
		v2, ok := m2[k]
		if !ok {
			return false
		}
		m2 := machine.New(m.Source, m.Pmachines, tk)
		m2.Push(v)
		m2.Push(v2)
		run(m2)
		tkr := m2.PopT(token.Bool)
		r, _ := tkr.B()
		if !r {
			return false
		}
	}
	return true
}

// Returns 'true' if mp1 and mp2 are equals with procedure 'p'.
func prEq(m *machine.T, run func(m *machine.T)) {
	tk := m.PopT(token.Procedure)
	mp1 := popMap(m)
	mp2 := popMap(m)
	m.Push(token.NewB(eq(m, run, tk, mp1, mp2), m.MkPos()))
}

// Returns 'true' if mp1 and mp2 are not equals with procedure 'p'.
func prNeq(m *machine.T, run func(m *machine.T)) {
	tk := m.PopT(token.Procedure)
	mp1 := popMap(m)
	mp2 := popMap(m)
	m.Push(token.NewB(!eq(m, run, tk, mp1, mp2), m.MkPos()))
}

// Returns the keys of 'mp' list.
//    m: Virtual machine.
func prKeys(m *machine.T) {
	mp := popMap(m)
	var a []*token.T
	for k := range mp {
		a = append(a, token.NewS(k, m.MkPos()))
	}
	m.Push(token.NewA(a, m.MkPos()))
}

// Returns the values of 'mp' list.
//    m: Virtual machine.
func prValues(m *machine.T) {
	mp := popMap(m)
	var a []*token.T
	for _, v := range mp {
		a = append(a, v)
	}
	m.Push(token.NewA(a, m.MkPos()))
}

// Returns a list with [key, value] elements.
//    m: Virtual machine.
func prPairs(m *machine.T) {
	pos := m.MkPos()
	mp := popMap(m)
	var a []*token.T
	for k, v := range mp {
		a = append(a, token.NewA([]*token.T{token.NewS(k, pos), v}, pos))
	}
	m.Push(token.NewA(a, m.MkPos()))
}

// Returns a list with a sequence of 'key, value', which can be reconverted
// with 'map.from'
//    m: Virtual machine.
func prTo(m *machine.T) {
	mp := popMap(m)
	var a []*token.T
	for k, v := range mp {
		a = append(a, token.NewS(k, m.MkPos()), v)
	}
	m.Push(token.NewA(a, m.MkPos()))
}

// Processes map procedures.
//    m   : Virtual machine.
//    proc: Procedure
//    run : Function which running a machine.
func Proc(m *machine.T, proc symbol.T, run func(m *machine.T)) {
	switch proc {
	case symbol.New("new"):
		prNew(m)
	case symbol.From:
		prFrom(m)
	case symbol.New("copy"):
		prCopy(m)
	case symbol.New("size"):
		prSize(m)
	case symbol.New("get"):
		prGet(m)
	case symbol.New("oget"):
		prOget(m)
	case symbol.New("has"):
		prHas(m)
	case symbol.New("put"):
		prPut(m)
	case symbol.New("up"):
		prUp(m, run)
	case symbol.New("remove"):
		prRemove(m)
	case symbol.New("eq"):
		prEq(m, run)
	case symbol.New("neq"):
		prNeq(m, run)
	case symbol.New("keys"):
		prKeys(m)
	case symbol.New("values"):
		prValues(m)
	case symbol.New("pairs"):
		prPairs(m)
	case symbol.New("to"):
		prTo(m)

	default:
		m.Failt("'map' does not contains the procedure '%v'.", proc.String())
	}
}

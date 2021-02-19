// Copyright 10-Jan-2021 ºDeme
// GNU General Public License - V3 <http://www.gnu.org/licenses/>

package js

import (
	"github.com/dedeme/dmstack/machine"
	"github.com/dedeme/dmstack/primitives/it"
	"github.com/dedeme/dmstack/token"
)

// Read a List
func prRlist(m *machine.T, run func(m *machine.T)) {
	tk := m.PopT(token.Procedure)
	prRa(m)
	tk2 := m.PopT(token.Array)
	ajs, _ := tk2.A()
	var a []*token.T
	for _, e := range ajs {
		m2 := machine.New(m.Source, m.Pmachines, tk)
		m2.Push(e)
		run(m2)
		a = append(a, m2.Pop())
	}
	m.Push(token.NewA(a, m.MkPos()))
}

// Read a Map
func prRmap(m *machine.T, run func(m *machine.T)) {
	tk := m.PopT(token.Procedure)
	prRo(m)
	tk2 := m.PopT(token.Array)
	mjs, _ := tk2.M()
	mp := map[string]*token.T{}
	for k, v := range mjs {
		m2 := machine.New(m.Source, m.Pmachines, tk)
		m2.Push(v)
		run(m2)
		mp[k] = m2.Pop()
	}
	m.Push(token.NewM(mp, m.MkPos()))
}

// Read an It
func prRit(m *machine.T, run func(m *machine.T)) {
	tk := m.PopT(token.Procedure)
	prRa(m)
	it.PrFrom(m)
	m.Push(tk)
	it.PrMap(m, run)
}

// Write a List
func prWlist(m *machine.T, run func(m *machine.T)) {
	tk := m.PopT(token.Procedure)
	tk2 := m.PopT(token.Array)
	a, _ := tk2.A()
	var ajs []*token.T
	for _, e := range a {
		m2 := machine.New(m.Source, m.Pmachines, tk)
		m2.Push(e)
		run(m2)
		ajs = append(ajs, m2.Pop())
	}
	m.Push(token.NewA(ajs, m.MkPos()))
	prWa(m)
}

// Write a Map
func prWmap(m *machine.T, run func(m *machine.T)) {
	tk := m.PopT(token.Procedure)
	tk2 := m.PopT(token.Map)
	mp, _ := tk2.M()
	mjs := map[string]*token.T{}
	for k, v := range mp {
		m2 := machine.New(m.Source, m.Pmachines, tk)
		m2.Push(v)
		run(m2)
		mjs[k] = m2.Pop()
	}
	m.Push(token.NewM(mjs, m.MkPos()))
	prWo(m)
}

// Write a It
func prWit(m *machine.T, run func(m *machine.T)) {
	it.PrMap(m, run)
	it.PrTo(m)
	prWa(m)
}

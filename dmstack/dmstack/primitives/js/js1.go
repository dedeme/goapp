// Copyright 27-Sep-2020 ºDeme
// GNU General Public License - V3 <http://www.gnu.org/licenses/>

// Json procedures.
package js

import (
	"github.com/dedeme/dmstack/machine"
	"github.com/dedeme/dmstack/primitives/it"
	"github.com/dedeme/dmstack/token"
)

// Read a List
func prRlist(m *machine.T, run func(m *machine.T)) {
	tk := m.PopT(token.Procedure)
	p, _ := tk.P()
	prRa(m)
	tk2 := m.PopT(token.List)
	ljs, _ := tk2.L()
	var l []*token.T
	for _, e := range ljs {
		m2 := machine.NewIsolate(m.SourceDir, m.Pmachines, p)
		m2.Push(e)
		run(m2)
		l = append(l, m2.Pop())
	}
	m.Push(token.NewL(l, m.MkPos()))
}

// Read a Map
func prRmap(m *machine.T, run func(m *machine.T)) {
	tk := m.PopT(token.Procedure)
	p, _ := tk.P()
	prRo(m)
	tk2 := m.PopT(token.List)
	mjs, _ := tk2.M()
	mp := map[string]*token.T{}
	for k, v := range mjs {
		m2 := machine.NewIsolate(m.SourceDir, m.Pmachines, p)
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
	p, _ := tk.P()
	tk2 := m.PopT(token.List)
	l, _ := tk2.L()
	var ljs []*token.T
	for _, e := range l {
		m2 := machine.NewIsolate(m.SourceDir, m.Pmachines, p)
		m2.Push(e)
		run(m2)
		ljs = append(ljs, m2.Pop())
	}
	m.Push(token.NewL(ljs, m.MkPos()))
	prWa(m)
}

// Write a Map
func prWmap(m *machine.T, run func(m *machine.T)) {
	tk := m.PopT(token.Procedure)
	p, _ := tk.P()
	tk2 := m.PopT(token.Map)
	mp, _ := tk2.M()
	mjs := map[string]*token.T{}
	for k, v := range mp {
		m2 := machine.NewIsolate(m.SourceDir, m.Pmachines, p)
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

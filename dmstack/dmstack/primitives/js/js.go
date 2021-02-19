// Copyright 10-Jan-2021 ºDeme
// GNU General Public License - V3 <http://www.gnu.org/licenses/>

// Json procedures.
package js

import (
	"github.com/dedeme/dmstack/machine"
	"github.com/dedeme/dmstack/operator"
	"github.com/dedeme/dmstack/symbol"
	"github.com/dedeme/dmstack/token"
)

// Auxiliar function
func popJs(m *machine.T) string {
	tk := m.PopT(token.Native)
	o, i, _ := tk.N()
	if o != operator.Js_ {
		m.Failt("\n  Expected: Json object.\n  Actual  : '%v'.", o)
	}
	return i.(string)
}

// Auxiliar function
func pushJs(m *machine.T, s string) {
	m.Push(token.NewN(operator.Js_, s, m.MkPos()))
}

// Processes js procedures.
//    m   : Virtual machine.
//    run : Function which running a machine.
func Proc(m *machine.T, proc symbol.T, run func(m *machine.T)) {
	switch proc {
	// js0.go ------------------------------------
	case symbol.New("from"):
		prFrom(m)
	case symbol.New("to"):
		prTo(m)
	case symbol.New("null"):
		prIsNull(m)
	case symbol.New("rb"):
		prRb(m)
	case symbol.New("ri"):
		prRi(m)
	case symbol.New("rf"):
		prRf(m)
	case symbol.New("rs"):
		prRs(m)
	case symbol.New("ra"):
		prRa(m)
	case symbol.New("ro"):
		prRo(m)
	case symbol.New("wn"):
		prWn(m)
	case symbol.New("wb"):
		prWb(m)
	case symbol.New("wi"):
		prWi(m)
	case symbol.New("wf"):
		prWf(m)
	case symbol.New("ws"):
		prWs(m)
	case symbol.New("wa"):
		prWa(m)
	case symbol.New("wo"):
		prWo(m)
	// js0.go ------------------------------------
	case symbol.New("rList"):
		prRlist(m, run)
	case symbol.New("rMap"):
		prRmap(m, run)
	case symbol.New("rIt"):
		prRit(m, run)
	case symbol.New("wList"):
		prWlist(m, run)
	case symbol.New("wMap"):
		prWmap(m, run)
	case symbol.New("wIt"):
		prWit(m, run)

	default:
		m.Failt("'js' does not contains the procedure '%v'.", proc.String())
	}
}

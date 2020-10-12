// Copyright 10-Aug-2020 ºDeme
// GNU General Public License - V3 <http://www.gnu.org/licenses/>

// Iterator constructor and procedures.
package it

import (
	"github.com/dedeme/dmstack/machine"
	"github.com/dedeme/dmstack/symbol"
	"github.com/dedeme/dmstack/token"
)

// It structure
type T struct {
	e    *token.T
	ok   bool
	next func() (tk *token.T, ok bool)
}

// Creates a native It token.
func New(next func() (*token.T, bool)) *T {
	e, ok := next()
	return &T{e, ok, next}
}

// Auxiliar function
func popIt(m *machine.T) *T {
	tk := m.PopT(token.Native)
	sym, i, _ := tk.N()
	if sym != symbol.Iterator_ {
		m.Failt("\n  Expected: Iterator object.\n  Actual  : '%v'.", sym)
	}
	return i.(*T)
}

// Auxiliar function
func pushIt(m *machine.T, i *T) {
	m.Push(token.NewN(symbol.Iterator_, i, m.MkPos()))
}

// Processes it procedures.
//    m: Virtual machine.
//    run: Function which running a machine.
func Proc(m *machine.T, run func(m *machine.T)) {
	tk, ok := m.PrgNext()
	if !ok {
		m.Failt("'it' procedure is missing")
	}
	sy, ok := tk.Sy()
	if !ok {
		m.Failt("\n  Expected: 'it' procedure.\n  Actual  : '%v'.", tk.StringDraft())
	}
	switch sy {
	// it0.go ------------------------------------
	case symbol.New("new"):
		prNew(m, run)
	case symbol.New("empty"):
		prEmpty(m)
	case symbol.New("unary"):
		prUnary(m)
	case symbol.From:
		PrFrom(m)
	case symbol.New("range"):
		prRange(m)
	case symbol.New("range0"):
		prRange0(m)
	case symbol.New("runes"):
		prRunes(m)
	case symbol.New("has"):
		prHas(m)
	case symbol.New("peek"):
		prPeek(m)
	case symbol.New("next"):
		prNext(m)
	// it1.go ------------------------------------
	case symbol.New("+"):
		prPlus(m)
	case symbol.New("++"):
		prPlus2(m)
	case symbol.New("drop"):
		prDrop(m)
	case symbol.New("dropf"):
		prDropf(m, run)
	case symbol.New("filter"):
		prFilter(m, run)
	case symbol.New("map"):
		PrMap(m, run)
	case symbol.New("push"):
		prPush(m)
	case symbol.New("push0"):
		prPush0(m)
	case symbol.New("take"):
		prTake(m)
	case symbol.New("takef"):
		prTakef(m, run)
	case symbol.New("zip"):
		prZip(m)
	case symbol.New("zip3"):
		prZip3(m)
	// it2.go ------------------------------------
	case symbol.New("all?"):
		prAll(m, run)
	case symbol.New("any?"):
		prAny(m, run)
	case symbol.New("contains?"):
		prContains(m, run)
	case symbol.New("each"):
		prEach(m, run)
	case symbol.New("eachIx"):
		prEachIx(m, run)
	case symbol.New("eq?"):
		prEq(m, run)
	case symbol.New("neq?"):
		prNeq(m, run)
	case symbol.Eq:
		prEquals(m)
	case symbol.Neq:
		prNequals(m)
	case symbol.New("find"):
		prFind(m, run)
	case symbol.New("index"):
		prIndex(m)
	case symbol.New("indexf"):
		prIndexF(m, run)
	case symbol.New("lastIndex"):
		prLastIndex(m)
	case symbol.New("lastIndexf"):
		prLastIndexF(m, run)
	case symbol.New("reduce"):
		prReduce(m, run)
	case symbol.New("to"):
		PrTo(m)
	case symbol.New("reverse"):
		prReverse(m)
	case symbol.New("shuffle"):
		prShuffle(m)
	case symbol.New("sort"):
		prSort(m, run)
	case symbol.New("box"):
		prBox(m)
	case symbol.New("box2"):
		prBox2(m)

	default:
		m.Failt("'it' does not contains the procedure '%v'.", sy.String())
	}
}

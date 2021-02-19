// Copyright 10-Jan-2021 ºDeme
// GNU General Public License - V3 <http://www.gnu.org/licenses/>

// Iterator constructor and procedures.
package it

import (
	"github.com/dedeme/dmstack/machine"
	"github.com/dedeme/dmstack/operator"
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
	o, i, _ := tk.N()
	if o != operator.Iterator_ {
		m.Failt("\n  Expected: Iterator object.\n  Actual  : '%v'.", o)
	}
	return i.(*T)
}

// Auxiliar function
func pushIt(m *machine.T, i *T) {
	m.Push(token.NewN(operator.Iterator_, i, m.MkPos()))
}

// Processes it procedures.
//    m   : Virtual machine.
//    proc: Procedure
//    run : Function which running a machine.
func Proc(m *machine.T, proc symbol.T, run func(m *machine.T)) {
	switch proc {
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
	case symbol.New("add"):
		prAdd(m)
	case symbol.New("join"):
		prJoin(m)
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
	case symbol.New("all"):
		prAll(m, run)
	case symbol.New("any"):
		prAny(m, run)
	case symbol.New("contains"):
		prContains(m, run)
	case symbol.New("each"):
		prEach(m, run)
	case symbol.New("eachIx"):
		prEachIx(m, run)
	case symbol.New("eqp"):
		prEqp(m, run)
	case symbol.New("neqp"):
		prNeqp(m, run)
	case symbol.New("eq"):
		prEq(m)
	case symbol.New("neq"):
		prNeq(m)
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
		m.Failt("'it' does not contains the procedure '%v'.", proc.String())
	}
}

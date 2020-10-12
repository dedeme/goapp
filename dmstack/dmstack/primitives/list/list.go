// Copyright 20-Jul-2020 ºDeme
// GNU General Public License - V3 <http://www.gnu.org/licenses/>

// List procedures.
package list

import (
	"github.com/dedeme/dmstack/machine"
	"github.com/dedeme/dmstack/symbol"
)

// Processes list procedures.
//    m: Virtual machine.
//    run: Function which running a machine.
func Proc(m *machine.T, run func(m *machine.T)) {
	tk, ok := m.PrgNext()
	if !ok {
		m.Failt("'list' procedure is missing")
	}
	sy, ok := tk.Sy()
	if !ok {
		m.Failt(
			"\n  Expected: 'list' procedure.\n  Actual  : '%v'.", tk.StringDraft(),
		)
	}
	switch sy {
	// list0.go ------------------------------------
	case symbol.New("new"):
		prNew(m)
	case symbol.New("make"):
		prMake(m)
	case symbol.New("size"):
		prSize(m)
	case symbol.New("push"):
		prPush(m)
	case symbol.New("push0"):
		prPush0(m)
	case symbol.New("pop"):
		prPop(m)
	case symbol.New("peek"):
		prPeek(m)
	case symbol.New("pop0"):
		prPop0(m)
	case symbol.New("peek0"):
		prPeek0(m)
	case symbol.New("insert"):
		prInsert(m)
	case symbol.New("insertList"):
		prInsertList(m)
	case symbol.New("remove"):
		prRemove(m)
	case symbol.New("removeRange"):
		prRemoveRange(m)
	case symbol.New("clear"):
		prClear(m)
	case symbol.New("reverse"):
		PrReverse(m)
	case symbol.New("shuffle"):
		PrShuffle(m)
	case symbol.New("sort"):
		PrSort(m, run)
	case symbol.New("get"):
		prGet(m)
	case symbol.New("set"):
		prSet(m)
	case symbol.New("up"):
		prUp(m, run)
	case symbol.New("fill"):
		prFill(m)
		// list1.go ----------------------------------
	case symbol.New("ref"):
		prRef(m)
	case symbol.New("none"):
		prNew(m) // list0.go
	case symbol.New("some"):
		prRef(m)
	case symbol.New("tp"):
		prTp(m)
	case symbol.New("tp3"):
		prTp3(m)
	case symbol.New("e1"):
		prE1(m)
	case symbol.New("e2"):
		prE2(m)
	case symbol.New("e3"):
		prE3(m)
	case symbol.New("left"):
		prLeft(m)
	case symbol.New("right"):
		prRight(m)
	case symbol.New("error"):
		prError(m)
	case symbol.New("ok"):
		prOk(m)
		// list2.go ------------------------------------
	case symbol.New("removeDup"):
		prRemoveDup(m, run)
	case symbol.New("all?"):
		prAll(m, run)
	case symbol.New("any?"):
		prAny(m, run)
	case symbol.New("each"):
		prEach(m, run)
	case symbol.New("eachIx"):
		prEachIx(m, run)
	case symbol.New("eq?"):
		prEq(m, run)
	case symbol.New("neq?"):
		prNeq(m, run)
	case symbol.New("index"):
		prIndex(m, run)
	case symbol.New("find"):
		prFind(m, run)
	case symbol.New("lastIndex"):
		prLastIndex(m, run)
	case symbol.New("reduce"):
		prReduce(m, run)
	case symbol.New("copy"):
		prCopy(m, run)
	case symbol.New("drop"):
		prDrop(m, run)
	case symbol.New("dropf"):
		prDropf(m, run)
	case symbol.New("filter"):
		prFilter(m, run)
	case symbol.New("flat"):
		prFlat(m, run)
	case symbol.New("map"):
		prMap(m, run)
	case symbol.New("sub"):
		prSub(m, run)
	case symbol.New("take"):
		prTake(m, run)
	case symbol.New("takef"):
		prTakef(m, run)
	case symbol.New("zip"):
		prZip(m, run)
	case symbol.New("zip3"):
		prZip3(m, run)
	case symbol.New("unzip"):
		prUnzip(m, run)
	case symbol.New("unzip3"):
		prUnzip3(m, run)

	default:
		m.Failt("'list' does not contains the procedure '%v'.", sy.String())
	}
}

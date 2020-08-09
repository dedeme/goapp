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
		m.Fail("'list' procedure is missing")
	}
	sy, ok := tk.Sy()
	if !ok {
		m.Failf("Expected: 'list' procedure.\nActual  : %v.", tk.StringDraft())
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
		prReverse(m)
	case symbol.New("shuffle"):
		prShuffle(m)
	case symbol.New("sort"):
		prSort(m, run)
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

	default:
		m.Failf("'list' does not contains the procedure '%v'.", sy.String())
	}
}

// Copyright 07-Jan-2021 ºDeme
// GNU General Public License - V3 <http://www.gnu.org/licenses/>

// Primitives symbols and operators
package primitives

import (
	"github.com/dedeme/dmstack/machine"
	"github.com/dedeme/dmstack/operator"
	"github.com/dedeme/dmstack/primitives/arr"
	"github.com/dedeme/dmstack/primitives/b64"
	"github.com/dedeme/dmstack/primitives/blob"
	"github.com/dedeme/dmstack/primitives/cryp"
	"github.com/dedeme/dmstack/primitives/date"
	"github.com/dedeme/dmstack/primitives/file"
	"github.com/dedeme/dmstack/primitives/float"
	"github.com/dedeme/dmstack/primitives/intp"
	"github.com/dedeme/dmstack/primitives/it"
	"github.com/dedeme/dmstack/primitives/js"
	"github.com/dedeme/dmstack/primitives/mapp"
	"github.com/dedeme/dmstack/primitives/mathp"
	"github.com/dedeme/dmstack/primitives/path"
	"github.com/dedeme/dmstack/primitives/str"
	"github.com/dedeme/dmstack/primitives/sys"
	"github.com/dedeme/dmstack/symbol"
	"github.com/dedeme/dmstack/token"
)

// Processes an Operator
//    m  : Virtual machine.
//    op : Operator to process.
//    run: Function to run a virtual machine.
func GlobalOperator(m *machine.T, op operator.T, run func(m *machine.T)) bool {
	switch op {
	// global2.go --------------------------------
	case operator.Plus:
		prPlus(m)
	case operator.Minus:
		prMinus(m)
	case operator.Mult:
		prMul(m)
	case operator.Div:
		prDiv(m)
	case operator.Mod:
		prMod(m)
	case operator.PlusPlus:
		prPlusPlus(m)
	case operator.MinusMinus:
		prMinusMinus(m)
	// global3.go --------------------------------
	case operator.Eq:
		prEq(m)
	case operator.Neq:
		prNeq(m)
	case operator.Less:
		prLess(m)
	case operator.LessEq:
		prLessEq(m)
	case operator.Greater:
		prGreater(m)
	case operator.GreaterEq:
		prGreaterEq(m)
	case operator.And:
		prAnd(m, run)
	case operator.Or:
		prOr(m, run)
	case operator.Not:
		prNot(m)
	// global4.go --------------------------------
	case operator.RefGet:
		prRefGet(m)
	case operator.RefSet:
		prRefSet(m)
	case operator.RefUp:
		prRefUp(m, run)
	// global5.go --------------------------------
	case operator.Assign:
		prAssign(m)
	case operator.ProcHeap:
		prProcHeap(m)
	case operator.StackCheck:
		prStackCheck(m)
	case operator.StackOpen:
		prStackOpen(m)
	case operator.StackClose:
		prStackClose(m)
	case operator.Stack:
		prStack(m)

	default:
		return false
	}
	return true
}

// Processes a symbol.
//    m  : Virtual machine.
//    sym: Symbol to process.
//    run: Function to run a virtual machine.
func GlobalSymbol(m *machine.T, sym symbol.T, run func(m *machine.T)) bool {
	switch sym {
	// global0.go ----------------------------------
	case symbol.Run:
		prRun(m, run)
	case symbol.Import:
		prImport(m, run)
	case symbol.If:
		prIf(m, run)
	case symbol.Else:
		m.Push(token.NewSy(sym, m.MkPos()))
	case symbol.Elif:
		prElif(m, run)
	case symbol.Loop:
		prLoop(m, run)
	case symbol.Break:
		m.Push(token.NewSy(sym, m.MkPos()))
	case symbol.Continue:
		m.Push(token.NewSy(sym, m.MkPos()))
	case symbol.While:
		prWhile(m, run)
	case symbol.For:
		prFor(m, run)
	case symbol.Sync:
		prSync(m, run)
	// global1.go --------------------------------
	case symbol.Puts:
		prPuts(m)
	case symbol.ToString:
		prToString(m)
	case symbol.Clone:
		prClone(m)
	case symbol.Assert:
		prAssert(m)
	case symbol.Expect:
		prExpect(m)
	case symbol.Fail:
		prFail(m)
	case symbol.Throw:
		prThrow(m)
	case symbol.Try:
		prTry(m, run)
	case symbol.Swap:
		prSwap(m)
	case symbol.Pop:
		prPop(m)
	case symbol.Dup:
		prDup(m)

	default:
		return false
	}
	return true
}

// Processes a symbol.
//
// It call functions in 'globalX.go'
//    m   : Virtual machine.
//    mod : Module.
//    proc: Procedure to run.
//    run : Function to run a virtual machine.
func Module(m *machine.T, mod, proc symbol.T, run func(m *machine.T)) bool {
	switch mod {
	case symbol.Arr:
		arr.Proc(m, proc, run)
	case symbol.Map:
		mapp.Proc(m, proc, run)
	case symbol.Sys:
		sys.Proc(m, proc, run)
	case symbol.Str:
		str.Proc(m, proc, run)
	case symbol.It:
		it.Proc(m, proc, run)
	case symbol.Int:
		intp.Proc(m, proc)
	case symbol.Float:
		float.Proc(m, proc)
	case symbol.Math:
		mathp.Proc(m, proc)
	case symbol.Date:
		date.Proc(m, proc)
	case symbol.Js:
		js.Proc(m, proc, run)
	case symbol.Blob:
		blob.Proc(m, proc, run)
	case symbol.B64:
		b64.Proc(m, proc)
	case symbol.Cryp:
		cryp.Proc(m, proc)
	case symbol.Path:
		path.Proc(m, proc)
	case symbol.File:
		file.Proc(m, proc)

	default:
		return false
	}
	return true
}

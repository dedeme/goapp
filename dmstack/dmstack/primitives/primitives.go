// Copyright 16-May-2020 ºDeme
// GNU General Public License - V3 <http://www.gnu.org/licenses/>

// Primitives main hub.
package primitives

import (
	"fmt"
	"github.com/dedeme/dmstack/machine"
	"github.com/dedeme/dmstack/primitives/b64"
	"github.com/dedeme/dmstack/primitives/blob"
	"github.com/dedeme/dmstack/primitives/cryp"
	"github.com/dedeme/dmstack/primitives/date"
	"github.com/dedeme/dmstack/primitives/file"
	"github.com/dedeme/dmstack/primitives/float"
	"github.com/dedeme/dmstack/primitives/intpk"
	"github.com/dedeme/dmstack/primitives/it"
	"github.com/dedeme/dmstack/primitives/js"
	"github.com/dedeme/dmstack/primitives/list"
	"github.com/dedeme/dmstack/primitives/mappk"
	"github.com/dedeme/dmstack/primitives/mathpk"
	"github.com/dedeme/dmstack/primitives/path"
	"github.com/dedeme/dmstack/primitives/str"
	"github.com/dedeme/dmstack/primitives/sys"
	"github.com/dedeme/dmstack/symbol"
	"github.com/dedeme/dmstack/token"
)

// Processes 'sym'.
//
// It call functions in 'globalX.go'
//    m: Virtual machine.
//    sym: Symbol to process.
//    run: Functions to run a virtual machine.
func Global(m *machine.T, sym symbol.T, run func(m *machine.T)) bool {
	defer func() {
		if err := recover(); err != nil {
			_, ok := err.(*machine.Error)
			if ok {
				panic(err)
			} else {
				panic(&machine.Error{
					Machine: m, Type: "Interpreter error",
					Message: fmt.Sprintf("%v", err),
				})
			}
		}
	}()

	switch sym {
	// global0.go ----------------------------------
	case symbol.Run:
		prRun(m, run)
	case symbol.Data:
		prData(m, run)
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
	case symbol.While:
		prWhile(m, run)
	case symbol.For:
		prFor(m, run)
	case symbol.Sync:
		prSync(m, run)
		// global1.go --------------------------------
	case symbol.Nop:
	case symbol.Puts:
		prPuts(m)
	case symbol.ToStr:
		prToStr(m)
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
	case symbol.StackCheck:
		prStackCheck(m)
	case symbol.StackOpen:
		prStackOpen(m)
	case symbol.StackClose:
		prStackClose(m)
	case symbol.Stack:
		prStack(m)
	// global2.go --------------------------------
	case symbol.Plus:
		prPlus(m)
	case symbol.Minus:
		prMinus(m)
	case symbol.Mult:
		prMul(m)
	case symbol.Div:
		prDiv(m)
	case symbol.Mod:
		prMod(m)
	case symbol.PlusPlus:
		prPlusPlus(m)
	case symbol.MinusMinus:
		prMinusMinus(m)
	// global3.go --------------------------------
	case symbol.Eq:
		prEq(m)
	case symbol.Neq:
		prNeq(m)
	case symbol.Less:
		prLess(m)
	case symbol.LessEq:
		prLessEq(m)
	case symbol.Greater:
		prGreater(m)
	case symbol.GreaterEq:
		prGreaterEq(m)
	case symbol.And:
		prAnd(m, run)
	case symbol.Or:
		prOr(m, run)
	case symbol.Not:
		prNot(m)
	// global4.go --------------------------------
	case symbol.RefGet:
		prRefGet(m)
	case symbol.RefSet:
		prRefSet(m)
	case symbol.RefUp:
		prRefUp(m, run)
	default:
		return false
	}
	return true
}

func Modules(m *machine.T, sym symbol.T, run func(m *machine.T)) bool {
	defer func() {
		if err := recover(); err != nil {
			_, ok := err.(*machine.Error)
			if ok {
				panic(err)
			} else {
				panic(&machine.Error{
					Machine: m, Type: "Machine error", Message: fmt.Sprintf("%v", err),
				})
			}
		}
	}()

	switch sym {
	case symbol.List:
		list.Proc(m, run)
		return true
	case symbol.Str:
		str.Proc(m, run)
		return true
	case symbol.Sys:
		sys.Proc(m, run)
		return true
	case symbol.Int:
		intpk.Proc(m)
		return true
	case symbol.Float:
		float.Proc(m)
		return true
	case symbol.Math:
		mathpk.Proc(m)
		return true
	case symbol.It:
		it.Proc(m, run)
		return true
	case symbol.Map:
		mappk.Proc(m, run)
		return true
	case symbol.Js:
		js.Proc(m, run)
		return true
	case symbol.Date:
		date.Proc(m, run)
		return true
	case symbol.Blob:
		blob.Proc(m, run)
		return true
	case symbol.B64:
		b64.Proc(m, run)
		return true
	case symbol.Cryp:
		cryp.Proc(m, run)
		return true
	case symbol.Path:
		path.Proc(m, run)
		return true
	case symbol.File:
		file.Proc(m, run)
		return true
	default:
		return false
	}
	return true
}

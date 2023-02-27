// Copyright 01-Feb-2022 ºDeme
// GNU General Public License - V3 <http://www.gnu.org/licenses/>

// Main file.
package main

import (
	"fmt"
	"github.com/dedeme/kut/checker"
	"github.com/dedeme/kut/expression"
	"github.com/dedeme/kut/fileix"
	"github.com/dedeme/kut/heap"
	"github.com/dedeme/kut/module"
	"github.com/dedeme/kut/modules"
	"github.com/dedeme/kut/reader"
	"github.com/dedeme/kut/reader/txReader"
	"github.com/dedeme/kut/runner"
	"github.com/dedeme/kut/runner/fail"
	"github.com/dedeme/kut/statement"
	"os"
	"path"
)

func help() {
	fmt.Println(
		"Usage: kut -v" +
			"\n     kut <file> [args]" +
			"\n       p.ej. kut -v" +
			"\n       p.ej. kut myprg" +
			"\n       p.ej. kut myprg arg1 arg2")
}

func main() {
	if len(os.Args) < 2 {
		help()
		return
	}

	check := false
	p := os.Args[1]
	if p == "-v" {
		fmt.Println("Kut version v2022.12")
		return
	} else if len(os.Args) == 3 && p == "-c" {
		check = true
		p = os.Args[2]
	}

	fileix.Root = path.Dir(p)
	fix := fileix.Add(-1, path.Base(p))
	var kutCode string
	var err error
	kutCode, err = fileix.Read(fix)
	if err == nil {
		modules.Add(fix)
		var mod *module.T
		mod, err = reader.ReadMainBlock(txReader.New(fix, kutCode))
		if err == nil {
			modules.Set(fix, mod)
			if check {
				checker.Run()
			} else {
				var bk, cont bool
				var stackTrace []*statement.T
				_, bk, cont, _, err, stackTrace = runner.Run(
					[]*statement.T{}, mod.Imports, mod.Heap0, []heap.T{mod.Heap}, mod.Code)
				if err == nil {
					if bk {
						err = fail.Mk("break' without 'while' or 'for'", stackTrace)
					} else if cont {
						err = fail.Mk("'continue' without 'while' or 'for'", stackTrace)
					}
				}
			}
		}
	}

	if err != nil {
		switch e := err.(type) {
		case *fail.SysErrorT:
			fn := e.Fn
			ex := expression.New(expression.ExPr, []interface{}{
				expression.MkFinal(fn),
				[]*expression.T{expression.MkFinal(e.Msg)}})
			_, err = runner.Solve(fn.Imports, fn.Hp0, fn.Hps, ex, []*statement.T{}) //runner.exSolver
			if err != nil {
				fmt.Printf("Error in custom function sys.fail:\n%v\n%v\n",
					expression.MkFinal(fn), err)
			}
		default:
			if runner.SysFffail != nil {
				e2 := fail.MkSysError(e.Error(), runner.SysFffail)
				fn := e2.Fn
				ex := expression.New(expression.ExPr, []interface{}{
					expression.MkFinal(fn),
					[]*expression.T{expression.MkFinal(e2.Msg)}})
				_, err = runner.Solve(fn.Imports, fn.Hp0, fn.Hps, ex, []*statement.T{}) //runner.exSolver
				if err != nil {
					fmt.Printf("Error in custom function sys.fail:\n%v\n%v\n",
						expression.MkFinal(fn), err)
				}
			} else {
				fmt.Println(err)
			}
		}
		return
	}
}

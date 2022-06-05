// Copyright 13-Mar-2022 ºDeme
// GNU General Public License - V3 <http://www.gnu.org/licenses/>

package runner

import (
	"fmt"
	"github.com/dedeme/kut/builtin/bfail"
	"github.com/dedeme/kut/builtin/bfunction"
	"github.com/dedeme/kut/expression"
	"github.com/dedeme/kut/function"
	"github.com/dedeme/kut/runner/fail"
	"github.com/dedeme/kut/statement"
	"sync"
)

type threadT struct {
	ch chan bool
}

var mtx sync.Mutex

// \* -> b
func threadCheckType(args []*expression.T) (ex *expression.T, err error) {
	switch (args[0].Value).(type) {
	case *threadT:
		ex = expression.MkFinal(true)
	default:
		ex = expression.MkFinal(false)
	}
	return
}

// \<thread> -> ()
func threadJoin(args []*expression.T) (ex *expression.T, err error) {
	switch th := (args[0].Value).(type) {
	case *threadT:
		<-th.ch
	default:
		err = bfail.Type(args[0], "<thread>")
	}
	return
}

// \\->() -> ()
func threadRun(args []*expression.T) (ex *expression.T, err error) {
	switch fn := (args[0].Value).(type) {
	case *function.T:
		if len(fn.Vars) != 0 {
			err = bfail.Mk(fmt.Sprintf(
				"Expected function with 0 parameters. Found %v", len(fn.Vars)))
			return
		}
		go func() {
			defer func() {
				if e := recover(); e != nil {
					err = fail.Mk(fmt.Sprintf("%v", e), []*statement.T{})
				}
			}()

			_, withBreak, withContinue, _, er, stackT :=
				RunStat([]*statement.T{}, fn.Imports, fn.Hp0, fn.Hps, fn.Stat)

			if er != nil {
			} else if withBreak {
				er = fail.Mk("break' without 'while' or 'for'", stackT)
			} else if withContinue {
				er = fail.Mk("'continue' without 'while' or 'for'", stackT)
			}

			if er != nil {
				if SysFffail != nil {
					e2 := fail.MkSysError(er.Error(), SysFffail)
					fn := e2.Fn
					ex := expression.New(expression.ExPr, []interface{}{
						expression.MkFinal(fn),
						[]*expression.T{expression.MkFinal(e2.Msg)}})
					_, er3 := Solve(fn.Imports, fn.Hp0, fn.Hps, ex, []*statement.T{})
					if er3 != nil {
						fmt.Printf("Error in custom function sys.fail:\n%v\n%v\n",
							expression.MkFinal(fn), er3)
					}
				} else {
					fmt.Println("[THREAD]:", er)
				}
			}
		}()
	default:
		err = bfail.Type(args[0], "function")
	}
	return
}

// \\->() -> <thread>
func threadStart(args []*expression.T) (ex *expression.T, err error) {
	switch fn := (args[0].Value).(type) {
	case *function.T:
		if len(fn.Vars) != 0 {
			err = bfail.Mk(fmt.Sprintf(
				"Expected function with 0 parameters. Found %v", len(fn.Vars)))
			return
		}
		th := &threadT{make(chan bool)}
		ex = expression.MkFinal(th)

		go func() {
			defer func() {
				if e := recover(); e != nil {
					err = fail.Mk(fmt.Sprintf("%v", e), []*statement.T{})
				}
				th.ch <- true
			}()

			_, withBreak, withContinue, _, er, stackT :=
				RunStat([]*statement.T{}, fn.Imports, fn.Hp0, fn.Hps, fn.Stat)

			if er != nil {
			} else if withBreak {
				er = fail.Mk("break' without 'while' or 'for'", stackT)
			} else if withContinue {
				er = fail.Mk("'continue' without 'while' or 'for'", stackT)
			}

			if er != nil {
				fmt.Println("[THREAD]:", er)
			}

			th.ch <- true
		}()
	default:
		err = bfail.Type(args[0], "function")
	}
	return
}

// \\->() -> ()
func threadSync(args []*expression.T) (ex *expression.T, err error) {
	switch fn := (args[0].Value).(type) {
	case *function.T:
		if len(fn.Vars) != 0 {
			err = bfail.Mk(fmt.Sprintf(
				"Expected function with 0 parameters. Found %v", len(fn.Vars)))
			return
		}
		mtx.Lock()
		_, err = solveIsolateFunction(fn, []*expression.T{})
		mtx.Unlock()
	default:
		err = bfail.Type(args[0], "function")
	}
	return
}

func threadGet(fname string) (fn *bfunction.T, ok bool) {
	ok = true
	switch fname {
	case "checkType":
		fn = bfunction.New(1, threadCheckType)
	case "join":
		fn = bfunction.New(1, threadJoin)
	case "run":
		fn = bfunction.New(1, threadRun)
	case "start":
		fn = bfunction.New(1, threadStart)
	case "sync":
		fn = bfunction.New(1, threadSync)
	default:
		ok = false
	}

	return
}

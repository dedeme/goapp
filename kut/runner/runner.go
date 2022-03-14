// Copyright 21-Feb-2022 ºDeme
// GNU General Public License - V3 <http://www.gnu.org/licenses/>

// Code runner.
package runner

import (
	"fmt"
	"github.com/dedeme/kut/expression"
	"github.com/dedeme/kut/fileix"
	"github.com/dedeme/kut/heap"
	"github.com/dedeme/kut/heap0"
	"github.com/dedeme/kut/runner/fail"
	"github.com/dedeme/kut/statement"
)

func freturn(
	imports map[string]int, hp0 heap0.T, hps []heap.T,
	st *statement.T, stackT []*statement.T,
) (exp *expression.T, err error) {
	if st.Value != nil {
		exp, err = Solve(imports, hp0, hps, st.Value.(*expression.T), stackT)
		if err == nil && exp.IsEmpty() {
			err = fail.Type(exp, stackT, "expression")
		}
	} else {
		exp = expression.MkEmpty()
	}
	return
}

func trace(
	imports map[string]int, hp0 heap0.T, hps []heap.T,
	st *statement.T, stackT []*statement.T,
) (err error) {
	var exp *expression.T
	exp, err = Solve(imports, hp0, hps, st.Value.(*expression.T), stackT)
	if err == nil {
		fmt.Printf("%v:%v: %v\n", fileix.Get(st.File), st.Nline, exp)
	}
	return
}

func assert(
	imports map[string]int, hp0 heap0.T, hps []heap.T,
	st *statement.T, stackT []*statement.T,
) (err error) {
	var exp *expression.T
	exp, err = Solve(imports, hp0, hps, st.Value.(*expression.T), stackT)
	if err == nil {
		switch v := exp.Value.(type) {
		case bool:
			if !v {
				err = fail.Mk("Assert failed", stackT)
			}
		default:
			err = fail.Type(exp, stackT, "bool")
		}
	}
	return
}

func assign(
	imports map[string]int, hp0 heap0.T, hps []heap.T,
	st *statement.T, stackT []*statement.T,
) (err error) {
	ps := st.Value.([]*expression.T)
	switch ps[0].Type {
	case expression.Sym:
		var ex *expression.T
		ex, err = Solve(imports, hp0, hps, ps[1], stackT)
		if err == nil {
			if !hps[0].Add(ps[0].Value.(string), ex) {
				err = fail.Mk(
					"Reassignation of variable '"+ps[0].Value.(string)+"'", stackT)
			}
			if ex.IsEmpty() {
				err = fail.Type(ex, stackT, "expression")
			}
		}
	case expression.ExSq:
		var ex, ex1, ex2 *expression.T
		ex1, err = Solve(imports, hp0, hps, (ps[0].Value).([]*expression.T)[0], stackT)
		if err == nil {
			switch v := (ex1.Value).(type) {
			case []*expression.T:
				ex2, err = Solve(imports, hp0, hps, (ps[0].Value).([]*expression.T)[1], stackT)
				if err == nil {
					switch i := (ex2.Value).(type) {
					case int64:
						ex, err = Solve(imports, hp0, hps, ps[1], stackT)
						if err == nil {
							if ex.IsEmpty() {
								err = fail.Type(ex, stackT, "expression")
							}
							v[i] = ex
						}
					default:
						err = fail.Type(ex2, stackT, "int")
					}
				}
			case map[string]*expression.T:
				ex2, err = Solve(imports, hp0, hps, (ps[0].Value).([]*expression.T)[1], stackT)
				if err == nil {
					switch i := (ex2.Value).(type) {
					case string:
						ex, err = Solve(imports, hp0, hps, ps[1], stackT)
						if err == nil {
							if ex.IsEmpty() {
								err = fail.Type(ex, stackT, "expression")
							}
							v[i] = ex
						}
					default:
						err = fail.Type(ex2, stackT, "string")
					}
				}
			default:
				err = fail.Type(ex1, stackT, "array", "map")
			}
		}
	case expression.ExPt:
		var ex, ex1, ex2 *expression.T
		ex1, err = Solve(imports, hp0, hps, (ps[0].Value).([]*expression.T)[0], stackT)
		if err == nil {
			switch v := (ex1.Value).(type) {
			case map[string]*expression.T:
				ex2 = (ps[0].Value).([]*expression.T)[1]
				if ex2.Type == expression.Sym {
					ex, err = Solve(imports, hp0, hps, ps[1], stackT)
					if err == nil {
						if ex.IsEmpty() {
							err = fail.Type(ex, stackT, "expression")
						}
						v[ex2.Value.(string)] = ex
					}
				} else {
					err = fail.Type(ex2, stackT, "symbol")
				}
			default:
				err = fail.Type(ex1, stackT, "map")
			}
		}
	default:
		err = fail.Type(ps[0], stackT, "symbol", "array", "map")
	}
	return
}

func xxxAs(
	tp expression.Type,
	imports map[string]int, hp0 heap0.T, hps []heap.T,
	st *statement.T, stackT []*statement.T,
) (err error) {
	ps := st.Value.([]*expression.T)
	switch ps[0].Type {
	case expression.ExSq:
		var ex, ex1, ex2 *expression.T
		ex1, err = Solve(imports, hp0, hps, (ps[0].Value).([]*expression.T)[0], stackT)
		if err == nil {
			switch v := (ex1.Value).(type) {
			case []*expression.T:
				ex2, err = Solve(imports, hp0, hps, (ps[0].Value).([]*expression.T)[1], stackT)
				if err == nil {
					switch i := (ex2.Value).(type) {
					case int64:
						ex, err = Solve(imports, hp0, hps, ps[1], stackT)
						if err == nil {
							ex, err = Solve(imports, hp0, hps,
								expression.New(tp, []*expression.T{v[i], ex}),
								stackT)
							if err == nil {
								v[i] = ex
							}
						}
					default:
						err = fail.Type(ex2, stackT, "int")
					}
				}
			case map[string]*expression.T:
				ex2, err = Solve(imports, hp0, hps, (ps[0].Value).([]*expression.T)[1], stackT)
				if err == nil {
					switch i := (ex2.Value).(type) {
					case string:
						ex, err = Solve(imports, hp0, hps, ps[1], stackT)
						if err == nil {
							ex, err = Solve(imports, hp0, hps,
								expression.New(tp, []*expression.T{v[i], ex}),
								stackT)
							if err == nil {
								v[i] = ex
							}
						}
					default:
						err = fail.Type(ex2, stackT, "string")
					}
				}
			default:
				err = fail.Type(ex1, stackT, "array", "map")
			}
		}
	case expression.ExPt:
		var ex, ex1, ex2 *expression.T
		ex1, err = Solve(imports, hp0, hps, (ps[0].Value).([]*expression.T)[0], stackT)
		if err == nil {
			switch v := (ex1.Value).(type) {
			case map[string]*expression.T:
				ex2 = (ps[0].Value).([]*expression.T)[1]
				if ex2.Type == expression.Sym {
					ex, err = Solve(imports, hp0, hps, ps[1], stackT)
					if err == nil {
						i := ex2.Value.(string)
						ex, err = Solve(imports, hp0, hps,
							expression.New(tp, []*expression.T{v[i], ex}),
							stackT)
						if err == nil {
							v[i] = ex
						}
					}
				} else {
					err = fail.Type(ex2, stackT, "symbol")
				}
			default:
				err = fail.Type(ex1, stackT, "array", "map")
			}
		}
	default:
		err = fail.Type(ps[0], stackT, "symbol", "array", "map")
	}
	return
}

// Run a list of statement.
//    stackTrace: Current stack trace when calling
//    imports: Module imports.
//    hp0: Module heap of not solved symbols.
//    hps: Module heaps of solved symbols.
//    code: Statements list.
func Run(
	stackTrace []*statement.T, imports map[string]int,
	hp0 heap0.T, hps []heap.T,
	code []*statement.T,
) (
	withReturn bool, withBreak bool, withContinue bool,
	ret *expression.T, err error,
	stackT []*statement.T,
) {
	hps = heap.New().AddTo(hps)
	for _, st := range code {
		withReturn, withBreak, withContinue, ret, err, stackT =
			RunStat(stackTrace, imports, hp0, hps, st)
		if withReturn || withBreak || withContinue || err != nil {
			break
		}
	}
	return
}

// Run a statement.
//    stackTrace: Current stack trace when calling
//    imports: Module imports.
//    hp0: Module heap of not solved symbols.
//    hps: Module heaps of solved symbols.
//    st: Statement to run.
// Returns:
//    withReturn: True if the statement is a return
//    withBreak: True if the statement is a break
//    withContinue: True if the statement is a continue
//    return: Expression returned
//    error: Error with stack trace.
//    stackT: New stackTrace, appending 'code' to 'stackTrace'.
func RunStat(
	stackTrace []*statement.T, imports map[string]int,
	hp0 heap0.T, hps []heap.T,
	st *statement.T,
) (
	withReturn bool, withBreak bool, withContinue bool,
	ret *expression.T, err error,
	stackT []*statement.T,
) {
	defer func() {
		if e := recover(); e != nil {
			err = fail.Mk(fmt.Sprintf("%v", e), stackT)
		}
	}()

	ret = expression.MkEmpty()
	stackT = append(stackTrace, st)
	switch st.Type {
	case statement.CloseBlock:
		err = fail.Mk("'}' without '{'", stackT)
	case statement.Break:
		withBreak = true
	case statement.Continue:
		withContinue = true
	case statement.Return:
		withReturn = true
		ret, err = freturn(imports, hp0, hps, st, stackT)
	case statement.Trace:
		err = trace(imports, hp0, hps, st, stackT)
	case statement.Assert:
		err = assert(imports, hp0, hps, st, stackT)
	case statement.Assign:
		err = assign(imports, hp0, hps, st, stackT)
	case statement.AddAs:
		err = xxxAs(expression.Add, imports, hp0, hps, st, stackT)
	case statement.SubAs:
		err = xxxAs(expression.Sub, imports, hp0, hps, st, stackT)
	case statement.MulAs:
		err = xxxAs(expression.Mul, imports, hp0, hps, st, stackT)
	case statement.DivAs:
		err = xxxAs(expression.Div, imports, hp0, hps, st, stackT)
	case statement.AndAs:
		err = xxxAs(expression.And, imports, hp0, hps, st, stackT)
	case statement.OrAs:
		err = xxxAs(expression.Or, imports, hp0, hps, st, stackT)
	case statement.Block:
		withReturn, withBreak, withContinue, ret, err, stackT =
			Run(stackT, imports, hp0, hps, st.Value.([]*statement.T))
	case statement.FunctionCalling:
		ret, err = Solve(imports, hp0, hps, st.Value.(*expression.T), stackT)
	case statement.If:
		withReturn, withBreak, withContinue, ret, err, stackT =
			runIf(stackT, imports, hp0, hps, st) // fluxRunner.go
	case statement.Else:
		err = fail.Mk("'else' without 'if'", stackT)
	case statement.While:
		withReturn, ret, err =
			runWhile(stackT, imports, hp0, hps, st) // fluxRunner.go
	case statement.For:
		withReturn, ret, err =
			runFor(stackT, imports, hp0, hps, st) // fluxRunner.go
	case statement.ForIx:
		withReturn, ret, err =
			runForIx(stackT, imports, hp0, hps, st) // fluxRunner.go
	case statement.ForR:
		withReturn, ret, err =
			runForR(stackT, imports, hp0, hps, st) // fluxRunner.go
	case statement.ForRI:
		withReturn, ret, err =
			runForRI(stackT, imports, hp0, hps, st) // fluxRunner.go
	case statement.Switch:
		withReturn, withBreak, withContinue, ret, err, stackT =
			runSwitch(stackT, imports, hp0, hps, st) // fluxRunner.go
	default:
		panic(fmt.Sprintf("Not valid statement '%v'", st))
	}
	return
}

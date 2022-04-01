// Copyright 02-Feb-2022 ºDeme
// GNU General Public License - V3 <http://www.gnu.org/licenses/>

package runner

import (
	"fmt"
	"github.com/dedeme/kut/builtin/bfunction"
	"github.com/dedeme/kut/expression"
	"github.com/dedeme/kut/function"
	"github.com/dedeme/kut/heap"
	"github.com/dedeme/kut/heap0"
	"github.com/dedeme/kut/iterator"
	"github.com/dedeme/kut/module"
	"github.com/dedeme/kut/modules"
	"github.com/dedeme/kut/runner/fail"
	"github.com/dedeme/kut/statement"
	"math"
	"strings"
)

func solveIsolateFunction(f *function.T, pars []*expression.T) (
	ex *expression.T, err error,
) {
	fex := expression.New(expression.ExPr, []interface{}{
		expression.MkFinal(f), pars})
	ex, err = Solve(f.Imports, f.Hp0, f.Hps, fex, []*statement.T{})
	return
}

func solveIterator(
	imports map[string]int, hp0 heap0.T, hps []heap.T,
	e *expression.T, stackT []*statement.T,
) (ex *expression.T, err error) {
	ps := e.Value.([]*expression.T)
	var startEx, endEx *expression.T
	startEx, err = Solve(imports, hp0, hps, ps[0], stackT)
	if err == nil {
		switch start := (startEx.Value).(type) {
		case int64:
			endEx, err = Solve(imports, hp0, hps, ps[1], stackT)
			if err == nil {
				switch end := (endEx.Value).(type) {
				case int64:
					var cmpType expression.Type
					var inc int64
					if end > start {
						inc = 1
						cmpType = expression.Less
					} else {
						inc = -1
						cmpType = expression.Greater
					}

					fhasNext := func() bool {
						if cmpType == expression.Less {
							if start < end {
								return true
							}
							return false
						} else {
							if start > end {
								return true
							}
							return false
						}
					}
					fnext := func() *expression.T {
						v := start
						start += inc
						return expression.MkFinal(v)
					}
					ex = expression.MkFinal(iterator.New(fhasNext, fnext))
				default:
					fail.Type(endEx, stackT, "int")
				}
			}
		default:
			fail.Type(startEx, stackT, "int")
		}
	}
	return
}

func solveSym(
	imports map[string]int, hp0 heap0.T, hps []heap.T,
	sym string, stackT []*statement.T,
) (ex *expression.T, err error) {
	var ok bool
	if ex, ok = heap.Get(hps, sym); ok {
		return
	}

	var e *heap0.EntryT
	if e, ok = hp0[sym]; ok {
		stk := append(stackT, statement.New(
			stackT[len(stackT)-1].File, e.Nline, statement.Assign,
			[]*expression.T{expression.New(expression.Sym, sym), e.Expr},
		))
		ex, err = Solve(imports, hp0, hps, e.Expr, stk)
		if err == nil {
			hps[len(hps)-1][sym] = ex
		}
		return
	}

	var mdix int
	if mdix, ok = imports[sym]; ok {
		ex = expression.MkFinal(modules.GetOk(mdix))
		return
	}

	var md *BModuleT
	if md, ok = GetModule(sym); ok { // mdIndex.go
		ex = expression.MkFinal(md)
		return
	}

	err = fail.Mk("Symbol '"+sym+"' not found.", stackT)
	return
}

func Solve(
	imports map[string]int, hp0 heap0.T, hps []heap.T,
	e *expression.T, stackT []*statement.T,
) (ex *expression.T, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = fail.Mk(fmt.Sprintf("%v", e), stackT)
		}
	}()

	var ex0, ex1, ex2 *expression.T

	switch e.Type {
	case expression.Final:
		ex = e
	case expression.Arr:
		var a []*expression.T
		for _, subE := range e.Value.([]*expression.T) {
			var ssubE *expression.T
			ssubE, err = Solve(imports, hp0, hps, subE, stackT)
			if err == nil {
				a = append(a, ssubE)
			} else {
				break
			}
		}
		if err == nil {
			ex = expression.MkFinal(a)
		}
	case expression.Map:
		m := map[string]*expression.T{}
		for k, v := range e.Value.(map[string]*expression.T) {
			var sv *expression.T
			sv, err = Solve(imports, hp0, hps, v, stackT)
			if err == nil {
				m[k] = sv
			} else {
				break
			}
		}
		if err == nil {
			ex = expression.MkFinal(m)
		}
	case expression.Func:
		fn := e.Value.(*function.T)
		fn.SetContext(imports, hp0, hps)
		ex = expression.MkFinal(fn)
	case expression.Sym:
		ex, err = solveSym(imports, hp0, hps, e.Value.(string), stackT)
	case expression.Range:
		ps := e.Value.([]*expression.T)
		if len(ps) == 2 {
			ex, err = solveIterator(imports, hp0, hps, e, stackT)
		} else {
			ex0, err = Solve(imports, hp0, hps, ps[0], stackT)
			if err == nil {
				ex1, err = Solve(imports, hp0, hps, ps[1], stackT)
				if err == nil {
					ex2, err = Solve(imports, hp0, hps, ps[2], stackT)
					if err == nil {
						if ex1.IsEmpty() {
							if ex2.IsEmpty() {
								switch s := (ex0.Value).(type) {
								case string:
									ex = expression.MkFinal(s)
								case []*expression.T:
									a := make([]*expression.T, len(s))
									copy(a, s)
									ex = expression.MkFinal(a)
								default:
									err = fail.Type(ex0, stackT, "string", "array")
								}
							} else {
								switch (ex0.Value).(type) {
								case string:
									fn := bfunction.New(2, strLeft)
									ex, err = fn.Run("left", []*expression.T{ex0, ex2}, stackT)
								case []*expression.T:
									fn := bfunction.New(2, arrLeft)
									ex, err = fn.Run("left", []*expression.T{ex0, ex2}, stackT)
								default:
									err = fail.Type(ex0, stackT, "string", "array")
								}
							}
						} else if ex2.IsEmpty() {
							switch (ex0.Value).(type) {
							case string:
								fn := bfunction.New(2, strRight)
								ex, err = fn.Run("right", []*expression.T{ex0, ex1}, stackT)
							case []*expression.T:
								fn := bfunction.New(2, arrRight)
								ex, err = fn.Run("right", []*expression.T{ex0, ex1}, stackT)
							default:
								err = fail.Type(ex0, stackT, "string", "array")
							}
						} else {
							switch (ex0.Value).(type) {
							case string:
								fn := bfunction.New(3, strSub)
								ex, err = fn.Run("sub", []*expression.T{ex0, ex1, ex2}, stackT)
							case []*expression.T:
								fn := bfunction.New(3, arrSub)
								ex, err = fn.Run("sub", []*expression.T{ex0, ex1, ex2}, stackT)
							default:
								err = fail.Type(ex0, stackT, "string", "array")
							}
						}
					}
				}
			}
		}
	case expression.ExPt:
		ps := e.Value.([]*expression.T)
		ex0, err = Solve(imports, hp0, hps, ps[0], stackT)
		if err == nil {
			switch v := (ex0.Value).(type) {
			case *module.T:
				ex, err = Solve(v.Imports, v.Heap0, []heap.T{v.Heap}, ps[1], stackT)
			case *BModuleT:
				ex1 = ps[1]
				if ex1.Type == expression.Sym {
					fname := (ex1.Value).(string)
					if fn, ok := GetFunction(v, fname); ok { //mdIndex.go
						ex = expression.MkFinal(fn)
					} else {
						err = fail.Mk("Function '"+v.Name+"."+fname+"' not found.", stackT)
					}
					return
				}
				if ex1.Type != expression.ExPr {
					err = fail.Type(ex1, stackT, "Function call")
					return
				}
				ps2 := ex1.Value.([]interface{})
				ex21 := ps2[0].(*expression.T)
				if ex21.Type != expression.Sym {
					err = fail.Type(ex21, stackT, "Symbol")
					return
				}
				fname := (ex21.Value).(string)

				var exs []*expression.T
				for _, ee := range ps2[1].([]*expression.T) {
					var see *expression.T
					see, err = Solve(imports, hp0, hps, ee, stackT)
					if err == nil {
						exs = append(exs, see)
					} else {
						break
					}
				}
				if err == nil {
					if fn, ok := GetFunction(v, fname); ok { // mdIndex.go
						ex, err = fn.Run(fname, exs, stackT)
						if err == nil && ex == nil {
							err = fail.Mk(
								"Function '"+fname+"' does not return any value", stackT)
						}
					} else {
						err = fail.Mk("Function '"+v.Name+"."+fname+"' not found.", stackT)
					}
				}
			case map[string]*expression.T:
				if ps[1].Type == expression.Sym {
					ex1 = expression.MkFinal(ps[1].Value)
					fn := bfunction.New(2, mapFget)
					ex, err = fn.Run("get", []*expression.T{ex0, ex1}, stackT)
				} else {
					err = fail.Type(ps[1], stackT, "symbol")
				}
			default:
				err = fail.Type(ex0, stackT, "module", "map")
			}
		}
	case expression.ExSq:
		ps := e.Value.([]*expression.T)
		ex0, err = Solve(imports, hp0, hps, ps[0], stackT)
		if err == nil {
			switch (ex0.Value).(type) {
			case []*expression.T:
				ex1, err = Solve(imports, hp0, hps, ps[1], stackT)
				if err == nil {
					fn := bfunction.New(2, arrFget)
					ex, err = fn.Run("get", []*expression.T{ex0, ex1}, stackT)
				}
			case string:
				ex1, err = Solve(imports, hp0, hps, ps[1], stackT)
				if err == nil {
					fn := bfunction.New(2, strAt)
					ex, err = fn.Run("at", []*expression.T{ex0, ex1}, stackT)
				}
			case map[string]*expression.T:
				ex1, err = Solve(imports, hp0, hps, ps[1], stackT)
				if err == nil {
					fn := bfunction.New(2, mapFget)
					ex, err = fn.Run("get", []*expression.T{ex0, ex1}, stackT)
				}
			default:
				err = fail.Type(ex0, stackT, "string", "array", "map")
			}
		}
	case expression.ExPr:
		ps := e.Value.([]interface{})
		ex0, err = Solve(imports, hp0, hps, ps[0].(*expression.T), stackT)
		if err == nil {
			switch fn := (ex0.Value).(type) {
			case *bfunction.T:
				var exs []*expression.T
				for _, ee := range ps[1].([]*expression.T) {
					var see *expression.T
					see, err = Solve(imports, hp0, hps, ee, stackT)
					if err == nil {
						exs = append(exs, see)
					} else {
						break
					}
				}
				if err == nil {
					ex, err = fn.Run(ps[0].(*expression.T).String(), exs, stackT)
					if err == nil && ex == nil {
						err = fail.Mk(
							"Function '"+ex0.String()+"' does not return any value", stackT)
					}
				}
			case *function.T:
				var exs []*expression.T
				for _, ee := range ps[1].([]*expression.T) {
					var see *expression.T
					see, err = Solve(imports, hp0, hps, ee, stackT)
					if err == nil {
						if see.IsEmpty() {
							err = fail.Type(see, stackT, "expression")
						}
						exs = append(exs, see)
					} else {
						break
					}
				}
				if err == nil {
					if len(exs) == len(fn.Vars) {
						hp2 := heap.New()
						for i, v := range fn.Vars {
							hp2.Add(v, exs[i])
						}
						var wrt, bk, cont bool
						var ret *expression.T
						var stackT2 []*statement.T
						wrt, bk, cont, ret, err, stackT2 = RunStat( // in runner.go
							stackT, fn.Imports, fn.Hp0, hp2.AddTo(fn.Hps), fn.Stat)
						if err == nil {
							if bk {
								err = fail.Mk("Unexpected 'break'", stackT2)
							} else if cont {
								err = fail.Mk("Unexpected 'continue'", stackT2)
							} else if wrt {
								ex = ret
							} else {
								ex = expression.MkEmpty()
							}
						}
					} else {
						var sexs []string
						for _, s := range exs {
							sexs = append(sexs, s.String())
						}
						err = fail.Mk(fmt.Sprintf(
							"Expected %v arguments but found %v in\n    %v(%v)",
							len(fn.Vars), len(exs), ex0, strings.Join(sexs, ", "),
						), stackT)
					}
				}
			default:
				err = fail.Type(ex0, stackT, "function")
			}
		}
	case expression.Switch:
		ps0 := e.Value.([]interface{})
		ex0, err = Solve(imports, hp0, hps, ps0[0].(*expression.T), stackT)

		if err == nil {
			ps1 := ps0[1].([][]*expression.T)
			var r *expression.T
			caseEntry := ""
			for _, ps := range ps1 {
				if ps[0].Type == expression.Sym && (ps[0].Value).(string) == "default" {
					r, err = Solve(imports, hp0, hps, ps[1], stackT)
					break
				}
				ex1, err = Solve(
					imports, hp0, hps,
					expression.New(expression.Eq, []*expression.T{ex0, ps[0]}),
					stackT)
				if err != nil {
					break
				}
				switch v := (ex1.Value).(type) {
				case bool:
					if v {
						caseEntry = ps[0].String()
						r, err = Solve(imports, hp0, hps, ps[1], stackT)
					}
				default:
					err = fail.Type(ex1, stackT, "bool")
				}
				if r != nil || err != nil {
					break
				}
			}
			if err == nil {
				if r != nil {
					if r.IsEmpty() {
						err = fail.Type(r, stackT, "expression in '"+caseEntry+"'")
					}
					ex = r
				} else {
					err = fail.Mk(fmt.Sprintf("switch did not catch '%v'", ex0), stackT)
				}
			}
		}
	case expression.Not:
		ex0, err = Solve(imports, hp0, hps, e.Value.(*expression.T), stackT)
		if err == nil {
			switch v := (ex0.Value).(type) {
			case bool:
				ex = expression.MkFinal(!v)
			default:
				err = fail.Type(ex0, stackT, "bool")
			}
		}
	case expression.Minus:
		ex0, err = Solve(imports, hp0, hps, e.Value.(*expression.T), stackT)
		if err == nil {
			switch v := (ex0.Value).(type) {
			case byte:
				ex = expression.MkFinal(-v)
			case int64:
				ex = expression.MkFinal(-v)
			case float64:
				ex = expression.MkFinal(-v)
			default:
				err = fail.Type(ex0, stackT, "int", "float")
			}
		}
	case expression.Add:
		ps := e.Value.([]*expression.T)
		ex0, err = Solve(imports, hp0, hps, ps[0], stackT)
		if err == nil {
			ex1, err = Solve(imports, hp0, hps, ps[1], stackT)
			if err == nil {
				switch v0 := (ex0.Value).(type) {
				case string:
					switch v1 := (ex1.Value).(type) {
					case string:
						ex = expression.MkFinal(v0 + v1)
					default:
						err = fail.Type(ex1, stackT, "string")
					}
				case int64:
					switch v1 := (ex1.Value).(type) {
					case int64:
						ex = expression.MkFinal(v0 + v1)
					default:
						err = fail.Type(ex1, stackT, "int")
					}
				case float64:
					switch v1 := (ex1.Value).(type) {
					case float64:
						ex = expression.MkFinal(v0 + v1)
					default:
						err = fail.Type(ex1, stackT, "float")
					}
				case []*expression.T:
					switch v1 := (ex1.Value).(type) {
					case []*expression.T:
						ex = expression.MkFinal(append(v0, v1...))
					default:
						err = fail.Type(ex1, stackT, "array")
					}
				default:
					err = fail.Type(ex0, stackT, "int", "float", "string", "array")
				}
			}
		}
	case expression.Sub, expression.Mul, expression.Div:
		ps := e.Value.([]*expression.T)
		ex0, err = Solve(imports, hp0, hps, ps[0], stackT)
		if err == nil {
			ex1, err = Solve(imports, hp0, hps, ps[1], stackT)
			if err == nil {
				switch v0 := (ex0.Value).(type) {
				case int64:
					switch v1 := (ex1.Value).(type) {
					case int64:
						switch e.Type {
						case expression.Sub:
							ex = expression.MkFinal(v0 - v1)
						case expression.Mul:
							ex = expression.MkFinal(v0 * v1)
						case expression.Div:
							ex = expression.MkFinal(v0 / v1)
						}
					default:
						err = fail.Type(ex1, stackT, "int")
					}
				case float64:
					switch v1 := (ex1.Value).(type) {
					case float64:
						switch e.Type {
						case expression.Sub:
							ex = expression.MkFinal(v0 - v1)
						case expression.Mul:
							ex = expression.MkFinal(v0 * v1)
						case expression.Div:
							ex = expression.MkFinal(v0 / v1)
						}
					default:
						err = fail.Type(ex1, stackT, "float")
					}
				default:
					err = fail.Type(ex0, stackT, "int", "float")
				}
			}
		}
	case expression.Mod:
		ps := e.Value.([]*expression.T)
		ex0, err = Solve(imports, hp0, hps, ps[0], stackT)
		if err == nil {
			ex1, err = Solve(imports, hp0, hps, ps[1], stackT)
			if err == nil {
				switch v0 := (ex0.Value).(type) {
				case int64:
					switch v1 := (ex1.Value).(type) {
					case int64:
						ex = expression.MkFinal(v0 % v1)
					default:
						err = fail.Type(ex1, stackT, "int")
					}
				default:
					err = fail.Type(ex0, stackT, "int")
				}
			}
		}
	case expression.And, expression.Or:
		ps := e.Value.([]*expression.T)
		ex0, err = Solve(imports, hp0, hps, ps[0], stackT)
		if err == nil {
			switch v0 := (ex0.Value).(type) {
			case bool:
				if v0 && e.Type == expression.Or {
					ex = expression.MkFinal(true)
				} else if !v0 && e.Type == expression.And {
					ex = expression.MkFinal(false)
				} else {
					ex1, err = Solve(imports, hp0, hps, ps[1], stackT)
					if err == nil {
						switch v1 := (ex1.Value).(type) {
						case bool:
							switch e.Type {
							case expression.And:
								ex = expression.MkFinal(v0 && v1)
							case expression.Or:
								ex = expression.MkFinal(v0 || v1)
							}
						default:
							err = fail.Type(ex1, stackT, "bool")
						}
					}
				}
			default:
				err = fail.Type(ex0, stackT, "bool")
			}
		}
	case expression.Greater, expression.GreaterEq, expression.Less,
		expression.LessEq, expression.Eq, expression.Neq:
		ps := e.Value.([]*expression.T)
		ex0, err = Solve(imports, hp0, hps, ps[0], stackT)
		if err == nil {
			ex1, err = Solve(imports, hp0, hps, ps[1], stackT)
			if err == nil {
				switch v0 := (ex0.Value).(type) {
				case string:
					switch v1 := (ex1.Value).(type) {
					case string:
						switch e.Type {
						case expression.Greater:
							ex = expression.MkFinal(v0 > v1)
						case expression.GreaterEq:
							ex = expression.MkFinal(v0 >= v1)
						case expression.Less:
							ex = expression.MkFinal(v0 < v1)
						case expression.LessEq:
							ex = expression.MkFinal(v0 <= v1)
						case expression.Eq:
							ex = expression.MkFinal(v0 == v1)
						case expression.Neq:
							ex = expression.MkFinal(v0 != v1)
						}
					default:
						err = fail.Type(ex1, stackT, "string")
					}
				case bool:
					switch v1 := (ex1.Value).(type) {
					case bool:
						switch e.Type {
						case expression.Greater:
							ex = expression.MkFinal(v0 && !v1)
						case expression.GreaterEq:
							ex = expression.MkFinal(v0 || !v1)
						case expression.Less:
							ex = expression.MkFinal(!v0 && v1)
						case expression.LessEq:
							ex = expression.MkFinal(!v0 || v1)
						case expression.Eq:
							ex = expression.MkFinal(v0 == v1)
						case expression.Neq:
							ex = expression.MkFinal(v0 != v1)
						}
					default:
						err = fail.Type(ex1, stackT, "bool")
					}
				case int64:
					switch v1 := (ex1.Value).(type) {
					case int64:
						switch e.Type {
						case expression.Greater:
							ex = expression.MkFinal(v0 > v1)
						case expression.GreaterEq:
							ex = expression.MkFinal(v0 >= v1)
						case expression.Less:
							ex = expression.MkFinal(v0 < v1)
						case expression.LessEq:
							ex = expression.MkFinal(v0 <= v1)
						case expression.Eq:
							ex = expression.MkFinal(v0 == v1)
						case expression.Neq:
							ex = expression.MkFinal(v0 != v1)
						}
					default:
						err = fail.Type(ex1, stackT, "int")
					}
				case float64:
					switch v1 := (ex1.Value).(type) {
					case float64:
						switch e.Type {
						case expression.Greater:
							ex = expression.MkFinal(v0 > v1)
						case expression.GreaterEq:
							ex = expression.MkFinal(v0 >= v1)
						case expression.Less:
							ex = expression.MkFinal(v0 < v1)
						case expression.LessEq:
							ex = expression.MkFinal(v0 <= v1)
						case expression.Eq:
							ex = expression.MkFinal(math.Abs(v0-v1) <= 0.0000001)
						case expression.Neq:
							ex = expression.MkFinal(math.Abs(v0-v1) > 0.0000001)
						}
					default:
						err = fail.Type(ex1, stackT, "float")
					}
				default:
					err = fail.Type(ex0, stackT, "bool", "int", "float", "string")
				}
			}
		}
	case expression.Ternary:
		ps := e.Value.([]*expression.T)
		ex0, err = Solve(imports, hp0, hps, ps[0], stackT)
		if err == nil {
			switch v := (ex0.Value).(type) {
			case bool:
				if v {
					ex1, err = Solve(imports, hp0, hps, ps[1], stackT)
				} else {
					ex1, err = Solve(imports, hp0, hps, ps[2], stackT)
				}
				if err == nil {
					if ex1.IsEmpty() {
						err = fail.Type(ex1, stackT, "expression")
					}
					ex = ex1
				}
			default:
				err = fail.Type(ex1, stackT, "bool")
			}
		}
	default:
		panic(fmt.Sprintf("runner.Solve: '%v'. Unknown expresion type", e))
	}
	return
}

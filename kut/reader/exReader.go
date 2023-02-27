// Copyright 02-Feb-2022 ºDeme
// GNU General Public License - V3 <http://www.gnu.org/licenses/>

package reader

import (
	"github.com/dedeme/kut/expression"
	"github.com/dedeme/kut/function"
	"github.com/dedeme/kut/reader/token"
	"github.com/dedeme/kut/reader/txReader"
	"github.com/dedeme/kut/statement"
)

func solveLevel1(exs []interface{}) []interface{} {
	if len(exs) == 1 {
		return exs
	}
	newExs := []interface{}{exs[0]}
	for i := 1; i < len(exs); i += 2 {
		tk := exs[i].(*token.T)
		if tk.IsBinary1() {
			len1 := len(newExs) - 1
			last := newExs[len1]
			if tk.Value.(string) == "*" {
				newExs = append(
					newExs[0:len1],
					expression.New(expression.Mul, []*expression.T{
						last.(*expression.T),
						exs[i+1].(*expression.T),
					}))
			} else if tk.Value.(string) == "/" {
				newExs = append(
					newExs[0:len1],
					expression.New(expression.Div, []*expression.T{
						last.(*expression.T),
						exs[i+1].(*expression.T),
					}))
			} else {
				newExs = append(
					newExs[0:len1],
					expression.New(expression.Mod, []*expression.T{
						last.(*expression.T),
						exs[i+1].(*expression.T),
					}))
			}
		} else {
			newExs = append(newExs, tk, exs[i+1])
		}
	}
	return newExs
}

func solveLevel2(exs []interface{}) []interface{} {
	if len(exs) == 1 {
		return exs
	}
	newExs := []interface{}{exs[0]}
	for i := 1; i < len(exs); i += 2 {
		tk := exs[i].(*token.T)
		if tk.IsBinary2() {
			len1 := len(newExs) - 1
			last := newExs[len1]
			if tk.Value.(string) == "+" {
				newExs = append(
					newExs[0:len1],
					expression.New(expression.Add, []*expression.T{
						last.(*expression.T),
						exs[i+1].(*expression.T),
					}))
			} else {
				newExs = append(
					newExs[0:len1],
					expression.New(expression.Sub, []*expression.T{
						last.(*expression.T),
						exs[i+1].(*expression.T),
					}))
			}
		} else {
			newExs = append(newExs, tk, exs[i+1])
		}
	}
	return newExs
}

func solveLevel3(exs []interface{}) []interface{} {
	if len(exs) == 1 {
		return exs
	}
	newExs := []interface{}{exs[0]}
	for i := 1; i < len(exs); i += 2 {
		tk := exs[i].(*token.T)
		if tk.IsBinary3() {
			len1 := len(newExs) - 1
			last := newExs[len1]
			if tk.Value.(string) == "==" {
				newExs = append(
					newExs[0:len1],
					expression.New(expression.Eq, []*expression.T{
						last.(*expression.T),
						exs[i+1].(*expression.T),
					}))
			} else if tk.Value.(string) == "!=" {
				newExs = append(
					newExs[0:len1],
					expression.New(expression.Neq, []*expression.T{
						last.(*expression.T),
						exs[i+1].(*expression.T),
					}))
			} else if tk.Value.(string) == ">" {
				newExs = append(
					newExs[0:len1],
					expression.New(expression.Greater, []*expression.T{
						last.(*expression.T),
						exs[i+1].(*expression.T),
					}))
			} else if tk.Value.(string) == ">=" {
				newExs = append(
					newExs[0:len1],
					expression.New(expression.GreaterEq, []*expression.T{
						last.(*expression.T),
						exs[i+1].(*expression.T),
					}))
			} else if tk.Value.(string) == "<" {
				newExs = append(
					newExs[0:len1],
					expression.New(expression.Less, []*expression.T{
						last.(*expression.T),
						exs[i+1].(*expression.T),
					}))
			} else {
				newExs = append(
					newExs[0:len1],
					expression.New(expression.LessEq, []*expression.T{
						last.(*expression.T),
						exs[i+1].(*expression.T),
					}))
			}
		} else {
			newExs = append(newExs, tk, exs[i+1])
		}
	}
	return newExs
}

func solveLevel4(exs []interface{}) *expression.T {
	ex := exs[0].(*expression.T)
	if len(exs) == 1 {
		return ex
	}

	for i := 1; i < len(exs); i += 2 {
		tk := exs[i].(*token.T)
		if tk.Value.(string) == "&" {
			ex = expression.New(expression.And, []*expression.T{
				ex,
				exs[i+1].(*expression.T),
			})
		} else {
			ex = expression.New(expression.Or, []*expression.T{
				ex,
				exs[i+1].(*expression.T),
			})
		}
	}
	return ex
}

func readExpression1(tk *token.T, tx *txReader.T) (
	empty bool, ex *expression.T, nextTk *token.T, err error,
) {
	if tk == nil {
		var eof bool
		tk, eof, err = tx.ReadToken()
		if err != nil {
			return
		}
		if eof {
			err = tx.Fail("Unexpected end of file")
			return
		}
	}

	if tk.IsOpenPar() {
		nline := tx.Nline
		empty, ex, nextTk, err = readExpression(tx)
		if err != nil {
			return
		}
		if !nextTk.IsClosePar() {
			err = tx.FailExpect(")", nextTk.String(), nline)
			return
		}

		var eof bool
		nextTk, eof, err = tx.ReadToken()
		if eof {
			err = tx.Fail("Unexpected end of file")
		}

		return
	}

	if tk.IsUnary() {
		empty, ex, nextTk, err = readExpression1(nil, tx)
		if err == nil {
			if empty {
				err = tx.FailExpect("Expresion", nextTk.String(), tx.Nline)
				return
			}
			ex, nextTk, err = readPtSqPr(ex, nextTk, tx) // ptSqPrReader.go
			if err != nil {
				return
			}

			if tk.Value.(string) == "!" {
				ex = expression.New(expression.Not, ex)
			} else {
				ex = expression.New(expression.Minus, ex)
			}
		}
		return
	}

	if tk.IsOpenSquare() {
		nline := tx.Nline
		var vs []*expression.T

		empty, ex, tk, err = readExpression(tx)
		if err != nil {
			return
		}

		exType := expression.Arr
		if !empty {
			if tk.Type == token.Operator && tk.Value.(string) == ":" {
				exType = expression.Range
				vs = append(vs, ex)

				empty, ex, tk, err = readExpression(tx)
				if err != nil {
					return
				}
				if empty {
					err = tx.FailExpect("Expresion", tk.String(), tx.Nline)
					return
				}
			} else {
				for {
					if !tk.IsComma() {
						break
					}
					vs = append(vs, ex)
					empty, ex, tk, err = readExpression(tx)
					if err != nil {
						return
					}
					if empty {
						err = tx.FailExpect("Expresion", tk.String(), tx.Nline)
						return
					}
				}
			}
			vs = append(vs, ex)
		}

		if tk.IsCloseSquare() {
			empty = false
			ex = expression.New(exType, vs)

			var eof bool
			nextTk, eof, err = tx.ReadToken()
			if eof {
				err = tx.Fail("Unexpected end of file")
			}

			return
		}

		err = tx.FailExpect("']'", tk.String(), tx.Nline)
		err = tx.FailLine(err.Error()+"\n'[' not closed", nline)
		return
	}

	if tk.IsOpenBracket() {
		nline := tx.Nline
		vs := map[string]*expression.T{}

		var kex *expression.T
		empty, kex, tk, err = readExpression(tx)
		if err != nil {
			return
		}

		if !empty {
			for {
				ex = kex
				if kex.Type == expression.Sym {
					kex = expression.New(expression.Final, kex.Value)
				} else {
					ex = nil
				}
				if !kex.IsString() {
					err = tx.FailExpect("Object key", kex.String(), tx.Nline)
					return
				}
				if !tk.IsColon() {
					if ex == nil {
						err = tx.FailExpect("':'", tk.String(), tx.Nline)
						return
					}
				} else {
					empty, ex, tk, err = readExpression(tx)
					if err != nil {
						return
					}
					if empty {
						err = tx.FailExpect("Expresion", tk.String(), tx.Nline)
						return
					}
				}

				if !tk.IsComma() {
					break
				}
        key := kex.Value.(string);
        _, ok := vs[key];
        if ok {
          err = tx.FailExpect(
            "A new key", "The duplicate key '" + key + "'", tx.Nline)
          return
        } else {
          vs[kex.Value.(string)] = ex
        }
				empty, kex, tk, err = readExpression(tx)
				if err != nil {
					return
				}
				if empty {
					err = tx.FailExpect("Expresion", tk.String(), tx.Nline)
					return
				}
			}
      key := kex.Value.(string);
      _, ok := vs[key];
      if ok {
        err = tx.FailExpect(
          "A new key", "The duplicate key '" + key + "'", tx.Nline)
        return
      } else {
        vs[key] = ex
      }
		}

		if tk.IsCloseBracket() {
			empty = false
			ex = expression.New(expression.Map, vs)

			var eof bool
			nextTk, eof, err = tx.ReadToken()
			if eof {
				err = tx.Fail("Unexpected end of file")
			}

			return
		}

		err = tx.FailExpect("'}'", tk.String(), tx.Nline)
		err = tx.FailLine(err.Error()+"\n'{' not closed", nline)
		return
	}

	if tk.IsBackSlash() {
		var eof bool
		var vars []string
		tk, eof, err = tx.ReadToken()
		if err != nil {
			return
		}
		if eof {
			err = tx.FailExpect("->", tk.String(), tx.Nline)
			return
		}

		if !tk.IsArrow() {
			for {
				if tk.Type != token.Symbol {
					err = tx.FailExpect("Parameter name", tk.String(), tx.Nline)
					return
				}
				vars = append(vars, tk.Value.(string))

				tk, eof, err = tx.ReadToken()
				if err != nil {
					return
				}
				if eof {
					err = tx.FailExpect("->", tk.String(), tx.Nline)
					return
				}

				if tk.IsArrow() {
					break
				}
				if !tk.IsComma() {
					err = tx.FailExpect("'->' or ','", tk.String(), tx.Nline)
					return
				}

				tk, eof, err = tx.ReadToken()
				if err != nil {
					return
				}
				if eof {
					err = tx.Fail("Unexpected end of file")
					return
				}
			}
		}

		var st *statement.T
		st, nextTk, eof, err = readStatement(nil, tx)
		if err != nil {
			return
		}
		if eof {
			err = tx.Fail("Unexpected end of file")
			return
		}

		empty = false
		ex = expression.New(expression.Func, function.New(vars, st))
		if nextTk == nil {
			nextTk, eof, err = tx.ReadToken()
			if err == nil && eof {
				err = tx.Fail("Unexpected end of file")
			}
		}

		return
	}

	switch tk.Type {
	case token.Comment, token.LineComment:
		empty, ex, nextTk, err = readExpression1(nil, tx)
	case token.Bool:
		ex = expression.New(expression.Final, tk.Value)
	case token.Int:
		ex = expression.New(expression.Final, tk.Value)
	case token.Float:
		ex = expression.New(expression.Final, tk.Value)
	case token.String:
		ex = expression.New(expression.Final, tk.Value)
	case token.Symbol:
		if tk.Value.(string) == "switch" {
			ex, err = readSwitchExpresion(tx) // switchReader.go
		} else {
			ex = expression.New(expression.Sym, tk.Value.(string))
		}
	default:
		nextTk = tk
		empty = true
	}

	if err == nil && nextTk == nil {
		var eof bool
		nextTk, eof, err = tx.ReadToken()
		if eof {
			err = tx.Fail("Unexpected end of file")
		}
	}

	return
}

func addExpressions(tk *token.T, iexs []interface{}, tx *txReader.T) (
	exs []interface{}, nextTk *token.T, err error,
) {
	exs = append(exs, iexs...)
	if tk == nil {
		var eof bool
		tk, eof, err = tx.ReadToken()
		if err != nil {
			return
		}
		if eof {
			err = tx.Fail("Unexpected end of file")
			return
		}
	}

	var empty bool
	var ex *expression.T
	empty, ex, nextTk, err = readExpression1(tk, tx)
	if err != nil {
		return
	}
	if empty {
		if len(exs) > 0 {
			err = tx.FailExpect("Expresion", nextTk.String(), tx.Nline)
		}
		return
	}

	ex, nextTk, err = readPtSqPr(ex, nextTk, tx) // ptSqPrReader.go
	if err != nil {
		return
	}
	exs = append(exs, ex)

	if nextTk.IsBinary() {
		exs = append(exs, nextTk)
		exs, nextTk, err = addExpressions(nil, exs, tx)
	}

	return
}

func readExpression2(tk *token.T, tx *txReader.T) (
	empty bool, ex *expression.T, nextTk *token.T, err error,
) {
	var exs []interface{}
	exs, nextTk, err = addExpressions(tk, []interface{}{}, tx)
	if err != nil {
		return
	}
	if len(exs) == 0 {
		empty = true
		return
	}

	exs = solveLevel1(exs) // * / %
	exs = solveLevel2(exs) // + ++ -
	exs = solveLevel3(exs) // == != > >= < <=
	ex = solveLevel4(exs)  // & |

	if nextTk.IsTernary() {
		var ex1 *expression.T
		empty, ex1, nextTk, err = readExpression(tx)
		if err != nil {
			return
		}
		if empty {
			err = tx.FailExpect("Expresion", nextTk.String(), tx.Nline)
			return
		}
		if !nextTk.IsColon() {
			err = tx.FailExpect(":", nextTk.String(), tx.Nline)
			return
		}

		var ex2 *expression.T
		empty, ex2, nextTk, err = readExpression(tx)
		if err != nil {
			return
		}
		if empty {
			err = tx.FailExpect("Expresion", nextTk.String(), tx.Nline)
			return
		}

		ex = expression.New(expression.Ternary, []*expression.T{ex, ex1, ex2})
	}

	return
}

func readExpression(tx *txReader.T) (
	empty bool, ex *expression.T, nextTk *token.T, err error,
) {
	empty, ex, nextTk, err = readExpression2(nil, tx)
	return
}

func readExpressionSeq(tx *txReader.T) (
	exs []*expression.T, nextTk *token.T, err error,
) {
	var empty bool
	var ex *expression.T
	empty, ex, nextTk, err = readExpression(tx)
	if err != nil || empty {
		return
	}

	exs = append(exs, ex)
	for {
		if nextTk.Type != token.Operator || nextTk.Value.(string) != "," {
			break
		}

		empty, ex, nextTk, err = readExpression(tx)
		if err != nil {
			break
		}
		if empty {
			err = tx.FailExpect("Expresion", nextTk.String(), tx.Nline)
			return
		}
		exs = append(exs, ex)
	}

	return
}

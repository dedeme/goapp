// Copyright 02-Feb-2022 ºDeme
// GNU General Public License - V3 <http://www.gnu.org/licenses/>

package reader

import (
	"github.com/dedeme/kut/expression"
	"github.com/dedeme/kut/reader/token"
	"github.com/dedeme/kut/reader/txReader"
)

func readPtSqPr(expr *expression.T, tk *token.T, tx *txReader.T) (
	ex *expression.T, nextTk *token.T, err error,
) {
	if tk == nil {
		var eof bool
		tk, eof, err = tx.ReadToken()
		if err != nil {
			return
		}
		if eof {
			err = tx.Fail("Unexpected end of file.")
			return
		}
	}

	var empty bool
	var ex2 *expression.T
	switch {
	case tk.Type == token.Operator && tk.Value.(string) == ".":
		empty, ex2, nextTk, err = readExpression1(nil, tx)
		if err == nil {
			if empty {
				err = tx.FailExpect("expresion", nextTk.String(), tx.Nline)
			} else {
				if ex2.Type != expression.Sym {
					err = tx.FailExpect("symbol", ex2.String(), tx.Nline)
				} else {
					ex, nextTk, err = readPtSqPr(expression.New(
						expression.ExPt, []*expression.T{expr, ex2}), nextTk, tx)
				}
			}
		}
	case tk.Type == token.Operator && tk.Value.(string) == "!":
		ex2 = expression.MkFinal(int64(0))
		ex, nextTk, err = readPtSqPr(expression.New(
			expression.ExSq, []*expression.T{expr, ex2}), nil, tx)
	case tk.IsOpenSquare():
		var eof bool
		tk, eof, err = tx.ReadToken()
		if err != nil {
			return
		}
		if eof {
			err = tx.Fail("Unexpected end of file")
			return
		}
		if tk.Type == token.Operator && tk.Value.(string) == ":" {
			tk, eof, err = tx.ReadToken()
			if err != nil {
				return
			}
			if eof {
				err = tx.Fail("Unexpected end of file")
				return
			}
			if tk.Type == token.Operator && tk.Value.(string) == "]" {
				ex, nextTk, err = readPtSqPr(expression.New(
					expression.Range,
					[]*expression.T{expr, expression.MkEmpty(), expression.MkEmpty()}),
					nil, tx)
			} else {
				empty, ex2, nextTk, err = readExpression2(tk, tx)
				if err == nil {
					if empty {
						err = tx.FailExpect("expresion", nextTk.String(), tx.Nline)
					} else {
						if !nextTk.IsCloseSquare() {
							err = tx.FailExpect("]", nextTk.String(), tx.Nline)
						} else {
							ex, nextTk, err = readPtSqPr(expression.New(
								expression.Range,
								[]*expression.T{expr, expression.MkEmpty(), ex2}), nil, tx)
						}
					}
				}
			}
		} else {
			empty, ex2, nextTk, err = readExpression2(tk, tx)
			if err == nil {
				if empty {
					err = tx.FailExpect("expresion", nextTk.String(), tx.Nline)
				} else {
					if !nextTk.IsCloseSquare() {
						if nextTk.Type == token.Operator && nextTk.Value.(string) == ":" {
							tk, eof, err = tx.ReadToken()
							if err != nil {
								return
							}
							if eof {
								err = tx.Fail("Unexpected end of file")
								return
							}
							if tk.Type == token.Operator && tk.Value.(string) == "]" {
								ex, nextTk, err = readPtSqPr(expression.New(
									expression.Range,
									[]*expression.T{expr, ex2, expression.MkEmpty()}),
									nil, tx)
							} else {
								var ex3 *expression.T
								empty, ex3, nextTk, err = readExpression2(tk, tx)
								if err == nil {
									if empty {
										err = tx.FailExpect("expresion", nextTk.String(), tx.Nline)
									} else {
										if !nextTk.IsCloseSquare() {
											err = tx.FailExpect("]", nextTk.String(), tx.Nline)
										} else {
											ex, nextTk, err = readPtSqPr(expression.New(
												expression.Range,
												[]*expression.T{expr, ex2, ex3}), nil, tx)
										}
									}
								}
							}
						} else {
							err = tx.FailExpect("']' or ':'", nextTk.String(), tx.Nline)
						}
					} else {
						ex, nextTk, err = readPtSqPr(expression.New(
							expression.ExSq, []*expression.T{expr, ex2}), nil, tx)
					}
				}
			}
		}
	case tk.IsOpenPar():
		var exs []*expression.T
		exs, nextTk, err = readExpressionSeq(tx) //exReader.go
		if err == nil {
			if !nextTk.IsClosePar() {
				err = tx.FailExpect(")", nextTk.String(), tx.Nline)
			} else {
				ex, nextTk, err = readPtSqPr(expression.New(
					expression.ExPr, []interface{}{expr, exs}), nil, tx)
			}
		}
	default:
		nextTk = tk
		ex = expr
	}

	return
}

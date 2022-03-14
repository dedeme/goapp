// Copyright 02-Feb-2022 ºDeme
// GNU General Public License - V3 <http://www.gnu.org/licenses/>

package reader

import (
	"fmt"
	"github.com/dedeme/kut/expression"
	"github.com/dedeme/kut/fileix"
	"github.com/dedeme/kut/reader/token"
	"github.com/dedeme/kut/reader/txReader"
	"github.com/dedeme/kut/statement"
	"path"
	"strings"
)

func readSymbol(sym string, tx *txReader.T) (
	st *statement.T, nextTk *token.T, err error,
) {
	nline := tx.Nline

	var tk *token.T
	var eof bool
	var empty bool
	var ex *expression.T
	switch sym {
	case "break":
		tk, eof, err = tx.ReadToken()
		if err == nil {
			if eof {
				err = tx.Fail("Unexpected end of file.")
			} else if !tk.IsSemicolon() {
				err = tx.FailExpect(";", tk.String(), tx.Nline)
			} else {
				st = statement.New(tx.File, nline, statement.Break, nil)
			}
		}
		return
	case "continue":
		tk, eof, err = tx.ReadToken()
		if err == nil {
			if eof {
				err = tx.Fail("Unexpected end of file.")
			} else if !tk.IsSemicolon() {
				err = tx.FailExpect(";", tk.String(), tx.Nline)
			} else {
				st = statement.New(tx.File, nline, statement.Continue, nil)
			}
		}
		return
	case "trace":
		empty, ex, tk, err = readExpression(tx) // exReader.go
		if err == nil {
			if empty {
				err = tx.FailExpect("Expresion", tk.String(), tx.Nline)
			} else if !tk.IsSemicolon() {
				err = tx.FailExpect(";", tk.String(), tx.Nline)
			} else {
				st = statement.New(tx.File, nline, statement.Trace, ex)
			}
		}
		return
	case "assert":
		empty, ex, tk, err = readExpression(tx) // exReader.go
		if err == nil {
			if empty {
				err = tx.FailExpect("Expresion", tk.String(), tx.Nline)
			} else if !tk.IsSemicolon() {
				err = tx.FailExpect(";", tk.String(), tx.Nline)
			} else {
				st = statement.New(tx.File, nline, statement.Assert, ex)
			}
		}
		return
	case "return":
		empty, ex, tk, err = readExpression(tx) // exReader.go
		if err == nil {
			if !tk.IsSemicolon() {
				err = tx.FailExpect(";", tk.String(), tx.Nline)
			} else {
				if empty {
					st = statement.New(tx.File, nline, statement.Return, nil)
				} else {
					st = statement.New(tx.File, nline, statement.Return, ex)
				}
			}
		}
		return
	case "while":
		st, err = readWhile(nline, tx) // fluxReader.go
		return
	case "if":
		st, err = readIf(nline, tx) // fluxReader.go
		if err == nil {
			nextTk, eof, err = tx.ReadToken()
			if err == nil {
				if eof {
					nextTk = nil
				}
				var st2 *statement.T
				if nextTk != nil && nextTk.IsElse() {
					st2, err = readElse(nline, tx) // fluxReader.go
				}
				if err == nil && st2 != nil {
					(st.Value.([]interface{}))[2] = st2.Value
					nextTk = nil
				}
			}
		}
		return
	case "else":
		st, err = readElse(nline, tx) // fluxReader.go
		return
	case "for":
		st, err = readFor(nline, tx) // fluxReader.go
		return
	case "switch":
		st, err = readSwitchStatement(nline, tx) // switchReader.go
		return
	case "import":
		tk, eof, err = tx.ReadToken()
		if err == nil {
			if eof {
				err = tx.Fail("Unexpected end of file.")
			} else if tk.Type != token.String {
				err = tx.FailExpect("String", tk.String(), tx.Nline)
			} else {
				fpath := tk.Value.(string)
				tk, eof, err = tx.ReadToken()
				if fpath == "" {
					err = tx.Fail("File path is empty")
				} else if fpath[len(fpath)-1] == '/' {
					err = tx.Fail("File path ends with '/'")
				} else if strings.IndexByte(fpath, '.') != -1 {
					err = tx.Fail("File path contains dots")
				}
				if err == nil {
					if eof {
						err = tx.Fail("Unexpected end of file.")
					} else if tk.IsSemicolon() {
						st = statement.New(tx.File, nline, statement.Import,
							[]interface{}{fileix.Add(fpath), path.Base(fpath)})
					} else if !tk.IsColon() {
						err = tx.FailExpect("':' or ';'", tk.String(), tx.Nline)
					} else {
						tk, eof, err = tx.ReadToken()
						if err == nil {
							if eof {
								err = tx.Fail("Unexpected end of file.")
							} else if tk.Type != token.Symbol {
								err = tx.FailExpect("Symbol", tk.String(), tx.Nline)
							} else {
								id := tk.Value.(string)
								tk, eof, err = tx.ReadToken()
								if err == nil {
									if eof {
										err = tx.Fail("Unexpected end of file.")
									} else if !tk.IsSemicolon() {
										err = tx.FailExpect(";", tk.String(), tx.Nline)
									} else {
										st = statement.New(tx.File, nline, statement.Import,
											[]interface{}{fileix.Add(fpath), id})
									}
								}
							}
						}
					}
				}
			}
		}
		return
	}

	// No resereved symbol

	ex, tk, err = readPtSqPr( // ptSqPrReader.go
		expression.New(expression.Sym, sym), nil, tx)

	if err != nil {
		return
	}

	if tk.Type == token.Operator {
		switch tk.Value.(string) {
		case ";":
			if !ex.IsFunctionCall() {
				err = tx.FailExpect("Function calling", ex.String(), nline)
				return
			}
			st = statement.New(
				tx.File, nline, statement.FunctionCalling, ex)
			return
		case "=":
			if ex.Type != expression.Sym && ex.Type != expression.ExSq &&
				ex.Type != expression.ExPt {
				err = tx.Fail("Unexpected '='")
			} else {
				st, err = readAssign(statement.Assign, nline, ex, tx) // assignReader.go
			}
			return
		case "+=", "-=", "*=", "/=", "&=", "|=":
			if ex.Type != expression.ExPt && ex.Type != expression.ExSq {
				err = tx.Fail("'" + tk.Value.(string) +
					"' only is applicable to array and map values")
			} else {
				tp := statement.OrAs
				switch tk.Value.(string) {
				case "+=":
					tp = statement.AddAs
				case "-=":
					tp = statement.SubAs
				case "*=":
					tp = statement.MulAs
				case "/=":
					tp = statement.DivAs
				case "&=":
					tp = statement.AndAs
				}
				st, err = readAssign(tp, nline, ex, tx) // assignReader.go
			}
			return
		}
	}

	err = tx.Fail(fmt.Sprintf("Unexpected '%v'", tk))
	return
}

func readStatement(tk *token.T, tx *txReader.T) (
	st *statement.T, nextTk *token.T, eof bool, err error,
) {
	if tk == nil {
		tk, eof, err = tx.ReadToken()
		if err != nil || eof {
			return
		}
	}

	switch tk.Type {
	case token.LineComment, token.Comment:
		st, _, eof, err = readStatement(nil, tx)
	case token.Symbol:
		st, nextTk, err = readSymbol(tk.Value.(string), tx)
	case token.Operator:
		switch tk.Value.(string) {
		case ";":
			st = statement.New(tx.File, tx.Nline, statement.Empty, nil)
		case "{":
			var stats []*statement.T
			stats, err = ReadBlock(tx)
			if err == nil {
				st = statement.New(tx.File, tx.Nline, statement.Block, stats)
			}
		case "}":
			st = statement.New(tx.File, tx.Nline, statement.CloseBlock, nil)
		default:
			err = tx.Fail(fmt.Sprintf("Unexpected '%v'", tk))
		}
	default:
		err = tx.Fail(fmt.Sprintf("Unexpected '%v'", tk))
	}

	return
}

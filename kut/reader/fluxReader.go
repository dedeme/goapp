// Copyright 02-Feb-2022 ºDeme
// GNU General Public License - V3 <http://www.gnu.org/licenses/>

package reader

import (
	"github.com/dedeme/kut/expression"
	"github.com/dedeme/kut/reader/token"
	"github.com/dedeme/kut/reader/txReader"
	"github.com/dedeme/kut/statement"
)

func readTry(nline int, tx *txReader.T) (
	st *statement.T, nextTk *token.T, err error,
) {
	var tk *token.T
	var varName string
	var eof bool
	var st1, st2, st3 *statement.T

	st1, nextTk, eof, err = readStatement(nil, tx) // stReader.go
	if err != nil {
		return
	}
	if eof {
		err = tx.Fail("Unexpected end of file.")
		return
	}

	tk, eof, err = tx.ReadToken()
	if err != nil {
		return
	}
	if eof {
		err = tx.Fail("Unexpected end of file.")
		return
	}
	if tk.Value.(string) != "catch" {
		err = tx.FailExpect("catch", tk.String(), tx.Nline)
		return
	}

	tk, eof, err = tx.ReadToken()
	if err != nil {
		return
	}
	if eof {
		err = tx.Fail("Unexpected end of file.")
		return
	}
	if tk.Type != token.Symbol {
		err = tx.FailExpect("Parameter name", tk.String(), tx.Nline)
		return
	}
	varName = tk.Value.(string)

	st2, nextTk, eof, err = readStatement(nil, tx) // stReader.go
	if err != nil {
		return
	}
	if eof {
		err = tx.Fail("Unexpected end of file.")
		return
	}

	tk, eof, err = tx.ReadToken()
	if err != nil {
		return
	}

	if !eof {
		if tk.Value.(string) == "finally" {
			st3, nextTk, eof, err = readStatement(nil, tx) // stReader.go
			if err != nil {
				return
			}
			if eof {
				err = tx.Fail("Unexpected end of file.")
				return
			}
		} else {
			nextTk = tk
		}
	}

	st = statement.New(tx.File, nline, statement.Try,
		[]interface{}{st1, varName, st2, st3})

	return
}

func readWhile(nline int, tx *txReader.T) (
	st *statement.T, nextTk *token.T, err error,
) {
	var tk *token.T
	var eof bool
	tk, eof, err = tx.ReadToken()

	if err != nil {
		return
	}
	if eof {
		err = tx.Fail("Unexpected end of file.")
		return
	}
	if !tk.IsOpenPar() {
		err = tx.FailExpect("(", tk.String(), tx.Nline)
		return
	}

	var ex *expression.T
	var empty bool
	empty, ex, tk, err = readExpression(tx) // exReader.go

	if err != nil {
		return
	}
	if !tk.IsClosePar() {
		err = tx.FailExpect(")", tk.String(), tx.Nline)
		return
	}
	if empty {
		ex = nil
	}

	st, nextTk, eof, err = readStatement(nil, tx) // stReader.go
	if err == nil {
		if eof {
			err = tx.Fail("Unexpected end of file")
		} else {
			st = statement.New(tx.File, nline, statement.While, []interface{}{ex, st})
		}
	}
	return
}

func readIf(nline int, tx *txReader.T) (
	st *statement.T, nextTk *token.T, err error,
) {
	var tk *token.T
	var eof bool
	tk, eof, err = tx.ReadToken()

	if err != nil {
		return
	}
	if eof {
		err = tx.Fail("Unexpected end of file.")
		return
	}
	if !tk.IsOpenPar() {
		err = tx.FailExpect("(", tk.String(), tx.Nline)
		return
	}

	var ex *expression.T
	var empty bool
	empty, ex, tk, err = readExpression(tx) // exReader.go

	if err != nil {
		return
	}
	if !tk.IsClosePar() {
		err = tx.FailExpect(")", tk.String(), tx.Nline)
		return
	}
	if empty {
		err = tx.FailExpect("If expresion", "()", tx.Nline)
		return
	}

	st, nextTk, eof, err = readStatement(nil, tx) // stReader.go
	if err == nil {
		if eof {
			err = tx.Fail("Unexpected end of file")
		} else {
			st = statement.New(tx.File, nline, statement.If, []interface{}{ex, st, nil})
		}
	}
	return
}

func readElse(nline int, tx *txReader.T) (
	st *statement.T, nextTk *token.T, err error,
) {
	var eof bool
	st, nextTk, eof, err = readStatement(nil, tx) // stReader.go
	if err == nil {
		if eof {
			err = tx.Fail("Unexpected end of file")
		} else {
			st = statement.New(tx.File, nline, statement.Else, st)
		}
	}
	return
}

func readFor(nline int, tx *txReader.T) (
	st *statement.T, nextTk *token.T, err error,
) {
	var tk *token.T
	var eof bool
	tk, eof, err = tx.ReadToken()

	if err != nil {
		return
	}
	if eof {
		err = tx.Fail("Unexpected end of file.")
		return
	}
	if !tk.IsOpenPar() {
		err = tx.FailExpect("(", tk.String(), tx.Nline)
		return
	}

	var var1, var2 string
	var ex1, ex2, ex3 *expression.T
	var stType statement.Type

	tk, eof, err = tx.ReadToken()
	if err != nil {
		return
	}
	if eof {
		err = tx.Fail("Unexpected end of file.")
		return
	}
	if tk.Type != token.Symbol {
		err = tx.FailExpect("Symbol", tk.String(), tx.Nline)
		return
	}
	var1 = tk.Value.(string)

	tk, eof, err = tx.ReadToken()
	if err != nil {
		return
	}
	if eof {
		err = tx.Fail("Unexpected end of file.")
		return
	}
	if tk.IsComma() {
		tk, eof, err = tx.ReadToken()
		if err != nil {
			return
		}
		if eof {
			err = tx.Fail("Unexpected end of file.")
			return
		}
		if tk.Type != token.Symbol {
			err = tx.FailExpect("Symbol", tk.String(), tx.Nline)
			return
		}
		var2 = tk.Value.(string)
		stType = statement.ForIx

		tk, eof, err = tx.ReadToken()
		if err != nil {
			return
		}
		if eof {
			err = tx.Fail("Unexpected end of file.")
			return
		}
	}

	if !tk.IsEquals() {
		err = tx.FailExpect("'='", tk.String(), tx.Nline)
		return
	}

	var empty bool
	empty, ex1, tk, err = readExpression(tx) // exReader.go
	if err != nil {
		return
	}
	if empty {
		err = tx.FailExpect("Expression", tk.String(), tx.Nline)
		return
	}
	if tk.Type == token.Operator && tk.Value.(string) == ":" {
		if stType == statement.ForIx {
			err = tx.Fail("Ranges are not allowed in 'for (i, e : ...)'")
			return
		}
		empty, ex2, tk, err = readExpression(tx) // exReader.go
		if err != nil {
			return
		}
		if empty {
			err = tx.FailExpect("Expression", tk.String(), tx.Nline)
			return
		}
		if tk.Type == token.Operator && tk.Value.(string) == ":" {
			empty, ex3, tk, err = readExpression(tx) // exReader.go
			if err != nil {
				return
			}
			if empty {
				err = tx.FailExpect("Expression", tk.String(), tx.Nline)
				return
			}
			stType = statement.ForRS
		} else {
			stType = statement.ForR
		}
	} else {
		if stType != statement.ForIx {
			stType = statement.For
		}
	}

	if !tk.IsClosePar() {
		err = tx.FailExpect(")", tk.String(), tx.Nline)
		return
	}

	st, nextTk, eof, err = readStatement(nil, tx) // stReader.go
	if err == nil {
		if eof {
			err = tx.Fail("Unexpected end of file")
		} else {
			if stType == statement.ForIx {
				st = statement.New(
					tx.File, nline, stType, []interface{}{var1, var2, ex1, st})
			} else if stType == statement.For {
				st = statement.New(
					tx.File, nline, stType, []interface{}{var1, ex1, st})
			} else if stType == statement.ForR {
				st = statement.New(
					tx.File, nline, stType, []interface{}{var1, ex1, ex2, st})
			} else {
				st = statement.New(
					tx.File, nline, stType, []interface{}{var1, ex1, ex2, ex3, st})
			}
		}
	}
	return
}

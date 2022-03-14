// Copyright 02-Feb-2022 ºDeme
// GNU General Public License - V3 <http://www.gnu.org/licenses/>

package reader

import (
	"github.com/dedeme/kut/expression"
	"github.com/dedeme/kut/reader/token"
	"github.com/dedeme/kut/reader/txReader"
	"github.com/dedeme/kut/statement"
)

func readSwitchExpresion(tx *txReader.T) (ex *expression.T, err error) {
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

	var ex0 *expression.T
	var empty bool
	empty, ex0, tk, err = readExpression(tx)
	if err != nil {
		return
	}
	if !tk.IsClosePar() {
		err = tx.FailExpect(")", tk.String(), tx.Nline)
		return
	}
	if empty {
		err = tx.FailExpect("(expression)", "()", tx.Nline)
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
	if !tk.IsOpenBracket() {
		err = tx.FailExpect("{", tk.String(), tx.Nline)
		return
	}

	var conds [][]*expression.T
	for {
		var cond1, cond2 *expression.T
		empty, cond1, tk, err = readExpression(tx)
		if err != nil {
			return
		}
		if empty {
			if !tk.IsCloseBracket() {
				err = tx.FailExpect("}", tk.String(), tx.Nline)
				return
			}
			break
		}
		if !tk.IsColon() {
			err = tx.FailExpect(":", tk.String(), tx.Nline)
			return
		}

		empty, cond2, tk, err = readExpression(tx)
		if err != nil {
			return
		}
		if empty {
			err = tx.FailExpect("Expresion", tk.String(), tx.Nline)
			return
		}
		if !tk.IsSemicolon() {
			err = tx.FailExpect(";", tk.String(), tx.Nline)
			return
		}
		conds = append(conds, []*expression.T{cond1, cond2})
	}

	ex = expression.New(expression.Switch, []interface{}{ex0, conds})
	return
}

func readSwitchStatement(nline int, tx *txReader.T) (
	st *statement.T, err error,
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
	empty, ex, tk, err = readExpression(tx)
	if err != nil {
		return
	}
	if !tk.IsClosePar() {
		err = tx.FailExpect(")", tk.String(), tx.Nline)
		return
	}
	if empty {
		err = tx.FailExpect("(expression)", "()", tx.Nline)
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
	if !tk.IsOpenBracket() {
		err = tx.FailExpect("{", tk.String(), tx.Nline)
		return
	}

	var conds [][]interface{}
	var nextTk *token.T
	for {
		var cond *expression.T
		empty, cond, tk, err = readExpression2(nextTk, tx)

		if err != nil {
			return
		}
		if empty {
			if !tk.IsCloseBracket() {
				err = tx.FailExpect("}", tk.String(), tx.Nline)
				return
			}
			break
		}
		if !tk.IsColon() {
			err = tx.FailExpect(":", tk.String(), tx.Nline)
			return
		}

		var stin *statement.T
		stin, nextTk, eof, err = readStatement(nil, tx) // stReader.go
		if err != nil {
			return
		}
		if eof {
			err = tx.Fail("Unexpected end of file.")
			return
		}

		conds = append(conds, []interface{}{cond, stin})
	}

	st = statement.New(tx.File, nline, statement.Switch, []interface{}{ex, conds})
	return
}

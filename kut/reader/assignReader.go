// Copyright 08-Feb-2022 ºDeme
// GNU General Public License - V3 <http://www.gnu.org/licenses/>

package reader

import (
	"github.com/dedeme/kut/expression"
	"github.com/dedeme/kut/reader/token"
	"github.com/dedeme/kut/reader/txReader"
	"github.com/dedeme/kut/statement"
)

func readAssign(
	tp statement.Type, nline int, symEx *expression.T, tx *txReader.T,
) (
	st *statement.T, err error,
) {
	var empty bool
	var ex *expression.T
	var nextTk *token.T
	empty, ex, nextTk, err = readExpression(tx) // exReader.go
	if err != nil {
		return
	}
	if empty {
		err = tx.FailExpect("Expresion", nextTk.String(), tx.Nline)
		return
	}

	if nextTk.Type != token.Operator || nextTk.Value.(string) != ";" {
		err = tx.FailExpect(";", nextTk.String(), tx.Nline)
		return
	}

	st = statement.New(
		tx.File, nline, tp, []*expression.T{symEx, ex},
	)
	return
}

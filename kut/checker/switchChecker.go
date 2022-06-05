// Copyright 31-Mar-2022 ºDeme
// GNU General Public License - V3 <http://www.gnu.org/licenses/>

// Statement checker.
package checker

import (
	"github.com/dedeme/kut/checker/cksym"
	"github.com/dedeme/kut/reader/token"
	"github.com/dedeme/kut/reader/txReader"
)

func checkSwitch(
	isExpression bool, layers [][]*cksym.T, tx *txReader.T,
) (errs []error) {
	tk := readToken(tx) // stChecker.go

	// tk IsOpenPar

	tk, errs = checkExpression(layers, tx) // exChecker.go

	// tk IsClosePar

	tk = readToken(tx) // stChecker.go

	// tk IsOpenBracket

	var ers []error
	tk = nil
	for {
		if tk == nil {
			tk = readToken(tx) // stChecker.go
		}

		if tk.IsCloseBracket() {
			break
		}

		if tk.Type == token.Symbol && tk.Value.(string) == "default" {
			tk = readToken(tx) // stChecker.go
		} else {
			tk, ers = checkExpression2(tk, layers, tx) // exChecker.go
			errs = append(errs, ers...)
		}

		// tk IsColon

		if isExpression {
			tk, ers = checkExpression(layers, tx) // exChecker.go
			// tk IsSemicolon
			tk = nil
		} else {
			_, tk, ers = checkStatement(false, false, nil, layers, tx) // stChecker.go
		}
		errs = append(errs, ers...)

	}

	return
}

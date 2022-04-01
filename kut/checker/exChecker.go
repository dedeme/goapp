// Copyright 30-Mar-2022 ºDeme
// GNU General Public License - V3 <http://www.gnu.org/licenses/>

// Statement checker.
package checker

import (
	"github.com/dedeme/kut/checker/cksym"
	"github.com/dedeme/kut/reader/token"
	"github.com/dedeme/kut/reader/txReader"
)

func checkExpression1(tk *token.T, layers [][]*cksym.T, tx *txReader.T) (
	sym *cksym.T, nextTk *token.T, errs []error,
) {
	if tk == nil {
		tk = readToken(tx) // stChecker.go
	}

	var ers []error
	if tk.IsOpenPar() {
		_, errs = checkExpression(layers, tx)

		nextTk = readToken(tx) // stChecker.go
		return
	}

	if tk.IsUnary() {
		sym, nextTk, errs = checkExpression1(nil, layers, tx)
		sym, nextTk, ers = checkPtSqPr(sym, nextTk, layers, tx) // ptSqPrChecker.go
		errs = append(errs, ers...)
		return
	}

	if tk.IsOpenSquare() {
		tk, errs = checkExpression(layers, tx)

		if tk.Type == token.Operator && tk.Value.(string) == ":" {
			tk, ers = checkExpression(layers, tx)
			errs = append(errs, ers...)
		} else {
			for {
				if !tk.IsComma() {
					break
				}
				tk, ers = checkExpression(layers, tx)
				errs = append(errs, ers...)
			}
		}

		// tk IsCloseSquare()
		nextTk = readToken(tx) // stChecker.go
		return
	}

	if tk.IsOpenBracket() {
		// Skip key
		tkKey := readToken(tx) //stChecker.go
    nline := tx.Nline;

		if tkKey.IsCloseBracket() {
			nextTk = readToken(tx) // stChecker.go
			return
		}

		// Read colon
		tk = readToken(tx) //stChecker.go

		for {
      if !tk.IsColon() {
        sym := cksym.New(tkKey.Value.(string), tx.File, nline)
        if e := cksym.ErrIfNotFound(layers, sym); e != nil {
          errs = append(errs, e)
        }
      } else {
        // tk IsColon

        tk, ers = checkExpression(layers, tx)
        errs = append(errs, ers...)
      }

			if !tk.IsComma() {
				break
			}

			// Skip key
			tkKey = readToken(tx) //stChecker.go
      nline = tx.Nline;
			// Read colon
			tk = readToken(tx) //stChecker.go
		}

		// tk IsCloseBracket
		nextTk = readToken(tx) // stChecker.go
		return
	}

	if tk.IsBackSlash() {
		tk = readToken(tx) // stChecker.go

		if !tk.IsArrow() {
			var newLayer []*cksym.T
			for {
				// tk is Symbol
				newLayer = append(newLayer, cksym.New(
					tk.Value.(string), tx.File, tx.Nline))

				tk = readToken(tx) // stChecker.go

				if tk.IsArrow() {
					break
				}

				// tk IsComma

				tk = readToken(tx) // stChecker.go
			}
			layers = append(layers, newLayer)
		}

		_, nextTk, ers = checkStatement(false, false, nil, layers, tx)
		errs = append(errs, ers...)

    if nextTk == nil {
      nextTk = readToken(tx) // stChecker.go
    }

		return
	}

	switch tk.Type {
	case token.Comment, token.LineComment:
		_, nextTk, errs = checkExpression1(nil, layers, tx)
	case token.Bool:
	case token.Int:
	case token.Float:
	case token.String:
	case token.Symbol:
		if tk.Value.(string) == "switch" {
			errs = checkSwitch(true, layers, tx) // switchChecker.go
		} else {
			sym = cksym.New(tk.Value.(string), tx.File, tx.Nline)
		}
	default:
		nextTk = tk
	}

	if nextTk == nil {
		nextTk = readToken(tx) // stChecker.go
	}
	return
}

func checkExpressions(tk *token.T, layers [][]*cksym.T, tx *txReader.T) (
	nextTk *token.T, errs []error,
) {
	if tk == nil {
		tk = readToken(tx) // stChecker.go
	}

	var sym *cksym.T
	sym, nextTk, errs = checkExpression1(tk, layers, tx)

	var ers []error
	sym, nextTk, ers = checkPtSqPr(sym, nextTk, layers, tx) // ptSqPrChecker.go
	errs = append(errs, ers...)

	if sym != nil {
		if e := cksym.ErrIfNotFound(layers, sym); e != nil {
			errs = append(errs, e)
		}
	}

	if nextTk.IsBinary() {
		nextTk, ers = checkExpressions(nil, layers, tx)
		errs = append(errs, ers...)
	}

	return
}

func checkExpression2(tk *token.T, layers [][]*cksym.T, tx *txReader.T) (
	nextTk *token.T, errs []error,
) {
	nextTk, errs = checkExpressions(tk, layers, tx)

	if nextTk.IsTernary() {
		_, ers := checkExpression(layers, tx)
		errs = append(errs, ers...)

		nextTk, ers = checkExpression(layers, tx)
		errs = append(errs, ers...)
	}
	return
}

func checkExpression(layers [][]*cksym.T, tx *txReader.T) (
	nextTk *token.T, errs []error,
) {
	nextTk, errs = checkExpression2(nil, layers, tx)
	return
}

func checkExpressionSeq(layers [][]*cksym.T, tx *txReader.T) (
	nextTk *token.T, errs []error,
) {
	nextTk, errs = checkExpression(layers, tx)

	var ers []error
	for {
		if nextTk.Type != token.Operator || nextTk.Value.(string) != "," {
			break
		}

		nextTk, ers = checkExpression(layers, tx)
		errs = append(errs, ers...)
	}

	return
}

// Copyright 30-Mar-2022 ºDeme
// GNU General Public License - V3 <http://www.gnu.org/licenses/>

// Statement checker.
package checker

import (
	"fmt"
	"github.com/dedeme/kut/checker/cksym"
	"github.com/dedeme/kut/reader/token"
	"github.com/dedeme/kut/reader/txReader"
)

func readToken(tx *txReader.T) *token.T {
	tk, eof, err := tx.ReadToken()
	if err != nil {
		panic(err)
	}
	if eof {
		panic(tx.Fail("Unexpected end of file"))
	}

	return tk
}

func checkMainSymbol(
	isTop bool, layers [][]*cksym.T, symName string, tx *txReader.T,
) (nextTk *token.T, errs []error) {
	var ers []error
	switch symName {
	case "break":
		readToken(tx)
		return
	case "continue":
		readToken(tx)
		return
	case "trace":
		_, errs = checkExpression(layers, tx) // exChecker.go
		return
	case "assert":
		_, errs = checkExpression(layers, tx) // exChecker.go
		return
	case "return":
		_, errs = checkExpression(layers, tx) // exChecker.go
		return
	case "while":
		readToken(tx)
		// read '('
		_, errs = checkExpression(layers, tx) // exChecker.go
		// read ')'
		_, nextTk, ers = checkStatement(false, false, nil, layers, tx)
		errs = append(errs, ers...)
		return
	case "if":
		readToken(tx)
		// read '('
		_, errs = checkExpression(layers, tx) // exChecker.go
		// read ')'
		_, nextTk, ers = checkStatement(false, false, nil, layers, tx)
		errs = append(errs, ers...)
    if nextTk == nil {
      var eof bool
      nextTk, eof, _ = tx.ReadToken()
      if eof {
        nextTk = nil
        return
      }
    }
    if nextTk.IsElse() {
      _, nextTk, errs = checkStatement(false, false, nil, layers, tx)
    }
		return
	case "for":
		var forLayer []*cksym.T
		readToken(tx)
		// Read '('

		tk := readToken(tx)
		// tk is Symbol
    sym := cksym.New(tk.Value.(string), tx.File, tx.Nline)
    sym.Used = true
		forLayer = append(forLayer, sym)

		tk = readToken(tx)
		if tk.IsComma() {
			tk = readToken(tx)
			// tk is Symbol
      sym = cksym.New(tk.Value.(string), tx.File, tx.Nline)
      sym.Used = true
			forLayer = append(forLayer, sym)

			tk = readToken(tx)
		}
		// tk is '='

		tk, errs = checkExpression(layers, tx) // exChecker.go
		if tk.Type == token.Operator && tk.Value.(string) == ":" {
			tk, ers = checkExpression(layers, tx) // exChecker.go
			errs = append(errs, ers...)
		}
		// tk is ')'

		_, nextTk, ers = checkStatement(false, true, nil, append(layers, forLayer), tx)
		errs = append(errs, ers...)
		return
	case "switch":
		errs = checkSwitch(false, layers, tx)
		return
	case "import":
		tk := readToken(tx)
    //tk is string
		tk = readToken(tx)
    if !tk.IsSemicolon() {
      readToken(tx)
      // Read Symbol
      readToken(tx)
      // Read ';'
    }
		return
	}

	// No resereved symbol

  nline := tx.Nline
	sym := cksym.New(symName, tx.File, tx.Nline)
	var tk *token.T
	sym, tk, ers = checkPtSqPr(sym, nil, layers, tx) // ptSqPrChecker.go
	errs = append(errs, ers...)
	switch tk.Value.(string) {
	case "=":
		if sym != nil {
      if isTop {
        for _, s := range layers[0] {
          if s.Name == sym.Name {
            s.Nline = nline
            s.Used = true
          }
        }
			} else {
				e := cksym.ErrIfFound(layers, sym)
				if e == nil {
					n := len(layers) - 1
					layers[n] = append(layers[n], sym)
				} else {
					errs = append(errs, e)
				}
			}
		}
		_, ers = checkExpression(layers, tx) // exChecker.go
		errs = append(errs, ers...)
	case "+=", "-=", "*=", "/=", "&=", "|=":
		_, ers = checkExpression(layers, tx) // exChecker.go
		errs = append(errs, ers...)
	}

	return
}

func checkStatement(
	isTop, isFor bool, tk *token.T, layers [][]*cksym.T, tx *txReader.T,
) (end bool, nextTk *token.T, errs []error) {
  if tk == nil {
    if isTop {
      var eof bool
      tk, eof, _ = tx.ReadToken()
      if eof {
        end = true
        return
      }
    } else {
      tk = readToken(tx)
    }
  }

	switch tk.Type {
	case token.LineComment, token.Comment:
	case token.Symbol:
		nextTk, errs = checkMainSymbol(isTop, layers, tk.Value.(string), tx)
	case token.Operator:
		switch tk.Value.(string) {
		case ";":
		case "{":
			errs = checkCode(false, isFor, layers, tx) //checker.go
		case "}":
			end = true
		default:
			panic(tx.Fail(fmt.Sprintf("Unexpected '%v'", tk)))
		}
	default:
		panic(tx.Fail(fmt.Sprintf("Unexpected '%v'", tk)))
	}

	return
}

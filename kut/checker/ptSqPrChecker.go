// Copyright 30-Mar-2022 ºDeme
// GNU General Public License - V3 <http://www.gnu.org/licenses/>

// Statement checker.
package checker

import (
	"fmt"
	"github.com/dedeme/kut/checker/cksym"
	"github.com/dedeme/kut/modules"
	"github.com/dedeme/kut/reader/token"
	"github.com/dedeme/kut/reader/txReader"
	"github.com/dedeme/kut/runner"
)

// syn and s can be nil
func checkPtSqPr(sym *cksym.T, tk *token.T, layers [][]*cksym.T, tx *txReader.T,
) (
	s *cksym.T, nextTk *token.T, errs []error,
) {
	var ers []error
	if tk == nil {
		tk = readToken(tx) // stChecker.go
	}

	switch {
	case tk.Type == token.Operator && tk.Value.(string) == ".":
		tk = readToken(tx)
		nline := tx.Nline
		var sym2 *cksym.T
		sym2, tk, errs = checkExpression1(tk, layers, tx) // exChecker.go
		// sym2 can not be nil
		md := modules.GetOk(tx.File)
		if sym != nil {
			if mdIx, ok := md.Imports[sym.Name]; ok {
				md = modules.GetOk(mdIx)
				if _, ok := md.Heap0[sym2.Name]; ok {
					_, nextTk, ers = checkPtSqPr(nil, tk, layers, tx)
					errs = append(errs, ers...)
					return
				}
				if mdIx, ok := md.Imports[sym2.Name]; ok {
					for {
						if tk.Type == token.Operator && tk.Value.(string) == "." {
							sym2, tk, ers = checkExpression1(tk, layers, tx) // exChecker.go
							errs = append(errs, ers...)
							if sym2 != nil {
								md = modules.GetOk(mdIx)
								if mdIx, ok = md.Imports[sym2.Name]; ok {
									continue
								}
							}
						}
						break
					}
					_, nextTk, ers = checkPtSqPr(nil, tk, layers, tx)
					errs = append(errs, ers...)
					return
				}
				errs = append(errs, tx.FailLine(fmt.Sprintf(
					"Symbol '%v' not found", sym2.Name), nline))
				_, nextTk, ers = checkPtSqPr(nil, tk, layers, tx)
				errs = append(errs, ers...)
				return
			}
			if md, ok := runner.GetModule(sym.Name); ok {
				_, ok = runner.GetFunction(md, sym2.Name)
				if !ok {
					errs = append(errs, tx.FailLine(fmt.Sprintf(
						"Symbol '%v' not found", sym2.Name), nline))
				}
				_, nextTk, ers = checkPtSqPr(nil, tk, layers, tx)
				errs = append(errs, ers...)
				return
			}

			// if map, do nothing
		}
		_, nextTk, ers = checkPtSqPr(nil, tk, layers, tx)
		errs = append(errs, ers...)
	case tk.Type == token.Operator && tk.Value.(string) == "!":
		_, nextTk, errs = checkPtSqPr(nil, nil, layers, tx)
	case tk.IsOpenSquare():
		tk = readToken(tx) // stChecker.go

		if tk.Type == token.Operator && tk.Value.(string) == ":" {
			tk = readToken(tx) // stChecker.go

			if tk.Type == token.Operator && tk.Value.(string) == "]" {
				_, nextTk, errs = checkPtSqPr(nil, nil, layers, tx)
			} else {
				nextTk, errs = checkExpression2(tk, layers, tx)

				// nextTk Is CloseSquare

				_, nextTk, ers = checkPtSqPr(nil, nil, layers, tx)
				errs = append(errs, ers...)
			}
		} else {
			nextTk, errs = checkExpression2(tk, layers, tx)
			if !nextTk.IsCloseSquare() {
				// nextTk  is ":"
				tk = readToken(tx) // stChecker.go

				if tk.Type == token.Operator && tk.Value.(string) == "]" {
					_, nextTk, ers = checkPtSqPr(nil, nil, layers, tx)
					errs = append(errs, ers...)
				} else {
					nextTk, ers = checkExpression2(tk, layers, tx)
					errs = append(errs, ers...)

					// nextTk IsCloseSquare

					_, nextTk, ers = checkPtSqPr(nil, nil, layers, tx)
					errs = append(errs, ers...)
				}
			} else {
				_, nextTk, ers = checkPtSqPr(nil, nil, layers, tx)
				errs = append(errs, ers...)
			}
		}
	case tk.IsOpenPar():
		nextTk, errs = checkExpressionSeq(layers, tx) //exChecker.go

		// nextTk IsClosePar
		_, nextTk, ers = checkPtSqPr(nil, nil, layers, tx)
		errs = append(errs, ers...)
	default:
		s = sym
		nextTk = tk
		sym = nil // prevent checking sym.
	}

	if sym != nil {
		e := cksym.ErrIfNotFound(layers, sym)
		if e != nil {
			errs = append(errs, e)
		}
	}

	return
}

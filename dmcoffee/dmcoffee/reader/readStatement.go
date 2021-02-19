// Copyright 15-Jan-2021 ºDeme
// GNU General Public License - V3 <http://www.gnu.org/licenses/>

package reader

import (
	"github.com/dedeme/dmcoffee/operator"
	"github.com/dedeme/dmcoffee/reader/lineReader"
	"github.com/dedeme/dmcoffee/statement"
	"github.com/dedeme/dmcoffee/symbol"
	"github.com/dedeme/dmcoffee/token"
	"path"
	"strings"
)

func updateImports(rd *T, imps [][]*token.T) {
	for _, tks := range imps {
		var ok bool
		var key symbol.T
		var pth string
		if len(tks) == 1 {
			pth, ok = tks[0].S()
			if !ok {
				rd.Fail("Expected String. Actual '%v'", tks[0])
			}
		} else {
			key, ok = tks[0].Sy()
			if !ok {
				rd.Fail("Expected Symbol. Actual '%v'", tks[0])
			}
			pth, ok = tks[1].S()
			if !ok {
				rd.Fail("Expected String. Actual '%v'", tks[1])
			}
		}

		pth = path.Clean(path.Join(path.Dir(rd.source.String()), pth))
		if len(tks) == 1 {
			key = symbol.New(path.Base(pth))
		}

		_, ok = rd.imps[key]
		if ok {
			rd.Fail("Import '%v' is duplicate", key)
		}
		pathSym := symbol.New(pth)
		rd.imps[key] = pathSym
	}
}

func interpolation(rd *T, s string, pos *token.PosT) []*token.T {
	var r []*token.T
	ix := strings.Index(s, "${")
	for ix != -1 {
		r = append(r, token.NewS(s[:ix], pos))
		rest := s[ix+2:]
		ix2 := strings.Index(rest, "}")
		if ix2 == -1 {
			rd.Fail("Interpolation not closed")
		}
		rd2 := New(-1, strings.TrimSpace(rest[:ix2]))
		proc := rd2.Process()
		if len(proc) > 0 {
			r = append(r,
				token.NewO(operator.Plus, pos),
				token.NewSy(symbol.ToString, pos),
				token.NewO(operator.Oparenthesis, pos),
			)
			r = append(r, proc[0].Tokens...)
			r = append(r, token.NewO(operator.Cparenthesis, pos))
		}

		s = rest[ix2+1:]
		if s != "" {
			r = append(r, token.NewO(operator.Plus, pos))
			ix = strings.Index(s, "${")
		} else {
			ix = -1
		}
	}
	if s != "" {
		r = append(r, token.NewS(s, pos))
	}
	return r
}

func (rd *T) readStatement() (r *statement.T, ok bool) {
	var l *lineReader.LineT
	if l, ok = rd.lreader.Peek(); ok {
		rd.lreader.Pop()
		rd.line = l
		rd.ix = 0
		r = statement.New(l.Depth(), []*token.T{})

		for {
			if tk, ok2 := rd.readToken(); ok2 {
				if o, ok2 := tk.O(); ok2 {
					if o == operator.Point {
						l := len(r.Tokens)
						if l == 0 || r.Tokens[l-1].Type() != token.Symbol {
							rd.Fail("Point operator must be after module identifier")
						}
						sym, _ := r.Tokens[l-1].Sy()
						if ref, ok2 := rd.imps[sym]; ok2 {
							r.Tokens[l-1] = token.NewSy(ref, r.Tokens[l-1].Pos)
							r.Tokens = append(r.Tokens, tk)
						}
						continue
					}
				}
				if sym, ok2 := tk.Sy(); ok2 && sym == symbol.Import {
					r.Tokens = append(r.Tokens, tk)
					if tk2, ok2 := rd.readToken(); ok2 {
						if o, ok2 := tk2.O(); ok2 && o == operator.Colon {
							panic("No implemented")
						}
						if _, ok2 := tk2.S(); ok2 {
							r.Tokens = append(r.Tokens, tk2)
							updateImports(rd, [][]*token.T{{tk2}})
							continue
						}
						if _, ok2 := tk2.Sy(); ok2 {
							if tk3, ok3 := rd.readToken(); ok3 {
								r.Tokens = append(r.Tokens, tk2, tk3)
								updateImports(rd, [][]*token.T{{tk2, tk3}})
								continue
							}
							rd.Fail("Import parameter path is missing")
						}
						rd.Fail("Expected Symbol, String, or Block. Actual '%v'", tk2)
					}
					rd.Fail("Import parameters are missing")
				}
				if s, ok2 := tk.S(); ok2 && ok2 {
					tks := interpolation(rd, s, tk.Pos)
					r.Tokens = append(r.Tokens, tks...)
					continue
				}
				r.Tokens = append(r.Tokens, tk)
				continue
			}
			break
		}
	}
	return
}

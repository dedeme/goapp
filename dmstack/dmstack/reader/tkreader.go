// Copyright 30-Apr-2020 ºDeme
// GNU General Public License - V3 <http://www.gnu.org/licenses/>

package reader

import (
	"fmt"
	"github.com/dedeme/dmstack/symbol"
	"github.com/dedeme/dmstack/token"
	"strconv"
	"strings"
)

func isBlank(ch byte) bool {
	return ch <= ' ' || ch == ';' || ch == ':' || ch == ','
}

func (rd *T) nextToken() (rtk *token.T, ok bool) {
	nLine := rd.nLine
	prg := rd.prg
	prgIx := rd.prgIx
	prgLen := len(prg)

	if prgIx >= prgLen {
		return
	}

	if prgIx == 0 && rd.isFile && nLine == 1 &&
		prgLen > 1 && prg[0] == '#' && prg[1] == '!' {

		ix := strings.IndexByte(prg, '\n')
		if ix == -1 {
			return
		}
		nLine = 2
		rd.nLine = nLine
		prgIx = ix + 1
		rd.prgIx = prgIx
	}

	if len(rd.nextsTk) != 0 {
		rtk = rd.nextsTk[0]
		rd.nextsTk = rd.nextsTk[1:]
		ok = true
		return
	}

	for prgIx < prgLen {
		ch := prg[prgIx]
		if isBlank(ch) {
			if ch == '\n' {
				nLine++
			}
			prgIx++
			continue
		}

		if ch == '/' && prgIx < prgLen-1 {
			ch2 := prg[prgIx+1]
			if ch2 == '/' {
				ix := strings.IndexByte(prg[prgIx:], '\n')
				if ix == -1 {
					return
				}
				nLine++
				prgIx += ix + 1
				continue
			} else if ch2 == '*' {
				ix := strings.Index(prg[prgIx:], "*/")
				if ix == -1 {
					rd.nLine = nLine
					rd.fail("Unclosed commentary")
				}
				start := prgIx
				prgIx += ix + 2

				ix = strings.IndexByte(prg[start:prgIx], '\n')
				for ix != -1 {
					nLine++
					start += ix + 1
					ix = strings.IndexByte(prg[start:prgIx], '\n')
				}

				continue
			}
		}

		break
	}
	if prgIx >= prgLen {
		return
	}

	start := prgIx
	ch := prg[start]

	if ch == '"' {
		prgIx++
		rd.nLine = nLine
		rd.prgIx = prgIx
		rtk = rd.processString()
		ok = true
		return
	}

	if ch == '`' {
		ix := strings.IndexByte(prg[prgIx:], '\n')
		prgIx = prgIx + ix + 1
		if ix == -1 || prgIx >= prgLen {
			rd.fail("Unclosed multiline string")
		}
		rd.prgIx = prgIx
		rd.nLine++
		rtk = rd.processString2(prg[start+1 : prgIx-1])
		ok = true
		return
	}

	if strings.IndexByte("([{", ch) != -1 {
		cch := ")"
		if ch == '[' {
			cch = "]"
		} else if ch == '{' {
			cch = "}"
		}
		closeTk := token.NewSy(symbol.New(cch), nil)

		rd.prgIx = prgIx + 1
		start = nLine
		rd.nLine = nLine
		var ls []*token.T
		for {
			tk, ok2 := rd.nextToken()
			if !ok2 {
				rd.nLine = start
				rd.fail(fmt.Sprintf("Unclosed '%v'", string(ch)))
			}
			if tk.Eq(closeTk) {
				ok = true
				pos := token.NewPos(rd.source, start)
				if ch == '[' {
					rd.nextsTk = []*token.T{
						token.NewSy(symbol.Data, pos),
					}
				} else if ch == '{' {
					rd.nextsTk = []*token.T{
						token.NewSy(symbol.Data, pos),
						token.NewSy(symbol.Map, pos),
						token.NewSy(symbol.From, pos),
					}
				}
				rtk = token.NewP(ls, pos)
				return
			}

			if tk.Type() == token.Symbol {
				for _, t := range rd.processSymbol(ls, tk) { // in tksymbol.go
					ls = append(ls, t)
				}
			} else if tk.Type() == token.String {
				for _, t := range rd.processInterpolation(tk) { // in reader.go
					ls = append(ls, t)
				}
			} else {
				ls = append(ls, tk)
			}
		}
	}

	if strings.IndexByte(")]}", ch) != -1 {
		rd.nLine = nLine
		prgIx++
		rd.prgIx = prgIx
		ok = true
		rtk = token.NewSy(symbol.New(string(ch)), token.NewPos(rd.source, nLine))
		return
	}

	// Until here prg[PrgIx] = ch
	prgIx++
	// From here prg[PrgIx] = next ch
	for prgIx < prgLen {
		ch := prg[prgIx]
		if isBlank(ch) || strings.IndexByte("()[]{}", ch) != -1 {
			break
		}
		prgIx++
	}

	sub := prg[start:prgIx]
	pos := token.NewPos(rd.source, nLine)

	if sub == "true" {
		rtk = token.NewB(true, pos)
	} else if sub == "false" {
		rtk = token.NewB(false, pos)
	} else if sub == "." {
		rd.nLine = nLine
		rd.prgIx = prgIx
		return rd.nextToken()
	} else {
		s0 := sub[0]
		if s0 == '0' && prgIx > start+2 && sub[1] == 'x' {
			if prgIx == start+3 {
				rd.nLine = nLine
				rd.fail(fmt.Sprintf("Reader: Wrong number (%v)", sub))
			} else {
				n, err := strconv.ParseInt(sub[2:], 16, 64)
				if err != nil {
					rd.nLine = nLine
					rd.fail(fmt.Sprintf("Reader: Wrong number (%v)", sub))
				}
				rtk = token.NewI(n, pos)
			}
		} else if (s0 >= '0' && s0 <= '9') ||
			(s0 == '-' && prgIx > start+1 && (sub[1] >= '0' && sub[1] <= '9')) {
			if strings.IndexByte(sub, '.') == -1 {
				n, err := strconv.ParseInt(sub, 10, 64)
				if err != nil {
					rd.nLine = nLine
					rd.fail(fmt.Sprintf("Reader: Wrong number (%v)", sub))
				}
				rtk = token.NewI(n, pos)
			} else {
				n, err := strconv.ParseFloat(sub, 64)
				if err != nil {
					rd.nLine = nLine
					rd.fail(fmt.Sprintf("Reader: Wrong number (%v)", sub))
				}
				rtk = token.NewF(n, pos)
			}
		} else {
			rtk = token.NewSy(symbol.New(sub), pos)
		}
	}
	ok = true

	rd.nLine = nLine
	rd.prgIx = prgIx
	return
}

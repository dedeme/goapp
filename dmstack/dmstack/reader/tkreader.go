// Copyright 04-Jan-2021 ºDeme
// GNU General Public License - V3 <http://www.gnu.org/licenses/>

package reader

import (
	"fmt"
	"github.com/dedeme/dmstack/args"
	"github.com/dedeme/dmstack/operator"
	"github.com/dedeme/dmstack/symbol"
	"github.com/dedeme/dmstack/token"
	"strconv"
	"strings"
)

func isBlank(ch byte) bool {
	return ch <= ' ' || ch == ';'
}

func isLetter(ch byte) bool {
	return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') || ch == '_'
}

func isDigit(ch byte) bool {
	return ch >= '0' && ch <= '9'
}

func isLetterOrDigit(ch byte) bool {
	return isLetter(ch) || isDigit(ch)
}

func (rd *T) nextToken() (rtk *token.T, ok bool) {
	nLine := rd.nLine
	prg := rd.prg
	prgIx := rd.prgIx
	prgLen := len(prg)

	if prgIx >= prgLen {
		rd.prgIx = prgLen
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

	// skips blanks and docs
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
					rd.Fail("Unclosed commentary")
				}
				nLine += strings.Count(prg[prgIx:prgIx+ix], "\n")
				prgIx += ix + 2
				continue
			}
		}

		break
	}

	if prgIx >= prgLen {
		rd.nLine = nLine
		rd.prgIx = prgLen
		return
	}

	start := prgIx // token start
	ch := prg[start]

	if strings.IndexByte(",)]}", ch) != -1 {
		rd.nLine = nLine
		rd.prgIx = prgIx
		return
	}

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
			rd.Fail(" Unclosed multiline string")
		}
		rd.nLine = nLine + 1
		rd.prgIx = prgIx
		rtk = rd.processString2(prg[start+1 : prgIx-1]) // in tkstring.go
		ok = true
		return
	}

	// Containers
	if strings.IndexByte("([{", ch) != -1 {
		cch := byte(')')
		if ch == '[' {
			cch = ']'
		} else if ch == '{' {
			cch = '}'
		}

		rd.nLine = nLine
		rd.prgIx = prgIx + 1
		var ls [][]*token.T
		var els []*token.T
		withComma := false
		pos := token.NewPos(rd.source, nLine)
		for {
			nextTk, rdOk := rd.nextToken()
			if rdOk {
				if nextTk.Type() == token.Operator {
					op, _ := nextTk.O()
					if op == operator.Point {
						if len(els) < 1 {
							rd.Fail("Module name is missing")
						}
						tk2 := els[len(els)-1]
						if tk2.Type() != token.Symbol {
							rd.Fail("Expected a name before '.'")
						}
						k, _ := tk2.Sy()
						v, ok := rd.syms[k]
						if ok {
							els = append(els[:len(els)-1], token.NewSy(v, tk2.Pos), nextTk)
						} else {
							els = append(els, nextTk)
						}
					} else {
						els = append(els, rd.expandOperator(nextTk)...) // in tkexpand
					}
				} else if nextTk.Type() == token.String {
					for _, t := range rd.processInterpolation(nextTk) { //tkstring
						els = append(els, t)
					}
				} else if nextTk.Type() == token.Procedure {
					els = append(els, nextTk, token.NewO(operator.ProcHeap, nextTk.Pos))
				} else if nextTk.Type() == token.Array {
					els = append(els, rd.expandArray(nextTk)...) // in tkexpand
				} else if nextTk.Type() == token.Map {
					els = append(els, rd.expandMap(nextTk)...) // in tkexpand
				} else {
					els = append(els, nextTk)
				}
				continue
			}

			if rd.prgIx == len(rd.prg) {
				rd.Fail(fmt.Sprintf("Expected ',' or '%c' are missing", cch))
			}

			last := rd.prg[rd.prgIx]

			if last == ',' {
				if cch == ')' {
					rd.Fail("Unexpected ','")
				}
				if len(els) == 0 {
					rd.Fail("Empty element in array or object")
				}
				withComma = true
				ls = append(ls, els)
				els = []*token.T{}

				rd.prgIx = rd.prgIx + 1
				continue
			}

			if last == cch {
				if withComma && len(els) == 0 {
					rd.Fail(fmt.Sprintf("Unexpected ',%c'", cch))
				}
				if len(els) > 0 {
					ls = append(ls, els)
				}

				if cch == ')' || cch == ']' {
					var ls2 []*token.T
					for _, tks := range ls {
						ls2 = append(ls2, token.NewP(tks, tks[0].Pos))
					}

					if cch == ')' {
						if len(ls2) == 0 {
							rtk = token.NewP([]*token.T{}, pos)
						} else {
							rtk = ls2[0]
						}
					} else {
						rtk = token.NewA(ls2, pos)
					}

					rd.prgIx = rd.prgIx + 1
					ok = true
					return
				}

				if cch == '}' {
					mp := map[string]*token.T{}
					for _, tks := range ls {
						if len(tks) < 3 {
							rd.Fail("Incomplete object")
						}
						key, rok := tks[0].S()
						if !rok {
							sym, rok := tks[0].Sy()
							if !rok {
								rd.Fail(fmt.Sprintf(
									"Expect object key. Found %v", tks[0].StringDraft(),
								))
							}
							key = sym.String()
						}
						sep, rok := tks[1].O()
						if !rok || sep != operator.Assign {
							rd.Fail(fmt.Sprintf(
								"Expected ':'. Found '$v'", tks[1].StringDraft(),
							))
						}
						value := token.NewP(tks[2:], tks[2].Pos)

						mp[key] = value
					}

					rtk = token.NewM(mp, pos)

					rd.prgIx = rd.prgIx + 1
					ok = true
					return
				}

				rd.prgIx = rd.prgIx + 1
				ok = true
				return
			}

			rd.Fail(fmt.Sprintf(
				"Expected ',' or '%c', but '%c' was found", cch, last,
			))
		}
	}

	// Until here prg[PrgIx] = ch
	prgIx++
	// From here prg[PrgIx] = prgLen or the first byte after token.
	for prgIx < prgLen {
		ch := prg[prgIx]
		if isBlank(ch) || strings.IndexByte(",()[]{}", ch) != -1 {
			break
		}
		prgIx++
	}

	sub := prg[start:prgIx]
	pos := token.NewPos(rd.source, nLine)
	s0 := sub[0]

	if s0 == '0' && prgIx > start+2 && sub[1] == 'x' {
		if prgIx == start+3 {
			rd.nLine = nLine
			rd.Fail(fmt.Sprintf("Wrong number (%v)", sub))
		} else {
			for i := 3; i < len(sub); i++ {
				ch := sub[i]
				if !isLetterOrDigit(ch) {
					prgIx = start + i
					sub = sub[:i]
					break
				}
			}
			n, err := strconv.ParseInt(sub[2:], 16, 64)
			if err != nil {
				rd.nLine = nLine
				rd.Fail(fmt.Sprintf("Wrong number (%v)", sub))
			}
			rtk = token.NewI(n, pos)
		}
	} else if isDigit(s0) ||
		(s0 == '-' && prgIx > start+1 && isDigit(sub[1])) {
		if strings.IndexByte(sub, '.') == -1 {
			for i := 1; i < len(sub); i++ {
				if isLetterOrDigit(sub[i]) {
					continue
				}
				prgIx = start + i
				sub = sub[:i]
				break
			}
			n, err := strconv.ParseInt(sub, 10, 64)
			if err != nil {
				rd.nLine = nLine
				rd.Fail(fmt.Sprintf("Wrong number (%v)", sub))
			}
			rtk = token.NewI(n, pos)
		} else {
			for i := 1; i < len(sub); i++ {
				ch := sub[i]
				if isLetterOrDigit(ch) || ch == '.' || ch == '-' {
					continue
				}
				prgIx = start + i
				sub = sub[:i]
				break
			}
			n, err := strconv.ParseFloat(sub, 64)
			if err != nil {
				rd.nLine = nLine
				rd.Fail(fmt.Sprintf("Wrong number (%v)", sub))
			}
			rtk = token.NewF(n, pos)
		}
	} else if isLetter(s0) {
		for i := 1; i < len(sub); i++ {
			if isLetterOrDigit(sub[i]) {
				continue
			}
			prgIx = start + i
			sub = sub[:i]
			break
		}
		if sub == "true" {
			rtk = token.NewB(true, pos)
		} else if sub == "false" {
			rtk = token.NewB(false, pos)
		} else if sub == "this" {
			rtk = token.NewSy(rd.source, pos)
		} else {
			rtk = token.NewSy(symbol.New(sub), pos)
		}
	} else if s0 == '@' {
		prgIx = start + len(sub)
		if !args.Production || (len(sub) > 1 && sub[1] == '?') {
			rtk = token.NewO(operator.New(sub), pos)
		} else {
			rd.nLine = nLine
			rd.prgIx = prgIx
			rtk, ok = rd.nextToken()
			return
		}
	} else if s0 == '.' &&
		(len(sub) == 2 && isDigit(sub[1]) ||
			len(sub) > 1 && sub[1] == '.') {
		rtk = token.NewO(operator.New(sub), pos)
	} else {
		for i := 1; i < len(sub); i++ {
			if isLetterOrDigit(sub[i]) ||
				(sub[i] == '-' && i < len(sub)-1 && isDigit(sub[i+1])) {
				prgIx = start + i
				sub = sub[:i]
				break
			}
			continue
		}

		rtk = token.NewO(operator.New(sub), pos)
	}

	ok = true

	rd.nLine = nLine
	rd.prgIx = prgIx
	return
}

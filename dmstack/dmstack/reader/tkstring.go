// Copyright 04-Jan-2021 ºDeme
// GNU General Public License - V3 <http://www.gnu.org/licenses/>

package reader

import (
	"fmt"
	"github.com/dedeme/dmstack/operator"
	"github.com/dedeme/dmstack/symbol"
	"github.com/dedeme/dmstack/token"
	"strconv"
	"strings"
)

func isHex(b byte) bool {
	return (b >= '0' && b <= '9') ||
		(b >= 'a' && b <= 'f') ||
		(b >= 'A' && b <= 'F')
}

func (rd *T) processString() *token.T {
	prg := rd.prg
	prgIx := rd.prgIx
	prgLen := len(prg)

	if prgIx >= prgLen {
		rd.Fail("Unclosed quotes")
	}

	rs := ""
	isEsc := false
	hex := 0
	closed := false
	for _, rn := range prg[prgIx:] {
		s := string(rn)
		l := len(s)
		prgIx += l
		if isEsc {
			isEsc = false
			if l == 1 {
				ch := s[0]
				switch ch {
				case '"':
					rs += "\""
				case '\\':
					rs += "\\"
				case 't':
					rs += "\t"
				case 'r':
					rs += "\r"
				case 'n':
					rs += "\n"
				case 'b':
					rs += "\b"
				case 'f':
					rs += "\f"
				case 'u':
					hex = 1
				}
			} else {
				rd.Fail(fmt.Sprintf("Bad escape sequence '\\%v'", rn))
			}
		} else if hex > 0 {
			ch := s[0]
			if l == 1 && isHex(ch) {
				if hex == 4 {
					s, _ := strconv.Unquote("'\\u" + prg[prgIx-4:prgIx] + "'")
					rs += s
					hex = 0
				} else {
					hex++
				}
			} else {
				rd.Fail(fmt.Sprintf(
					"Bad unicode sequence '\\u%v'", prg[prgIx-hex:prgIx],
				))
			}
		} else {
			if l == 1 {
				ch := s[0]
				if ch < ' ' {
					rd.Fail(fmt.Sprintf("Control code '%v' don't allowed", ch))
				} else if ch == '"' {
					closed = true
					break
				} else if ch == '\\' {
					isEsc = true
					continue
				}
			}
			rs += s
		}
	}

	if isEsc {
		rd.Fail("'\\' at the end of string")
	}
	if !closed {
		rd.Fail("Unclosed quotes")
	}

	rd.prgIx = prgIx
	return token.NewS(rs, token.NewPos(rd.source, rd.nLine))
}

// Expected rd postion after '\n' and with characters
//    key : "here doc" text value. (e.g. ` -> "", `abc -> abc)
func (rd *T) processString2(key string) *token.T {
	lineStart := rd.nLine - 1

	prg := rd.prg
	prgIx := rd.prgIx
	ix := strings.Index(prg[prgIx:], key+"`")
	if ix == -1 {
		rd.nLine = lineStart
		rd.Fail("Unclosed multiline string")
	}
	start := prgIx
	end := prgIx + ix

	for _, rn := range prg[start:end] {
		if string(rn) == "\n" {
			rd.nLine++
		}
	}

	rd.prgIx = end + len(key) + 1
	return token.NewS(prg[start:end], token.NewPos(rd.source, rd.nLine))
}

// Process a String interpolation
func (rd *T) processInterpolation(tk *token.T) []*token.T {
	fpos := func(nline int) *token.PosT {
		return token.NewPos(tk.Pos.Source, nline)
	}

	s, _ := tk.S()
	nline := tk.Pos.Nline

	pos := 0
	var tks []*token.T
	ix := strings.Index(s, "${")
	for ix != -1 {
		tks = append(tks, token.NewS(s[pos:pos+ix], fpos(nline)))
		if pos > 0 {
			tks = append(tks, token.NewO(operator.Plus, fpos(nline)))
		}
		nline += strings.Count(s[pos:pos+ix], "\n")
		pos += ix + 2

		subr := newFromReader(rd, s[pos:], nline)
		prg, _ := subr.Process().P()

		if subr.LastChar() != "}" {
			subr.Fail("Interpolation not closed")
		}
		ix2 := subr.prgIx
		nline += strings.Count(s[pos:pos+ix2], "\n")

		lprg := len(prg)
		if lprg == 0 {
			tks = append(tks, token.NewS("", fpos(nline)))
			tks = append(tks, token.NewO(operator.Plus, fpos(nline)))
		} else {
			for _, tk2 := range prg {
				tks = append(tks, tk2)
			}
			tks = append(tks, token.NewSy(symbol.ToString, fpos(nline)))
			tks = append(tks, token.NewO(operator.Plus, fpos(nline)))
		}

		pos += ix2 + 1
		ix = strings.Index(s[pos:], "${")
	}

	tks = append(tks, token.NewS(s[pos:], fpos(nline)))
	if len(tks) > 1 {
		nline += strings.Count(s[pos:], "\n")
		tks = append(tks, token.NewO(operator.Plus, fpos(nline)))
	}

	return tks
}

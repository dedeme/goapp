// Copyright 15-Jan-2021 ºDeme
// GNU General Public License - V3 <http://www.gnu.org/licenses/>

package reader

import (
	"github.com/dedeme/dmcoffee/fail"
	"github.com/dedeme/dmcoffee/operator"
	"github.com/dedeme/dmcoffee/symbol"
	"github.com/dedeme/dmcoffee/token"
	"strconv"
	"strings"
)

func isLetter(ch byte) bool {
	return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') || ch == '_'
}

func isDigit(ch byte) bool {
	return ch >= '0' && ch <= '9'
}

func isLetterOrDigit(ch byte) bool {
	return isLetter(ch) || isDigit(ch)
}

func skipBlanks(rd *T) {
	l := rd.line.Tx()
	for rd.ix < len(l) {
		if l[rd.ix] <= ' ' {
			rd.ix++
			continue
		}
		if l[rd.ix] == '/' && len(l)-rd.ix > 1 {
			if l[rd.ix+1] == '/' {
				rd.ix = len(l)
				continue
			} else if l[rd.ix+1] == '*' {
				startNline := rd.line.N()
				cut := rd.ix + 2
				for {
					l2 := l[cut:]
					ix := strings.Index(l2, "*/")
					if ix == -1 {
						if ln, ok := rd.lreader.Pop(); ok {
							rd.line = ln
							l = ln.Tx()
							cut = 0
							continue
						}
						panic(fail.New(
							token.NewPos(rd.source, startNline),
							fail.ESyntax(),
							"Commentary not closed",
						))
					}
					rd.ix = cut + ix + 2
					break
				}
				continue
			}
		}
		break
	}
}

func (rd *T) readToken() (tk *token.T, ok bool) {
	if rd.ix == len(rd.line.Tx()) {
		return
	}

	skipBlanks(rd)
	l := rd.line.Tx()
	start := rd.ix
	if start == len(rd.line.Tx()) {
		return
	}
	pos := rd.getPos()

	//start := ix
	ch := l[start]
	ch2 := byte(0)
	if start+1 < len(l) {
		ch2 = l[start+1]
	}

	if ch == '"' {
		rd.ix++
		ok = true
		tk = rd.readString()
		return
	} else if ch == '`' {
		rd.ix++
		ok = true
		tk = rd.readHereDoc()
		return
	} else if ch == '0' && ch2 == 'x' {
		i := rd.ix + 2
		for ; i < len(l); i++ {
			if isLetterOrDigit(l[i]) {
				continue
			}
		}
		rd.ix = i
		n, err := strconv.ParseInt(l[start+2:i], 16, 64)
		if err != nil {
			rd.Fail("Wrong number (%v)", l[start:i])
		}
		ok = true
		tk = token.NewI(n, pos)
		return
	} else if isDigit(ch) || (ch == '-' && isDigit(ch2)) {
		i := rd.ix + 1
		exponent := false
		for ; i < len(l); i++ {
			ch := l[i]
			if ch == 'e' || ch == 'E' && i < len(l)-1 && !exponent {
				ch2 := l[i+1]
				if isDigit(ch2) || ch2 == '-' || ch2 == '+' {
					exponent = true
					i++
					continue
				}
			}
			if isDigit(ch) {
				continue
			}
			break
		}

		if i < len(l) && l[i] == '.' && !exponent {
			i++
			for ; i < len(l); i++ {
				ch := l[i]
				if ch == 'e' || ch == 'E' && i < len(l)-1 && !exponent {
					ch2 := l[i+1]
					if isDigit(ch2) || ch2 == '-' || ch2 == '+' {
						exponent = true
						i++
						continue
					}
				}
				if isDigit(ch) {
					continue
				}
				break
			}
			rd.ix = i
			n, err := strconv.ParseFloat(l[start:i], 64)
			if err != nil {
				rd.Fail("Wrong number (%v)", l[start:i])
			}
			ok = true
			tk = token.NewF(n, pos)
			return
		}

		rd.ix = i

		if exponent {
			n, err := strconv.ParseFloat(l[start:i], 64)
			if err != nil {
				rd.Fail("Wrong number (%v)", l[start:i])
			}
			ok = true
			tk = token.NewF(n, pos)
			return
		}

		n, err := strconv.ParseInt(l[start:i], 10, 64)
		if err != nil {
			rd.Fail("Wrong number (%v)", l[start:i])
		}
		ok = true
		tk = token.NewI(n, pos)
		return
	} else if isLetter(ch) {
		i := rd.ix + 1
		for ; i < len(l); i++ {
			ch := l[i]
			if isLetterOrDigit(ch) {
				continue
			}
			break
		}

		rd.ix = i
		ok = true
		sub := l[start:i]
		if sub == "true" {
			tk = token.NewB(true, pos)
		} else if sub == "false" {
			tk = token.NewB(false, pos)
		} else {
			tk = token.NewSy(symbol.New(sub), pos)
		}
		return
	} else if o, ok2 := operator.GetO2(ch, ch2); ok2 {
		ok = true
		if o == operator.MinusMinus && len(l) > rd.ix+2 && isDigit(l[rd.ix+2]) {
			rd.ix++
			tk = token.NewO(operator.Minus, pos)
			return
		}
		rd.ix += 2
		tk = token.NewO(o, pos)
		return
	} else if o, ok2 := operator.GetO1(ch); ok2 {
		ok = true
		rd.ix++
		tk = token.NewO(o, pos)
		return
	}

	rd.Fail("Unknown symbol '%c'", ch)
	// Unreachable
	return
}

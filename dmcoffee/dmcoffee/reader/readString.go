// Copyright 15-Jan-2021 ºDeme
// GNU General Public License - V3 <http://www.gnu.org/licenses/>

package reader

import (
	"github.com/dedeme/dmcoffee/token"
	"strconv"
	"strings"
)

func isHex(b byte) bool {
	return (b >= '0' && b <= '9') ||
		(b >= 'a' && b <= 'f') ||
		(b >= 'A' && b <= 'F')
}

func (rd *T) readString() *token.T {
	l := rd.line.Tx()
	if rd.ix == len(l) {
		rd.Fail("Unclosed quotes")
	}

	var bf strings.Builder
	isEsc := false
	hex := 0
	closed := false
	i := rd.ix
	for ; i < len(l); i++ {
		ch := l[i]
		if isEsc {
			isEsc = false
			switch ch {
			case '"', '\\', 't', 'r', 'n', 'b', 'f':
				bf.WriteByte(ch)
			case 'u':
				hex = 1
			default:
				rd.Fail("Bad escape sequence '\\%c'", ch)
			}
		} else if hex > 0 {
			if isHex(ch) {
				if hex == 4 {
					s, _ := strconv.Unquote("'\\u" + l[i-3:i+1] + "'")
					bf.WriteString(s)
					hex = 0
				} else {
					hex++
				}
			} else {
				rd.Fail("Bad unicode sequence '\\%v'", l[i-hex:i])
			}
		} else {
			if ch == '"' {
				closed = true
				break
			} else if ch == '\\' {
				isEsc = true
				continue
			}
			bf.WriteByte(ch)
		}
	}

	if isEsc {
		rd.Fail("'\\' at the end of string")
	}
	if hex > 0 {
		rd.Fail("Bad unicode escape (%v)", l[i-hex:])
	}
	if !closed {
		rd.Fail("Unclosed quotes")
	}

	rd.ix = i + 1
	return token.NewS(bf.String(), rd.getPos())
}

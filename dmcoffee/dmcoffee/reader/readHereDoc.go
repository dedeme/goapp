// Copyright 15-Jan-2021 ºDeme
// GNU General Public License - V3 <http://www.gnu.org/licenses/>

package reader

import (
	"fmt"
	"github.com/dedeme/dmcoffee/fail"
	"github.com/dedeme/dmcoffee/token"
	"strings"
)

// Label is leters or digits.
func (rd *T) readHereDoc() *token.T {
	l := rd.line.Tx()
	i := rd.ix
	start := i
	for ; i < len(l); i++ {
		ch := l[i]
		if isLetterOrDigit(ch) {
			continue
		}
		break
	}

	if i != len(l) {
		rd.Fail("HereDoc label must be at end of line (`%v)", l[start:i])
	}

	label := l[start:] + "`"
	pos := rd.getPos()
	var bf strings.Builder
	first := true
	for {
		if l, ok := rd.lreader.Pop(); ok {
			tx := l.Tx()
			if tx == label {
				rd.line = l
				rd.ix = len(tx)
				break
			}
			if first {
				first = false
			} else {
				bf.WriteByte('\n')
			}
			bf.WriteString(tx)
			continue
		}
		panic(fail.New(
			pos, fail.ESyntax(),
			fmt.Sprintf("HereDoc label whitout closing (%v)", label),
		))
	}

	return token.NewS(bf.String(), pos)
}

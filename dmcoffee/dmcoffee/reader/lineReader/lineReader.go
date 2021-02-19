// Copyright 14-Jan-2021 ºDeme
// GNU General Public License - V3 <http://www.gnu.org/licenses/>

// Line reader.

package lineReader

import (
	"fmt"
	"github.com/dedeme/dmcoffee/cts"
	"github.com/dedeme/dmcoffee/fail"
	"github.com/dedeme/dmcoffee/symbol"
	"github.com/dedeme/dmcoffee/token"
	"strings"
)

// Line structure
type LineT struct {
	n     int
	depth int
	tx    string
}

func (l *LineT) N() int {
	return l.n
}
func (l *LineT) Depth() int {
	return l.depth
}
func (l *LineT) Tx() string {
	return l.tx
}

// Reader
type T struct {
	ix    int
	lines []*LineT
}

// Constructor
// If tabulation is wrong returns ok = false
func New(source symbol.T, tx string) *T {
	var lines []*LineT
	for ln, l := range strings.Split(tx, "\n") {
		if ln == 0 && strings.HasPrefix(l, "#!") {
			continue
		}

		i := 0
		for ; i < len(l); i++ {
			if l[i] != ' ' {
				break
			}
		}
		if i%cts.TabSpaces != 0 {
			spaces := "spaces"
			if i == 1 {
				spaces = "space"
			}
			panic(fail.New(
				token.NewPos(source, ln+1),
				fail.ESyntax(),
				fmt.Sprintf("Odd tabulation (%v %v)", i, spaces),
			))
		}
		l = strings.TrimSpace(l)
		if len(l) > 0 {
			lines = append(lines, &LineT{ln + 1, i / cts.TabSpaces, l})
		}
	}

	return &T{0, lines}
}

// Returns the current line or "ok=false" if there are no more lines.
func (rd *T) Peek() (line *LineT, ok bool) {
	if rd.ix < len(rd.lines) {
		ok = true
		line = rd.lines[rd.ix]
	}
	return
}

// Returns the current line or "ok=false" if there are no more lines, advancing
// one position.
func (rd *T) Pop() (line *LineT, ok bool) {
	if rd.ix < len(rd.lines) {
		ok = true
		line = rd.lines[rd.ix]
		rd.ix++
	}
	return
}

// Returns the number of lines read.
func (rd *T) Nlines() int {
	return len(rd.lines)
}

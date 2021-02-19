// Copyright 20-Jan-2021 ºDeme
// GNU General Public License - V3 <http://www.gnu.org/licenses/>

// Statement data.
package statement

import (
	"github.com/dedeme/dmcoffee/token"
)

type T struct {
	Depth  int
	Tokens []*token.T
}

func New(depth int, tokens []*token.T) *T {
	return &T{depth, tokens}
}

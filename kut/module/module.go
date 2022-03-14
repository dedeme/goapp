// Copyright 20-Feb-2022 ºDeme
// GNU General Public License - V3 <http://www.gnu.org/licenses/>

// Module data.
package module

import (
	"github.com/dedeme/kut/heap"
	"github.com/dedeme/kut/heap0"
	"github.com/dedeme/kut/statement"
)

type T struct {
	Imports map[string]int
	Heap0   heap0.T
	Heap    heap.T
	Code    []*statement.T
}

func New(imports map[string]int, hp heap0.T, code []*statement.T) *T {
	return &T{imports, hp, heap.New(), code}
}

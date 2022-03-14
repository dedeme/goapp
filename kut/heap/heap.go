// Copyright 02-Feb-2022 ºDeme
// GNU General Public License - V3 <http://www.gnu.org/licenses/>

// Heap data.
package heap

import (
	"github.com/dedeme/kut/expression"
)

type T map[string]*expression.T

func New() T {
	return map[string]*expression.T{}
}

// Adds symbol to heap.
// If symbol already has been added, returns false.
func (h T) Add(symbol string, expr *expression.T) bool {
	if _, ok := h[symbol]; ok {
		return false
	}
	h[symbol] = expr
	return true
}

// Returns the expresion matching 'symbol'.
func Get(heaps []T, symbol string) (expr *expression.T, ok bool) {
	for _, h := range heaps {
		if expr, ok = h[symbol]; ok {
			return
		}
	}
	return
}

// Returns a new []heap.T with 'h' added on top of 'heaps'
func (h T) AddTo(heaps []T) []T {
	return append([]T{h}, heaps...)
}

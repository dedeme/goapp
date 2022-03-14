// Copyright 02-Feb-2022 ºDeme
// GNU General Public License - V3 <http://www.gnu.org/licenses/>

// File heap data.
package heap0

import (
	"github.com/dedeme/kut/expression"
)

type EntryT struct {
	Nline int
	Expr  *expression.T
}

type T map[string]*EntryT

func New() T {
	return map[string]*EntryT{}
}

// Adds symbol to heap0.
// If symbol already has been added, returns false.
func (h T) Add(symbol string, nline int, expr *expression.T) bool {
	if _, ok := h[symbol]; ok {
		return false
	}
	h[symbol] = &EntryT{nline, expr}
	return true
}

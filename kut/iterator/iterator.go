// Copyright 25-Feb-2022 ºDeme
// GNU General Public License - V3 <http://www.gnu.org/licenses/>

// Function data
package iterator

import (
	"github.com/dedeme/kut/expression"
)

type T struct {
	HasNext func() bool
	// A function with 0 arguments with returns a expression.T
	Next func() *expression.T
}

// Constructor
//    hasNext: A function with 0 arguments with returns a bool
//    next: A function with 0 arguments with returns a expression.T
func New(HasNext func() bool, Next func() *expression.T) *T {
	return &T{HasNext, Next}
}

func (s *T) String() string {
	return "<iterator>"
}

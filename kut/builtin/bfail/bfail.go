// Copyright 21-Feb-2022 ºDeme
// GNU General Public License - V3 <http://www.gnu.org/licenses/>

// Built-in function fail.
package bfail

import (
	"errors"
	"fmt"
	"github.com/dedeme/kut/expression"
	"strings"
)

func Mk(msg string) error {
	return errors.New(msg)
}

func Type(expr *expression.T, expected ...string) error {
	msg := "Type error:"
	if len(expected) == 1 {
		msg += fmt.Sprintf(
			"\n    Expected: %v\n    Found   : %T (%v)", expected[0], expr.Value, expr)
	} else {
		msg += fmt.Sprintf(
			"\n    Expected: %v\n    Found   : %T (%v)",
			strings.Join(expected, ", "), expr.Value, expr)
	}
	return Mk(expression.ReplaceGoName(msg))
}

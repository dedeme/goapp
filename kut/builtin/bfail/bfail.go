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
	return Mk(strings.ReplaceAll(
		strings.ReplaceAll(
			strings.ReplaceAll(
				strings.ReplaceAll(
					strings.ReplaceAll(
						strings.ReplaceAll(
							strings.ReplaceAll(
								strings.ReplaceAll(msg, "map[string]*expression.T", "map"),
								"[]*expression.T", "array"),
							"*expression.emptyT", "()"),
						"*runner.BModuleT", "module"),
					"*bfunction.T", "buit-in function"),
				"*function.T", "function"),
			"float64", "float"),
		"int64", "int"))
}

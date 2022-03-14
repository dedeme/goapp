// Copyright 21-Feb-2022 ºDeme
// GNU General Public License - V3 <http://www.gnu.org/licenses/>

// Runtime fail.
package fail

import (
	"errors"
	"fmt"
	"github.com/dedeme/kut/expression"
	"github.com/dedeme/kut/fileix"
	"github.com/dedeme/kut/function"
	"github.com/dedeme/kut/statement"
	"strconv"
	"strings"
)

type SysErrorT struct {
	Msg string
	Fn  *function.T
}

func (e *SysErrorT) Error() string {
	return e.Msg
}

func MkSysError(msg string, fn *function.T) *SysErrorT {
	return &SysErrorT{msg, fn}
}

func Mk(msg string, stackTrace []*statement.T) error {
	ix := len(stackTrace) - 1
	for i := 0; i < 15; i++ {
		if ix < 0 {
			break
		}
		st := stackTrace[ix]
		ix--
		msg += "\n  " + fileix.Get(st.File) + ":" + strconv.Itoa(st.Nline)
	}
	return errors.New(msg)
}

func Type(expr *expression.T, stackTrace []*statement.T, expected ...string) error {
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
		"int64", "int"), stackTrace)
}

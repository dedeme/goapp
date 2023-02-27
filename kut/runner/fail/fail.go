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
		msg += "\n  " + fileix.GetFail(st.File) + ":" + strconv.Itoa(st.Nline) + ":"
	}
	return errors.New(msg)
}

func Type(expr *expression.T, stackTrace []*statement.T, expected ...string) error {
	msg := "Type error:" + fmt.Sprintf(
		"\n    Expected: %v\n    Found   : %T (%v)",
		strings.Join(expected, ", "), expr.Value, expr)
	return Mk(expression.ReplaceGoName(msg), stackTrace)
}

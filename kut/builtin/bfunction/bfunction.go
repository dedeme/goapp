// Copyright 23-Feb-2022 ºDeme
// GNU General Public License - V3 <http://www.gnu.org/licenses/>

// Built-in function data.
package bfunction

import (
	"fmt"
	"github.com/dedeme/kut/expression"
	"github.com/dedeme/kut/heap"
	"github.com/dedeme/kut/heap0"
	"github.com/dedeme/kut/runner/fail"
	"github.com/dedeme/kut/statement"
)

type T struct {
	// Variables number
	nvars int
	// Function (varsArrayOfFinalExpressions, stackTrace)(finalExpression, error)
	//    'finalExpression' is nil for procedures.
	//    'error' must be construct with 'errors.New(string)'
	fn func([]*expression.T) (*expression.T, error)
}

// Constructor
//    nvars: Variables number of 'fn'
//    fn: Function (varsArrayOfFinalExpressions, stackTrace)(finalExpression, error)
func New(nvar int, fn func([]*expression.T) (*expression.T, error)) *T {
	return &T{nvar, fn}
}

// Executes a bfunction.
//    name  : Function name (for debug)
//    vars  : Function variables. Must be final expresions.
//    stackT: Stack for debug.
// Returns an error or an expression.
//    if 'bf' returns <nil> 'Run' returns a empty expression.
func (bf *T) Run(name string, vars []*expression.T, stackT []*statement.T) (
	ex *expression.T, err error,
) {
	if len(vars) != bf.nvars {
		err = fail.Mk(fmt.Sprintf(
			"Function %v:\n    expect %v argument(s), found %v",
			name, bf.nvars, len(vars),
		), stackT)
		return
	}
	ex, err = bf.fn(vars)
	if err != nil {
		switch e := err.(type) {
		case *fail.SysErrorT:
			e.Msg = fail.Mk(err.Error(), stackT).Error()
		default:
			err = fail.Mk("Function "+name+":\n    "+err.Error(), stackT)
		}
	}
	if ex == nil {
		ex = expression.MkEmpty()
	}
	return
}

// Implements 'function.I' interface.
// Implements 'I' interface.
func (f *T) MkExPr(pars []*expression.T) (
	ex *expression.T, imports map[string]int, hp0 heap0.T, hps []heap.T,
) {
	ex = expression.New(expression.ExPr, []interface{}{
		expression.MkFinal(f), pars})
	imports = map[string]int{}
	hp0 = heap0.New()
	return
}

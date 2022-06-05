// Copyright 15-Feb-2022 ºDeme
// GNU General Public License - V3 <http://www.gnu.org/licenses/>

// Function data
package function

import (
	"github.com/dedeme/kut/expression"
	"github.com/dedeme/kut/heap"
	"github.com/dedeme/kut/heap0"
	"github.com/dedeme/kut/statement"
	"strings"
)

// Converter to 'expression.ExPr'. Used in 'runner.solveIsolateFunction'
// (exSolver.go)
type I interface {
	MkExPr(pars []*expression.T) (
		ex *expression.T, imports map[string]int, hp0 heap0.T, hps []heap.T,
	)
}

type T struct {
	Imports map[string]int
	Hp0     heap0.T
	Hps     []heap.T
	Vars    []string
	Stat    *statement.T
}

func New(
	vars []string, stat *statement.T,
) *T {
	return &T{map[string]int{}, heap0.New(), []heap.T{}, vars, stat}
}

func NewWithHeaps(
	hps []heap.T, vars []string, stat *statement.T,
) *T {
	return &T{map[string]int{}, heap0.New(), hps, vars, stat}
}

// Implements 'I' interface.
func (f *T) MkExPr(pars []*expression.T) (
	ex *expression.T, imports map[string]int, hp0 heap0.T, hps []heap.T,
) {
	ex = expression.New(expression.ExPr, []interface{}{
		expression.MkFinal(f), pars})
	imports = f.Imports
	hp0 = f.Hp0
	hps = f.Hps
	return
}

func (f *T) SetContext(imports map[string]int, hp0 heap0.T, hps []heap.T) *T {
	return &T{
		imports,
		hp0,
		hps,
		f.Vars,
		f.Stat,
	}
}

func (f *T) String() string {
	return "(\\" + strings.Join(f.Vars, ", ") + " -> " + f.Stat.String() + ")"
}

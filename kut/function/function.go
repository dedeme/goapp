// Copyright 15-Feb-2022 ºDeme
// GNU General Public License - V3 <http://www.gnu.org/licenses/>

// Function data
package function

import (
	"github.com/dedeme/kut/heap"
	"github.com/dedeme/kut/heap0"
	"github.com/dedeme/kut/statement"
	"strings"
)

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

func (f *T) SetContext(imports map[string]int, hp0 heap0.T, hps []heap.T) {
	f.Imports = imports
	f.Hp0 = hp0
	f.Hps = hps
}

func (f *T) String() string {
	return "(\\" + strings.Join(f.Vars, ", ") + " -> " + f.Stat.String() + ")"
}

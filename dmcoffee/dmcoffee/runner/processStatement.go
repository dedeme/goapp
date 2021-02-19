// Copyright 08-Feb-2021 ºDeme
// GNU General Public License - V3 <http://www.gnu.org/licenses/>

package runner

import (
	"fmt"
	"github.com/dedeme/dmcoffee/machine"
	"github.com/dedeme/dmcoffee/symbol"
	"github.com/dedeme/dmcoffee/operator"
)

func puts(m *machine.T) {
	if tk1, ok := m.NextToken(); ok {
		if op1, ok := tk1.O(); ok && op1 == operator.Oparenthesis {
    }
    m.Faile(operator.Oparenthesis, tk1)
	}
  m.Faile(operator.Oparenthesis, "Nothing")
}

func processStatement(m *machine.T) {
  fmt.Println("here")
	if tk, ok := m.NextToken(); ok {
    fmt.Println(tk)
		if sym, ok := tk.Sy(); ok {
			switch sym {
			case symbol.Puts:
				fmt.Println("abc")
			default:
				m.Fail(machine.EUnexpected(), "(%v) Symbol unexpected", sym)
			}
		}
    m.Faile("Symbol", tk)
	}
}

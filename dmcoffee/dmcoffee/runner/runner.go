// Copyright 22-Jan-2021 ºDeme
// GNU General Public License - V3 <http://www.gnu.org/licenses/>

// Code runner.
package runner

import (
	"fmt"
	"github.com/dedeme/dmcoffee/machine"
	"os"
	"strings"
)

func Run(m *machine.T) {
	defer func() {
		if err := recover(); err != nil {
			e, ok := err.(*machine.Error)
			if ok {
				fmt.Printf("%v: %v\n", e.Type, e.Message)
				fmt.Println(strings.Join(e.Machine.StackTrace(), "\n"))
				os.Exit(0)
			} else {
				panic(err)
			}
		}
	}()

	for {
		if st, ok := m.NextStatement(); ok {
			if st.Depth != m.Depth() {
				m.Fail(
					machine.ERuntime(), "Unexpected spaces (Expected '%v'. Actual '%v'.)",
					m.Depth()*2, st.Depth*2,
				)
			}
			processStatement(m) // In processStatement.go
			continue
		}
		break
	}

}

// Copyright 16-May-2020 ºDeme
// GNU General Public License - V3 <http://www.gnu.org/licenses/>

// Primitives main hub.
package primitives

import (
	"github.com/dedeme/dmstack/machine"
	"github.com/dedeme/dmstack/symbol"
)

// Processes 'sym'.
//
// It call functions in 'global0.go'
//    m: Virtual machine.
//    sym: Symbol to process.
//    run: Functions to run a virtual machine.
func Global(m *machine.T, sym symbol.T, run func(m *machine.T)) bool {
	switch sym {
	case symbol.Nop:
	case symbol.Puts:
		prPuts(m)
	case symbol.Assert:
		prAssert(m)
	case symbol.Expect:
		prExpect(m)
	case symbol.Run:
		prRun(m, run)
	case symbol.Data:
		prData(m, run)
	case symbol.Import:
		prImport(m, run)
	default:
		return false
	}
	return true
}

func Modules(m *machine.T, sym symbol.T, run func(m *machine.T)) bool {
	switch sym {
	default:
		return false
	}
	return true
}

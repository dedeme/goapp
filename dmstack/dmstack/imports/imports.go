// Copyright 04-Jan-2021 ºDeme
// GNU General Public License - V3 <http://www.gnu.org/licenses/>

// Imports management.
package imports

import (
	"fmt"
	"github.com/dedeme/dmstack/symbol"
	"github.com/dedeme/dmstack/token"
)

var imported = map[symbol.T]map[symbol.T]*token.T{}

func Add(id symbol.T) {
	imported[id] = nil
}

func Initialize(id symbol.T) {
	imported[id] = map[symbol.T]*token.T{}
}

// Returns the heap of source 'id'.
// If source has not been loaded yeat, returns heap == null, ok == true.
// If source has not a heap with 'id', returns ok == false.
func Get(id symbol.T) (heap map[symbol.T]*token.T, ok bool) {
	heap, ok = imported[id]
	return
}

func AddKey(id, key symbol.T, value *token.T) {
	if heap, ok := Get(id); ok {
		if heap != nil {
			heap[key] = value
			return
		}
		panic(fmt.Sprintf("Module '%v' not initialized", id))
	}
	panic(fmt.Sprintf("Module '%v' not found", id))
}

// Reads an import symbol.
//  'tk' can be:
//    - A String (e.g. "lib/person" import):
//        Returns symbol.Kv{-1, string symbol}
//    - A Procedure with two elements (e.g. ("lib/person" per) import):
//        Returns symbol.Kv{symbol, string symbol}
//    - Otherwise returns an error.
func ReadSymbol(tk *token.T) (symMap *symbol.Kv, err error) {
	if s, ok := tk.S(); ok {
		symMap = &symbol.Kv{-1, symbol.New(s)}
		return
	}

	if a, ok := tk.P(); ok {
		if len(a) == 2 {
			if s, ok := a[1].S(); ok {
				if sy, ok := a[0].Sy(); ok {
					symMap = &symbol.Kv{sy, symbol.New(s)}
					return
				}
				err = fmt.Errorf(
					"\n  Expected: First element is Symbol."+
						"\n  Actual  : First element is '%v'",
					a[0],
				)
				return
			}
			err = fmt.Errorf(
				"\n  Expected: Second element is String."+
					"\n  Actual  : Second element is '%v'",
				a[1],
			)
			return
		}
		err = fmt.Errorf(
			"\n  Expected: Procedure with 2 elements.\n  Actual  : '%v'", len(a),
		)
		return
	}

	err = fmt.Errorf(
		"\n  Expected:  Procedure or String.\n  Actual  : '%v'",
		tk,
	)
	return
}

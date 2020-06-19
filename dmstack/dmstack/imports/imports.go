// Copyright 27-Apr-2020 ºDeme
// GNU General Public License - V3 <http://www.gnu.org/licenses/>

// Imports management.
package imports

import (
	"fmt"
	"github.com/dedeme/dmstack/symbol"
	"github.com/dedeme/dmstack/token"
)

var onWay []symbol.T
var imported = map[symbol.T]map[symbol.T]*token.T{}

// Anotates an 'on way' import.
func PutOnWay(imp symbol.T) {
	onWay = append(onWay, imp)
}

// Removes an 'on way' import.
func QuitOnWay(imp symbol.T) {
	var tmp []symbol.T
	for _, e := range onWay {
		if e != imp {
			tmp = append(tmp, e)
		}
	}
	onWay = tmp
}

// Returns 'true' if the importation of'imp' has started but not finished.
func IsOnWay(imp symbol.T) bool {
	for _, e := range onWay {
		if e == imp {
			return true
		}
	}
	return false
}

func Add(id symbol.T, heap map[symbol.T]*token.T) {
  imported[id] = heap
}

func Get(id symbol.T) (heap map[symbol.T]*token.T, ok bool) {
  heap, ok = imported[id]
  return
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

	if l, ok := tk.P(); ok {
		if len(l) == 2 {
			if s, ok := l[1].S(); ok {
				if sy, ok := l[0].Sy(); ok {
					symMap = &symbol.Kv{sy, symbol.New(s)}
					return
				}
				err = fmt.Errorf(
					"First element of import: Expected Symbol. Actual %v",
					l[0],
				)
				return
			}
			err = fmt.Errorf(
				"Second element of import: Expected String. Actual %v",
				l[1],
			)
			return
		}
		err = fmt.Errorf(
			"Import: Expected a Procedure with 2 elements, but it has %v",
			len(l),
		)
		return
	}

	err = fmt.Errorf(
    "Import: Expected Procedure or String. Actual: %v",
    tk.Type(),
  )
	return
}

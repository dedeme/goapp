// Copyright 14-Jan-2021 ºDeme
// GNU General Public License - V3 <http://www.gnu.org/licenses/>

// Imports management.
package imports

import (
	"fmt"
	"github.com/dedeme/dmcoffee/symbol"
	"github.com/dedeme/dmcoffee/token"
)

var imported = map[symbol.T]map[symbol.T]*token.T{}

// Adds a new module.
//    id: Module identifier.
func Add(id symbol.T) {
	imported[id] = map[symbol.T]*token.T{}
}

// Returns the heap of source 'id'.
// If source has not been loaded yeat, returns heap == null, ok == true.
// If source has not a heap with 'id', returns ok == false.
//    id: Module identifier.
func Get(id symbol.T) (heap map[symbol.T]*token.T, ok bool) {
	heap, ok = imported[id]
	return
}

// Adds a key to imported map.
//    id   : Module identifier.
//    key  : Key
//    value: Value
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

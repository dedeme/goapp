// Copyright 20-Feb-2022 ºDeme
// GNU General Public License - V3 <http://www.gnu.org/licenses/>

// Modules regiter.
package modules

import (
	"github.com/dedeme/kut/module"
)

var list []*module.T

// Add an empty module to module list.
//  ix: Index of module file.
func Add(ix int) {
	for len(list) < ix+1 {
		list = append(list, nil)
	}
	list[ix] = module.New(nil, nil, nil)
}

// Set the value of module 'ix'.
//  ix : Index of module (and module file)
//  mod: Module
func Set(ix int, mod *module.T) {
	list[ix] = mod
}

// Returns the module 'ix' or false if such module does not exist.
//   ix: Index of mudule (and module file)
// NOTE: Can return empty modules.
func Get(ix int) (mod *module.T, ok bool) {
	if len(list) > ix {
		mod = list[ix]
		ok = mod != nil
	}
	return
}

// Returns the module 'ix' or raise a fail if module does not exist.
//   ix: Index of mudule (and module file)
// NOTE: Can return empty modules.
func GetOk(ix int) *module.T {
	return list[ix]
}

// Returns list of modules.
// Each index match the file index of module.
// If a module is empty, its index value is 'nil'.
func List() (l []*module.T) {
	for _, e := range list {
		if e.Code == nil {
			l = append(l, nil)
		} else {
			l = append(l, e)
		}
	}
	return
}

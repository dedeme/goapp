// Copyright 20-Feb-2022 ºDeme
// GNU General Public License - V3 <http://www.gnu.org/licenses/>

// Modules regiter.
package modules

import (
	"github.com/dedeme/kut/module"
)

var list []*module.T

func Add(ix int) {
	for len(list) < ix+1 {
		list = append(list, nil)
	}
	list[ix] = module.New(nil, nil, nil)
}

func Set(ix int, mod *module.T) {
	list[ix] = mod
}

func Get(ix int) (mod *module.T, ok bool) {
	if len(list) > ix {
		mod = list[ix]
		ok = mod != nil
	}
	return
}

func GetOk(ix int) *module.T {
	return list[ix]
}

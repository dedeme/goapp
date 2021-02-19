// Copyright 22-Jan-2021 ºDeme
// GNU General Public License - V3 <http://www.gnu.org/licenses/>

// Symbol-value map.
package heap

import (
	"github.com/dedeme/dmcoffee/symbol"
	"github.com/dedeme/dmcoffee/value"
)

type T map[symbol.T]value.T

func New() T {
	return T{}
}

func From(hp map[symbol.T]value.T) T {
	return hp
}

func (hp T) Cast() map[symbol.T]value.T {
	return hp
}

func (hp T) Copy() T {
	r := New()
	for k, v := range hp {
		r[k] = v
	}
	return r
}

func (hp T) Get(key symbol.T) (v value.T, ok bool) {
	v, ok = hp[key]
	return
}

func (hp T) Put(key symbol.T, val value.T) bool {
	if _, ok := hp[key]; ok {
		return false
	}
	hp[key] = val
	return true
}

// Copyright 23-Feb-2022 ºDeme
// GNU General Public License - V3 <http://www.gnu.org/licenses/>

package runner

import (
	"github.com/dedeme/kut/builtin/bfunction"
)

type BModuleT struct {
	Name string
}

func getModule(mod string) (md *BModuleT, ok bool) {
	switch mod {
	case "arr", "b64", "bytes", "cryp", "iter", "js", "file", "map",
		"math", "path", "str", "tcp", "thread", "time", "sys":
		md = &BModuleT{mod}
		ok = true
	}
	return
}

// Returns the built-in function 'fname' of module 'mod'.
//    mod  : Built-in module.
//    fname: Function name.
func getFunction(mod *BModuleT, fname string) (f *bfunction.T, ok bool) {
	switch mod.Name {
	case "arr":
		f, ok = arrGet(fname)
	case "b64":
		f, ok = b64Get(fname)
	case "bytes":
		f, ok = bytesGet(fname)
	case "cryp":
		f, ok = crypGet(fname)
	case "file":
		f, ok = fileGet(fname)
	case "iter":
		f, ok = iterGet(fname)
	case "js":
		f, ok = jsGet(fname)
	case "map":
		f, ok = mapGet(fname)
	case "math":
		f, ok = mathGet(fname)
	case "path":
		f, ok = pathGet(fname)
	case "str":
		f, ok = strGet(fname)
	case "tcp":
		f, ok = tcpGet(fname)
	case "thread":
		f, ok = threadGet(fname)
	case "time":
		f, ok = timeGet(fname)
	case "sys":
		f, ok = sysGet(fname)
	}
	return
}

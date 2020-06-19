// Copyright 08-Apr-2020 ºDeme
// GNU General Public License - V3 <http://www.gnu.org/licenses/>

package directs

import (
	format "fmt"
	"strings"
)

func Process(l string) (r string, err string) {
	key, value := Kv(l)

	switch key {
	case "for":
		r, err = pfor(value)
	case "each":
		r, err = peach(value)
	case "map":
		r, err = pmap(value)
	case "mapTo":
		r, err = pmapTo(value)
	case "mapFrom":
		r, err = pmapFrom(value)
	case "filter":
		r, err = pfilter(value)
	case "sort":
		r, err = psort(value)
	case "props":
		r, err = pprops(value)
	case "toJs":
		r, err = ptoJs(value)
	case "fromJs":
		r, err = pfromJs(value)
	case "copyright":
		r, err = pcopyright(value)
	case "yes":
		r, err = pyes(value)
	case "not":
		r, err = pnot(value)
	case "eqs":
		r, err = peqs(value)
	case "eqi":
		r, err = peqi(value)
	default:
		err = format.Sprintf("Key '%v' is unknown", key)
	}
	//fmt.Println(r)
	return
}

func fmt(tpl string, ps ...interface{}) string {
	return format.Sprintf(tpl, ps...)
}

// Split functon.
func Kv(tx string) (k, v string) {
	k = tx
	ix := strings.Index(tx, "·")
	if ix != -1 {
		k = strings.TrimSpace(tx[:ix])
		v = strings.TrimSpace(tx[ix+2:])
	}
	return
}

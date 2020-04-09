// Copyright 08-Apr-2020 ºDeme
// GNU General Public License - V3 <http://www.gnu.org/licenses/>

//  Go templates.
//
//  General form
//
//  Templates have the following form:
//    ·<key> <par1> <par2> ...
//  For example:
//    ·for i len(x)
//
//  Concrete forms
//
//    for each map mapTo mapFrom filter sort props toJs fromJs copyright
//
//  for
//    ·for·end
//      ·for·len(x)
//    ----------
//    ·for·var·end
//      ·for·i·len(x)
//    ----------
//    ·for·var·start·end
//      ·for·i·3·len(x)
//    ----------
//    ·for·var·start·end·step
//      ·for·i·3·len(x)·4
//
//  each
//    ·each·collection
//      ·each·ls
//    ----------
//    ·each·elem·collection
//      ·each·e·ls
//    ----------
//    ·each·index·elem·collection
//      ·for·i·e·ls
//
//  map
//    ·map·source·taget·targetType·function
//      ·map·src·tg·string·fcopy(e + 4)
//
//  mapTo
//    ·mapTo·source·fn
//      ·mapTo·src·json.Ws(e)
//
//  mapFrom
//    ·mapFrom·index·targetType·function
//      ·mapFrom·2·string·e.Rs()
//
//  filter
//    ·filter·source·taget·type·function
//      ·filter·src·tg·string·e > 33
//
//  sort
//    ·sort·iname·type
//      ·sort·sortInt·string
//
//  props
//    ·props·receiver·recType·[@]property·propType·[@]property·propType...
//      ·props·person·Person·id·int·@name·string
//
//  toJs
//    ·toJs·receiver·recType
//      ·toJs·person·Person
//
//  fromJs
//    ·toJs·type
//      ·toJs·Person
//
//  copyright
//    ·copyright
//      ·copyright
package tpls

import (
	format "fmt"
	"strings"
)

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

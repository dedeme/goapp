// Copyright 08-Apr-2020 ºDeme
// GNU General Public License - V3 <http://www.gnu.org/licenses/>

package tpls

import (
	"time"
)

// copyright template.
//
// Forms:
//    ·copyright
//      ·copyright
//      -----------
//      func (person *Person) ToJs() json.T {
//        return json.Wa([]json.T{})
//      }
func Pcopyright(pars string) (r string, err string) {
	r = fmt(
"// Copyright %v ºDeme\n"+
"// GNU General Public License - V3 <http://www.gnu.org/licenses/>\n\n",
		time.Now().Format("02-Jan-2006"),
	)
	return
}

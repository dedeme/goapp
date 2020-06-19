// Copyright 08-Apr-2020 ºDeme
// GNU General Public License - V3 <http://www.gnu.org/licenses/>

package directs

// each template.
//
// Forms:
//    ·each·collection
//      ·each·ls
//      ----------
//      for _, e := range ls {
//
//      }
//      --------
//    ·each·elem·collection
//      ·each·e·ls
//      ----------
//      for _, e := range ls {
//
//      }
//      --------
//    ·each·index·elem·collection
//      ·for·i·e·ls
//      -----------
//      for i, e := range ls {
//
//      }
func peach(pars string) (r string, err string) {
	p1, rs := Kv(pars)
	if p1 == "" {
		err = fmt("Parameters missings in '%v'", pars)
		return
	}
	p2, rs := Kv(rs)
	if p2 == "" {
		r = fmt("for _, e := range %v {\n\n}", p1)
		return
	}
	p3, rs := Kv(rs)
	if p3 == "" {
		r = fmt("for _, %v := range %v {\n\n}", p1, p2)
		return
	}
	r = fmt("for %v, %v := range %v {\n\n}", p1, p2, p3)

	return
}

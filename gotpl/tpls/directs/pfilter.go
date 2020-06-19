// Copyright 08-Apr-2020 ºDeme
// GNU General Public License - V3 <http://www.gnu.org/licenses/>

package directs

// map template.
//
// Forms:
//    ·filter·source·taget·type·function
//      ·filter·src·tg·string·e > 33
//      -----------
//      var tg []string
//      for _, e := range src {
//        if e > 33 {
//          tg = append(tg, e)
//        }
//      }
func pfilter(pars string) (r string, err string) {
	src, rs := Kv(pars)
	if src == "" {
		err = fmt("Parameters missings in '%v'", pars)
		return
	}
	tg, rs := Kv(rs)
	if tg == "" {
		err = fmt("Parameters missings in '%v'", pars)
		return
	}
	tp, rs := Kv(rs)
	if tp == "" {
		err = fmt("Parameters missings in '%v'", pars)
		return
	}
	fn, rs := Kv(rs)
	if fn == "" {
		err = fmt("Parameters missings in '%v'", pars)
		return
	}
	r = fmt(
		"var %v []%v\nfor _, e := range %v {\nif %v {\n"+
			"%v = append(%v, e)\n}\n}",
		tg, tp, src, fn, tg, tg,
	)

	return
}

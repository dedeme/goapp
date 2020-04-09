// Copyright 08-Apr-2020 ºDeme
// GNU General Public License - V3 <http://www.gnu.org/licenses/>

package tpls

// map template.
//
// Forms:
//    ·map·source·taget·targetType·function
//      ·map·src·tg·string·fcopy(e + 4)
//      -----------
//      var tg []string
//      for _, e := range src {
//        tg = append(tg, fcopy(e+4))
//      }
func Pmap(pars string) (r string, err string) {
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
		"var %v []%v\nfor _, e := range %v {\n"+
			"%v = append(%v, %v)\n}\n",
		tg, tp, src, tg, tg, fn,
	)

	return
}

// map to []json.T template.
//
// Forms:
//    ·mapTo·source·fn
//      ·mapTo·src·json.Ws(e)
//      -----------
//      var namesJs []json.T
//      for _, e := range names {
//        namesJs = append(namesJs, json.Ws(e))
//      }
func PmapTo(pars string) (r string, err string) {
	p1, rs := Kv(pars)
	if p1 == "" {
		err = fmt("Parameters missings in '%v'", pars)
		return
	}
	p2, rs := Kv(rs)
	if p2 == "" {
		err = fmt("Parameters missings in '%v'", pars)
		return
	}
	r = fmt(
		"var %vJs []json.T\nfor _, e := range %v {\n"+
			"%vJs = append(%vJs, %v)\n}\n",
		p1, p1, p1, p1, p2,
	)

	return
}

// map from Json template.
//
// Forms:
//    ·mapFrom·index·targetType·function
//      ·mapFrom·2·string·e.Rs()
//      -----------
//      var tg []string
//      for _, e := range src {
//        tg = append(tg, fcopy(e+4))
//      }
func PmapFrom(pars string) (r string, err string) {
	p1, rs := Kv(pars)
	if p1 == "" {
		err = fmt("Parameters missings in '%v'", pars)
		return
	}
	p2, rs := Kv(rs)
	if p2 == "" {
		err = fmt("Parameters missings in '%v'", pars)
		return
	}
	p3, rs := Kv(rs)
	if p3 == "" {
		err = fmt("Parameters missings in '%v'", pars)
		return
	}
	r = fmt(
		"var a%v []%v\nfor _, e := range a[%v] {\n"+
			"a%v = append(a%v, %v)\n}\n",
		p1, p2, p1, p1, p1, p3,
	)

	return
}

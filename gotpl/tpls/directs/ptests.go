// Copyright 24-Apr-2020 ºDeme
// GNU General Public License - V3 <http://www.gnu.org/licenses/>

package directs

// yes template.
//
// Forms:
//    ·yes
//      ·yes·cond
//      --------
//      if !(cond) {
//        t.Fatal(fail)
//      }
func pyes(pars string) (r string, err string) {
	v, _ := Kv(pars)
	if v == "" {
		err = fmt("Parameters missings in '%v'", pars)
		return
	}

	r = fmt(
		"if !(%v) {\n"+
			"  t.Fatal(fail)\n"+
			"}",
		v,
	)
	return
}

// not template.
//
// Forms:
//    ·not
//      ·not·cond
//      --------
//      if cond {
//        t.Fatal(failNot)
//      }
func pnot(pars string) (r string, err string) {
	v, _ := Kv(pars)
	if v == "" {
		err = fmt("Parameters missings in '%v'", pars)
		return
	}

	r = fmt(
		"if %v {\n"+
			"  t.Fatal(failNot)\n"+
			"}",
		v,
	)
	return
}

// eqs template.
//
// Forms:
//    ·eqs
//      ·eqs·str1·"abc"
//      --------
//      if r := eqs(str1, "abc"); r != "" {
//        t.Fatal(r)
//      }
func peqs(pars string) (r string, err string) {
	v1, rs := Kv(pars)
	if v1 == "" {
		err = fmt("Parameters missings in '%v'", pars)
		return
	}

	v2, _ := Kv(rs)
	if v2 == "" {
		err = fmt("Parameters missings in '%v'", pars)
		return
	}

	r = fmt(
		"if r := eqs(\n"+
			"  %v,\n"+
			"  %v,\n"+
			"); r != \"\" {\n"+
			"  t.Fatal(r)\n"+
			"}",
		v1, v2,
	)
	return
}

// eqi template.
//
// Forms:
//    ·eqi
//      ·eqi·n1·33
//      --------
//      if r := eqi(n1, 33); r != "" {
//        t.Fatal(r)
//      }
func peqi(pars string) (r string, err string) {
	v1, rs := Kv(pars)
	if v1 == "" {
		err = fmt("Parameters missings in '%v'", pars)
		return
	}

	v2, _ := Kv(rs)
	if v2 == "" {
		err = fmt("Parameters missings in '%v'", pars)
		return
	}

	r = fmt(
		"if r := eqi(\n"+
			"  %v,\n"+
			"  %v,\n"+
			"); r != \"\" {\n"+
			"  t.Fatal(r)\n"+
			"}",
		v1, v2,
	)
	return
}

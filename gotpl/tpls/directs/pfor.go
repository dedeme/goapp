// Copyright 08-Apr-2020 ºDeme
// GNU General Public License - V3 <http://www.gnu.org/licenses/>

package directs

// for template.
//
// Forms:
//    ·for·end
//      ·for·len(x)
//      --------
//      for i := 0; i < len(x); i++ {
//
//      }
//    ·for·var·end
//      ·for·i·len(x)
//      --------
//      for i := 0; i < len(x); i++ {
//
//      }
//      --------
//    ·for·var·start·end
//      ·for·i·3·len(x)
//      ---------------
//      for i := 3; i < len(x); i++ {
//
//      }
//    ·for·var·start·end·step
//      ·for·i·3·len(x)·4
//      -----------------
//      for i := 3; i < len(x); i += 4 {
//
//      }
func pfor(pars string) (r string, err string) {
	v, rs := Kv(pars)
	if v == "" {
		err = fmt("Parameters missings in '%v'", pars)
		return
	}
	n1, rs := Kv(rs)
	if n1 == "" {
		r = fmt("for i := 0; i < %v; i++ {\n\n}", v)
		return
	}
	n2, rs := Kv(rs)
	if n2 == "" {
		r = fmt("for %v := 0; %v < %v; %v++ {\n\n}", v, v, n1, v)
		return
	}
	n3, rs := Kv(rs)
	if n3 == "" {
		r = fmt("for %v := %v; %v < %v; %v++ {\n\n}", v, n1, v, n2, v)
		return
	}
	r = fmt("for %v := %v; %v < %v; %v += %v {\n\n}", v, n1, v, n2, v, n3)

	return
}

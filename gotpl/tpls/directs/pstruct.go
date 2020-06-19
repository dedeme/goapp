// Copyright 08-Apr-2020 ºDeme
// GNU General Public License - V3 <http://www.gnu.org/licenses/>

package directs

import (
	"strings"
)

func cap(s string) string {
	return strings.ToUpper(string(s[0])) + s[1:]
}

// Poperties template.
//
// Forms:
//    ·props·receiver·recType·[@]property·propType·[@]property·propType...
//      ·props·person·Person·id·int·@name·string
//      -----------
//      func (person *Person) Id() int {
//        return person.id
//      }
//      func (person *Person) Name() string {
//        return person.name
//      }
//      func (person *Person) SetName(name string) {
//        person.name = name
//      }
func pprops(pars string) (r string, err string) {
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
	p4, rs := Kv(rs)
	if p4 == "" {
		err = fmt("Parameters missings in '%v'", pars)
		return
	}

	for {
		setter := false
		if strings.HasPrefix(p3, "@") {
			p3 = p3[1:]
			setter = true
		}
		if p3 == "" {
			err = fmt("Parameter '@' not valid in '%v'", pars)
			return
		}

		r += fmt(
			"func (%v *%v) %v() %v {\n"+
				"  return %v.%v\n"+
				"}\n",
			p1, p2, cap(p3), p4, p1, p3,
		)
		if setter {
			r += fmt(
				"func (%v *%v) Set%v(%v %v) {\n"+
					"  %v.%v = %v\n"+
					"}\n",
				p1, p2, cap(p3), p3, p4, p1, p3, p3,
			)
		}

		p3, rs = Kv(rs)
		if p3 == "" {
			break
		}
		p4, rs = Kv(rs)
		if p4 == "" {
			err = fmt("Parameter missing at end of '%v'", pars)
			return
		}
	}

	return
}

// toJs template.
//
// Forms:
//    ·toJs·receiver·recType
//      ·toJs·person·Person
//      -----------
//      func (person *Person) ToJs() json.T {
//        return json.Wa([]json.T{})
//      }
func ptoJs(pars string) (r string, err string) {
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
		"func (%v *%v) ToJs() json.T {\n"+
			"  return json.Wa([]json.T{})\n"+
			"}",
		p1, p2,
	)
	return
}

// fromJs template.
//
// Forms:
//    ·toJs·type
//      ·toJs·Person
//      -----------
//      func (person *Person) ToJs() json.T {
//        return json.Wa([]json.T{})
//      }
func pfromJs(pars string) (r string, err string) {
	p1, _ := Kv(pars)
	if p1 == "" {
		err = fmt("Parameters missings in '%v'", pars)
		return
	}
	r = fmt(
		"func FromJs(js json.T) *%v {\n"+
			"  a := js.Ra()\n"+
			"  return &%v {}\n"+
			"}",
		p1, p1,
	)
	return
}

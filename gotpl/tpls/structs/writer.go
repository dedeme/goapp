// Copyright 13-Apr-2020 ºDeme
// GNU General Public License - V3 <http://www.gnu.org/licenses/>

package structs

import (
	"github.com/dedeme/gotpl/tpls/funcs"
	"github.com/dedeme/gotpl/tpls/model"
	"strings"
)

func capitalize(s string) string {
	return strings.ToUpper(s[:1]) + s[1:]
}

func uncapitalize(s string) string {
	return strings.ToLower(s[:1]) + s[1:]
}

func appendDoc(ss []string, doc []string) []string {
	if len(doc) > 0 {
		for _, e := range doc {
			e = "//" + e[1:]
			ss = append(ss, e)
		}
	}
	return ss
}

func wStruct(m *model.T) (result []string) {
	if !m.IsBasicType() {
		result = appendDoc(result, m.Doc)
		result = append(result, fmt("type %v struct {", m.Id))
		for _, e := range m.Params {
			result = appendDoc(result, e.Doc)
			id := cutSignum(e.Id)
			result = append(result, fmt("  %v %v", id, e.Type))
		}
		for _, e := range m.Vars {
			result = appendDoc(result, e.Doc)
			id := cutSignum(e.Id)
			result = append(result, fmt("  %v %v", id, e.Type))
		}
		result = append(result, "}")
	}
	return
}

func wConstructor(m *model.T) (result []string) {
	if !m.IsBasicType() {
		id := "New"
		if !m.IsMainType() {
			id = id + capitalize(m.Id)
			if uncapitalize(m.Id) == m.Id {
				id = uncapitalize(id)
			}
		} else if uncapitalize(m.Id) == m.Id {
			id = "new"
		}

		var pls string
		if len(m.Params) < 4 {
			var ls []string
			for _, e := range m.Params {
				ls = append(ls, fmt("%v %v", cutSignum(e.Id), e.Type))
			}
			pls = strings.Join(ls, ", ")
		} else {
			ls := []string{""}
			for _, e := range m.Params {
				ls = append(ls, fmt("  %v %v,", cutSignum(e.Id), e.Type))
			}
			ls = append(ls, "")
			pls = strings.Join(ls, "\n")
		}

		// Head ----------------------------------------------
		result = append(result, fmt("func "+id+"("+pls+") *"+m.Id+"{"))

		var allFields []string
		for _, e := range m.Params {
			allFields = append(allFields, cutSignum(e.Id))
		}

		var ls []string
		for _, e := range m.Vars {
			allFields = append(allFields, cutSignum(e.Id))
			if e.Value == "" {
				ls = append(ls, fmt("  var %v %v", cutSignum(e.Id), e.Type))
			} else {
				ls = append(ls, fmt("  %v := %v", cutSignum(e.Id), e.Value))
			}
		}

		// Asignations --------------
		result = append(result, ls...)

		pls = ""
		if len(allFields) < 7 {
			var ls []string
			for _, e := range allFields {
				ls = append(ls, fmt("%v", e))
			}
			pls = strings.Join(ls, ", ")
		} else {
			c := 0
			lns := []string{""}
			ls := []string{}
			for _, e := range allFields {
				ls = append(ls, fmt("%v", e))
				c++
				if c == 6 {
					c = 0
					lns = append(lns, fmt("    %v,", strings.Join(ls, ", ")))
					ls = []string{}
				}
			}
			if len(ls) > 0 {
				lns = append(lns, fmt("    %v,", strings.Join(ls, ", ")))
			}
			lns = append(lns, "")
			pls = strings.Join(lns, "\n  ")
		}

		result = append(result, "  return &"+m.Id+"{"+pls+"}")

		result = append(result, "}")
	}
	return
}

func wSetter(mType, id, tp string, doc []string) (result []string) {
	result = appendDoc(result, doc)
	result = append(result, fmt(
		"func (this *%v) Set%v(value %v) {", mType, capitalize(id), tp,
	))
	result = append(result, fmt("  this.%v = value", id))
	result = append(result, "}")
	return
}

func wGetter(mType, id, tp string, doc []string) (result []string) {
	result = appendDoc(result, doc)
	result = append(result, fmt(
		"func (this *%v) %v() %v {", mType, capitalize(id), tp,
	))
	result = append(result, fmt("  return this.%v", id))
	result = append(result, "}")
	return
}

func wGetSet(m *model.T) (result []string) {
	fn := func(id, tp string, doc []string) (result []string) {
		if id[0] != '-' {
			withSetter := false
			if id[0] == '+' {
				withSetter = true
				id = id[1:]
			}

			if uncapitalize(id) == id {
				if withSetter {
					tmp := wSetter(m.Id, id, tp, doc)
					result = append(result, tmp...)
				}
				if id[0] != '-' {
					tmp := wGetter(m.Id, id, tp, doc)
					result = append(result, tmp...)
				}
			}
		}
		return
	}

	for _, e := range m.Params {
		rs := fn(e.Id, e.Type, e.Doc)
		result = append(result, rs...)
	}
	for _, e := range m.Vars {
		rs := fn(e.Id, e.Type, e.Doc)
		result = append(result, rs...)
	}
	return
}

func write(m *model.T) (result []string, err error) {
	result = append(result, wStruct(m)...)
	result = append(result, wConstructor(m)...)
	result = append(result, wGetSet(m)...)
	tmp, err := funcs.Process(m)
	if err == nil {
		result = append(result, tmp...)
	}
	return
}

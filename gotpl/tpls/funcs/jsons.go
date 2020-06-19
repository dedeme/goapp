// Copyright 13-Apr-2020 ºDeme
// GNU General Public License - V3 <http://www.gnu.org/licenses/>

package funcs

import (
	"github.com/dedeme/gotpl/tpls/model"
	"strings"
)

func removeSpaces(s string) (r string) {
	for _, e := range s {
		if e > ' ' {
			r += string(e)
		}
	}
	return
}

func isSlice(tp string) bool {
	return strings.HasPrefix(tp, "[]")
}

func isMap(tp string) bool {
	return strings.HasPrefix(tp, "map[string]")
}

func splitType(tp string) (rs []string, err error) {
	rs = strings.Split(tp, ".")
	if len(rs) == 1 {
		rs = strings.Split("."+tp, ".")
	} else if len(rs) > 2 {
		err = efmt("Bad type (%v dots)", len(rs))
	}
	return
}

func wToJs(receiver string, m *model.T) (result []string, err error) {
	var apres []string
	var mpres []string
	var posts []string

	simplyTo := func(id, tp string) string {
		switch tp {
		case "bool":
			return fmt("json.Wb(%v)", id)
		case "int":
			return fmt("json.Wi(%v)", id)
		case "int32":
			return fmt("json.Wi(int(%v))", id)
		case "int64":
			return fmt("json.Wl(%v)", id)
		case "float32":
			return fmt("json.Wf(%v)", id)
		case "float64":
			return fmt("json.Wd(%v)", id)
		case "string":
			return fmt("json.Ws(%v)", id)
		default:
			return fmt("%v.ToJs()", id)
		}
	}

	fn := func(id, tp string) (err error) {
		id = cutSignum(id)
		tp = removeSpaces(tp)
		if isSlice(tp) {
			if len(apres) == 0 {
				apres = append(apres, "  var _tmp []json.T")
			}
			tp2 := tp[2:]
			if tp2 == "" {
				err = efmt("Type basic is missing in '%v'", tp)
				return
			}
			apres = append(apres, fmt("  for _, e := range this.%v {", id))
			apres = append(apres, fmt(
				"    _tmp = append(_tmp, %v)", simplyTo("e", tp2),
			))
			apres = append(apres, "  }")
			apres = append(apres, fmt("  %v := json.Wa(_tmp)", id))
			apres = append(apres, "  _tmp = []json.T{}")

			posts = append(posts, fmt("    %v,", id))
		} else if isMap(tp) {
			if len(mpres) == 0 {
				mpres = append(mpres, "  var __tmp map[string]json.T")
			}
			tp2 := tp[11:]
			if tp2 == "" {
				err = efmt("Type basic is missing in '%v'", tp)
				return
			}
			mpres = append(mpres, fmt("  for k, v := range this.%v {", id))
			mpres = append(mpres, fmt(
				"    __tmp[k] = %v", simplyTo("v", tp2)))
			mpres = append(mpres, "  }")
			mpres = append(mpres, fmt("  %v := json.Wo(__tmp)", id))
			mpres = append(mpres, "  __tmp = map[string]json.T{}")

			posts = append(posts, fmt("    %v,", id))
		} else {
			id = "this." + id
			posts = append(posts, fmt("    %v,", simplyTo(id, tp)))
		}

		return
	}

	// main function -------------
	for _, e := range m.Params {
		err = fn(e.Id, e.Type)
		if err != nil {
			break
		}
	}
	if err == nil {
		for _, e := range m.Vars {
			err = fn(e.Id, e.Type)
			if err != nil {
				break
			}
		}
	}
	if err == nil {
		result = append(result, fmt(
			"func %v ToJs () json.T {", receiver,
		))
		result = append(result, apres...)
		result = append(result, mpres...)
		result = append(result, "  return json.Wa([]json.T {")
		result = append(result, posts...)
		result = append(result, "  })")
		result = append(result, "}")
	}
	return
}

func wFromJs(m *model.T) (result []string, err error) {
	var apres []string
	var mpres []string
	var posts []string

	funcName := func(tp string) string {
		if tp[0] == '*' {
			tp = tp[1:]
		}
		prefix := ""
		ix := strings.Index(tp, ".")
		if ix != -1 {
			prefix = tp[0 : ix+1]
			tp = tp[ix+1:]
		}
		if tp == "T" || tp == "t" {
			return prefix + "FromJs"
		}
		return prefix + capitalize(tp) + "FromJs"
	}

	simplyFrom := func(id, tp string) string {
		switch tp {
		case "bool":
			return fmt("%v.Rb()", id)
		case "int":
			return fmt("%v.Ri()", id)
		case "int32":
			return fmt("int32(%v.Ri())", id)
		case "int64":
			return fmt("%v.Rl()", id)
		case "float32":
			return fmt("%v.Rf()", id)
		case "float64":
			return fmt("%v.Rd()", id)
		case "string":
			return fmt("%v.Rs()", id)
		default:
			return fmt("%v(%v)", funcName(tp), id)
		}
	}

	fn := func(ix int, id, tp string) (err error) {
		id = cutSignum(id)
		tp = removeSpaces(tp)

		if isSlice(tp) {
			tp2 := tp[2:]
			if tp2 == "" {
				err = efmt("Type basic is missing in '%v'", tp)
				return
			}

			apres = append(apres, fmt("  var %v %v", id, tp))
			apres = append(apres, fmt(
				"  for _, e := range %v.Ra() {", fmt("a[%v]", ix),
			))
			apres = append(apres, fmt(
				"    %v = append(%v, %v)", id, id, simplyFrom("e", tp2),
			))
			apres = append(apres, "  }")

			posts = append(posts, fmt("    %v,", id))
		} else if isMap(tp) {
			tp2 := tp[11:]
			if tp2 == "" {
				err = efmt("Type basic is missing in '%v'", tp)
				return
			}
			mpres = append(mpres, fmt("  var %v %v", id, tp))
			mpres = append(mpres, fmt(
				"  for k, v := range %v.Ro() {", fmt("a[%v]", ix),
			))
			mpres = append(mpres, fmt(
				"    %v[k] = %v", id, simplyFrom("v", tp2)))
			mpres = append(mpres, "  }")

			posts = append(posts, fmt("    %v,", id))
		} else {
			posts = append(posts, fmt("    %v,", simplyFrom(fmt("a[%v]", ix), tp)))
		}
		return
	}

	// main function -------------
	ix := 0
	for _, e := range m.Params {
		err = fn(ix, e.Id, e.Type)
		ix++
		if err != nil {
			break
		}
	}
	if err == nil {
		for _, e := range m.Vars {
			err = fn(ix, e.Id, e.Type)
			ix++
			if err != nil {
				break
			}
		}
	}
	if err == nil {
		result = append(result, fmt(
			"func %v (js json.T) *%v {", funcName(m.Id), m.Id,
		))
		result = append(result, "  a := js.Ra()")
		result = append(result, apres...)
		result = append(result, mpres...)
		result = append(result, fmt("  return &%v{", m.Id))
		result = append(result, posts...)
		result = append(result, "  }")
		result = append(result, "}")
	}
	return
}

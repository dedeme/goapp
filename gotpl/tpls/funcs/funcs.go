// Copyright 13-Apr-2020 ºDeme
// GNU General Public License - V3 <http://www.gnu.org/licenses/>

package funcs

import (
	format "fmt"
	"github.com/dedeme/gotpl/tpls/model"
	"strings"
)

func fmt(tpl string, ps ...interface{}) string {
	return format.Sprintf(tpl, ps...)
}

func efmt(tpl string, ps ...interface{}) error {
	return format.Errorf(tpl, ps...)
}

func println(s string) {
	format.Println(s)
}

func cutSignum(s string) string {
	if s[0] == '-' || s[0] == '+' {
		return s[1:]
	}
	return s
}

func capitalize(s string) string {
	return strings.ToUpper(s[:1]) + s[1:]
}

func Process(m *model.T) (result []string, err error) {
	if len(m.Funcs) == 0 {
		return
	}

	for _, e := range m.Funcs {
		receiver := fmt("(this *%v)", m.Id)
		if m.IsBasicType() {
			receiver = fmt("(this %v)", m.Id)
		}

		var tmp []string
		switch e.Id {
		case "toJs":
			tmp, err = wToJs(receiver, m)
		case "fromJs":
			tmp, err = wFromJs(m)
		case "filter":
			tmp = wFilter(m)
		case "map":
      if len(e.Params) != 2 {
        err = efmt("'map' need two arguments: functionId and targetType")
      } else {
        tmp = wMap(m, e.Params[0], e.Params[1])
      }
		case "take":
			tmp = wTake(m)
		case "drop":
			tmp = wDrop(m)
		default:
			err = efmt("Unkown function '%s'", e.Id)
		}
		if err != nil {
			err = efmt("Line %v: %v", e.Nline, err.Error())
			break
		}
		result = append(result, tmp...)
	}
	return
}

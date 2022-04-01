// Copyright 02-Mar-2022 ºDeme
// GNU General Public License - V3 <http://www.gnu.org/licenses/>

package runner

import (
	"github.com/dedeme/kut/builtin/bfail"
	"github.com/dedeme/kut/builtin/bfunction"
	"github.com/dedeme/kut/expression"
	"regexp"
)

// \s, s -> [s...]
func regexMatches(args []*expression.T) (ex *expression.T, err error) {
	switch s := (args[0].Value).(type) {
	case string:
		switch rg := (args[1].Value).(type) {
		case string:
			re := regexp.MustCompile(rg)
			rs := re.FindAllString(s, -1)
			var a []*expression.T
			if rs != nil {
				for _, e := range rs {
					a = append(a, expression.MkFinal(e))
				}
			}
			ex = expression.MkFinal(a)
		default:
			err = bfail.Type(args[1], "string")
		}
	default:
		err = bfail.Type(args[0], "string")
	}
	return
}

// \s, s, s -> s
func regexReplace(args []*expression.T) (ex *expression.T, err error) {
	switch s := (args[0].Value).(type) {
	case string:
		switch rg := (args[1].Value).(type) {
		case string:
			switch sub := (args[2].Value).(type) {
			case string:
				re := regexp.MustCompile(rg)
				ex = expression.MkFinal(re.ReplaceAllString(s, sub))
			default:
				err = bfail.Type(args[2], "string")
			}
		default:
			err = bfail.Type(args[1], "string")
		}
	default:
		err = bfail.Type(args[0], "string")
	}
	return
}

func regexGet(fname string) (fn *bfunction.T, ok bool) {
	ok = true
	switch fname {
	case "matches":
		fn = bfunction.New(2, regexMatches)
	case "replace":
		fn = bfunction.New(3, regexReplace)
	default:
		ok = false
	}

	return
}

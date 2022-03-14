// Copyright 03-Mar-2022 ºDeme
// GNU General Public License - V3 <http://www.gnu.org/licenses/>

package runner

import (
	"github.com/dedeme/kut/builtin/bfail"
	"github.com/dedeme/kut/builtin/bfunction"
	"github.com/dedeme/kut/expression"
	"path"
)

/// \s -> s
func pathBase(args []*expression.T) (ex *expression.T, err error) {
	switch p := (args[0].Value).(type) {
	case string:
		ex = expression.MkFinal(path.Base(p))
	default:
		err = bfail.Type(args[0], "string")
	}
	return
}

/// \s -> s
func pathCanonical(args []*expression.T) (ex *expression.T, err error) {
	switch p := (args[0].Value).(type) {
	case string:
		ex = expression.MkFinal(path.Clean(p))
	default:
		err = bfail.Type(args[0], "string")
	}
	return
}

/// \[s...] -> s
func pathCat(args []*expression.T) (ex *expression.T, err error) {
	switch a := (args[0].Value).(type) {
	case []*expression.T:
		var ss []string
		for _, e := range a {
			switch v := (e.Value).(type) {
			case string:
				ss = append(ss, v)
			default:
				err = bfail.Type(e, "string")
			}
		}
		if err == nil {
			ex = expression.MkFinal(path.Join(ss...))
		}
	default:
		err = bfail.Type(args[0], "array")
	}
	return
}

/// \s -> s
func pathExtension(args []*expression.T) (ex *expression.T, err error) {
	switch p := (args[0].Value).(type) {
	case string:
		ex = expression.MkFinal(path.Ext(p))
	default:
		err = bfail.Type(args[0], "string")
	}
	return
}

/// \s -> s
func pathParent(args []*expression.T) (ex *expression.T, err error) {
	switch p := (args[0].Value).(type) {
	case string:
		ex = expression.MkFinal(path.Dir(p))
	default:
		err = bfail.Type(args[0], "string")
	}
	return
}

func pathGet(fname string) (fn *bfunction.T, ok bool) {
	ok = true
	switch fname {
	case "base":
		fn = bfunction.New(1, pathBase)
	case "canonical":
		fn = bfunction.New(1, pathCanonical)
	case "cat":
		fn = bfunction.New(1, pathCat)
	case "extension":
		fn = bfunction.New(1, pathExtension)
	case "parent":
		fn = bfunction.New(1, pathParent)
	default:
		ok = false
	}

	return
}

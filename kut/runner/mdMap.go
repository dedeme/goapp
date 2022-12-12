// Copyright 23-Feb-2022 ºDeme
// GNU General Public License - V3 <http://www.gnu.org/licenses/>

package runner

import (
	"fmt"
	"github.com/dedeme/kut/builtin/bfail"
	"github.com/dedeme/kut/builtin/bfunction"
	"github.com/dedeme/kut/expression"
	"github.com/dedeme/kut/iterator"
)

// \m -> m
func mapCopy(args []*expression.T) (ex *expression.T, err error) {
	switch a := (args[0].Value).(type) {
	case map[string]*expression.T:
		m := map[string]*expression.T{}
		for k, v := range a {
			m[k] = v
		}
		ex = expression.MkFinal(m)
	default:
		err = bfail.Type(args[0], "map")
	}
	return
}

// \[[s, *]...] -> m
func mapFromArr(args []*expression.T) (ex *expression.T, err error) {
	switch a := (args[0].Value).(type) {
	case []*expression.T:
		m := map[string]*expression.T{}
		for _, ae := range a {
			switch a2 := (ae.Value).(type) {
			case []*expression.T:
				if len(a2) != 2 {
					err = bfail.Mk(fmt.Sprintf(
						"Expected array with 2 elements. Found %v", len(a2)))
					return
				}
				switch v := (a2[0].Value).(type) {
				case string:
					m[v] = a2[1]
				default:
					err = bfail.Type(a2[0], "string")
					return
				}
			default:
				err = bfail.Type(ae, "array")
				return
			}
		}
		ex = expression.MkFinal(m)
	default:
		err = bfail.Type(args[0], "array")
	}
	return
}

// \<iterator>[[s, *]...] -> m
func mapFromIter(args []*expression.T) (ex *expression.T, err error) {
	switch it := (args[0].Value).(type) {
	case *iterator.T:
		m := map[string]*expression.T{}
		for it.HasNext() {
			ae := it.Next()
			switch a2 := (ae.Value).(type) {
			case []*expression.T:
				if len(a2) != 2 {
					err = bfail.Mk(fmt.Sprintf(
						"Expected array with 2 elements. Found %v", len(a2)))
					return
				}
				switch v := (a2[0].Value).(type) {
				case string:
					m[v] = a2[1]
				default:
					err = bfail.Type(a2[0], "string")
					return
				}
			default:
				err = bfail.Type(ae, "array")
				return
			}
		}
		ex = expression.MkFinal(m)
	default:
		err = bfail.Type(args[0], "<iterator>")
	}
	return
}

func mapFget(args []*expression.T) (ex *expression.T, err error) {
	switch m := (args[0].Value).(type) {
	case map[string]*expression.T:
		switch s := (args[1].Value).(type) {
		case string:
			var ok bool
			if ex, ok = m[s]; !ok {
				err = bfail.Mk("Key '" + s + "' not found")
			}
		default:
			err = bfail.Type(args[1], "string")
		}
	default:
		err = bfail.Type(args[0], "map")
	}
	return
}

// \m, s -> [*]|[]
func mapGetOp(args []*expression.T) (ex *expression.T, err error) {
	switch m := (args[0].Value).(type) {
	case map[string]*expression.T:
		switch s := (args[1].Value).(type) {
		case string:
			ex0, ok := m[s]
      if ok {
				ex = expression.MkFinal([]*expression.T{ex0})
			} else {
        ex = expression.MkFinal([]*expression.T{})
      }
		default:
			err = bfail.Type(args[1], "string")
		}
	default:
		err = bfail.Type(args[0], "map")
	}
	return
}

// \m, s -> b
func mapHasKey(args []*expression.T) (ex *expression.T, err error) {
	switch m := (args[0].Value).(type) {
	case map[string]*expression.T:
		switch s := (args[1].Value).(type) {
		case string:
			_, ok := m[s]
			ex = expression.MkFinal(ok)
		default:
			err = bfail.Type(args[1], "string")
		}
	default:
		err = bfail.Type(args[0], "map")
	}
	return
}

// \m -> [s...]
func mapKeys(args []*expression.T) (ex *expression.T, err error) {
	switch m := (args[0].Value).(type) {
	case map[string]*expression.T:
		var exs []*expression.T
		for k := range m {
			exs = append(exs, expression.MkFinal(k))
		}
		ex = expression.MkFinal(exs)
	default:
		err = bfail.Type(args[0], "map")
	}
	return
}

// \m -> i
func mapSize(args []*expression.T) (ex *expression.T, err error) {
	switch m := (args[0].Value).(type) {
	case map[string]*expression.T:
		ex = expression.MkFinal(int64(len(m)))
	default:
		err = bfail.Type(args[0], "map")
	}
	return
}

// \m, s -> ()
func mapRemove(args []*expression.T) (ex *expression.T, err error) {
	switch m := (args[0].Value).(type) {
	case map[string]*expression.T:
		switch s := (args[1].Value).(type) {
		case string:
			delete(m, s)
		default:
			err = bfail.Type(args[1], "string")
		}
	default:
		err = bfail.Type(args[0], "map")
	}
	return
}

// \m -> [[s, *]...]
func mapToArr(args []*expression.T) (ex *expression.T, err error) {
	switch m := (args[0].Value).(type) {
	case map[string]*expression.T:
		var exs []*expression.T
		for k, v := range m {
			e := expression.MkFinal([]*expression.T{
				expression.MkFinal(k), v})
			exs = append(exs, e)
		}
		ex = expression.MkFinal(exs)
	default:
		err = bfail.Type(args[0], "map")
	}
	return
}

// \m -> <iterator>[[s, *]...]
func mapToIter(args []*expression.T) (ex *expression.T, err error) {
	switch m := (args[0].Value).(type) {
	case map[string]*expression.T:
		var keys []string
		for k := range m {
			keys = append(keys, k)
		}
		l := len(keys)
		ix := 0
		hasNext := func() bool {
			return ix < l
		}
		next := func() *expression.T {
			k := keys[ix]
			e := expression.MkFinal([]*expression.T{
				expression.MkFinal(k), m[k]})
			ix++
			return e
		}
		ex = expression.MkFinal(iterator.New(hasNext, next))
	default:
		err = bfail.Type(args[0], "map")
	}
	return
}

// \m -> a
func mapValues(args []*expression.T) (ex *expression.T, err error) {
	switch m := (args[0].Value).(type) {
	case map[string]*expression.T:
		var exs []*expression.T
		for _, v := range m {
			exs = append(exs, v)
		}
		ex = expression.MkFinal(exs)
	default:
		err = bfail.Type(args[0], "map")
	}
	return
}

func mapGet(fname string) (fn *bfunction.T, ok bool) {
	ok = true
	switch fname {
	case "copy":
		fn = bfunction.New(1, mapCopy)
	case "fromArr":
		fn = bfunction.New(1, mapFromArr)
	case "fromIter":
		fn = bfunction.New(1, mapFromIter)
	case "get":
		fn = bfunction.New(2, mapGetOp)
	case "hasKey":
		fn = bfunction.New(2, mapHasKey)
	case "keys":
		fn = bfunction.New(1, mapKeys)
	case "size":
		fn = bfunction.New(1, mapSize)
	case "remove":
		fn = bfunction.New(2, mapRemove)
	case "toArr":
		fn = bfunction.New(1, mapToArr)
	case "toIter":
		fn = bfunction.New(1, mapToIter)
	case "values":
		fn = bfunction.New(1, mapValues)
	default:
		ok = false
	}
	return
}

// Copyright 23-Feb-2022 ºDeme
// GNU General Public License - V3 <http://www.gnu.org/licenses/>

package runner

import (
	"fmt"
	"github.com/dedeme/kut/builtin/bfail"
	"github.com/dedeme/kut/builtin/bfunction"
	"github.com/dedeme/kut/expression"
	"github.com/dedeme/kut/function"
	"github.com/dedeme/kut/iterator"
	"math/rand"
	"sort"
	"strings"
)

// \a, \*->b -> b
func arrAll(args []*expression.T) (ex *expression.T, err error) {
	switch a := (args[0].Value).(type) {
	case []*expression.T:
		switch fn := (args[1].Value).(type) {
		case *function.T:
			rt := true
			for _, e := range a {
				var r *expression.T
				r, err = solveIsolateFunction(fn, []*expression.T{e}) // exSolver.go
				if err != nil {
					break
				}
				switch v := (r.Value).(type) {
				case bool:
					rt = v
				default:
					err = bfail.Type(r, "bool")
				}
				if err != nil || !rt {
					break
				}
			}
			if err == nil {
				ex = expression.MkFinal(rt)
			}
		default:
			err = bfail.Type(args[1], "function")
		}
	default:
		err = bfail.Type(args[0], "array")
	}
	return
}

// \a, \*->b -> b
func arrAny(args []*expression.T) (ex *expression.T, err error) {
	switch a := (args[0].Value).(type) {
	case []*expression.T:
		switch fn := (args[1].Value).(type) {
		case *function.T:
			rt := false
			for _, e := range a {
				var r *expression.T
				r, err = solveIsolateFunction(fn, []*expression.T{e}) // exSolver.go
				if err != nil {
					break
				}
				switch v := (r.Value).(type) {
				case bool:
					rt = v
				default:
					err = bfail.Type(r, "bool")
				}
				if err != nil || rt {
					break
				}
			}
			if err == nil {
				ex = expression.MkFinal(rt)
			}
		default:
			err = bfail.Type(args[1], "function")
		}
	default:
		err = bfail.Type(args[0], "array")
	}
	return
}

// a -> ()
func arrClear(args []*expression.T) (ex *expression.T, err error) {
	switch (args[0].Value).(type) {
	case []*expression.T:
		args[0].Value = []*expression.T{}
	default:
		err = bfail.Type(args[0], "array")
	}
	return
}

// a -> a
func arrCopy(args []*expression.T) (ex *expression.T, err error) {
	switch a := (args[0].Value).(type) {
	case []*expression.T:
		ex = expression.MkFinal(append([]*expression.T{}, a...))
	default:
		err = bfail.Type(args[0], "array")
	}
	return
}

// \a, i -> a
func arrDrop(args []*expression.T) (ex *expression.T, err error) {
	switch a := (args[0].Value).(type) {
	case []*expression.T:
		switch n := (args[1].Value).(type) {
		case int64:
			l := int64(len(a))
			if n < 0 {
				n = 0
			}
			var rt []*expression.T
			for c := n; c < l; c++ {
				rt = append(rt, a[c])
			}
			ex = expression.MkFinal(rt)
		default:
			err = bfail.Type(args[1], "int")
		}
	default:
		err = bfail.Type(args[0], "array")
	}
	return
}

// \a, \*->b -> a
func arrDropWhile(args []*expression.T) (ex *expression.T, err error) {
	switch a := (args[0].Value).(type) {
	case []*expression.T:
		switch fn := (args[1].Value).(type) {
		case *function.T:
			begin := 0
			for _, e := range a {
				var r *expression.T
				r, err = solveIsolateFunction(fn, []*expression.T{e}) // exSolver.go
				if err != nil {
					return
				}
				stop := false
				switch v := (r.Value).(type) {
				case bool:
					if !v {
						stop = true
					} else {
						begin++
					}
				default:
					err = bfail.Type(r, "bool")
					return
				}
				if stop {
					break
				}
			}
			var exs []*expression.T
			for i := begin; i < len(a); i++ {
				exs = append(exs, a[i])
			}
			ex = expression.MkFinal(exs)
		default:
			err = bfail.Type(args[1], "function")
		}
	default:
		err = bfail.Type(args[0], "array")
	}
	return
}

// \a -> [a, a]
func arrDuplicates(args []*expression.T) (ex *expression.T, err error) {
	switch a := (args[0].Value).(type) {
	case []*expression.T:
		switch fn := (args[1].Value).(type) {
		case *function.T:
			contains := func(es []*expression.T, e *expression.T) bool {
				for _, ees := range es {
					r, er := solveIsolateFunction(fn, []*expression.T{e, ees}) // exSolver.go
					if er != nil {
						panic(er.Error())
					}
					switch v := (r.Value).(type) {
					case bool:
						if v {
							return true
						}
					default:
						panic(bfail.Type(r, "bool"))
					}
				}
				return false
			}
			var dup []*expression.T
			var rest []*expression.T
			for _, e := range a {
				if contains(rest, e) {
					if !contains(dup, e) {
						dup = append(dup, e)
					}
				} else {
					rest = append(rest, e)
				}
			}
			ex = expression.MkFinal([]*expression.T{
				expression.MkFinal(rest), expression.MkFinal(dup),
			})
		default:
			err = bfail.Type(args[0], "function")
		}
	default:
		err = bfail.Type(args[0], "array")
	}
	return
}

// \a -> b
func arrEmpty(args []*expression.T) (ex *expression.T, err error) {
	switch a := (args[0].Value).(type) {
	case []*expression.T:
		ex = expression.MkFinal(len(a) == 0)
	default:
		err = bfail.Type(args[0], "array")
	}
	return
}

// <iterator> -> a
func arrFromIter(args []*expression.T) (ex *expression.T, err error) {
	switch it := (args[0].Value).(type) {
	case *iterator.T:
		var a []*expression.T
		for it.HasNext() {
			a = append(a, it.Next())
		}
		ex = expression.MkFinal(a)
	default:
		err = bfail.Type(args[0], "<iterator>")
	}
	return
}

// \a, \*->b -> a
func arrFilter(args []*expression.T) (ex *expression.T, err error) {
	switch a := (args[0].Value).(type) {
	case []*expression.T:
		switch fn := (args[1].Value).(type) {
		case *function.T:
			var exs []*expression.T
			for _, e := range a {
				var r *expression.T
				r, err = solveIsolateFunction(fn, []*expression.T{e}) // exSolver.go
				if err != nil {
					break
				}
				switch v := (r.Value).(type) {
				case bool:
					if v {
						exs = append(exs, e)
					}
				default:
					err = bfail.Type(r, "bool")
				}
				if err != nil {
					break
				}
			}
			if err == nil {
				ex = expression.MkFinal(exs)
			}
		default:
			err = bfail.Type(args[1], "function")
		}
	default:
		err = bfail.Type(args[0], "array")
	}
	return
}

// \a, \*->b -> ()
func arrFilterIn(args []*expression.T) (ex *expression.T, err error) {
	switch a := (args[0].Value).(type) {
	case []*expression.T:
		switch fn := (args[1].Value).(type) {
		case *function.T:
			var exs []*expression.T
			for _, e := range a {
				var r *expression.T
				r, err = solveIsolateFunction(fn, []*expression.T{e}) // exSolver.go
				if err != nil {
					break
				}
				switch v := (r.Value).(type) {
				case bool:
					if v {
						exs = append(exs, e)
					}
				default:
					err = bfail.Type(r, "bool")
				}
				if err != nil {
					break
				}
			}
			if err == nil {
				args[0].Value = exs
			}
		default:
			err = bfail.Type(args[1], "function")
		}
	default:
		err = bfail.Type(args[0], "array")
	}
	return
}

// \a, \*->b -> [*]
func arrFind(args []*expression.T) (ex *expression.T, err error) {
	switch a := (args[0].Value).(type) {
	case []*expression.T:
		switch fn := (args[1].Value).(type) {
		case *function.T:
			var exfund *expression.T
			for _, e := range a {
				var r *expression.T
				r, err = solveIsolateFunction(fn, []*expression.T{e}) // exSolver.go
				if err != nil {
					break
				}
				switch v := (r.Value).(type) {
				case bool:
					if v {
						exfund = e
					}
				default:
					err = bfail.Type(r, "bool")
				}
				if err != nil || exfund != nil {
					break
				}
			}
			if err == nil {
				if exfund == nil {
					ex = expression.MkFinal([]*expression.T{})
				} else {
					ex = expression.MkFinal([]*expression.T{exfund})
				}
			}
		default:
			err = bfail.Type(args[1], "function")
		}
	default:
		err = bfail.Type(args[0], "array")
	}
	return
}

// \a, i -> *
func arrFget(args []*expression.T) (ex *expression.T, err error) {
	switch a := (args[0].Value).(type) {
	case []*expression.T:
		switch ix := (args[1].Value).(type) {
		case int64:
			ex = a[ix]
		default:
			err = bfail.Type(args[1], "int")
		}
	default:
		err = bfail.Type(args[0], "array")
	}
	return
}

// \a, s -> s
func arrJoin(args []*expression.T) (ex *expression.T, err error) {
	switch a := (args[0].Value).(type) {
	case []*expression.T:
		switch sep := (args[1].Value).(type) {
		case string:
			var ss []string
			for _, e := range a {
				switch s := (e.Value).(type) {
				case string:
					ss = append(ss, s)
				default:
					err = bfail.Type(e, "string")
				}
			}
			if err == nil {
				ex = expression.MkFinal(strings.Join(ss, sep))
			}
		default:
			err = bfail.Type(args[1], "string")
		}
	default:
		err = bfail.Type(args[0], "array")
	}
	return
}

// \a, \*->* -> a
func arrMp(args []*expression.T) (ex *expression.T, err error) {
	switch a := (args[0].Value).(type) {
	case []*expression.T:
		switch fn := (args[1].Value).(type) {
		case *function.T:
			var exs []*expression.T
			for _, e := range a {
				var r *expression.T
				r, err = solveIsolateFunction(fn, []*expression.T{e}) // exSolver.go
				if err != nil {
					break
				}
				exs = append(exs, r)
			}
			if err == nil {
				ex = expression.MkFinal(exs)
			}
		default:
			err = bfail.Type(args[1], "function")
		}
	default:
		err = bfail.Type(args[0], "array")
	}
	return
}

// \i, * -> a
func arrNew(args []*expression.T) (ex *expression.T, err error) {
	switch n := (args[0].Value).(type) {
	case int64:
		if n < 0 {
			n = 0
		}
		a := make([]*expression.T, n)
		e := args[1]
		switch v := (e.Value).(type) {
		case bool, int64, float64, string:
			for i := int64(0); i < n; i++ {
				a[i] = e
			}
			ex = expression.MkFinal(a)
		case []*expression.T:
			for i := int64(0); i < n; i++ {
				var newa []*expression.T
				for _, ve := range v {
					newa = append(newa, ve)
				}
				a[i] = expression.MkFinal(newa)
			}
			ex = expression.MkFinal(a)
		default:
			err = bfail.Type(args[1], "bool", "int", "float", "string", "array")
		}
	default:
		err = bfail.Type(args[0], "int")
	}
	return
}

// \a -> *
func arrPeek(args []*expression.T) (ex *expression.T, err error) {
	switch a := (args[0].Value).(type) {
	case []*expression.T:
		l := len(a)
		if l < 1 {
			err = bfail.Mk("Empty array")
			return
		}
		ex = a[l-1]
	default:
		err = bfail.Type(args[0], "array")
	}
	return
}

// \a -> *
func arrPop(args []*expression.T) (ex *expression.T, err error) {
	switch a := (args[0].Value).(type) {
	case []*expression.T:
		l := len(a)
		if l < 1 {
			err = bfail.Mk("Empty array")
			return
		}
		ex = a[l-1]
		args[0].Value = a[:l-1]
	default:
		err = bfail.Type(args[0], "array")
	}
	return
}

// \a, * -> ()
func arrPush(args []*expression.T) (ex *expression.T, err error) {
	switch a := (args[0].Value).(type) {
	case []*expression.T:
		args[0].Value = append(a, args[1])
	default:
		err = bfail.Type(args[0], "array")
	}
	return
}

// \a, i -> ()
func arrRemove(args []*expression.T) (ex *expression.T, err error) {
	switch a := (args[0].Value).(type) {
	case []*expression.T:
		switch ix := (args[1].Value).(type) {
		case int64:
			args[0].Value = append(a[:ix], a[ix+1:]...)
		default:
			err = bfail.Type(args[1], "int")
		}
	default:
		err = bfail.Type(args[0], "array")
	}
	return
}

// \a, i, i -> ()
func arrRemoveRange(args []*expression.T) (ex *expression.T, err error) {
	switch a := (args[0].Value).(type) {
	case []*expression.T:
		switch begin := (args[1].Value).(type) {
		case int64:
			switch end := (args[2].Value).(type) {
			case int64:
				args[0].Value = append(a[:begin], a[end:]...)
			default:
				err = bfail.Type(args[2], "int")
			}
		default:
			err = bfail.Type(args[1], "int")
		}
	default:
		err = bfail.Type(args[0], "array")
	}
	return
}

// \a -> ()
func arrReverse(args []*expression.T) (ex *expression.T, err error) {
	switch a := (args[0].Value).(type) {
	case []*expression.T:
		right := len(a) - 1
		for left := 0; left < right; left++ {
			a[left], a[right] = a[right], a[left]
			right--
		}
	default:
		err = bfail.Type(args[0], "array")
	}
	return
}

// \a, i, * -> ()
func arrSet(args []*expression.T) (ex *expression.T, err error) {
	switch a := (args[0].Value).(type) {
	case []*expression.T:
		switch ix := (args[1].Value).(type) {
		case int64:
			a[ix] = args[2]
		default:
			err = bfail.Type(args[1], "int")
		}
	default:
		err = bfail.Type(args[0], "array")
	}
	return
}

// \a -> *
func arrShift(args []*expression.T) (ex *expression.T, err error) {
	switch a := (args[0].Value).(type) {
	case []*expression.T:
		l := len(a)
		if l < 1 {
			err = bfail.Mk("Empty array")
			return
		}
		ex = a[0]
		args[0].Value = a[1:]
	default:
		err = bfail.Type(args[0], "array")
	}
	return
}

// \a -> ()
func arrShuffle(args []*expression.T) (ex *expression.T, err error) {
	switch a := (args[0].Value).(type) {
	case []*expression.T:
		for i := len(a); i > 1; {
			n := rand.Intn(i)
			i--
			a[n], a[i] = a[i], a[n]
		}
	default:
		err = bfail.Type(args[0], "array")
	}
	return
}

// \a -> i
func arrSize(args []*expression.T) (ex *expression.T, err error) {
	switch a := (args[0].Value).(type) {
	case []*expression.T:
		ex = expression.MkFinal(int64(len(a)))
	default:
		err = bfail.Type(args[0], "array")
	}
	return
}

// \a, \*, *->b -> ()
func arrSort(args []*expression.T) (ex *expression.T, err error) {
	switch a := (args[0].Value).(type) {
	case []*expression.T:
		switch fnLess := (args[1].Value).(type) {
		case *function.T:
			sort.Slice(a, func(i, j int) bool {
				r, er := solveIsolateFunction(fnLess, []*expression.T{a[i], a[j]}) // exSolver.go
				if er != nil {
					panic(er.Error())
				}
				switch v := (r.Value).(type) {
				case bool:
					return v
				default:
					panic(bfail.Type(r, "bool"))
				}
			})
		default:
			err = bfail.Type(args[1], "function")
		}
	default:
		err = bfail.Type(args[0], "array")
	}
	return
}

// \a, i -> a
func arrTake(args []*expression.T) (ex *expression.T, err error) {
	switch a := (args[0].Value).(type) {
	case []*expression.T:
		switch n := (args[1].Value).(type) {
		case int64:
			if n > int64(len(a)) {
				n = int64(len(a))
			}
			var rt []*expression.T
			for c := int64(0); c < n; c++ {
				rt = append(rt, a[c])
			}
			ex = expression.MkFinal(rt)
		default:
			err = bfail.Type(args[1], "int")
		}
	default:
		err = bfail.Type(args[0], "array")
	}
	return
}

// \a, \*->b -> a
func arrTakeWhile(args []*expression.T) (ex *expression.T, err error) {
	switch a := (args[0].Value).(type) {
	case []*expression.T:
		switch fn := (args[1].Value).(type) {
		case *function.T:
			var exs []*expression.T
			for _, e := range a {
				var r *expression.T
				r, err = solveIsolateFunction(fn, []*expression.T{e}) // exSolver.go
				if err != nil {
					return
				}
				stop := false
				switch v := (r.Value).(type) {
				case bool:
					if !v {
						stop = true
					} else {
						exs = append(exs, e)
					}
				default:
					err = bfail.Type(r, "bool")
					return
				}
				if stop {
					break
				}
			}
			ex = expression.MkFinal(exs)
		default:
			err = bfail.Type(args[1], "function")
		}
	default:
		err = bfail.Type(args[0], "array")
	}
	return
}

// \a -> <iterator>
func arrToIter(args []*expression.T) (ex *expression.T, err error) {
	switch a := (args[0].Value).(type) {
	case []*expression.T:
		l := len(a)
		ix := 0
		hasNext := func() bool {
			return ix < l
		}
		next := func() *expression.T {
			e := a[ix]
			ix++
			return e
		}
		ex = expression.MkFinal(iterator.New(hasNext, next))
	default:
		err = bfail.Type(args[0], "array")
	}
	return
}

// \a -> s
func arrToStr(args []*expression.T) (ex *expression.T, err error) {
	switch a := (args[0].Value).(type) {
	case []*expression.T:
		ex = expression.MkFinal(fmt.Sprint(
			expression.MkFinal(a)))
	default:
		err = bfail.Type(args[0], "array")
	}
	return
}

// \a, * -> ()
func arrUnshift(args []*expression.T) (ex *expression.T, err error) {
	switch a := (args[0].Value).(type) {
	case []*expression.T:
		args[0].Value = append([]*expression.T{args[1]}, a...)
	default:
		err = bfail.Type(args[0], "array")
	}
	return
}

func arrGet(fname string) (fn *bfunction.T, ok bool) {
	ok = true
	switch fname {
	case "all":
		fn = bfunction.New(2, arrAll)
	case "any":
		fn = bfunction.New(2, arrAny)
	case "clear":
		fn = bfunction.New(1, arrClear)
	case "copy":
		fn = bfunction.New(1, arrCopy)
	case "drop":
		fn = bfunction.New(2, arrDrop)
	case "dropWhile":
		fn = bfunction.New(2, arrDropWhile)
	case "duplicates":
		fn = bfunction.New(2, arrDuplicates)
	case "empty":
		fn = bfunction.New(1, arrEmpty)
	case "filter":
		fn = bfunction.New(2, arrFilter)
	case "filterIn":
		fn = bfunction.New(2, arrFilterIn)
	case "find":
		fn = bfunction.New(2, arrFind)
	case "fromIter":
		fn = bfunction.New(1, arrFromIter)
	case "join":
		fn = bfunction.New(2, arrJoin)
	case "mp":
		fn = bfunction.New(2, arrMp)
	case "new":
		fn = bfunction.New(2, arrNew)
	case "peek":
		fn = bfunction.New(1, arrPeek)
	case "pop":
		fn = bfunction.New(1, arrPop)
	case "push":
		fn = bfunction.New(2, arrPush)
	case "remove":
		fn = bfunction.New(2, arrRemove)
	case "removeRange":
		fn = bfunction.New(3, arrRemoveRange)
	case "reverse":
		fn = bfunction.New(1, arrReverse)
	case "shift":
		fn = bfunction.New(1, arrShift)
	case "shuffle":
		fn = bfunction.New(1, arrShuffle)
	case "size":
		fn = bfunction.New(1, arrSize)
	case "sort":
		fn = bfunction.New(2, arrSort)
	case "take":
		fn = bfunction.New(2, arrTake)
	case "takeWhile":
		fn = bfunction.New(2, arrTakeWhile)
	case "toIter":
		fn = bfunction.New(1, arrToIter)
	case "toStr":
		fn = bfunction.New(1, arrToStr)
	case "unshift":
		fn = bfunction.New(2, arrUnshift)
	default:
		ok = false
	}

	return
}

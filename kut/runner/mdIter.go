// Copyright 25-Feb-2022 ºDeme
// GNU General Public License - V3 <http://www.gnu.org/licenses/>

package runner

import (
	"fmt"
	"github.com/dedeme/kut/builtin/bfail"
	"github.com/dedeme/kut/builtin/bfunction"
	"github.com/dedeme/kut/expression"
	"github.com/dedeme/kut/function"
	"github.com/dedeme/kut/iterator"
)

// \<iterator>, \*->b -> b
func iterAll(args []*expression.T) (ex *expression.T, err error) {
	switch it := (args[0].Value).(type) {
	case *iterator.T:
		switch fn := (args[1].Value).(type) {
		case *function.T:
			for it.HasNext() {
				r, er := solveIsolateFunction(fn, []*expression.T{it.Next()})
				if er != nil {
					panic(er.Error())
				}
				switch v := (r.Value).(type) {
				case bool:
					if !v {
						ex = r
						return
					}
				default:
					panic(bfail.Type(r, "bool"))
				}
			}
			ex = expression.MkFinal(true)
		default:
			err = bfail.Type(args[1], "function")
		}
	default:
		err = bfail.Type(args[0], "<iterator>")
	}
	return
}

// \<iterator>, \*->b -> b
func iterAny(args []*expression.T) (ex *expression.T, err error) {
	switch it := (args[0].Value).(type) {
	case *iterator.T:
		switch fn := (args[1].Value).(type) {
		case *function.T:
			for it.HasNext() {
				r, er := solveIsolateFunction(fn, []*expression.T{it.Next()})
				if er != nil {
					panic(er.Error())
				}
				switch v := (r.Value).(type) {
				case bool:
					if v {
						ex = r
						return
					}
				default:
					panic(bfail.Type(r, "bool"))
				}
			}
			ex = expression.MkFinal(false)
		default:
			err = bfail.Type(args[1], "function")
		}
	default:
		err = bfail.Type(args[0], "<iterator>")
	}
	return
}

// \<iterator>, <iterator> -> <iterator>
func iterCat(args []*expression.T) (ex *expression.T, err error) {
	switch it1 := (args[0].Value).(type) {
	case *iterator.T:
		switch it2 := (args[1].Value).(type) {
		case *iterator.T:
			hasNext := func() bool {
				return it1.HasNext() || it2.HasNext()
			}
			next := func() *expression.T {
				if it1.HasNext() {
					return it1.Next()
				} else {
					return it2.Next()
				}
			}
			ex = expression.MkFinal(iterator.New(hasNext, next))
		default:
			err = bfail.Type(args[1], "<iterator>")
		}
	default:
		err = bfail.Type(args[0], "<iterator>")
	}
	return
}

// \<iterator> -> i
func iterCount(args []*expression.T) (ex *expression.T, err error) {
	switch it := (args[0].Value).(type) {
	case *iterator.T:
		n := int64(0)
		for it.HasNext() {
			it.Next()
			n++
		}
		ex = expression.MkFinal(n)
	default:
		err = bfail.Type(args[0], "<iterator>")
	}
	return
}

// \<iterator>, i -> <iterator>
func iterDrop(args []*expression.T) (ex *expression.T, err error) {
	switch it := (args[0].Value).(type) {
	case *iterator.T:
		switch end := (args[1].Value).(type) {
		case int64:
			ix := int64(0)
			for it.HasNext() && ix < end {
				it.Next()
				ix++
			}
			ex = args[0]
		default:
			err = bfail.Type(args[1], "int")
		}
	default:
		err = bfail.Type(args[0], "<iterator>")
	}
	return
}

// \<iterator>, \*->b -> <iterator>
func iterDropWhile(args []*expression.T) (ex *expression.T, err error) {
	switch it := (args[0].Value).(type) {
	case *iterator.T:
		switch fn := (args[1].Value).(type) {
		case *function.T:
			var nx *expression.T
			for it.HasNext() {
				nx = it.Next()
				r, er := solveIsolateFunction(fn, []*expression.T{nx})
				if er != nil {
					err = er
					return
				}
				switch v := (r.Value).(type) {
				case bool:
					if v {
						nx = nil
					}
				default:
					panic(bfail.Type(r, "bool"))
				}
				if nx != nil {
					break
				}
			}
			hasNext := func() bool {
				return nx != nil
			}
			next := func() *expression.T {
				rt := nx
				nx = nil
				if it.HasNext() {
					nx = it.Next()
				}
				return rt
			}
			ex = expression.MkFinal(iterator.New(hasNext, next))
		default:
			err = bfail.Type(args[1], "function")
		}
	default:
		err = bfail.Type(args[0], "<iterator>")
	}
	return
}

// \<iterator>, \*->() -> ()
func iterEach(args []*expression.T) (ex *expression.T, err error) {
	switch it := (args[0].Value).(type) {
	case *iterator.T:
		switch fn := (args[1].Value).(type) {
		case *function.T:
			for it.HasNext() {
				_, err = solveIsolateFunction(fn, []*expression.T{it.Next()})
				if err != nil {
					break
				}
			}
			if err == nil {
				ex = expression.MkEmpty()
			}
		default:
			err = bfail.Type(args[1], "function")
		}
	default:
		err = bfail.Type(args[0], "<iterator>")
	}
	return
}

// \-> <iterator>
func iterEmpty(args []*expression.T) (ex *expression.T, err error) {
	hasNext := func() bool {
		return false
	}
	next := func() *expression.T {
		return expression.MkEmpty()
	}
	ex = expression.MkFinal(iterator.New(hasNext, next))
	return
}

// \<iterator>, \*->b -> <iterator>
func iterFilter(args []*expression.T) (ex *expression.T, err error) {
	switch it := (args[0].Value).(type) {
	case *iterator.T:
		switch fn := (args[1].Value).(type) {
		case *function.T:
			var nx *expression.T
			for it.HasNext() {
				nx = it.Next()
				r, er := solveIsolateFunction(fn, []*expression.T{nx})
				if er != nil {
					err = er
					return
				}
				switch v := (r.Value).(type) {
				case bool:
					if !v {
						nx = nil
					}
				default:
					panic(bfail.Type(r, "bool"))
				}
				if nx != nil {
					break
				}
			}

			hasNext := func() bool {
				return nx != nil
			}
			next := func() *expression.T {
				rt := nx
				nx = nil
				for it.HasNext() {
					nx = it.Next()
					r, er := solveIsolateFunction(fn, []*expression.T{nx})
					if er != nil {
						panic(er.Error())
					}
					switch v := (r.Value).(type) {
					case bool:
						if !v {
							nx = nil
						}
					default:
						panic(bfail.Type(r, "bool"))
					}
					if nx != nil {
						break
					}
				}
				return rt
			}
			ex = expression.MkFinal(iterator.New(hasNext, next))
		default:
			err = bfail.Type(args[1], "function")
		}
	default:
		err = bfail.Type(args[0], "<iterator>")
	}
	return
}

// \<iterator>, \*->b -> [*]
func iterFind(args []*expression.T) (ex *expression.T, err error) {
	switch it := (args[0].Value).(type) {
	case *iterator.T:
		switch fn := (args[1].Value).(type) {
		case *function.T:
			var rt *expression.T
			for it.HasNext() {
				nx := it.Next()
				r, er := solveIsolateFunction(fn, []*expression.T{nx})
				if er != nil {
					err = er
					return
				}
				switch v := (r.Value).(type) {
				case bool:
					if v {
						rt = nx
					}
				default:
					panic(bfail.Type(r, "bool"))
				}
				if rt != nil {
					break
				}
			}
			if rt == nil {
				ex = expression.MkFinal([]*expression.T{})
			} else {
				ex = expression.MkFinal([]*expression.T{rt})
			}
		default:
			err = bfail.Type(args[1], "function")
		}
	default:
		err = bfail.Type(args[0], "<iterator>")
	}
	return
}

// \<iterator>, \*->b -> i
func iterIndex(args []*expression.T) (ex *expression.T, err error) {
	switch it := (args[0].Value).(type) {
	case *iterator.T:
		switch fn := (args[1].Value).(type) {
		case *function.T:
			n := int64(-1)
			c := int64(0)
			for it.HasNext() {
				r, er := solveIsolateFunction(fn, []*expression.T{it.Next()})
				if er != nil {
					err = er
					return
				}
				stop := false
				switch v := (r.Value).(type) {
				case bool:
					if v {
						stop = true
					}
				default:
					panic(bfail.Type(r, "bool"))
				}
				if stop {
					n = c
					break
				}
				c++
			}
			ex = expression.MkFinal(n)
		default:
			err = bfail.Type(args[1], "function")
		}
	default:
		err = bfail.Type(args[0], "<iterator>")
	}
	return
}

// \<iterator> -> b
func iterHasNext(args []*expression.T) (ex *expression.T, err error) {
	switch it := (args[0].Value).(type) {
	case *iterator.T:
		ex = expression.MkFinal(it.HasNext())
	default:
		err = bfail.Type(args[0], "<iterator>")
	}
	return
}

// \<iterator>, \*->* -> <iterator>
func iterMp(args []*expression.T) (ex *expression.T, err error) {
	switch it := (args[0].Value).(type) {
	case *iterator.T:
		switch fn := (args[1].Value).(type) {
		case *function.T:
			next := func() *expression.T {
				r, er := solveIsolateFunction(fn, []*expression.T{it.Next()})
				if er == nil {
					return r
				}
				panic(er.Error())
			}
			ex = expression.MkFinal(iterator.New(it.HasNext, next))
		default:
			err = bfail.Type(args[1], "function")
		}
	default:
		err = bfail.Type(args[0], "<iterator>")
	}
	return
}

// \<iterator> -> *
func iterNext(args []*expression.T) (ex *expression.T, err error) {
	switch it := (args[0].Value).(type) {
	case *iterator.T:
		ex = it.Next()
	default:
		err = bfail.Type(args[0], "<iterator>")
	}
	return
}

// \\->b, \->* -> <iterator>
func iterNew(args []*expression.T) (ex *expression.T, err error) {
	switch fhasNext := (args[0].Value).(type) {
	case *function.T:
		if len(fhasNext.Vars) != 0 {
			err = bfail.Mk(fmt.Sprintf(
				"'hasNext' function of 'iter.new' must have 0 arguments, but it has %v",
				len(fhasNext.Vars),
			))
		} else {
			switch fnext := (args[1].Value).(type) {
			case *function.T:
				if len(fnext.Vars) != 0 {
					err = bfail.Mk(fmt.Sprintf(
						"'next' function of 'iter.new' must have 0 arguments, but it has %v",
						len(fnext.Vars),
					))
				} else {
					fhasNext2 := func() bool {
						var rex *expression.T
						rex, err = solveIsolateFunction(fhasNext, []*expression.T{})
						if err == nil {
							switch r := (rex.Value).(type) {
							case bool:
								return r
							default:
								panic(bfail.Type(rex, "bool"))
							}
						}
						panic(err.Error())
					}
					fnext2 := func() *expression.T {
						var rex *expression.T
						rex, err = solveIsolateFunction(fnext, []*expression.T{})
						if err == nil {
							return rex
						}
						panic(err.Error())
					}

					ex = expression.MkFinal(iterator.New(fhasNext2, fnext2))
				}
			default:
				err = bfail.Type(args[1], "function")
			}
		}
	default:
		err = bfail.Type(args[0], "function")
	}
	return
}

// \<iterator>, *1, \*1, *2->*1 -> *1
func iterReduce(args []*expression.T) (ex *expression.T, err error) {
	switch it := (args[0].Value).(type) {
	case *iterator.T:
		seed := args[1]
		switch fn := (args[2].Value).(type) {
		case *function.T:
			for it.HasNext() {
				var r *expression.T
				r, err = solveIsolateFunction(fn, []*expression.T{seed, it.Next()})
				if err != nil {
					return
				}
				seed = r
			}
			ex = seed
		default:
			err = bfail.Type(args[2], "function")
		}
	default:
		err = bfail.Type(args[0], "<iterator>")
	}
	return
}

// \<iterator>, i -> <iterator>
func iterTake(args []*expression.T) (ex *expression.T, err error) {
	switch it := (args[0].Value).(type) {
	case *iterator.T:
		switch end := (args[1].Value).(type) {
		case int64:
			ix := int64(0)
			hasNext := func() bool {
				return it.HasNext() && ix < end
			}
			next := func() *expression.T {
				ix++
				return it.Next()
			}
			ex = expression.MkFinal(iterator.New(hasNext, next))
		default:
			err = bfail.Type(args[1], "int")
		}
	default:
		err = bfail.Type(args[0], "<iterator>")
	}
	return
}

// \<iterator>, \*->b -> <iterator>
func iterTakeWhile(args []*expression.T) (ex *expression.T, err error) {
	switch it := (args[0].Value).(type) {
	case *iterator.T:
		switch fn := (args[1].Value).(type) {
		case *function.T:
			more := false
			var nx *expression.T
			if it.HasNext() {
				nx = it.Next()
				r, er := solveIsolateFunction(fn, []*expression.T{nx})
				if er != nil {
					err = er
					return
				}
				switch v := (r.Value).(type) {
				case bool:
					more = v
				default:
					err = bfail.Type(r, "bool")
					return
				}
			}

			hasNext := func() bool {
				return more
			}
			next := func() *expression.T {
				rt := nx
				more = false
				if it.HasNext() {
					nx = it.Next()
					r, er := solveIsolateFunction(fn, []*expression.T{nx})
					if er != nil {
						panic(er.Error())
					}
					switch v := (r.Value).(type) {
					case bool:
						more = v
					default:
						panic(bfail.Type(r, "bool"))
					}
				}
				return rt
			}
			ex = expression.MkFinal(iterator.New(hasNext, next))
		default:
			err = bfail.Type(args[1], "function")
		}
	default:
		err = bfail.Type(args[0], "<iterator>")
	}
	return
}

// \<iterator> -> s
func iterToStr(args []*expression.T) (ex *expression.T, err error) {
	switch it := (args[0].Value).(type) {
	case *iterator.T:
		var a []*expression.T
		for it.HasNext() {
			a = append(a, it.Next())
		}
		ex = expression.MkFinal("<iterator>" + fmt.Sprint(
			expression.MkFinal(a)))
	default:
		err = bfail.Type(args[0], "<iterator>")
	}
	return
}

// \* -> <iterator>
func iterUnary(args []*expression.T) (ex *expression.T, err error) {
	has := true
	hasNext := func() bool {
		return has
	}
	next := func() *expression.T {
		has = false
		return args[0]
	}
	ex = expression.MkFinal(iterator.New(hasNext, next))
	return
}

func iterGet(fname string) (fn *bfunction.T, ok bool) {
	ok = true
	switch fname {
	case "all":
		fn = bfunction.New(2, iterAll)
	case "any":
		fn = bfunction.New(2, iterAny)
	case "cat":
		fn = bfunction.New(2, iterCat)
	case "count":
		fn = bfunction.New(1, iterCount)
	case "drop":
		fn = bfunction.New(2, iterDrop)
	case "dropWhile":
		fn = bfunction.New(2, iterDropWhile)
	case "each":
		fn = bfunction.New(2, iterEach)
	case "empty":
		fn = bfunction.New(0, iterEmpty)
	case "filter":
		fn = bfunction.New(2, iterFilter)
	case "find":
		fn = bfunction.New(2, iterFind)
	case "hasNext":
		fn = bfunction.New(1, iterHasNext)
	case "index":
		fn = bfunction.New(2, iterIndex)
	case "mp":
		fn = bfunction.New(2, iterMp)
	case "next":
		fn = bfunction.New(1, iterNext)
	case "new":
		fn = bfunction.New(2, iterNew)
	case "reduce":
		fn = bfunction.New(3, iterReduce)
	case "take":
		fn = bfunction.New(2, iterTake)
	case "takeWhile":
		fn = bfunction.New(2, iterTakeWhile)
	case "toStr":
		fn = bfunction.New(1, iterToStr)
	case "unary":
		fn = bfunction.New(1, iterUnary)
	default:
		ok = false
	}

	return
}

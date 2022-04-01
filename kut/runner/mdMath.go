// Copyright 25-Feb-2022 ºDeme
// GNU General Public License - V3 <http://www.gnu.org/licenses/>

package runner

import (
	"github.com/dedeme/kut/builtin/bfail"
	"github.com/dedeme/kut/builtin/bfunction"
	"github.com/dedeme/kut/expression"
	"github.com/dedeme/kut/reader/token"
	"math"
	"math/rand"
	"strconv"
	"strings"
)

/// Round n with 'dec' decimals.
/// dec is set between 0 and 9, both inclusive.
func mathFRound(n float64, dec int64) (r float64) {
	switch {
	case dec <= 0:
		r = math.Round(n)
	case dec == 1:
		r = math.Round(n*10.0) / 10.0
	case dec == 2:
		r = math.Round(n*100.0) / 100.0
	case dec == 3:
		r = math.Round(n*1000.0) / 1000.0
	case dec == 4:
		r = math.Round(n*10000.0) / 10000.0
	case dec == 5:
		r = math.Round(n*100000.0) / 100000.0
	case dec == 6:
		r = math.Round(n*1000000.0) / 1000000.0
	case dec == 7:
		r = math.Round(n*10000000.0) / 10000000.0
	case dec == 8:
		r = math.Round(n*100000000.0) / 100000000.0
	default:
		r = math.Round(n*1000000000.0) / 1000000000.0
	}
	return
}

// \f -> f
func mathAbs(exs []*expression.T) (ex *expression.T, err error) {
	switch n := (exs[0].Value).(type) {
	case float64:
		ex = expression.MkFinal(math.Abs(n))
	default:
		err = bfail.Type(exs[0], "float")
	}
	return
}

// \f -> f
func mathAcos(exs []*expression.T) (ex *expression.T, err error) {
	switch n := (exs[0].Value).(type) {
	case float64:
		ex = expression.MkFinal(math.Acos(n))
	default:
		err = bfail.Type(exs[0], "float")
	}
	return
}

// \f -> f
func mathAcosh(exs []*expression.T) (ex *expression.T, err error) {
	switch n := (exs[0].Value).(type) {
	case float64:
		ex = expression.MkFinal(math.Acosh(n))
	default:
		err = bfail.Type(exs[0], "float")
	}
	return
}

// \f -> f
func mathAsin(exs []*expression.T) (ex *expression.T, err error) {
	switch n := (exs[0].Value).(type) {
	case float64:
		ex = expression.MkFinal(math.Asin(n))
	default:
		err = bfail.Type(exs[0], "float")
	}
	return
}

// \f -> f
func mathAsinh(exs []*expression.T) (ex *expression.T, err error) {
	switch n := (exs[0].Value).(type) {
	case float64:
		ex = expression.MkFinal(math.Asinh(n))
	default:
		err = bfail.Type(exs[0], "float")
	}
	return
}

// \f -> f
func mathAtan(exs []*expression.T) (ex *expression.T, err error) {
	switch n := (exs[0].Value).(type) {
	case float64:
		ex = expression.MkFinal(math.Atan(n))
	default:
		err = bfail.Type(exs[0], "float")
	}
	return
}

// \f -> f
func mathAtanh(exs []*expression.T) (ex *expression.T, err error) {
	switch n := (exs[0].Value).(type) {
	case float64:
		ex = expression.MkFinal(math.Atanh(n))
	default:
		err = bfail.Type(exs[0], "float")
	}
	return
}

// \f -> f
func mathCeil(exs []*expression.T) (ex *expression.T, err error) {
	switch n := (exs[0].Value).(type) {
	case float64:
		ex = expression.MkFinal(math.Ceil(n))
	default:
		err = bfail.Type(exs[0], "float")
	}
	return
}

// \f -> f
func mathCos(exs []*expression.T) (ex *expression.T, err error) {
	switch n := (exs[0].Value).(type) {
	case float64:
		ex = expression.MkFinal(math.Cos(n))
	default:
		err = bfail.Type(exs[0], "float")
	}
	return
}

// \f -> f
func mathCosh(exs []*expression.T) (ex *expression.T, err error) {
	switch n := (exs[0].Value).(type) {
	case float64:
		ex = expression.MkFinal(math.Cosh(n))
	default:
		err = bfail.Type(exs[0], "float")
	}
	return
}

// \f, f, f -> b
func mathEq(exs []*expression.T) (ex *expression.T, err error) {
	switch n1 := (exs[0].Value).(type) {
	case float64:
		switch n2 := (exs[1].Value).(type) {
		case float64:
			switch df := (exs[2].Value).(type) {
			case float64:
				ex = expression.MkFinal(math.Abs(n1-n2) <= df)
			default:
				err = bfail.Type(exs[2], "float")
			}
		default:
			err = bfail.Type(exs[1], "float")
		}
	default:
		err = bfail.Type(exs[0], "float")
	}
	return
}

// \f -> f
func mathExp(exs []*expression.T) (ex *expression.T, err error) {
	switch n := (exs[0].Value).(type) {
	case float64:
		ex = expression.MkFinal(math.Exp(n))
	default:
		err = bfail.Type(exs[0], "float")
	}
	return
}

// \f -> f
func mathExp2(exs []*expression.T) (ex *expression.T, err error) {
	switch n := (exs[0].Value).(type) {
	case float64:
		ex = expression.MkFinal(math.Exp2(n))
	default:
		err = bfail.Type(exs[0], "float")
	}
	return
}

// \f -> f
func mathFloor(exs []*expression.T) (ex *expression.T, err error) {
	switch n := (exs[0].Value).(type) {
	case float64:
		ex = expression.MkFinal(math.Floor(n))
	default:
		err = bfail.Type(exs[0], "float")
	}
	return
}

// \s -> [f] | []
func mathFromEn(exs []*expression.T) (ex *expression.T, err error) {
	switch s := (exs[0].Value).(type) {
	case string:
		s = strings.ReplaceAll(s, ",", "")
		f, er := strconv.ParseFloat(s, 64)
		if er == nil {
			ex = expression.MkFinal([]*expression.T{expression.MkFinal(f)})
		} else {
			ex = expression.MkFinal([]*expression.T{})
		}
	default:
		err = bfail.Type(exs[0], "string")
	}
	return
}

// \s -> [f] | []
func mathFromIso(exs []*expression.T) (ex *expression.T, err error) {
	switch s := (exs[0].Value).(type) {
	case string:
		s = strings.ReplaceAll(strings.ReplaceAll(s, ".", ""), ",", ".")
		f, er := strconv.ParseFloat(s, 64)
		if er == nil {
			ex = expression.MkFinal([]*expression.T{expression.MkFinal(f)})
		} else {
			ex = expression.MkFinal([]*expression.T{})
		}
	default:
		err = bfail.Type(exs[0], "string")
	}
	return
}

// \s -> [f] | []
func mathFromStr(exs []*expression.T) (ex *expression.T, err error) {
	switch s := (exs[0].Value).(type) {
	case string:
		f, er := strconv.ParseFloat(s, 64)
		if er == nil {
			ex = expression.MkFinal([]*expression.T{expression.MkFinal(f)})
		} else {
			ex = expression.MkFinal([]*expression.T{})
		}
	default:
		err = bfail.Type(exs[0], "string")
	}
	return
}

// \f -> f
func mathLog(exs []*expression.T) (ex *expression.T, err error) {
	switch n := (exs[0].Value).(type) {
	case float64:
		ex = expression.MkFinal(math.Log(n))
	default:
		err = bfail.Type(exs[0], "float")
	}
	return
}

// \f -> f
func mathLog10(exs []*expression.T) (ex *expression.T, err error) {
	switch n := (exs[0].Value).(type) {
	case float64:
		ex = expression.MkFinal(math.Log10(n))
	default:
		err = bfail.Type(exs[0], "float")
	}
	return
}

// \f -> f
func mathLog2(exs []*expression.T) (ex *expression.T, err error) {
	switch n := (exs[0].Value).(type) {
	case float64:
		ex = expression.MkFinal(math.Log2(n))
	default:
		err = bfail.Type(exs[0], "float")
	}
	return
}

// \f, f -> f
func mathMax(exs []*expression.T) (ex *expression.T, err error) {
	switch n1 := (exs[0].Value).(type) {
	case float64:
		switch n2 := (exs[1].Value).(type) {
		case float64:
			ex = expression.MkFinal(math.Max(n1, n2))
		default:
			err = bfail.Type(exs[1], "float")
		}
	default:
		err = bfail.Type(exs[0], "float")
	}
	return
}

// \f, f -> f
func mathMin(exs []*expression.T) (ex *expression.T, err error) {
	switch n1 := (exs[0].Value).(type) {
	case float64:
		switch n2 := (exs[1].Value).(type) {
		case float64:
			ex = expression.MkFinal(math.Min(n1, n2))
		default:
			err = bfail.Type(exs[1], "float")
		}
	default:
		err = bfail.Type(exs[0], "float")
	}
	return
}

// \f, f -> f
func mathPow(exs []*expression.T) (ex *expression.T, err error) {
	switch base := (exs[0].Value).(type) {
	case float64:
		switch exp := (exs[1].Value).(type) {
		case float64:
			ex = expression.MkFinal(math.Pow(base, exp))
		default:
			err = bfail.Type(exs[1], "float")
		}
	default:
		err = bfail.Type(exs[0], "float")
	}
	return
}

// \f -> f
func mathPow10(exs []*expression.T) (ex *expression.T, err error) {
	switch n := (exs[0].Value).(type) {
	case int64:
		ex = expression.MkFinal(math.Pow10(int(n)))
	default:
		err = bfail.Type(exs[0], "float")
	}
	return
}

// \-> f
func mathRnd(exs []*expression.T) (ex *expression.T, err error) {
	ex = expression.MkFinal(rand.Float64())
	return
}

// \i -> f
func mathRndi(exs []*expression.T) (ex *expression.T, err error) {
	switch n := (exs[0].Value).(type) {
	case int64:
		ex = expression.MkFinal(rand.Int63n(n))
	default:
		err = bfail.Type(exs[0], "int")
	}
	return
}

// \f, i -> f
func mathRound(exs []*expression.T) (ex *expression.T, err error) {
	switch n := (exs[0].Value).(type) {
	case float64:
		switch decimals := (exs[1].Value).(type) {
		case int64:
			ex = expression.MkFinal(mathFRound(n, decimals))
		default:
			err = bfail.Type(exs[1], "float")
		}
	default:
		err = bfail.Type(exs[0], "float")
	}
	return
}

// \f -> f
func mathSin(exs []*expression.T) (ex *expression.T, err error) {
	switch n := (exs[0].Value).(type) {
	case float64:
		ex = expression.MkFinal(math.Sin(n))
	default:
		err = bfail.Type(exs[0], "float")
	}
	return
}

// \f -> f
func mathSinh(exs []*expression.T) (ex *expression.T, err error) {
	switch n := (exs[0].Value).(type) {
	case float64:
		ex = expression.MkFinal(math.Sinh(n))
	default:
		err = bfail.Type(exs[0], "float")
	}
	return
}

// \f -> f
func mathSqrt(exs []*expression.T) (ex *expression.T, err error) {
	switch n := (exs[0].Value).(type) {
	case float64:
		ex = expression.MkFinal(math.Sqrt(n))
	default:
		err = bfail.Type(exs[0], "float")
	}
	return
}

// \f -> f
func mathTan(exs []*expression.T) (ex *expression.T, err error) {
	switch n := (exs[0].Value).(type) {
	case float64:
		ex = expression.MkFinal(math.Tan(n))
	default:
		err = bfail.Type(exs[0], "float")
	}
	return
}

// \f -> f
func mathTanh(exs []*expression.T) (ex *expression.T, err error) {
	switch n := (exs[0].Value).(type) {
	case float64:
		ex = expression.MkFinal(math.Tanh(n))
	default:
		err = bfail.Type(exs[0], "float")
	}
	return
}

// \i | f | s -> i | f | s
func mathTo(tpTo token.Type) func([]*expression.T) (*expression.T, error) {
	return func(exs []*expression.T) (ex *expression.T, err error) {
		e := exs[0]
		switch v := (e.Value).(type) {
		case int64:
			switch tpTo {
			case token.Int:
				ex = e
			case token.Float:
				ex = expression.MkFinal(float64(v))
			case token.String:
				ex = expression.MkFinal(strconv.FormatInt(v, 10))
			}
		case float64:
			switch tpTo {
			case token.Int:
				ex = expression.MkFinal(int64(v))
			case token.Float:
				ex = expression.MkFinal(v)
			case token.String:
				ex = expression.MkFinal(strconv.FormatFloat(v, 'g', -1, 64))
			}
		case string:
			var er error
			switch tpTo {
			case token.Int:
				var i int64
				i, er = strconv.ParseInt(v, 10, 64)
				if er == nil {
					ex = expression.MkFinal(i)
				}
			case token.Float:
				var f float64
				f, er = strconv.ParseFloat(v, 64)
				if er == nil {
					ex = expression.MkFinal(f)
				}
			case token.String:
				_, err = strconv.ParseFloat(v, 64)
				if err == nil {
					ex = e
				}
			}
			if er != nil {
				err = bfail.Mk(er.Error())
			}
		default:
			err = bfail.Type(e, "int", "float", "string")
		}
		return
	}
}

// \f -> f
func mathTrunc(exs []*expression.T) (ex *expression.T, err error) {
	switch n := (exs[0].Value).(type) {
	case float64:
		ex = expression.MkFinal(math.Trunc(n))
	default:
		err = bfail.Type(exs[0], "float")
	}
	return
}

func mathGet(fname string) (fn *bfunction.T, ok bool) {
	ok = true
	switch fname {
	case "abs":
		fn = bfunction.New(1, mathAbs)
	case "acos":
		fn = bfunction.New(1, mathAcos)
	case "acosh":
		fn = bfunction.New(1, mathAcosh)
	case "asin":
		fn = bfunction.New(1, mathAsin)
	case "asinh":
		fn = bfunction.New(1, mathAsinh)
	case "atan":
		fn = bfunction.New(1, mathAtan)
	case "atanh":
		fn = bfunction.New(1, mathAtanh)
	case "ceil":
		fn = bfunction.New(1, mathCeil)
	case "cos":
		fn = bfunction.New(1, mathCos)
	case "cosh":
		fn = bfunction.New(1, mathCosh)
	case "eq":
		fn = bfunction.New(3, mathEq)
	case "exp":
		fn = bfunction.New(1, mathExp)
	case "exp2":
		fn = bfunction.New(1, mathExp2)
	case "floor":
		fn = bfunction.New(1, mathFloor)
	case "fromEn":
		fn = bfunction.New(1, mathFromEn)
	case "fromIso":
		fn = bfunction.New(1, mathFromIso)
	case "fromStr":
		fn = bfunction.New(1, mathFromStr)
	case "log":
		fn = bfunction.New(1, mathLog)
	case "log10":
		fn = bfunction.New(1, mathLog10)
	case "log2":
		fn = bfunction.New(1, mathLog2)
	case "max":
		fn = bfunction.New(2, mathMax)
	case "min":
		fn = bfunction.New(2, mathMin)
	case "pow":
		fn = bfunction.New(2, mathPow)
	case "pow10":
		fn = bfunction.New(1, mathPow10)
	case "rnd":
		fn = bfunction.New(0, mathRnd)
	case "rndi":
		fn = bfunction.New(1, mathRndi)
	case "round":
		fn = bfunction.New(2, mathRound)
	case "sin":
		fn = bfunction.New(1, mathSin)
	case "sinh":
		fn = bfunction.New(1, mathSinh)
	case "sqrt":
		fn = bfunction.New(1, mathSqrt)
	case "tan":
		fn = bfunction.New(1, mathTan)
	case "tanh":
		fn = bfunction.New(1, mathTanh)
	case "toFloat":
		fn = bfunction.New(1, mathTo(token.Float))
	case "toInt":
		fn = bfunction.New(1, mathTo(token.Int))
	case "toStr":
		fn = bfunction.New(1, mathTo(token.String))
	case "trunc":
		fn = bfunction.New(1, mathTrunc)
	default:
		ok = false
	}

	return
}

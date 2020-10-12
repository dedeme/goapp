// Copyright 02-Aug-2020 ºDeme
// GNU General Public License - V3 <http://www.gnu.org/licenses/>

// Integer procedures.
package intpk

import (
	"github.com/dedeme/dmstack/machine"
	"github.com/dedeme/dmstack/symbol"
	"github.com/dedeme/dmstack/token"
	"math/rand"
	"strconv"
)

// Auxiliar function
func popInt(m *machine.T) (i int64) {
	tk := m.PopT(token.Int)
	i, _ = tk.I()
	return
}

// Auxiliar function
func pushInt(m *machine.T, i int64) {
	m.Push(token.NewI(i, m.MkPos()))
}

// Returns an integer from a string. If it fails, throws an "Int error".
//    m  : Virtual machine.
func prFromStr(m *machine.T) {
	tk := m.PopT(token.String)
	s, _ := tk.S()
	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		m.Fail("Int error", "'%v' bad integer.", tk.StringDraft())
	}
	pushInt(m, i)
}

// Returns the absolute value of i.
//    m  : Virtual machine.
func prAbs(m *machine.T) {
	i := popInt(m)
	if i < 0 {
		i = -i
	}
	pushInt(m, i)
}

// Returns a random number from 0 (inclusive) to i (exclusive). If i <= 0
// throws an "Int error".
//    m  : Virtual machine.
func prRnd(m *machine.T) {
	tk := m.PopT(token.Int)
	i, _ := tk.I()
	if i <= 0 {
		m.Fail("Int error", "'%v' is <= 0.", tk.StringDraft())
	}
	pushInt(m, rand.Int63n(i))
}

// Returns division and module of num / den. It raises a "Math error" if
// den == 0.
//    m  : Virtual machine.
func prDiv(m *machine.T) {
	den := popInt(m)
	if den == 0 {
		m.Fail("Math error", "Division by 0")
	}
	num := popInt(m)
	pushInt(m, num/den)
	pushInt(m, num%den)
}

// Returns arithmetic 'and'.
//    m  : Virtual machine.
func prAnd(m *machine.T) {
	n2 := popInt(m)
	n1 := popInt(m)
	pushInt(m, n1&n2)
}

// Returns arithmetic 'or'
//    m  : Virtual machine.
func prOr(m *machine.T) {
	n2 := popInt(m)
	n1 := popInt(m)
	pushInt(m, n1|n2)
}

// Returns arithmetic 'xor'
//    m  : Virtual machine.
func prXor(m *machine.T) {
	n2 := popInt(m)
	n1 := popInt(m)
	pushInt(m, n1^n2)
}

// Returns arithmetic 'not'.
//    m  : Virtual machine.
func prNot(m *machine.T) {
	n := popInt(m)
	pushInt(m, ^n)
}

// Returns shift left of n1, n2 times.
//    m  : Virtual machine.
func prLeft(m *machine.T) {
	n2 := popInt(m)
	n1 := popInt(m)
	pushInt(m, n1<<n2)
}

// Returns shift right of n1, n2 times.
//    m  : Virtual machine.
func prRight(m *machine.T) {
	n2 := popInt(m)
	n1 := popInt(m)
	pushInt(m, n1>>n2)
}

// Returns the maximum integer.
//    m  : Virtual machine.
func prMax(m *machine.T) {
	pushInt(m, 9223372036854775807)
}

// Returns the minimum integer.
//    m  : Virtual machine.
func prMin(m *machine.T) {
	pushInt(m, -9223372036854775808)
}

// Returns n converted to float.
//    m  : Virtual machine.
func prToFloat(m *machine.T) {
	n := popInt(m)
	m.Push(token.NewF(float64(n), m.MkPos()))
}

// Auxiliar function.
func toIsoEn(m *machine.T, sep string) {
	n := popInt(m)
	neg := false
	if n < 0 {
		neg = true
		n = -n
	}
	s := strconv.FormatInt(n, 10)
	ix := 3
	for {
		if len(s) <= ix {
			break
		}
		pos := len(s) - ix
		s = s[:pos] + sep + s[pos:]
		ix += 4
	}
	if neg {
		s = "-" + s
	}
	m.Push(token.NewS(s, m.MkPos()))
}

// Returns an integer converted to string in Iso format.
//    m  : Virtual machine.
func prToIso(m *machine.T) {
	toIsoEn(m, ".")
}

// Returns an integer converted to string in En format.
//    m  : Virtual machine.
func prToEn(m *machine.T) {
	toIsoEn(m, ",")
}

// Processes string procedures.
//    m  : Virtual machine.
func Proc(m *machine.T) {
	tk, ok := m.PrgNext()
	if !ok {
		m.Failt("'int' procedure is missing")
	}
	sy, ok := tk.Sy()
	if !ok {
		m.Failt(
			"\n  Expected: 'int' procedure.\n  Actual  : '%v'.", tk.StringDraft(),
		)
	}
	switch sy {
	case symbol.New("fromStr"):
		prFromStr(m)
	case symbol.New("abs"):
		prAbs(m)
	case symbol.New("rnd"):
		prRnd(m)
	case symbol.New("div"):
		prDiv(m)

	case symbol.New("and"):
		prAnd(m)
	case symbol.New("or"):
		prOr(m)
	case symbol.New("xor"):
		prXor(m)
	case symbol.New("not"):
		prNot(m)
	case symbol.New("<<"):
		prLeft(m)
	case symbol.New(">>"):
		prRight(m)

	case symbol.New("max"):
		prMax(m)
	case symbol.New("min"):
		prMin(m)
	case symbol.New("toFloat"):
		prToFloat(m)
	case symbol.New("toIso"):
		prToIso(m)
	case symbol.New("toEn"):
		prToEn(m)

	default:
		m.Failt("'int' does not contains the procedure '%v'.", sy.String())
	}
}

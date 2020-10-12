// Copyright 07-Aug-2020 ºDeme
// GNU General Public License - V3 <http://www.gnu.org/licenses/>

// Float-numeric functions.
package mathpk

import (
	"github.com/dedeme/dmstack/machine"
	"github.com/dedeme/dmstack/symbol"
	"github.com/dedeme/dmstack/token"
	"math"
)

// Auxiliar function
func popFloat(m *machine.T) (f float64) {
	tk := m.PopT(token.Float)
	f, _ = tk.F()
	return
}

// Auxiliar function
func pushFloat(m *machine.T, f float64) {
	m.Push(token.NewF(f, m.MkPos()))
}

// Returns n round with 0 decimals.
//    m  : Virtual machine.
func prRound(m *machine.T) {
	n := popFloat(m)
	pushFloat(m, math.Round(n))
}

// Returns n round with scale decimals. If 'scale' is out of [0-9) throws a
// "Math error".
//    m  : Virtual machine.
func prRoundn(m *machine.T) {
	tk := m.PopT(token.Int)
	scale, _ := tk.I()
	if scale < 0 || scale > 9 {
		m.Fail("Math error", "Expected: scale [0-9).\nActual  : '%v'.", scale)
	}
	n := popFloat(m)
	sc := math.Pow10(int(scale))
	pushFloat(m, math.Round(n*sc)/sc)
}

// Returns the number PI
//    m  : Virtual machine.
func prPi(m *machine.T) {
	pushFloat(m, math.Pi)
}

// Returns sin(n)
//    m  : Virtual machine.
func prSin(m *machine.T) {
	n := popFloat(m)
	pushFloat(m, math.Sin(n))
}

// Returns cos(n)
//    m  : Virtual machine.
func prCos(m *machine.T) {
	n := popFloat(m)
	pushFloat(m, math.Cos(n))
}

// Returns tan(n)
//    m  : Virtual machine.
func prTan(m *machine.T) {
	n := popFloat(m)
	pushFloat(m, math.Tan(n))
}

// Returns Asin(n)
//    m  : Virtual machine.
func prAsin(m *machine.T) {
	n := popFloat(m)
	pushFloat(m, math.Asin(n))
}

// Returns Acos(n)
//    m  : Virtual machine.
func prAcos(m *machine.T) {
	n := popFloat(m)
	pushFloat(m, math.Acos(n))
}

// Returns Atan(n)
//    m  : Virtual machine.
func prAtan(m *machine.T) {
	n := popFloat(m)
	pushFloat(m, math.Atan(n))
}

// Returns Atan(n1/n2)
//    m  : Virtual machine.
func prAtan2(m *machine.T) {
	n2 := popFloat(m)
	n1 := popFloat(m)
	pushFloat(m, math.Atan2(n1, n2))
}

// Returns sinh(n)
//    m  : Virtual machine.
func prSinh(m *machine.T) {
	n := popFloat(m)
	pushFloat(m, math.Sinh(n))
}

// Returns cosh(n)
//    m  : Virtual machine.
func prCosh(m *machine.T) {
	n := popFloat(m)
	pushFloat(m, math.Cosh(n))
}

// Returns tanh(n)
//    m  : Virtual machine.
func prTanh(m *machine.T) {
	n := popFloat(m)
	pushFloat(m, math.Tanh(n))
}

// Returns Asinh(n)
//    m  : Virtual machine.
func prAsinh(m *machine.T) {
	n := popFloat(m)
	pushFloat(m, math.Asinh(n))
}

// Returns Acosh(n)
//    m  : Virtual machine.
func prAcosh(m *machine.T) {
	n := popFloat(m)
	pushFloat(m, math.Acosh(n))
}

// Returns Atanh(n)
//    m  : Virtual machine.
func prAtanh(m *machine.T) {
	n := popFloat(m)
	pushFloat(m, math.Atanh(n))
}

// Returns the number E
//    m  : Virtual machine.
func prE(m *machine.T) {
	pushFloat(m, math.E)
}

// Returns e^n
//    m  : Virtual machine.
func prExp(m *machine.T) {
	n := popFloat(m)
	pushFloat(m, math.Exp(n))
}

// Returns 2^n
//    m  : Virtual machine.
func prExp2(m *machine.T) {
	n := popFloat(m)
	pushFloat(m, math.Exp2(n))
}

// Returns 10^n)
//    m  : Virtual machine.
func prExp10(m *machine.T) {
	n := popFloat(m)
	pushFloat(m, math.Pow(10.0, n))
}

// Returns lg(n)
//    m  : Virtual machine.
func prLog(m *machine.T) {
	n := popFloat(m)
	pushFloat(m, math.Log(n))
}

// Returns lg2(n)
//    m  : Virtual machine.
func prLog2(m *machine.T) {
	n := popFloat(m)
	pushFloat(m, math.Log2(n))
}

// Returns lg10(n)
//    m  : Virtual machine.
func prLog10(m *machine.T) {
	n := popFloat(m)
	pushFloat(m, math.Log10(n))
}

// Returns n1^n2
//    m  : Virtual machine.
func prPow(m *machine.T) {
	n2 := popFloat(m)
	n1 := popFloat(m)
	pushFloat(m, math.Pow(n1, n2))
}

// Returns n^(1/2)
//    m  : Virtual machine.
func prSqrt(m *machine.T) {
	n := popFloat(m)
	pushFloat(m, math.Sqrt(n))
}

// Processes string procedures.
//    m  : Virtual machine.
func Proc(m *machine.T) {
	tk, ok := m.PrgNext()
	if !ok {
		m.Failt("'math' procedure is missing")
	}
	sy, ok := tk.Sy()
	if !ok {
		m.Failt("Expected: 'math' procedure.\nActual  : '%v'.", tk.StringDraft())
	}
	switch sy {
	case symbol.New("round"):
		prRound(m)
	case symbol.New("roundn"):
		prRoundn(m)

	case symbol.New("pi"):
		prPi(m)
	case symbol.New("sin"):
		prSin(m)
	case symbol.New("cos"):
		prCos(m)
	case symbol.New("tan"):
		prTan(m)
	case symbol.New("asin"):
		prAsin(m)
	case symbol.New("acos"):
		prAcos(m)
	case symbol.New("atan"):
		prAtan(m)
	case symbol.New("atan2"):
		prAtan2(m)

	case symbol.New("sinh"):
		prSinh(m)
	case symbol.New("cosh"):
		prCosh(m)
	case symbol.New("tanh"):
		prTanh(m)
	case symbol.New("asinh"):
		prAsinh(m)
	case symbol.New("acosh"):
		prAcosh(m)
	case symbol.New("atanh"):
		prAtanh(m)

	case symbol.New("e"):
		prE(m)
	case symbol.New("exp"):
		prExp(m)
	case symbol.New("exp2"):
		prExp2(m)
	case symbol.New("exp10"):
		prExp10(m)
	case symbol.New("log"):
		prLog(m)
	case symbol.New("log2"):
		prLog2(m)
	case symbol.New("log10"):
		prLog10(m)
	case symbol.New("pow"):
		prPow(m)
	case symbol.New("sqrt"):
		prSqrt(m)

	default:
		m.Failt("'math' does not contains the procedure '%v'.", sy.String())
	}
}

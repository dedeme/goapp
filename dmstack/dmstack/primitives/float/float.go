// Copyright 10-Jan-2021 ºDeme
// GNU General Public License - V3 <http://www.gnu.org/licenses/>

// Float procedures.
package float

import (
	"fmt"
	"github.com/dedeme/dmstack/machine"
	"github.com/dedeme/dmstack/symbol"
	"github.com/dedeme/dmstack/token"
	"math/rand"
	"strconv"
	"strings"
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

// Returns a float from a string. If it fails, throws a "Float error".
//    m  : Virtual machine.
func prFromStr(m *machine.T) {
	tk := m.PopT(token.String)
	s, _ := tk.S()
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		m.Failt("'%v' bad Float.", tk.StringDraft())
	}
	pushFloat(m, f)
}

// Returns 'true' if n2 <= n1 + gap and n2 >= n1 - gap. If gap < 0 returns
// a "Float error".
//    m  : Virtual machine.
func prEq(m *machine.T) {
	gap := popFloat(m)
	if gap < 0 {
		m.Fail(machine.ERange(), "Gap is < 0 (%v)", gap)
	}
	n2 := popFloat(m)
	n1 := popFloat(m)
	m.Push(token.NewB(n2 <= n1+gap && n2 >= n1-gap, m.MkPos()))
}

// Returns 'false' if n2 <= n1 + gap and n2 >= n1 - gap. If gap < 0 returns
// a "Float error".
//    m  : Virtual machine.
func prNeq(m *machine.T) {
	gap := popFloat(m)
	if gap < 0 {
		m.Fail(machine.ERange(), "Gap is < 0 (%v)", gap)
	}
	n2 := popFloat(m)
	n1 := popFloat(m)
	m.Push(token.NewB(n2 > n1+gap || n2 < n1-gap, m.MkPos()))
}

// Returns the absolute value of f.
//    m  : Virtual machine.
func prAbs(m *machine.T) {
	f := popFloat(m)
	if f < 0 {
		f = -f
	}
	pushFloat(m, f)
}

// Returns a random number from 0 (inclusive) to 1 (exclusive)
//    m  : Virtual machine.
func prRnd(m *machine.T) {
	pushFloat(m, rand.Float64())
}

// Returns the maximum Float.
//    m  : Virtual machine.
func prMax(m *machine.T) {
	pushFloat(m, 1.797693134862315708145274237317043567981e+308)
}

// Returns the minimum Float.
//    m  : Virtual machine.
func prMin(m *machine.T) {
	pushFloat(m, -1.797693134862315708145274237317043567981e+308)
}

// Returns n truncated to Int.
//    m  : Virtual machine.
func prToInt(m *machine.T) {
	f := popFloat(m)
	m.Push(token.NewI(int64(f), m.MkPos()))
}

// Auxiliar function.
func toIsoEn(m *machine.T, isIso bool) {
	// Auxiliar function
	thousands := func(n string, isNeg bool) string {
		sep := ","
		if isIso {
			sep = "."
		}

		if isNeg {
			n = n[1:]
		}
		ix := 3
		for {
			if len(n) <= ix {
				break
			}
			pos := len(n) - ix
			n = n[:pos] + sep + n[pos:]
			ix += 4
		}
		if isNeg {
			n = "-" + n
		}
		return n
	}

	tk := m.PopT(token.Int)
	scale, _ := tk.I()
	if scale < 0 || scale > 9 {
		m.Fail(
			machine.ERange(),
			"\n  Expected: scale [0-9).\n  Actual  : '%v'.", scale,
		)
	}
	n := popFloat(m)
	fm := "%." + string(48+scale) + "f"
	ns := fmt.Sprintf(fm, n)
	if scale > 0 {
		sep := "."
		if isIso {
			sep = ","
		}
		ps := strings.Split(ns, ".")
		ns = thousands(ps[0], n < 0) + sep + ps[1]
	} else {
		ns = thousands(ns, n < 0)
	}
	m.Push(token.NewS(ns, m.MkPos()))
}

// Returns a Float converted to string in Iso format. If 'scale' is out of
// [0-9) throws a "Float error".
//    m  : Virtual machine.
func prToIso(m *machine.T) {
	toIsoEn(m, true)
}

// Returns a Float converted to string in En format. If 'scale' is out of
// [0-9) throws a "Float error".
//    m  : Virtual machine.
func prToEn(m *machine.T) {
	toIsoEn(m, false)
}

// Processes string procedures.
//    m   : Virtual machine.
//    proc: Procedure
func Proc(m *machine.T, proc symbol.T) {
	switch proc {
	case symbol.New("fromStr"):
		prFromStr(m)
	case symbol.New("eq"):
		prEq(m)
	case symbol.New("neq"):
		prNeq(m)
	case symbol.New("abs"):
		prAbs(m)
	case symbol.New("rnd"):
		prRnd(m)
	case symbol.New("max"):
		prMax(m)
	case symbol.New("min"):
		prMin(m)
	case symbol.New("toInt"):
		prToInt(m)
	case symbol.New("toIso"):
		prToIso(m)
	case symbol.New("toEn"):
		prToEn(m)

	default:
		m.Failt("'float' does not contains the procedure '%v'.", proc.String())
	}
}

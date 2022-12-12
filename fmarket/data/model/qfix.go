// Copyright 03-Dic-2022 ºDeme
// GNU General Public License - V3 <http://www.gnu.org/licenses/>

package model

import (
	"github.com/dedeme/ktlib/math"
)

func qfixCalc(
	closes [][]float64,
	params []float64,
	action func(closes, refs []float64),
) {
	jmp := params[0] + 1.0
	lgJmp := math.Log(jmp)

	downGap := func(q float64) float64 {
		return math.Pow(jmp, math.Round(math.Log(q)/lgJmp, 0)-1.0)
	}
	upGap := func(q float64) float64 {
		return math.Pow(jmp, math.Round(math.Log(q)/lgJmp, 0)+1.0)
	}
	downGap2 := func(q, ref float64) float64 {
		for {
			ref2 := ref * jmp
			if ref2*math.Sqrt(jmp) >= q {
				return ref
			}
			ref = ref2
		}
	}
	upGap2 := func(q, ref float64) float64 {
		for {
			ref2 := ref / jmp
			if ref2/math.Sqrt(jmp) <= q {
				return ref
			}
			ref = ref2
		}
	}

	nCos := len(closes[0])
	pvrow := make([]float64, nCos)
	refs := make([]float64, nCos)
	for i, cl := range closes[0] {
		pvrow[i] = cl
		refs[i] = downGap(cl) / jmp
	}

	for _, row := range closes {
		var newRefs []float64

		for i, q := range row {
			q0 := pvrow[i]
			ref := refs[i]

			if q0 < ref {
				if q < q0 {
					newRefs = append(newRefs, upGap2(q, ref))
				} else if q > ref {
					newRefs = append(newRefs, downGap(q))
				} else {
					newRefs = append(newRefs, ref)
				}
			} else {
				if q > q0 {
					newRefs = append(newRefs, downGap2(q, ref))
				} else if q < ref {
					newRefs = append(newRefs, upGap(q))
				} else {
					newRefs = append(newRefs, ref)
				}
			}
		}

		pvrow = row
		refs = newRefs
		action(row, refs)
	}
}

// Model Fix Quantum.
func QfixNew() *T {
	return &T{
		"QFIJO",
		[]string{
			"Intervalo",
		},
		[]int{
			4,
		},
		[]float64{
			0.1,
		},
		[]float64{
			0.3,
		},
		qfixCalc,
	}
}

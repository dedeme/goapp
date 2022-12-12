// Copyright 01-Jul-2022 ºDeme
// GNU General Public License - V3 <http://www.gnu.org/licenses/>

package model

func apprCalc(
	closes [][]float64,
	params []float64,
	action func(closes, refs []float64),
) {
	start := params[0]
	incr := params[1]
	nCos := len(closes[0])

	isSolds := make([]bool, nCos)
	refs := make([]float64, nCos)
	for i, c := range closes[0] {
		refs[i] = c * (1 - start)
	}

	for _, cs := range closes {
		newRefs := make([]float64, nCos)
		newIsSolds := make([]bool, nCos)

		for i, c := range cs {
			rf := refs[i]
			isSold := isSolds[i]

			if isSold {
				if c > rf {
					// newIsSolds[i] = false
					newRefs[i] = c * (1 - start)
				} else {
					newIsSolds[i] = true
					newRefs[i] = rf - (rf-c)*incr
				}
			} else {
				if c < rf {
					newIsSolds[i] = true
					newRefs[i] = c * (1 + start)
				} else {
					// newIsSolds[i] = false
					newRefs[i] = rf + (c-rf)*incr
				}
			}

		}
		refs = newRefs
		isSolds = newIsSolds

		action(cs, refs)
	}

}

// Model Approximations.
func ApprNew() *T {
	return &T{
		"APRX",
		[]string{
			"Inicio",
			"Aproximación",
		},
		[]int{
			4,
			4,
		},
		[]float64{
			0.02,
			0.002,
		},
		[]float64{
			0.45,
			0.05,
		},
		apprCalc,
	}
}

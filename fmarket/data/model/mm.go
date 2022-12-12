// Copyright 04-Jul-2022 ºDeme
// GNU General Public License - V3 <http://www.gnu.org/licenses/>

package model

func mmCalc(
	closes [][]float64,
	params []float64,
	action func(closes, refs []float64),
) {
	days := int(params[0] + 0.4)
	strip := params[1]
	nCos := len(closes[0])

	isSolds := make([]bool, nCos)
	refs := make([]float64, nCos)

	for day, cs := range closes {
		newIsSolds := make([]bool, nCos)
		newRefs := make([]float64, nCos)

		if day < days {
			newIsSolds = isSolds
			for i, c := range cs {
				newRefs[i] = c * (1 - strip)
			}
		} else if day == days {
			for i, c := range cs {
				oldC := closes[0][i]
				if oldC > c {
					newIsSolds[i] = true
					newRefs[i] = oldC * (1 + strip)
				} else {
					newRefs[i] = oldC * (1 - strip)
				}
			}
		} else {
			for i, c := range cs {
				isSold := isSolds[i]
				rf := refs[i]
				oldC := closes[day-days][i]

				if isSold {
					if c > rf {
						// newIsSolds[i] = false
						newRefs[i] = oldC * (1 - strip)
					} else {
						newIsSolds[i] = true
						newRf := oldC * (1 + strip)
						if newRf < rf {
							newRefs[i] = newRf
						} else {
							newRefs[i] = rf
						}
					}
				} else {
					if c < rf {
						newIsSolds[i] = true
						newRefs[i] = oldC * (1 + strip)
					} else {
						// newIsSolds[i] = false
						newRf := oldC * (1 - strip)
						if newRf > rf {
							newRefs[i] = newRf
						} else {
							newRefs[i] = rf
						}
					}
				}
			}
		}

		isSolds = newIsSolds
		refs = newRefs

		action(cs, refs)
	}

}

// Maximum-Minimum.
func MmNew() *T {
	return &T{
		"MX_MN",
		[]string{
			"Días",
			"Banda",
		},
		[]int{
			0,
			4,
		},
		[]float64{
			5.0,
			0.002,
		},
		[]float64{
			250.0,
			0.1,
		},
		mmCalc,
	}
}

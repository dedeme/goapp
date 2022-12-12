// Copyright 04-Jul-2022 ºDeme
// GNU General Public License - V3 <http://www.gnu.org/licenses/>

package model

func maCalc(
	closes [][]float64,
	params []float64,
	action func(closes, refs []float64),
) {
	days := int(params[0] + 0.4)
	strip := params[1]
	nCos := len(closes[0])

	isSolds := make([]bool, nCos)
	sums := make([]float64, nCos)
	refs := make([]float64, nCos)

	for day, cs := range closes {
		newSums := make([]float64, nCos)
		newIsSolds := make([]bool, nCos)
		newRefs := make([]float64, nCos)

		if day < days {
			newIsSolds = isSolds
			for i, c := range cs {
				newSums[i] = sums[i] + c
				newRefs[i] = c * (1 - strip)
			}
		} else if day == days {
			for i, c := range cs {
				sum := sums[i] + c - closes[day-days][i]
				avg := sum / float64(days)
				newSums[i] = sum
				if avg > c {
					newIsSolds[i] = true
					newRefs[i] = avg * (1 + strip)
				} else {
					newRefs[i] = avg * (1 - strip)
				}
			}
		} else {
			for i, c := range cs {
				isSold := isSolds[i]
				rf := refs[i]
				sum := sums[i] + c - closes[day-days][i]
				avg := sum / float64(days)
				newSums[i] = sum

				if isSold {
					if c > rf {
						// newIsSolds[i] = false
						newRefs[i] = avg * (1 - strip)
					} else {
						newIsSolds[i] = true
						newRf := avg * (1 + strip)
						if newRf < rf {
							newRefs[i] = newRf
						} else {
							newRefs[i] = rf
						}
					}
				} else {
					if c < rf {
						newIsSolds[i] = true
						newRefs[i] = avg * (1 + strip)
					} else {
						// newIsSolds[i] = false
						newRf := avg * (1 - strip)
						if newRf > rf {
							newRefs[i] = newRf
						} else {
							newRefs[i] = rf
						}
					}
				}
			}
		}

		sums = newSums
		isSolds = newIsSolds
		refs = newRefs

		action(cs, refs)
	}

}

// Mobil average.
func MaNew() *T {
	return &T{
		"MM",
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
		maCalc,
	}
}

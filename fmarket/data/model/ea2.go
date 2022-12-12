// Copyright 04-Jul-2022 ºDeme
// GNU General Public License - V3 <http://www.gnu.org/licenses/>

package model

func ea2Calc(
	closes [][]float64,
	params []float64,
	action func(closes, refs []float64),
) {
	days := int(params[0] + 0.4)
	strip := params[1]
	nCos := len(closes[0])

	isSolds := make([]bool, nCos)
	avgs := make([]float64, nCos)
	refs := make([]float64, nCos)

	for day, cs := range closes {
		newAvgs := make([]float64, nCos)
		newIsSolds := make([]bool, nCos)
		newRefs := make([]float64, nCos)

		if day < days {
			newIsSolds = isSolds
			for i, c := range cs {
				newAvgs[i] = avgs[i] + c
				newRefs[i] = c * (1 - strip)
			}
		} else if day == days {
			for i, c := range cs {
				avg := avgs[i] / float64(days)
				newAvg := avg + (c-avg)/float64(days)

				newAvgs[i] = newAvg
				if newAvg > c {
					newIsSolds[i] = true
					newRefs[i] = newAvg * (1 + strip)
				} else {
					newRefs[i] = newAvg * (1 - strip)
				}
			}
		} else {
			for i, c := range cs {
				isSold := isSolds[i]
				avg := avgs[i]
				rf := refs[i]
				newAvg := avg + (c-avg)/float64(days)
				newAvgs[i] = newAvg

				if isSold {
					if c > rf {
						// newIsSolds[i] = false
						newRefs[i] = newAvg * (1 - strip)
					} else {
						newIsSolds[i] = true
						newRf := newAvg * (1 + strip)
						if newRf < rf {
							newRefs[i] = newRf
						} else {
							newRefs[i] = rf
						}
					}
				} else {
					if c < rf {
						newIsSolds[i] = true
						newRefs[i] = newAvg * (1 + strip)
					} else {
						// newIsSolds[i] = false
						newRf := newAvg * (1 - strip)
						if newRf > rf {
							newRefs[i] = newRf
						} else {
							newRefs[i] = rf
						}
					}
				}
			}
		}

		avgs = newAvgs
		isSolds = newIsSolds
		refs = newRefs

		action(cs, refs)
	}

}

// Exponential mobil average (2).
func Ea2New() *T {
	return &T{
		"ME2",
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
		ea2Calc,
	}
}

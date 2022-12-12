// Copyright 04-Dic-2022 ºDeme
// GNU General Public License - V3 <http://www.gnu.org/licenses/>

// Quotes table data
package quotes

import (
	"github.com/dedeme/ktlib/arr"
	"github.com/dedeme/ktlib/js"
)

// Quotes table data
type T struct {
	// Date in which data was read.
	Date string
	//Company nicks.
	Cos []string
	// Dates in format YYYYMMDD, from before to after.
	Dates []string
	// Matrix of normalized opens quotes (without -1).
	// Its rows match 'dates' and is columns 'cos'.
	Opens [][]float64
	// Matrix of normalized closes quotes (without -1).
	// Its rows match 'dates' and is columns 'cos'.
	Closes [][]float64
	// Matrix of normalized maximum quotes (without -1).
	// Its rows match 'dates' and is columns 'cos'.
	Maxs [][]float64
}

// Constructor
//   date  : Date in which data was read.
//   cos   : Company nicks.
//   dates : Dates in format YYYYMMDD, from before to after.
//   opens : Matrix of normalized opens quotes (without -1).
//           Its rows match 'dates' and is columns 'cos'.
//   closes: Matrix of normalized closes quotes (without -1).
//           Its rows match 'dates' and is columns 'cos'.
//   maxs  : Matrix of normalized maximum quotes (without -1).
//           Its rows match 'dates' and is columns 'cos'.
func New(date string, cos, dates []string, opens, closes, maxs [][]float64) *T {
	return &T{date, cos, dates, opens, closes, maxs}
}

// Returns the index of one company from its nick, or -1 if nick is not found.
func (qs *T) CompanyIndex(nick string) int {
	return arr.Index(qs.Cos, nick)
}

// Returns data of the company with index 'coIx' in format *T
func (qs *T) GetSingle(coIx int) *T {
	return New(
		qs.Date,
		[]string{qs.Cos[coIx]},
		qs.Dates,
		arr.Map(qs.Opens, func(row []float64) []float64 {
			return []float64{row[coIx]}
		}),
		arr.Map(qs.Closes, func(row []float64) []float64 {
			return []float64{row[coIx]}
		}),
		arr.Map(qs.Maxs, func(row []float64) []float64 {
			return []float64{row[coIx]}
		}),
	)
}

func ToJs(qs *T) string {
	return js.Wa([]string{
		js.Ws(qs.Date),
		js.Wa(arr.Map(qs.Cos, js.Ws)),
		js.Wa(arr.Map(qs.Dates, js.Ws)),
		js.Wa(arr.Map(qs.Opens, func(cos []float64) string {
			return js.Wa(arr.Map(cos, js.Wd))
		})),
		js.Wa(arr.Map(qs.Closes, func(cos []float64) string {
			return js.Wa(arr.Map(cos, js.Wd))
		})),
		js.Wa(arr.Map(qs.Maxs, func(cos []float64) string {
			return js.Wa(arr.Map(cos, js.Wd))
		})),
	})
}

func FromJs(j string) *T {
	a := js.Ra(j)
	return New(
		js.Rs(a[0]),
		arr.Map(js.Ra(a[1]), js.Rs),
		arr.Map(js.Ra(a[2]), js.Rs),
		arr.Map(js.Ra(a[3]), func(j2 string) []float64 {
			return arr.Map(js.Ra(j2), js.Rd)
		}),
		arr.Map(js.Ra(a[4]), func(j2 string) []float64 {
			return arr.Map(js.Ra(j2), js.Rd)
		}),
		arr.Map(js.Ra(a[5]), func(j2 string) []float64 {
			return arr.Map(js.Ra(j2), js.Rd)
		}),
	)
}

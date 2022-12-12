// Copyright 05-Dic-2022 ºDeme
// GNU General Public License - V3 <http://www.gnu.org/licenses/>

// Model evaluation data.
package modelEval

import (
	"github.com/dedeme/ktlib/arr"
	"github.com/dedeme/ktlib/js"
)

type T struct {
	// Params evaluated.
	Params []float64
	// Weeks weight to calculate historic values.
	Weeks int
	// Historic value.
	Hvalue float64
	// Historic sales.
	Hsales float64
	// Last value.
	Value float64
	// Last sales.
	Sales float64
}

func New(params []float64, weeks int, hvalue, hsales, value, sales float64) *T {
	return &T{params, weeks, hvalue, hsales, value, sales}
}

func ToJs(me *T) string {
	return js.Wa([]string{
		js.Wa(arr.Map(me.Params, func(n float64) string {
			return js.WdDec(n, 6)
		})),
		js.Wi(me.Weeks),
		js.Wd(me.Hvalue),
		js.Wd(me.Hsales),
		js.Wd(me.Value),
		js.Wd(me.Sales),
	})
}

func FromJs(j string) *T {
	a := js.Ra(j)
	return New(
		arr.Map(js.Ra(a[0]), js.Rd),
		js.Ri(a[1]),
		js.Rd(a[2]),
		js.Rd(a[3]),
		js.Rd(a[4]),
		js.Rd(a[5]),
	)
}

type TbT struct {
	// Date of last Sunday evaluation.
	Date string
	// Evaluations
	Evals []*T
}

func NewTb(date string, evals []*T) *TbT {
	return &TbT{date, evals}
}

func TbToJs(t *TbT) string {
	return js.Wa([]string{
		js.Ws(t.Date),
		js.Wa(arr.Map(t.Evals, ToJs)),
	})
}

func TbFromJs(j string) *TbT {
	a := js.Ra(j)
	return NewTb(
		js.Rs(a[0]),
		arr.Map(js.Ra(a[1]), FromJs),
	)
}

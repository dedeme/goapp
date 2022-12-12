// Copyright 02-Dic-2022 ºDeme
// GNU General Public License - V3 <http://www.gnu.org/licenses/>

package tests

import (
	"github.com/dedeme/fmarket/data/model"
	"github.com/dedeme/fmarket/data/modelEval"
	"github.com/dedeme/fmarket/db"
	"github.com/dedeme/ktlib/arr"
	"github.com/dedeme/ktlib/str"
	"github.com/dedeme/ktlib/sys"
)

func ModelTests() {
	sys.Println("Model Tests")

	qs := db.QuotesTb().Read()
	for _, md := range model.List() {
		evals := db.EvalsDb(md.Id()).Read().Evals
		arr.Sort(evals, func(e1, e2 *modelEval.T) bool {
			return e1.Hvalue > e2.Hvalue // Reverse order
		})
		first := evals[0]
		params := first.Params

		//    if md.Id() == "QFIJO" {
		//      params = []float64{0.185}
		//    }

		profitsPc := md.Simulation(qs, params)

		if len(params) == 2 {
			sys.Println(str.Fmt("%v [%v, %v]: %v",
				md.Id(), params[0], params[1], profitsPc,
			))
		} else {
			sys.Println(str.Fmt("%v [%v]: %v", md.Id(), params[0], profitsPc))
		}
	}

	sys.Println("  finished")
}

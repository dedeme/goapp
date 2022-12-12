// Copyright 02-Dic-2022 ºDeme
// GNU General Public License - V3 <http://www.gnu.org/licenses/>

package tests

import (
	"github.com/dedeme/fmarket/data/cts"
	"github.com/dedeme/fmarket/data/flea"
	"github.com/dedeme/fmarket/data/model"
	"github.com/dedeme/fmarket/db"
	"github.com/dedeme/ktlib/sys"
)

func FleaTests() {
	sys.Println("Flea Tests")
	f := flea.New(
		int64(123),
		1,
		true,
		[]*flea.FmodelT{
			flea.NewFmodel("EA", []float64{20, 0.05}),
			flea.NewFmodel("EA2", []float64{21, 0.06}),
			flea.NewFmodel("MA", []float64{22, 0.07}),
		},
		300000.0,
	)
	fJs := flea.ToJs(f)
	eq(flea.ToJs(flea.FromJs(fJs)), fJs)
	fleas := db.FleasTb().Read().Fleas()
	eq(len(fleas), cts.FleasPerModel*len(model.List()))
	sys.Println("  finished")
}

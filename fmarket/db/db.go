// Copyright 05-Dic-2022 ºDeme
// GNU General Public License - V3 <http://www.gnu.org/licenses/>

// Data base of jstb's.
package db

import (
	"github.com/dedeme/fmarket/data/cts"
	"github.com/dedeme/fmarket/data/flea"
	"github.com/dedeme/fmarket/data/mdStats"
	"github.com/dedeme/fmarket/data/model"
	"github.com/dedeme/fmarket/data/modelEval"
	"github.com/dedeme/fmarket/data/quotes"
	"github.com/dedeme/fmarket/data/ranking"
	"github.com/dedeme/ktlib/arr"
	"github.com/dedeme/ktlib/file"
	"github.com/dedeme/ktlib/js"
	"github.com/dedeme/ktlib/jstb"
	"github.com/dedeme/ktlib/path"
)

// Initialize data base
func Initialize() {
	if !file.IsDirectory(cts.FMarketPath) {
		file.Mkdir(cts.FMarketPath)
	}
}

// Returns KtMMarket 'quotes.tb'
func QuotesTb() *jstb.T[*quotes.T] {
	return jstb.New(
		path.Cat(cts.MMarketPath, "quotes.tb"),
		nil,
		quotes.ToJs,
		quotes.FromJs,
	)
}

// Returns KtMMarket 'evals/' + modelId + '.tb' (Used for testing)
func EvalsDb(modelId string) *jstb.T[*modelEval.TbT] {
	return jstb.New(
		path.Cat(cts.MMarketPath, "evals", modelId+".tb"),
		nil,
		modelEval.TbToJs,
		modelEval.TbFromJs,
	)
}

// Returns FMarket 'activity.tb'
func ActivityTb() *jstb.T[string] {
	return jstb.New(
		path.Cat(cts.FMarketPath, "activity.tb"),
		func() string { return "stopped" },
		js.Ws,
		js.Rs,
	)
}

// Returns FMarket 'fleas.tb'
func FleasTb() *jstb.T[*flea.TbT] {
	return jstb.New(
		path.Cat(cts.FMarketPath, "fleas.tb"),
		func() *flea.TbT {
			qs := QuotesTb().Read()
			models := model.List()
			nFleas := len(models) * cts.FleasPerModel
			fleas := make([]*flea.T, nFleas)
			for i := range fleas {
				mds := arr.Copy(models)
				arr.Shuffle(mds)
				fmodels := make([]*flea.FmodelT, 3)
				assets := 0.0
				for mdIx := 0; mdIx < 3; mdIx++ {
					md := mds[mdIx]
					params := md.Mutation()
					fmodels[mdIx] = flea.NewFmodel(md.Id(), params)
					assets += md.Simulation(qs, params)
				}
				if i%2 == 0 {
					fleas[i] = flea.New(int64(i), 0, true, fmodels, assets)
				} else {
					fleas[i] = flea.New(int64(i), 0, false, fmodels, assets)
				}
			}
			flea.Sort(fleas)
			return flea.NewTbFromSorted(int64(nFleas), 1, fleas)
		},
		flea.TbToJs,
		flea.TbFromJs,
	)
}

// Returns FMarket 'rankings.tb"
func RankingsTb() *jstb.T[*ranking.DbT] {
	return jstb.New(
		path.Cat(cts.FMarketPath, "rankings.tb"),
		func() *ranking.DbT {
			return ranking.MakeDbFromSorted(FleasTb().Read().Fleas())
		},
		ranking.DbToJs,
		ranking.DbFromJs,
	)
}

// Returns FMarket 'mdStats.tb"
func MdStatsTb() *jstb.T[*mdStats.DbT] {
	return jstb.New(
		path.Cat(cts.FMarketPath, "mdStats.tb"),
		func() *mdStats.DbT {
			return mdStats.MakeDbFromSorted(FleasTb().Read().Fleas())
		},
		mdStats.DbToJs,
		mdStats.DbFromJs,
	)
}

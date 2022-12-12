// Copyright 06-Dic-2022 ºDeme
// GNU General Public License - V3 <http://www.gnu.org/licenses/>

// Ranking data
package mdStats

import (
	"github.com/dedeme/fmarket/data/cts"
	"github.com/dedeme/fmarket/data/flea"
	"github.com/dedeme/fmarket/data/model"
	"github.com/dedeme/ktlib/arr"
	"github.com/dedeme/ktlib/js"
)

type OrderFleaT struct {
	// Order
	ix int
	// Flea
	fl *flea.T
}

func newOrderFlea(ix int, fl *flea.T) *OrderFleaT {
	return &OrderFleaT{ix, fl}
}

func orderFleaToJs(of *OrderFleaT) string {
	return js.Wa([]string{
		js.Wi(of.ix),
		flea.ToJs(of.fl),
	})
}

func orderFleaFromJs(j string) *OrderFleaT {
	a := js.Ra(j)
	return newOrderFlea(
		js.Ri(a[0]),
		flea.FromJs(a[1]),
	)
}

type T struct {
	// Model identifier
	modelId string
	// The three bests fleas. It can be less than three.
	bests []*OrderFleaT
	// The three worst fleas. It can be less than three.
	worsts []*OrderFleaT
	// Number of active fleas
	nFleas int
	// Position average
	position int
	// Assets average
	assets float64
}

func new(
	modelId string, bests, worsts []*OrderFleaT,
	nFleas int, position int, assets float64,
) *T {
	return &T{modelId, bests, worsts, nFleas, position, assets}
}

func toJs(m *T) string {
	return js.Wa([]string{
		js.Ws(m.modelId),
		js.Wa(arr.Map(m.bests, orderFleaToJs)),
		js.Wa(arr.Map(m.worsts, orderFleaToJs)),
		js.Wi(m.nFleas),
		js.Wi(m.position),
		js.WdDec(m.assets, 2),
	})
}

func fromJs(j string) *T {
	a := js.Ra(j)
	return new(
		js.Rs(a[0]),
		arr.Map(js.Ra(a[1]), orderFleaFromJs),
		arr.Map(js.Ra(a[2]), orderFleaFromJs),
		js.Ri(a[3]),
		js.Ri(a[4]),
		js.Rd(a[5]),
	)
}

type GroupT struct {
	modelIds   []string
	duplicates int
}

func newGroup(modelIds []string, duplicates int) *GroupT {
	return &GroupT{modelIds, duplicates}
}

func (g1 *GroupT) eqModel(mds []string) bool {
	mds1 := g1.modelIds
	return mds1[0] == mds[0] && mds1[1] == mds[1] && mds1[2] == mds[2]
}

func groupToJs(g *GroupT) string {
	return js.Wa([]string{
		js.Wa(arr.Map(g.modelIds, js.Ws)),
		js.Wi(g.duplicates),
	})
}

func groupFromJs(j string) *GroupT {
	a := js.Ra(j)
	return newGroup(
		arr.Map(js.Ra(a[0]), js.Rs),
		js.Ri(a[1]),
	)
}

type DbT struct {
	models        []*T
	groupsRanking []*GroupT
}

func newDb(ranks []*T, groupsRanking []*GroupT) *DbT {
	return &DbT{ranks, groupsRanking}
}

// Create a models statitics database from a sorted slice of fleas
func MakeDbFromSorted(fleas []*flea.T) *DbT {
	var models []*T
	for _, md := range model.List() {
		mdId := md.Id()
		var fs []*OrderFleaT
		sumPositions := 0
		sumAssets := 0.0
		for i, f := range fleas {
			if f.HasModel(mdId) {
				fs = append(fs, newOrderFlea(i, f))
				sumPositions += i
				sumAssets += f.Assets()
			}
		}
		lenFs := len(fs)
		position := sumPositions / lenFs
		assets := sumAssets / float64(lenFs)
		if lenFs < 4 {
			models = append(models, new(mdId, fs, fs, lenFs, position, assets))
		} else {
			models = append(models, new(mdId, fs[:3], fs[lenFs-3:], lenFs, position, assets))
		}
	}

	var groups []*GroupT
	for _, f := range fleas {
		ids := f.ModelIds()
		ix := arr.Indexf(groups, func(g *GroupT) bool {
			return g.eqModel(ids)
		})
		if ix == -1 {
			groups = append(groups, newGroup(ids, 1))
		} else {
			groups[ix].duplicates++
		}
	}
	arr.Sort(groups, func(g1, g2 *GroupT) bool {
		return g1.duplicates > g2.duplicates // Reverse order
	})

	return newDb(models, groups[:cts.GroupsInModelRanking])
}

func DbToJs(rdb *DbT) string {
	return js.Wa([]string{
		js.Wa(arr.Map(rdb.models, toJs)),
		js.Wa(arr.Map(rdb.groupsRanking, groupToJs)),
	})
}

func DbFromJs(j string) *DbT {
	a := js.Ra(j)
	return newDb(
		arr.Map(js.Ra(a[0]), fromJs),
		arr.Map(js.Ra(a[1]), groupFromJs),
	)
}

// Copyright 06-Dic-2022 ºDeme
// GNU General Public License - V3 <http://www.gnu.org/licenses/>

// Ranking data
package ranking

import (
	"github.com/dedeme/fmarket/data/cts"
	"github.com/dedeme/fmarket/data/flea"
	"github.com/dedeme/ktlib/arr"
	"github.com/dedeme/ktlib/js"
	"github.com/dedeme/ktlib/time"
)

// Ranking structure
type T struct {
	date  string
	fleas []*flea.T
}

func new(date string, fleas []*flea.T) *T {
	return &T{date, fleas}
}

// Create a ranking from a sorted slice of fleas
func makeFromSorted(fleas []*flea.T) *T {
	return new(time.ToStr(time.Now()), fleas[:cts.FleasInRanking])
}

func toJs(r *T) string {
	return js.Wa([]string{
		js.Ws(r.date),
		js.Wa(arr.Map(r.fleas, flea.ToJs)),
	})
}

func fromJs(j string) *T {
	a := js.Ra(j)
	return new(
		js.Rs(a[0]),
		arr.Map(js.Ra(a[1]), flea.FromJs),
	)
}

// Array of rankings from after to before
type DbT struct {
	ranks []*T
}

func newDb(ranks []*T) *DbT {
	return &DbT{ranks}
}

// Create a ranking data base from a sorted slice of fleas
func MakeDbFromSorted(fleas []*flea.T) *DbT {
	var r []*T
	rk := makeFromSorted(fleas)
	for i := 0; i < cts.RankingsInDatabase; i++ {
		r = append(r, rk)
	}
	return newDb(r)
}

// Modify rdb in place, adding the sorted slice 'fleas'
func (rdb *DbT) UpdateFromSorted(fleas []*flea.T) {
	rk := makeFromSorted(fleas)
	oldDate := rdb.ranks[0].date
	newDate := rk.date
	if newDate < oldDate {
		panic("Old date " + oldDate + " is after new date " + newDate)
	}
	if newDate == oldDate {
		rdb.ranks[0] = rk
	} else {
		arr.Pop(&rdb.ranks)
		arr.Unshift(&rdb.ranks, rk)
	}
}

func DbToJs(rdb *DbT) string {
	return js.Wa(arr.Map(rdb.ranks, toJs))
}

func DbFromJs(j string) *DbT {
	return newDb(arr.Map(js.Ra(j), fromJs))
}

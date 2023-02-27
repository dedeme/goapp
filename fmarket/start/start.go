// Copyright 01-Dic-2022 ºDeme
// GNU General Public License - V3 <http://www.gnu.org/licenses/>

// Start run.
package start

import (
	"github.com/dedeme/fmarket/data/cts"
	"github.com/dedeme/fmarket/data/flea"
	"github.com/dedeme/fmarket/data/mdStats"
	"github.com/dedeme/fmarket/data/model"
	"github.com/dedeme/fmarket/db"
	"github.com/dedeme/ktlib/arr"
	"github.com/dedeme/ktlib/math"
	"github.com/dedeme/ktlib/thread"
)

func Run() {
	activityTb := db.ActivityTb()
	if activityTb.Read() != "stopped" {
		return
	}

	activityTb.Write("Generating")
	models := model.List()
	fleasTb := db.FleasTb()
	fleasTbData := fleasTb.Read()
	id := fleasTbData.NextId()
	cycle := fleasTbData.NextCycle()
	fleas := fleasTbData.Fleas()
	fleas, _ = arr.Duplicatesf(fleas, func(f1, f2 *flea.T) bool {
		return f1.IsMale() == f2.IsMale() && f1.EqModel(f2)
	})

	var males []*flea.T
	var females []*flea.T
	for _, f := range fleas {
		if f.IsMale() {
			males = append(males, f)
		} else {
			if !arr.Anyf(females, func(f2 *flea.T) bool {
				return f.EqModel(f2)
			}) {
				females = append(females, f)
			}
		}
	}
	malesLen := len(males)
	females = females[:len(females)/2]
	for {
		for _, fm := range females {
			fleas = append(fleas, flea.Generate(id, cycle, fm, males[math.Rndi(malesLen)]))
			id++
		}
		if len(fleas) > cts.FleasPerModel*len(models) {
			break
		}
	}

	activityTb.Write("Evaluating")
	qs := db.QuotesTb().Read()
	var chs []chan bool
	for _, f := range fleas {
		f2 := f
		chs = append(chs, thread.Start(func() {
			f2.Update(qs)
		}))
	}
	for _, ch := range chs {
		thread.Join(ch)
	}

	activityTb.Write("Selecting")
	flea.Sort(fleas)
	fleas = fleas[:cts.FleasPerModel*len(models)]

	activityTb.Write("Saving")
	fleasTb.Write(flea.NewTbFromSorted(id, cycle+1, fleas))
	rankingsTb := db.RankingsTb()
	rankings := rankingsTb.Read()
	rankings.UpdateFromSorted(fleas)
	rankingsTb.Write(rankings)
	db.MdStatsTb().Write(mdStats.MakeDbFromSorted(fleas))

	activityTb.Write("stopped")
}

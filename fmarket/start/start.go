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
	"github.com/dedeme/ktlib/math"
	"github.com/dedeme/ktlib/thread"
)

func Run() {
	activityTb := db.ActivityTb()
	if activityTb.Read() != "stopped" {
		return
	}

	activityTb.Write("Generating")
	fleasTb := db.FleasTb()
	fleasTbData := fleasTb.Read()
	id := fleasTbData.NextId()
	cycle := fleasTbData.NextCycle()
	fleas := fleasTbData.Fleas()

	var males []*flea.T
	var females []*flea.T
	for _, f := range fleas {
		if f.IsMale() {
			males = append(males, f)
		} else {
			females = append(females, f)
		}
	}
	malesLen := len(males)
	females = females[:len(females)/2]
	for _, fm := range females {
		fleas = append(fleas, flea.Generate(id, cycle, fm, males[math.Rndi(malesLen)]))
		id++
	}

	activityTb.Write("Evaluating")
	qs := db.QuotesTb().Read()
  var chs []chan bool
  for _, f := range fleas {
    f2 := f
    chs = append(chs, thread.Start(func () {
      f2.Update(qs)
    }))
	}
  for _, ch := range chs {
    thread.Join(ch)
  }

	models := model.List()
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

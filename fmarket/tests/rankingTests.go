// Copyright 06-Dic-2022 ºDeme
// GNU General Public License - V3 <http://www.gnu.org/licenses/>

package tests

import (
	"github.com/dedeme/fmarket/data/cts"
	"github.com/dedeme/fmarket/db"
	"github.com/dedeme/ktlib/js"
	"github.com/dedeme/ktlib/sys"
)

func RankingTests() {
	sys.Println("Ranking Tests")
	eq(len(js.Ra(db.RankingsTb().ReadJs())), cts.RankingsInDatabase)
	sys.Println("  finished")
}

// Copyright 08-Dic-2022 ºDeme
// GNU General Public License - V3 <http://www.gnu.org/licenses/>

package tests

import (
	"github.com/dedeme/fmarket/db"
	"github.com/dedeme/ktlib/sys"
)

func ActivityTests() {
	sys.Println("Activity Tests")

	var act string
	activityTb := db.ActivityTb()
	act = activityTb.Read()
	activityTb.Write(act)
	eq(activityTb.Read(), act)
	sys.Println("  finished")
}

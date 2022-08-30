// Copyright 31-May-2022 ºDeme
// GNU General Public License - V3 <http://www.gnu.org/licenses/>

package clock

import (
	"github.com/dedeme/barbery/cts"
	"github.com/dedeme/ktlib/sys"
)

var TimeOver = false

func Run() {
	sys.Sleep(cts.ClockTime)
	sys.Println("Clock: Time is over")
	TimeOver = true
}

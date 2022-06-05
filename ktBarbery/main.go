// Copyright 31-May-2022 ºDeme
// GNU General Public License - V3 <http://www.gnu.org/licenses/>

package main

import (
	"github.com/dedeme/ktBarbery/barber"
	"github.com/dedeme/ktBarbery/client"
	"github.com/dedeme/ktBarbery/clock"
	"github.com/dedeme/ktlib/sys"
	"github.com/dedeme/ktlib/thread"
)

func main() {
	sys.Rand()

	th1 := thread.Start(clock.Run)
	th2 := thread.Start(client.MakerRun)
	th3 := thread.Start(barber.Run)

	thread.Join(th1)
	thread.Join(th2)
	thread.Join(th3)

	sys.Println("Program end.")
}

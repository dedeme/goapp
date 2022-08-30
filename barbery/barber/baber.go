// Copyright 31-May-2022 ºDeme
// GNU General Public License - V3 <http://www.gnu.org/licenses/>

package barber

import (
	"github.com/dedeme/barbery/clock"
	"github.com/dedeme/barbery/cts"
	"github.com/dedeme/barbery/shop"
	"github.com/dedeme/ktlib/str"
	"github.com/dedeme/ktlib/sys"
	"github.com/dedeme/ktlib/thread"
)

func Run() {
	sys.Println("Barber: Go out from home")
	sys.Sleep(1000)
	sys.Println("Barber: Open barbery")
	shop.Open()

	end := false
	sleeping := false
	for !end {
		cl := ""
		thread.Sync(func() {
			if clock.TimeOver && shop.IsOpen {
				sys.Println("Barber: Close barbery")
				shop.Close()
			}

			if shop.IsEmpty() {
				if !shop.IsOpen {
					end = true
				} else if !sleeping {
					sleeping = true
					sys.Println("Barber: Barber is sleeping")
				}
			} else {
				sleeping = false
				cl = shop.TakeAClient()
			}
		})
		if cl != "" {
			cutHair(cl)
		} else if !end {
			sys.Sleep(10)
		}
	}
}

func cutHair(cl string) {
	sys.Println(str.Fmt("Barber: Cutting hair to %s", cl))
	sys.Println(str.Fmt("Shop: %s", shop.SitsToStr()))
	sys.Sleep(cts.BarberTime)
}

// Copyright 31-May-2022 ºDeme
// GNU General Public License - V3 <http://www.gnu.org/licenses/>

package client

import (
	"github.com/dedeme/ktBarbery/cts"
	"github.com/dedeme/ktBarbery/shop"
	"github.com/dedeme/ktlib/math"
	"github.com/dedeme/ktlib/str"
	"github.com/dedeme/ktlib/sys"
	"github.com/dedeme/ktlib/thread"
)

func MakerRun() {
	i := 0
	for shop.IsPreopen || shop.IsOpen {
		i++
		thread.Run(func() {
			run(str.Fmt("%d", i))
		})
		sys.Sleep(1 + math.Rndi(cts.ClientMakerTime))
	}
}

func run(cl string) {
	msg := func(m string) {
		sys.Println(str.Fmt("Client %s: %s", cl, m))
	}

	msg("Arrive to barbery")
	if !shop.IsOpen {
		msg("Go out because barbery is close")
	} else if shop.IsFull() {
		msg("Go out because barbery is full")
	} else {
		thread.Sync(func() {
			shop.TakeASit(cl)
		})
	}
}

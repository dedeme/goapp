// Copyright 02-Dic-2022 ºDeme
// GNU General Public License - V3 <http://www.gnu.org/licenses/>

package tests

import (
	"github.com/dedeme/fmarket/data/broker"
	"github.com/dedeme/ktlib/sys"
)

func BrokerTests() {
	sys.Println("Broker Tests")

	fees80000 := 80000.00*0.00103 + 0.11
	eq(broker.Fees(80000.00), fees80000)
	eq(broker.Buy(1000, 80.00), 80000.00+fees80000+80000.00*0.002)
	eq(broker.Sell(1000, 80.00), 80000.00-fees80000)

	fees5000 := 5000.00*0.00003 + 9.86
	eq(broker.Fees(5000.00), fees5000)
	eq(broker.Buy(1000, 5.00), 5000.00+fees5000+5000.00*0.002)
	eq(broker.Sell(1000, 5.00), 5000.00-fees5000)

	sys.Println("  finished")
}

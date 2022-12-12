// Copyright 02-Dic-2022 ºDeme
// GNU General Public License - V3 <http://www.gnu.org/licenses/>

// Tests.
package tests

import (
	"github.com/dedeme/ktlib/str"
	"github.com/dedeme/ktlib/sys"
)

func eq[T comparable](actual, expected T) {
	if actual != expected {
		panic(sys.Fail(str.Fmt("%v != %v\n", actual, expected)))
	}
}

func All() {
	BrokerTests()
	ActivityTests()
	ModelTests()
	FleaTests()
	RankingTests()
	sys.Println("All tests finised")
}

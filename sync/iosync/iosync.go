// Copyright 10-Apr-2020 ºDeme
// GNU General Public License - V3 <http://www.gnu.org/licenses/>

// To synchronize
package iosync

import (
	"fmt"
	"github.com/dedeme/sync/synchronizer"
	"time"
)

func Print(lk *synchronizer.Lock, i int) {
  synchronizer.Check(lk)
	fmt.Printf("%va\n", i)
	time.Sleep(time.Duration(2000) * time.Millisecond)
	fmt.Printf("%vb\n", i)
}

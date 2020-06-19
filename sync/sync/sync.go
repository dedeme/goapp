// Copyright 10-Apr-2020 ºDeme
// GNU General Public License - V3 <http://www.gnu.org/licenses/>

package main

import (
	"fmt"
	"github.com/dedeme/sync/synchronizer"
	"github.com/dedeme/sync/iosync"
	"time"
)

var count int

func launcher(i int) {
	go func() {
		synchronizer.Run(func(lk *synchronizer.Lock) {
      iosync.Print(lk, i)
      count++
		})
	}()
}

func main() {
  limit := 10
	for i := 0; i < limit; i++ {
		launcher(i)
	}

  for {
    if count >= limit {
      time.Sleep(time.Duration(5000) * time.Millisecond)
      break;
    }
    time.Sleep(time.Duration(5000) * time.Millisecond)
  }
  fmt.Println("end application")

}

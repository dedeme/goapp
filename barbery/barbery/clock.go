// Copyright 28-Mar-2020 ºDeme
// GNU General Public License - V3 <http://www.gnu.org/licenses/>

package main

import (
	"fmt"
	"time"
)

var clockTimeOver = false

func clockRun(ch chan *Client) {
	time.Sleep(CLOCK_TIME * time.Millisecond)
	fmt.Println("Clock: Time is over")
	clockTimeOver = true
	ch <- nil
}

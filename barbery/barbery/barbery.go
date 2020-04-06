// Copyright 28-Mar-2020 ºDeme
// GNU General Public License - V3 <http://www.gnu.org/licenses/>

package main

import (
	"fmt"
	"math/rand"
	"time"
)

const (
	CLOCK_TIME        = 40000
	CLIENT_MAKER_TIME = 5000
	BARBER_TIME       = 3000
)

func main() {
	rand.Seed(time.Now().UTC().UnixNano())

	mainCh := make(chan bool)
	barberCh := make(chan *Client)

	go clockRun(barberCh)
	go clientMakerRun(mainCh, barberCh)
	go barberRun(mainCh, barberCh)

	<-mainCh
	<-mainCh
	fmt.Println("Program end")
}

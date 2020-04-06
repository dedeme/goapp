// Copyright 28-Mar-2020 ºDeme
// GNU General Public License - V3 <http://www.gnu.org/licenses/>

package main

import (
	"math/rand"
	"time"
)

func clientMakerRun(mainCh chan bool, barberCh chan *Client) {
	var i = 0
	for shopPreopen || shopOpenV {
		i++
		go clientRun(clientMake(i), barberCh)
		time.Sleep(time.Duration(1+rand.Intn(CLIENT_MAKER_TIME)) * time.Millisecond)
	}

	mainCh <- true
}

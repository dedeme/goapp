// Copyright 28-Mar-2020 ºDeme
// GNU General Public License - V3 <http://www.gnu.org/licenses/>

package main

import (
	"fmt"
	"time"
)

func _cutHair(cl *Client) {
	shopTakeClient(cl)
	fmt.Printf("Barber: Cutting hair to %v\n", cl.id)
	time.Sleep(BARBER_TIME * time.Millisecond)
}

func barberRun(mainCh chan bool, barberCh chan *Client) {
	fmt.Println("Barber: Go out from home")
	time.Sleep(1000 * time.Millisecond)
	fmt.Println("Barber: Open barbery")
	shopOpen()

	for {
		if clockTimeOver && shopOpenV {
			fmt.Println("Barber: Close barbery")
			shopClose()
		}
		if shopIsEmpty() {
			fmt.Println("Barber: Barber is sleeping")
		}
		who := <-barberCh
		if who == nil {
			if clockTimeOver && shopOpenV {
				fmt.Println("Barber: Close barbery")
				shopClose()
			}

		lastClients:
			for {
				select {
				case cl := <-barberCh:
					_cutHair(cl)
				default:
					break lastClients
				}
			}

			time.Sleep(10 * time.Millisecond) // Wait for last client going out
			mainCh <- true
			break
		}
		_cutHair(who)
	}
}

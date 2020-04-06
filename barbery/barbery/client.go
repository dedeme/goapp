// Copyright 28-Mar-2020 ºDeme
// GNU General Public License - V3 <http://www.gnu.org/licenses/>

package main

import (
	"fmt"
)

// A Client and such
//    id : Is such
type Client struct {
	id string
}

func _msg(cl Client, m string) {
	fmt.Printf("Client %v: %v\n", cl.id, m)
}

func clientMake(id int) Client {
	return Client{fmt.Sprintf("%d", id)}
}

func clientRun(cl Client, ch chan *Client) {
	msg := func(m string) {
		fmt.Printf("Client %v: %v\n", cl.id, m)
	}

	msg("Arrive to barbery")
	switch {
	case !shopOpenV:
		msg("Go out because barbery is close")
	case shopIsFull():
		msg("Go out because barbery is full")
	default:
		{
			shopTakeASit(cl)
			ch <- &cl
		}
	}
}

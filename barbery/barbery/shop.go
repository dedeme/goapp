// Copyright 28-Mar-2020 ºDeme
// GNU General Public License - V3 <http://www.gnu.org/licenses/>

package main

import (
	"fmt"
	"math/rand"
	"strings"
)

var shopPreopen = true
var shopOpenV = false

func shopOpen() {
	shopOpenV = true
	shopPreopen = false
}

func shopClose() {
	shopOpenV = false
}

var _clients = 0
var _sits [4]*Client

func _showSits() string {
	var b strings.Builder
	b.WriteString("|")
	for i := 0; i < len(_sits); i++ {
		cl := _sits[i]
		if cl == nil {
			b.WriteString("-|")
		} else {
			fmt.Fprint(&b, _sits[i].id, "|")
		}
	}
	return b.String()
}

func shopIsFull() bool {
	return _clients == 4
}

func shopIsEmpty() bool {
	return _clients == 0
}

func shopTakeASit(cl Client) {
	if _clients >= 4 {
		panic("Barbery is full")
	}

	n := rand.Intn(4 - _clients)
	j := 0
	r := -1
	for i := 0; i < len(_sits); i++ {
		if _sits[i] == nil {
			if j == n {
				r = i
				break
			}
			j++
		}
	}
	if r == -1 {
		panic("Sit not found")
	}

	_clients++
	_sits[r] = &cl

	fmt.Printf("Shop: %v\n", _showSits())
}

func shopTakeClient(cl *Client) {
	for i := 0; i < len(_sits); i++ {
		s := _sits[i]
		if s != nil && s.id == cl.id {
			_sits[i] = nil
			_clients--
			return
		}
	}
}

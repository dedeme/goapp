// Copyright 31-May-2022 ºDeme
// GNU General Public License - V3 <http://www.gnu.org/licenses/>

package shop

import (
	"github.com/dedeme/ktlib/arr"
	"github.com/dedeme/ktlib/math"
	"github.com/dedeme/ktlib/str"
	"github.com/dedeme/ktlib/sys"
)

var IsPreopen = true
var IsOpen = false
var clients []string
var sits = make([]string, 4)

func SitsToStr() string {
	b := []string{""}
	for _, cl := range sits {
		if cl == "" {
			arr.Push(&b, "-")
		} else {
			arr.Push(&b, cl)
		}
	}
	arr.Push(&b, "")
	return arr.Join(b, "|")
}

func IsFull() bool {
	return len(clients) == 4
}

func IsEmpty() bool {
	return len(clients) == 0
}

func Open() {
	IsPreopen = false
	IsOpen = true
}

func Close() {
	IsOpen = false
}

func TakeASit(cl string) {
	nClients := len(clients)
	if nClients >= 4 {
		panic("Barbery is full")
	}

	n := math.Rndi(4 - nClients)
	j := 0
	r := -1
	for i := 0; i < len(sits); i++ {
		if sits[i] == "" {
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

	arr.Push(&clients, cl)
	sits[r] = cl

	sys.Println(str.Fmt("Shop: %s", SitsToStr()))
}

func TakeAClient() string {
	if arr.Empty(clients) {
		panic("Barbery is empty")
	}

	cl := arr.Shift(&clients)
	for i, c := range sits {
		if c == cl {
			sits[i] = ""
			break
		}
	}

	return cl
}

// Copyright 08-Apr-2020 ºDeme
// GNU General Public License - V3 <http://www.gnu.org/licenses/>

// tpls example
package _tests

import (
	"fmt"
	"github.com/dedeme/golib/json"
)

/* ·
# A record
#   R = 3 + 7
T
# Identifier
Id string
# Of your own property
car []string
-
# The real age
+age int
# Identifier
Id2 []string
# The most new
car2 map[]string
--
toJs ·· filter ·· map·IdMap·string
take ·· drop
*/
// A record
//   R = 3 + 7
type T struct {
	// Identifier
	Id string
	// Of your own property
	car []string
	// The real age
	age int
	// Identifier
	Id2 []string
	// The most new
	car2 map[string]string
}

func New(Id string, car []string) *T {
	var age int
	var Id2 []string
	var car2 map[string]string
	return &T{Id, car, age, Id2, car2}
}

// Of your own property
func (this *T) Car() []string {
	return this.car
}

// The real age
func (this *T) SetAge(value int) {
	this.age = value
}

// The real age
func (this *T) Age() int {
	return this.age
}

// The most new
func (this *T) Car2() map[string]string {
	return this.car2
}
func (this *T) ToJs() json.T {
	var _tmp []json.T
	for _, e := range this.car {
		_tmp = append(_tmp, json.Ws(e))
	}
	car := json.Wa(_tmp)
	_tmp = []json.T{}
	for _, e := range this.Id2 {
		_tmp = append(_tmp, json.Ws(e))
	}
	Id2 := json.Wa(_tmp)
	_tmp = []json.T{}
	var __tmp map[string]json.T
	for k, v := range this.car2 {
		__tmp[k] = json.Ws(v)
	}
	car2 := json.Wo(__tmp)
	__tmp = map[string]json.T{}
	return json.Wa([]json.T{
		json.Ws(this.Id),
		car,
		json.Wi(this.age),
		Id2,
		car2,
	})
}
func Filter(sl []*T, fn func(*T) bool) (rs []*T) {
	for _, e := range sl {
		if fn(e) {
			rs = append(rs, e)
		}
	}
	return
}
func IdMap(sl []*T, fn func(*T) string) (rs []string) {
	for _, e := range sl {
		rs = append(rs, fn(e))
	}
	return
}
func Take(sl []*T, n int) (rs []*T) {
	if n < 1 {
		return
	}
	if n > len(sl) {
		n = len(sl)
	}
	return sl[0:n]
}
func Drop(sl []*T, n int) (rs []*T) {
	if n < 1 {
		return sl
	}
	if n >= len(sl) {
		return
	}
	return sl[n:]
}

//===

func main() {
	fmt.Println("hello")

	x := "Hello world"

	fmt.Println(x)
}

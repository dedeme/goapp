// Copyright 13-Apr-2020 ºDeme
// GNU General Public License - V3 <http://www.gnu.org/licenses/>

package _tests

import (
	"fmt"
	"testing"
)

func TestFilter(t *testing.T) {
  var car []string
  c1 := New("a", car)
  c2 := New("b", car)
  c3 := New("ab", car)
  cs := []*T{c1, c2, c3}
  rs := Filter(cs, func (e *T) bool {
    return e.Id[0] == 'a'
  })

  //for _, e := range rs {
  //  fmt.Printf("%v\n", e.ToJs())
  //}
  if len(rs) != 2 {
    t.Fatal(fail)
  }
}

func TestMap(t *testing.T) {
  var car []string
  c1 := New("a", car)
  c2 := New("b", car)
  c3 := New("ab", car)
  cs := []*T{c1, c2, c3}
  rs := IdMap(cs, func (e *T) string {
    return e.Id
  })

  for _, e := range rs {
    fmt.Printf("%v\n", e)
  }
  if len(rs) != 3 {
    t.Fatal(fail)
  }
}

func TestTake(t *testing.T) {
  var car []string
  c1 := New("a", car)
  c2 := New("b", car)
  c3 := New("ab", car)
  cs := []*T{c1, c2, c3}
  if len(Take(cs, -2)) != 0 {
    t.Fatal(fail)
  }
  if len(Take(cs, 1)) != 1 {
    t.Fatal(fail)
  }
  if len(Take(cs, 100)) != 3 {
    t.Fatal(fail)
  }
}

func TestDrop(t *testing.T) {
  var car []string
  c1 := New("a", car)
  c2 := New("b", car)
  c3 := New("ab", car)
  cs := []*T{c1, c2, c3}
  if len(Drop(cs, -2)) != 3 {
    t.Fatal(fail)
  }
  if len(Drop(cs, 1)) != 2 {
    t.Fatal(fail)
  }
  if len(Drop(cs, 100)) != 0 {
    t.Fatal(fail)
  }
}

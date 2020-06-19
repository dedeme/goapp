// Copyright 10-Apr-2020 ºDeme
// GNU General Public License - V3 <http://www.gnu.org/licenses/>

// Synchronizer
package synchronizer

import (
  "sync"
)

var mutex sync.Mutex

type Lock struct {
  value string
}

func Check(lk *Lock) {
  if lk == nil || lk.value == "" {
    panic ("Invalid lock")
  }
}

func Run(fn func(lk *Lock)) {
  mutex.Lock()
  l := Lock{"a"}
  fn(&l)
  l.value = ""
  mutex.Unlock()
}

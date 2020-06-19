// Copyright 08-May-2020 ºDeme
// GNU General Public License - V3 <http://www.gnu.org/licenses/>

// Management of []*token.T as stack.
// NOTE: Function push is missing because you can use the standard 'append'.
package stack

import (
	"github.com/dedeme/dmstack/token"
)

// Returns the last element of 'st' and a new stack without this element.
//    Values returned:
//    newSt: 'st' without its last token.
//    tk: The last token of 'st'.
//    ok: 'false' if 'st' is empty.
func Pop (st []*token.T) (newSt []*token.T, tk *token.T, ok bool) {
  ix := len(st) - 1
  if ix >= 0 {
    newSt = st[0:ix]
    tk = st[ix]
    ok = true
  }
  return
}

// Returns the last element of 'st'.
//    Values returned:
//    tk: The last token of 'st'.
//    ok: 'false' if 'st' is empty.
func Peek (st []*token.T) (tk *token.T, ok bool) {
  ix := len(st) - 1
  if ix >= 0 {
    tk = st[ix]
    ok = true
  }
  return
}


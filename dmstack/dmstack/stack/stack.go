// Copyright 08-May-2020 ºDeme
// GNU General Public License - V3 <http://www.gnu.org/licenses/>

// Management of []*token.T as stack.
// NOTE: Function push is missing because you can use the standard 'append'.
package stack

import (
	"errors"
	"fmt"
	"github.com/dedeme/dmstack/symbol"
	"github.com/dedeme/dmstack/token"
	"strings"
)

// Returns the last element of 'st' and a new stack without this element.
//    Values returned:
//    newSt: 'st' without its last token.
//    tk: The last token of 'st'.
//    ok: 'false' if 'st' is empty.
func Pop(st []*token.T) (newSt []*token.T, tk *token.T, ok bool) {
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
func Peek(st []*token.T) (tk *token.T, ok bool) {
	ix := len(st) - 1
	if ix >= 0 {
		tk = st[ix]
		ok = true
	}
	return
}

// Returns the number of types in 'types' or an error if stack types do not
// match 'types'.
func TypesOk(st []*token.T, types string) (n int, err error) {
	errMsg := func() string {
		var b strings.Builder
		for _, tk := range st {
			b.WriteString(tk.TypeCode())
		}
		return fmt.Sprintf(
			"\n  Expected: '%s'.\n  Actual  : '%s'.",
			types, b.String(),
		)
	}

	stIx := len(st) - 1
	tIx := len(types) - 1
	if tIx == -1 {
		if stIx != -1 {
			err = errors.New(errMsg())
		}
		return
	}
	for {
		if tIx < 0 {
			n = len(st) - stIx - 1
			return
		}
		if stIx < 0 {
			err = errors.New(errMsg())
			return
		}
		tpCode := ""
		by := types[tIx]
		if by == 62 { // ">"
			tIx2 := tIx + 1
			tIx--
			for {
				if tIx < 0 {
					err = errors.New("Character '>' without '<'")
					return
				}
				if types[tIx] == 60 { // "<"
					break
				}
				tIx--
			}
			tpCode = types[tIx:tIx2]
		} else if by != 60 {
			tpCode = string(by)
		} else {
			err = errors.New("Unclosed character '<'")
			return
		}

		stCode := st[stIx].TypeCode()
		if tpCode != "*" && tpCode != stCode {
			err = fmt.Errorf(
				"\n  Expected: '...%s'.\n  Actual  : '...%s%s'.",
				types[tIx:], stCode, types[tIx+1:],
			)
			return
		}

		tIx--
		stIx--
	}
}

// Returns an error if stack types from stopt mark do not match 'types'.
func StopTypesOk(st []*token.T, types string) (n int, err error) {
	var st2 []*token.T
	st2 = nil
	for _, tk := range st {
		sym, ok := tk.Sy()
		if ok && sym == symbol.StackStop {
			st2 = []*token.T{}
			continue
		}
		if st2 != nil {
			st2 = append(st2, tk)
		}
	}
	if st2 == nil {
		err = errors.New("Stack stop not found")
		return
	}
	n, err = TypesOk(st2, types)
	if err == nil && n < len(st2) {
		err = fmt.Errorf(
			"\n  Expected: '<.= @!.>%s'.\n  Actual  : '<.= @!.>...%s%s'.",
			types, st2[len(st2)-n-1].TypeCode(), types,
		)
	}
	return
}

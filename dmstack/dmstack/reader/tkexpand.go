// Copyright 04-Jan-2021 ºDeme
// GNU General Public License - V3 <http://www.gnu.org/licenses/>

package reader

import (
	"fmt"
	"github.com/dedeme/dmstack/args"
	"github.com/dedeme/dmstack/operator"
	"github.com/dedeme/dmstack/symbol"
	"github.com/dedeme/dmstack/token"
	"strconv"
	"strings"
)

func (rd *T) processAt(opStr string, pos *token.PosT) []*token.T {
	var tks []*token.T
	if strings.HasPrefix(opStr, "?") {
		t := token.NewS(opStr[1:], pos)
		tks = append(tks, t)
		return append(tks, token.NewO(operator.StackCheck, pos))
	}

	if !args.Production {
		if strings.HasPrefix(opStr, "+") {
			rd.stackCounter++
			t := token.NewS(opStr[1:], pos)
			tks = append(tks, t)
			return append(tks, token.NewO(operator.StackOpen, pos))
		}
		if strings.HasPrefix(opStr, "-") {
			rd.stackCounter--
			t := token.NewS(opStr[1:], pos)
			tks = append(tks, t)
			return append(tks, token.NewO(operator.StackClose, pos))
		}
		t := token.NewS(opStr, pos)
		tks = append(tks, t)
		return append(tks, token.NewO(operator.Stack, pos))
	}

	return tks
}

func (rd *T) expandOperator(tk *token.T) (r []*token.T) {
	op, _ := tk.O()
	opStr := op.String()
	o0 := opStr[0]

	if o0 == '.' {
		if len(opStr) == 2 && isDigit(opStr[1]) {
			i, err := strconv.ParseUint(opStr[1:], 10, 64)
			if err != nil {
				rd.Fail(fmt.Sprintf("Expected number >= 0 after '.' (%v)", opStr))
			}
			r = append(r, token.NewI(int64(i), tk.Pos))
			r = append(r, token.NewSy(symbol.Arr, tk.Pos))
			r = append(r, token.NewO(operator.Point, tk.Pos))
			r = append(r, token.NewSy(symbol.Get, tk.Pos))
		} else if len(opStr) > 1 && opStr[1] == '.' {
			if len(opStr) > 2 {
				r = append(r, token.NewS(opStr[2:], tk.Pos))
				r = append(r, token.NewSy(symbol.Map, tk.Pos))
				r = append(r, token.NewO(operator.Point, tk.Pos))
				r = append(r, token.NewSy(symbol.Get, tk.Pos))
			} else {
				rd.Fail(fmt.Sprintf("Expected number name after '..' (%v)", opStr))
			}
		}

	} else if o0 == '@' {
		for _, t := range rd.processAt(opStr[1:], tk.Pos) {
			r = append(r, t)
		}
	} else {
		r = append(r, tk)
	}

	return
}

func (rd *T) expandArray(tk *token.T) (r []*token.T) {
	r = append(r,
		token.NewSy(symbol.Arr, tk.Pos),
		token.NewO(operator.Point, tk.Pos),
		token.NewSy(symbol.Newp, tk.Pos),
	)

	a, _ := tk.A()
	for _, tka := range a {
		r = append(r,
			token.NewSy(symbol.Dup, tka.Pos),
			tka,
			token.NewO(operator.ProcHeap, tka.Pos),
			token.NewSy(symbol.Run, tka.Pos),
			token.NewSy(symbol.Arr, tka.Pos),
			token.NewO(operator.Point, tka.Pos),
			token.NewSy(symbol.Push, tka.Pos),
		)
	}

	return
}

func (rd *T) expandMap(tk *token.T) (r []*token.T) {
	r = append(r,
		token.NewSy(symbol.Map, tk.Pos),
		token.NewO(operator.Point, tk.Pos),
		token.NewSy(symbol.Newp, tk.Pos),
	)

	m, _ := tk.M()
	for k, tkv := range m {
		r = append(r,
			token.NewSy(symbol.Dup, tkv.Pos),
			token.NewS(k, tkv.Pos),
			tkv,
			token.NewO(operator.ProcHeap, tkv.Pos),
			token.NewSy(symbol.Run, tkv.Pos),
			token.NewSy(symbol.Map, tkv.Pos),
			token.NewO(operator.Point, tkv.Pos),
			token.NewSy(symbol.Put, tkv.Pos),
		)
	}

	return
}

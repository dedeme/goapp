// Copyright 04-Jan-2021 ºDeme
// GNU General Public License - V3 <http://www.gnu.org/licenses/>

// Code reader.
package reader

import (
	"fmt"
	"github.com/dedeme/dmstack/imports"
	"github.com/dedeme/dmstack/operator"
	"github.com/dedeme/dmstack/stack"
	"github.com/dedeme/dmstack/symbol"
	"github.com/dedeme/dmstack/token"
	"os"
	"path"
)

type T struct {
	isFile       bool
	source       symbol.T //file path without extension.
	nLine        int
	prg          string
	prgIx        int
	syms         map[symbol.T]symbol.T
	stackCounter int
}

// Creates a new Reader from a file.
//    module: Symbol of file path without extension.
//    prg   : Contents of file.
func New(module symbol.T, prg string) *T {
	return &T{
		true,
		module,
		1,
		prg,
		0,
		map[symbol.T]symbol.T{},
		0,
	}
}

// Constructor used in string interpolation.
func newFromReader(rd *T, prg string, nline int) *T {
	return &T{
		false,
		rd.source,
		nline,
		prg,
		0,
		rd.syms,
		0,
	}
}

// Returns a token type Procedure after parsing 'rd.Prg()'.
func (rd *T) Process() (procedure *token.T) {
	var r []*token.T
	tk, tkOk := rd.nextToken() // in tkreader.go
	for tkOk {
		if tk.Type() == token.Symbol {
			sym, _ := tk.Sy()

			if sym == symbol.Import {
				r2, tk2, ok := stack.Pop(r)
				if !ok {
					rd.Fail("Argument of import is missing")
				}
        if tk2.Type() == token.Operator {
          o, _ := tk2.O()
          if o == operator.ProcHeap {
            r2, tk2, ok = stack.Pop(r2)
            if !ok {
              rd.Fail("Argument of import is missing")
            }
          }
        }

				kv, err := imports.ReadSymbol(tk2)
				if err != nil {
					rd.Fail(err.Error())
				}
				mod := path.Clean(
					path.Join(path.Dir(rd.source.String()), kv.Value.String()),
				)
				modSym := symbol.New(mod)
				imports.Add(modSym)

				key := kv.Key
				if key == -1 {
					key = symbol.New(path.Base(mod))
				}
				rd.syms[key] = modSym

				r = append(r2, token.NewSy(modSym, tk2.Pos), tk)

			} else {
				r = append(r, tk)
			}
		} else if tk.Type() == token.Operator {
			op, _ := tk.O()
			if op == operator.Point {
				if len(r) < 1 {
					rd.Fail("Module name is missing")
				}
				tk2 := r[len(r)-1]
				if tk2.Type() != token.Symbol {
					rd.Fail("Expected a name before '.'")
				}
				k, _ := tk2.Sy()
				v, ok := rd.syms[k]
				if ok {
					r = append(r[:len(r)-1], token.NewSy(v, tk2.Pos), tk)
				} else {
					r = append(r, tk)
				}
			} else {
				r = append(r, rd.expandOperator(tk)...) // in tkexpand
			}
		} else if tk.Type() == token.String {
			for _, t := range rd.processInterpolation(tk) { //tkstring
				r = append(r, t)
			}
		} else if tk.Type() == token.Procedure {
			r = append(r, tk, token.NewO(operator.ProcHeap, tk.Pos))
		} else if tk.Type() == token.Array {
			r = append(r, rd.expandArray(tk)...) // in tkexpand
		} else if tk.Type() == token.Map {
			r = append(r, rd.expandMap(tk)...) // in tkexpand
		} else {
			r = append(r, tk)
		}
		tk, tkOk = rd.nextToken() // in tkreader.go
	}

	if rd.stackCounter > 0 {
		rd.Fail(fmt.Sprintf("%d @+ entries without @- removes.", rd.stackCounter))
	}

	return token.NewP(r, token.NewPos(rd.source, 1))
}

// Returns "" if reader read all the text. otherwise returns the last
// character read.
func (rd *T) LastChar() string {
	if rd.prgIx >= len(rd.prg) {
		return ""
	}
	return string(rd.prg[rd.prgIx])
}

// Print a message and exit from program.
func (rd *T) Fail(msg string) {
	fmt.Printf("%v.dms:%d: Syntax error.\n  %s\n", rd.source, rd.nLine, msg)
	os.Exit(1)
}

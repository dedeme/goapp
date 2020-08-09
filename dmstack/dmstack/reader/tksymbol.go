// Copyright 30-Apr-2020 ºDeme
// GNU General Public License - V3 <http://www.gnu.org/licenses/>

package reader

import (
	"fmt"
	"github.com/dedeme/dmstack/args"
	"github.com/dedeme/dmstack/imports"
	"github.com/dedeme/dmstack/stack"
	"github.com/dedeme/dmstack/symbol"
	"github.com/dedeme/dmstack/token"
	"path"
	"strconv"
	"strings"
)

func (rd *T) processSymbol(prg []*token.T, tk *token.T) []*token.T {
	sym, _ := tk.Sy()
	if sym == symbol.New(")") ||
		sym == symbol.New("}") ||
		sym == symbol.New("}") {
		rd.fail(fmt.Sprintf("Unexpected '%v'", sym))
	}

	var tks []*token.T
	switch sym {
	// ------------------------------------------------------------------- Import
	case symbol.Import:
		tk2, ok := stack.Peek(prg)
		if !ok {
			rd.fail("Argument of import is missing")
		}

		smap, err := imports.ReadSymbol(tk2)
		if err != nil {
			rd.fail(err.Error())
		}
		mod := path.Clean(
			path.Join(path.Dir(rd.source.String()), smap.Value.String()),
		)
		key := smap.Key
		if key == -1 {
			key = symbol.New(path.Base(mod))
		}
		rd.syms = append(rd.syms, &symbol.Kv{key, symbol.New(mod)})

		return append(tks, tk)
	// --------------------------------------------------------------------- this
	case symbol.This:
		return append(tks, token.NewSy(rd.source, tk.Pos))
	default:
		symStr := sym.String()

		if strings.Contains(symStr, "{}[]()") {
			rd.fail(fmt.Sprintf("%v contains one of '{}[]()' characters", symStr))
		}

		// -------------------------------------------------------------------- Dot
		ix := strings.IndexByte(symStr, '.')
		if ix != -1 {
			left := symStr[0:ix]
			right := symStr[ix+1:]
			if left == "this" {
				ix2 := strings.IndexByte(right, '.')
				if ix2 != -1 {
					l1 := len(right) - 1
					if l1 == 0 || ix2 != l1 {
						rd.fail(fmt.Sprintf("Too many dots in '%v'", symStr))
					}
				}
				if len(right) > 0 {
					tks = append(tks, token.NewSy(rd.source, tk.Pos))
					return append(tks, token.NewSy(symbol.New(right), tk.Pos))
				} else {
					rd.fail("Syntax error: 'this.'")
				}
			} else if right != "" {
				if strings.IndexByte(right, '.') != -1 {
					rd.fail(fmt.Sprintf("Too many dots in '%v'", symStr))
				}
				if left == "" {
					tks = append(tks, token.NewS(right, tk.Pos))
					tks = append(tks, token.NewSy(symbol.Map, tk.Pos))
					return append(tks, token.NewSy(symbol.Get, tk.Pos))
				} else {
					leftSym := symbol.New(left)
					if leftSym == symbol.This {
						leftSym = rd.source
					} else {
						for _, sy := range rd.syms {
							if leftSym == sy.Key {
								leftSym = sy.Value
								break
							}
						}
					}
					tks = append(tks, token.NewSy(leftSym, tk.Pos))
					return append(tks, token.NewSy(symbol.New(right), tk.Pos))
				}
			} else {
				rd.fail(fmt.Sprintf(
					"Only 'this.xxx.' expresions can finish at dot (%v).", symStr,
				))
			}
		}

		// ------------------------------------------------------- Exclamation mark
		if strings.HasPrefix(symStr, "!") && len(symStr) > 1 {
			n, err := strconv.ParseInt(symStr[1:], 10, 64)
			if err == nil && n >= 0 {
				tks = append(tks, token.NewI(n, tk.Pos))
				tks = append(tks, token.NewSy(symbol.List, tk.Pos))
				return append(tks, token.NewSy(symbol.Get, tk.Pos))
			}
			// if dgs is not an Int constinue as normal symbol
		}

		// --------------------------------------------------------------------- At
		if strings.HasPrefix(symStr, "@") {
			if strings.HasPrefix(symStr, "@?") {
				t := token.NewS(symStr[2:], tk.Pos)
				tks = append(tks, t)
				return append(tks, token.NewSy(symbol.StackCheck, tk.Pos))
			}
			if !args.Production {
				if strings.HasPrefix(symStr, "@+") {
					rd.stackCounter++
					t := token.NewS(symStr[2:], tk.Pos)
					tks = append(tks, t)
					return append(tks, token.NewSy(symbol.StackOpen, tk.Pos))
				}
				if strings.HasPrefix(symStr, "@-") {
					rd.stackCounter--
					t := token.NewS(symStr[2:], tk.Pos)
					tks = append(tks, t)
					return append(tks, token.NewSy(symbol.StackClose, tk.Pos))
				}
				t := token.NewS(symStr[1:], tk.Pos)
				tks = append(tks, t)
				return append(tks, token.NewSy(symbol.Stack, tk.Pos))
			} else {
				return tks
			}
		}

		// ----------------------------------------------- Default in rd.syms index
		for _, sy := range rd.syms {
			if sym == sy.Key {
				return append(tks, token.NewSy(sy.Value, tk.Pos))
			}
		}

		// ---------------------------------------------------------------- Default
		return append(tks, tk)
	}
}

// Copyright 14-May-2020 ºDeme
// GNU General Public License - V3 <http://www.gnu.org/licenses/>

// Manchine runner
package runner

import (
	"github.com/dedeme/dmstack/imports"
	"github.com/dedeme/dmstack/machine"
	"github.com/dedeme/dmstack/primitives"
	"github.com/dedeme/dmstack/symbol"
	"github.com/dedeme/dmstack/token"
)

func primitiveModule(m *machine.T, sym symbol.T) bool {
	return primitives.Modules(m, sym, Run)
}

func module(m *machine.T, sym symbol.T) bool {
	if m.InImports(sym) {
		tk, ok := m.PrgNext()
		if !ok {
			m.Fail("Machine error", "Unexpected end of procedure")
		}
		sym2, ok := tk.Sy()
		if !ok {
			m.Failt(
				"\n Expected: '%v'.\n  Actual  : %v.",
				token.Symbol, tk.StringDraft(),
			)
		}
		heap, _ := imports.Get(sym)
		tk2, ok := heap[sym2]
		if !ok {
			m.Fail(
				"Machine error", "Module '%v' does not defines the symbol '%v'",
				sym, sym2,
			)
		}
		m.Push(tk2)
		if tk2.Type() == token.Procedure {
			tk, ok = m.PrgPeek()
			if ok {
				sym, ok = tk.Sy()
			}
			if ok && sym == symbol.Ampersand {
				m.PrgNext()
				return true
			}
			primitives.Global(m, symbol.Run, Run)
		}
		return true
	}
	return false
}

func heap(m *machine.T, sym symbol.T) bool {
	if tk, ok := m.HeapGet(sym); ok {
		m.Push(tk)
		if tk.Type() == token.Procedure {
			tk, ok = m.PrgPeek()
			if ok {
				sym, ok = tk.Sy()
			}
			if ok && sym == symbol.Ampersand {
				m.PrgNext()
				return true
			}
			primitives.Global(m, symbol.Run, Run)
		}
		return true
	}
	return false
}

func equals(m *machine.T, sym symbol.T) bool {
	if tk, ok := m.PrgPeek(); ok {
		sym2, _ := tk.Sy()
		if sym2 == symbol.Equals {
			m.HeapAdd(sym, m.Pop())
			m.PrgSkip()
			return true
		}
	}
	return false
}

// Runs a machine
//    m: Machine to run.
func Run(m *machine.T) {
	for {
		if tk, ok := m.PrgNext(); ok {
			if tk.Type() == token.Symbol {
				sym, _ := tk.Sy()

				if sym == symbol.Equals {
					m.Fail(
						"Machine error", "Unexpected '%v' (Possible redefinition)", sym,
					)
				}

				if ok := primitives.Global(m, sym, Run); ok {
					continue
				}

				if ok := primitiveModule(m, sym); ok {
					continue
				}

				if ok := module(m, sym); ok {
					continue
				}

				if ok := heap(m, sym); ok {
					continue
				}

				if ok := equals(m, sym); ok {
					continue
				}

				m.Fail("Machine error", "Unknown symbol '%v'", sym)
			}

			m.Push(tk)
			continue
		}

		break
	}
}

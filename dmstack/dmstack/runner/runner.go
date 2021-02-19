// Copyright 04-Jan-2021 ºDeme
// GNU General Public License - V3 <http://www.gnu.org/licenses/>

// Code runner.
package runner

import (
	"github.com/dedeme/dmstack/imports"
	"github.com/dedeme/dmstack/machine"
	"github.com/dedeme/dmstack/operator"
	"github.com/dedeme/dmstack/primitives"
	"github.com/dedeme/dmstack/symbol"
	"github.com/dedeme/dmstack/token"
)

// Module case
func module(m *machine.T, sym symbol.T) bool {
	if tk, ok := m.PrgPeek(); ok {
		if sym2, ok := tk.Sy(); ok {
			if sym2 == symbol.Import {
				m.Push(token.NewSy(sym, tk.Pos))
				return true
			}
		} else if op, ok := tk.O(); ok {
			if op == operator.Point {
				m.PrgNext()
				if tk2, ok := m.PrgNext(); ok {
					if sym2, ok := tk2.Sy(); ok {
						if primitives.Module(m, sym, sym2, Run) {
							return true
						}
						if heap, ok := imports.Get(sym); ok {
							if heap != nil {
								if tk3, ok := heap[sym2]; ok {
									if _, ok := tk3.P(); ok {
										m2 := machine.New(m.Source, m.Pmachines, tk3)
										Run(m2)
										return true
									}
									m.Push(tk3)
									return true
								}
							}
							m.Failt("Module '%v' has not procedure '%v'", sym, sym2)
						}
						m.Failt("Module '%v' not found", sym)
					}
					m.Failt(
						"\n  Expected: Procedure name."+
							"\n  Actual  : '%v'.",
						tk2.StringDraft(),
					)
				}
				m.Failt(
					"\n  Expected: Procedure name." +
						"\n  Actual  : End of code.",
				)
			}
		}
	}
	return false
}

// General symbol case
func heap(m *machine.T, sym symbol.T) bool {
	if tk, ok := m.Heap[sym]; ok {
		if _, ok := tk.P(); ok {
			m2 := machine.New(m.Source, m.Pmachines, tk)
			Run(m2)
			return true
		}
		m.Push(tk)
		return true
	}
	return false
}

// Runs a machine
//    m: Machine to run.
func Run(m *machine.T) {
	for {
		if tk, ok := m.PrgNext(); ok {
			if tk.Type() == token.Operator {
				op, _ := tk.O()
				if ok := primitives.GlobalOperator(m, op, Run); ok {
					continue
				}
				if op == operator.Point {
					m.Fail(
						machine.EMachine(),
						"Operator '.' only can be used after a module name",
					)
				}
				m.Fail(machine.EMachine(), "Unknown operator '%v'", op)
			}

			if tk.Type() == token.Symbol {
				sym, _ := tk.Sy()
				if ok := module(m, sym); ok {
					continue
				}
				if ok := primitives.GlobalSymbol(m, sym, Run); ok {
					continue
				}
				if ok := heap(m, sym); ok {
					continue
				}
				m.Fail(machine.EMachine(), "Unknown symbol '%v'", sym)
			}

			m.Push(tk)
			continue
		}
		break
	}
}

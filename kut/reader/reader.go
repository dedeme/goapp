// Copyright 01-Feb-2022 ºDeme
// GNU General Public License - V3 <http://www.gnu.org/licenses/>

// Kut text Reader.
package reader

import (
	"github.com/dedeme/kut/expression"
	"github.com/dedeme/kut/fileix"
	"github.com/dedeme/kut/heap0"
	"github.com/dedeme/kut/module"
	"github.com/dedeme/kut/modules"
	"github.com/dedeme/kut/reader/token"
	"github.com/dedeme/kut/reader/txReader"
	"github.com/dedeme/kut/statement"
)

// Translate main Kut text to a statements block.
//    tx: Reader with text and file reference
func ReadMainBlock(tx *txReader.T) (
	mod *module.T, err error,
) {
	imports := map[string]int{}
	heap := heap0.New()
	var stats []*statement.T
	var tk *token.T
	for {
		var st *statement.T
		var eof bool
		st, tk, eof, err = readStatement(tk, tx) // in stReader.go
		if err != nil {
			break
		}
		if eof {
			break
		}

		if st.Type != statement.Empty {
			if st.Type == statement.Import {
				ps := st.Value.([]interface{})
				fix := ps[0].(int)
				id := ps[1].(string)

				_, ok := imports[id]
				if ok {
					err = tx.FailLine("Import '"+id+"' already exists", st.Nline)
					break
				}
				imports[id] = fix

				_, ok = modules.Get(fix)
				if !ok {
					var kutCode string
					kutCode, err = fileix.Read(fix)
					if err != nil {
						err = tx.FailLine("Import '"+id+"':\n    "+err.Error(), st.Nline)
						break
					}
					modules.Add(fix)
					var md *module.T
					md, err = ReadMainBlock(txReader.New(fix, kutCode))
					if err != nil {
						break
					}
					modules.Set(fix, md)
				}

				continue
			}

			if st.Type == statement.CloseBlock {
				err = tx.FailLine("Unexpected '}'", st.Nline)
				break
			}

			if st.Type == statement.Assign {
				ps := st.Value.([]*expression.T)
				ex := ps[0]
				if ex.Type == expression.Sym {
					if ok := heap.Add(ex.Value.(string), st.Nline, ps[1]); !ok {
						err = tx.FailLine(
							"Duplicate assignation to symbol '"+ex.Value.(string)+"'",
							st.Nline)
						break
					}
				}
				stats = append(stats, st)
				continue
			}

			stats = append(stats, st)
		}
	}

	mod = module.New(imports, heap, stats)
	return
}

// Translate no-main Kut text to a statements block.
//    tx: Reader with text and file reference
func ReadBlock(tx *txReader.T) (
	stats []*statement.T, err error,
) {
	var tk *token.T
	for {
		var st *statement.T
		var eof bool
		st, tk, eof, err = readStatement(tk, tx) // in stReader.go
		if err != nil {
			break
		}
		if eof {
			err = tx.Fail("Unexpected end of file")
			break
		}

		if st.Type != statement.Empty {
			if st.Type == statement.Import {
				err = tx.FailLine("'import' out of main block", st.Nline)
				return
			}

			if st.Type == statement.CloseBlock {
				break
			}

			stats = append(stats, st)
		}
	}
	return
}

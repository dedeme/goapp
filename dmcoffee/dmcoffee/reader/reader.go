// Copyright 14-Jan-2021 ºDeme
// GNU General Public License - V3 <http://www.gnu.org/licenses/>

// Code reader.

package reader

import (
	"fmt"
	"github.com/dedeme/dmcoffee/fail"
	"github.com/dedeme/dmcoffee/reader/lineReader"
	"github.com/dedeme/dmcoffee/statement"
	"github.com/dedeme/dmcoffee/symbol"
	"github.com/dedeme/dmcoffee/token"
)

type T struct {
	isFile  bool
	source  symbol.T //file path without extension.
	lreader *lineReader.T
	line    *lineReader.LineT     // Current line
	ix      int                   // Position in line
	imps    map[symbol.T]symbol.T // Imports map
}

// If code came from a string, 'module' is -1.
func New(module symbol.T, code string) *T {
	isFile := true
	if module == -1 {
		isFile = false
		module = symbol.CodeEvaluation
	}
	return &T{
		isFile, module, lineReader.New(module, code),
		nil, 0, map[symbol.T]symbol.T{},
	}
}

func (rd *T) getPos() *token.PosT {
	return token.NewPos(rd.source, rd.line.N())
}

// Process the reader code and returns an array of statements.
func (rd *T) Process() []*statement.T {
	var statements []*statement.T
	var currentSt *statement.T
	for {
		if st, ok := rd.readStatement(); ok {
			if len(st.Tokens) == 0 {
				continue
			}
			if currentSt == nil {
				currentSt = st
				continue
			}
			if st.Depth == currentSt.Depth+2 {
				currentSt.Tokens = append(currentSt.Tokens, st.Tokens...)
				continue
			}
			statements = append(statements, currentSt)
			currentSt = st
			continue
		}
		break
	}
	if currentSt != nil {
		statements = append(statements, currentSt)
	}

	return statements
}

func (rd *T) Fail(template string, values ...interface{}) {
	panic(fail.New(rd.getPos(), fail.ESyntax(), fmt.Sprintf(template, values...)))
}

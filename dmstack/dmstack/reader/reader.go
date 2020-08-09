// Copyright 27-Apr-2020 ºDeme
// GNU General Public License - V3 <http://www.gnu.org/licenses/>

// Code reader.
package reader

import (
	"fmt"
	"github.com/dedeme/dmstack/symbol"
	"github.com/dedeme/dmstack/token"
	"os"
	"strings"
)

type T struct {
	isFile       bool
	source       symbol.T
	nLine        int
	prg          string
	prgIx        int
	nextsTk      []*token.T
	stackCounter int
	syms         []*symbol.Kv
}

// Creates a new Reader from a file.
//    module - File path without extension.
//    prg - Contents of file.
func New(module, prg string) *T {
	return &T{
		true,
		symbol.New(module),
		1,
		prg,
		0,
		nil,
		0,
		[]*symbol.Kv{},
	}
}

func newFromReader(rd *T, prg string, nline int) *T {
	return &T{
		false,
		rd.source,
		nline,
		prg,
		0,
		nil,
		0,
		rd.syms,
	}
}

// Returns a token type Procedure after parsing 'rd.Prg()'
func (rd *T) Process() *token.T {
	var r []*token.T
	tk, tkOk := rd.nextToken() // in tkreader.go
	for tkOk {
		if tk.Type() == token.Symbol {
			for _, t := range rd.processSymbol(r, tk) { // in tksymbol.go
				r = append(r, t)
			}
		} else if tk.Type() == token.String {
			for _, t := range rd.processInterpolation(tk) {
				r = append(r, t)
			}
		} else {
			r = append(r, tk)
		}
		tk, tkOk = rd.nextToken() // in tkreader.go -> tkstring -> tkcontainers
	}

	if rd.stackCounter > 0 {
		rd.fail(fmt.Sprintf(
			"Expected 0 'stack stops'. Actual %d.", rd.stackCounter,
		))
	}

	return token.NewP(r, token.NewPos(rd.source, rd.nLine))
}

// Process a String interpolation
func (rd *T) processInterpolation(tk *token.T) []*token.T {
	fpos := func(nline int) *token.PosT {
		return token.NewPos(tk.Pos.Source, nline)
	}

	s, _ := tk.S()
	nline := tk.Pos.Nline

	pos := 0
	var tks []*token.T
	ix := strings.Index(s, "${")
	for ix != -1 {
		tks = append(tks, token.NewS(s[pos:pos+ix], fpos(nline)))
		if pos > 0 {
			tks = append(tks, token.NewSy(symbol.Plus, fpos(nline)))
		}
		nline += strings.Count(s[pos:pos+ix], "\n")
		pos += ix + 2

		ix2 := strings.IndexByte(s[pos:], '}')
		if ix2 == -1 {
			rd.nLine = nline
			rd.fail("Interpolation not closed")
		}
		subr := newFromReader(rd, s[pos:pos+ix2], nline)
		prg, _ := subr.Process().P()
		nline += strings.Count(s[pos:pos+ix2], "\n")

		lprg := len(prg)
		if lprg == 0 {
			tks = append(tks, token.NewS("", fpos(nline)))
			tks = append(tks, token.NewSy(symbol.Plus, fpos(nline)))
		} else if len(prg) == 1 {
			tks = append(tks, token.NewL(prg, fpos(nline)))
			tks = append(tks, token.NewSy(symbol.Data, fpos(nline)))
			tks = append(tks, token.NewI(0, fpos(nline)))
			tks = append(tks, token.NewSy(symbol.List, fpos(nline)))
			tks = append(tks, token.NewSy(symbol.Get, fpos(nline)))
			tks = append(tks, token.NewSy(symbol.ToStr, fpos(nline)))
			tks = append(tks, token.NewSy(symbol.Plus, fpos(nline)))
		} else {
			rd.fail(fmt.Sprintf(
				"Interpolation '%v' yield more tha one value (%v)", subr, lprg,
			))
		}

		pos += ix2 + 1
		ix = strings.Index(s[pos:], "${")
	}

	tks = append(tks, token.NewS(s[pos:], fpos(nline)))
	if len(tks) > 1 {
		nline += strings.Count(s[pos:], "\n")
		tks = append(tks, token.NewSy(symbol.Plus, fpos(nline)))
	}

	return tks
}

// Print a message and exit from program.
func (rd *T) fail(msg string) {
	fmt.Printf("%v.dms:%d: %s\n", rd.source, rd.nLine, msg)
	os.Exit(1)
}

// Copyright 14-May-2020 ºDeme
// GNU General Public License - V3 <http://www.gnu.org/licenses/>

// Virtual machine
package machine

import (
	"fmt"
	"github.com/dedeme/dmstack/cts"
	"github.com/dedeme/dmstack/symbol"
	"github.com/dedeme/dmstack/token"
)

// Dmstack error
type Error struct {
	Machine *T
	Type    string
	Message string
}

// Machine structure
type T struct {
	SourceDir string // Path of parent to source file.
	Pmachines []*T   // Parent machines including this.
	Stack     *[]*token.T
	imports   []symbol.T
	Heap      map[symbol.T]*token.T
	prg       []*token.T
	ix        int
}

// Creates a new virtual machine.
//    sourceDir: Path of source directory.
//    pmachines : Parent machines.
//    prg       : Token list read from the procedure takes out from sourceFile.
func New(sourceDir string, pmachines []*T, prg []*token.T) *T {
	if sourceDir == "" {
		panic("sourceFile is missing")
	}

	stack := &[]*token.T{}
	imports := []symbol.T{}
	heap := map[symbol.T]*token.T{}
	l := len(pmachines) - 1
	if l >= 0 {
		stack = pmachines[l].Stack
		imports = pmachines[l].imports
	}

	m := &T{
		sourceDir,
		nil,
		stack,
		imports,
		heap,
		prg,
		0,
	}
	m.Pmachines = append(pmachines, m)

	return m
}

// Creates a new virtual machine.
//
// An isolate machine differs with a normal one on that the latter inherits
// the stack from pmachines and the former has a new stack.
//    sourceDir: Path of source directory.
//    pmachines : Parent machines.
//    prg       : Token list read from the procedure takes out from sourceFile.
func NewIsolate(sourceDir string, pmachines []*T, prg []*token.T) *T {
	if sourceDir == "" {
		panic("sourceDir is missing")
	}

	imports := []symbol.T{}
	heap := map[symbol.T]*token.T{}
	l := len(pmachines) - 1
	if l >= 0 {
		imports = pmachines[l].imports
	}

	m := &T{
		sourceDir,
		nil,
		&[]*token.T{},
		imports,
		heap,
		prg,
		0,
	}
	m.Pmachines = append(pmachines, m)

	return m
}

// Returns the stack trace of 'm'
func (m *T) StackTrace() []string {
	var r []string
	l := len(m.Pmachines) - 1
	for i := range m.Pmachines {
		mch := m.Pmachines[l-i]
		ix := mch.ix - 1
		tk := mch.prg[ix]
		pos := tk.Pos
		r = append(r, fmt.Sprintf("%v%v:%v: %v",
			pos.Source, cts.SourceExtension, pos.Nline, tk.StringDraft(),
		))
	}

	return r
}

// Returns next program token and removes it. If there are no more tokens
// returns ok = false.
func (m *T) PrgNext() (tk *token.T, ok bool) {
	if m.ix < len(m.prg) {
		tk = m.prg[m.ix]
		m.ix++
		ok = true
	}
	return
}

// Returns next program token but does not remove it. If there are no more
// tokens returns ok = false.
func (m *T) PrgPeek() (tk *token.T, ok bool) {
	if m.ix < len(m.prg) {
		tk = m.prg[m.ix]
		ok = true
	}
	return
}

// Skips next token. If there are no more tokens, produces a panic error.
func (m *T) PrgSkip() {
	if m.ix < len(m.prg) {
		m.ix++
	}
}

// Pushes a token in machine stack.
//    tk: Token to save.
func (m *T) Push(tk *token.T) {
	*m.Stack = append(*m.Stack, tk)
}

// Pops a token from machine stack.  If there are no more tokens, it raises
// a "Stack error"
func (m *T) Pop() (tk *token.T) {
	ix := len(*m.Stack) - 1
	if ix >= 0 {
		tk = (*m.Stack)[ix]
		*m.Stack = (*m.Stack)[0:ix]
	} else {
		m.Fail("Stack error", "Stack is empty")
	}
	return
}

// Pops a token of type 'tp'. If there are no more tokens or the last token
// is not of type 'tp', it raises a "Stack error".
func (m *T) PopT(tp token.TypeT) (tk *token.T) {
	tk = m.Pop()
	if tk.Type() != tp {
		m.Failt(
			"\n  Expected: '%v'.\n  Actual  : '%v'.",
			tp, tk.StringDraft(),
		)
	}
	return
}

// Adds an import. If 'sym' has already been defined, it raises an
// "Import error".
//    sym: Import symbol.
func (m *T) ImportsAdd(sym symbol.T) {
	if m.InImports(sym) {
		m.Fail("Import error", "Redefinition of import '%v'", sym)
	}
	m.imports = append(m.imports, sym)
}

func (m *T) InImports(sym symbol.T) bool {
	for _, s := range m.imports {
		if s == sym {
			return true
		}
	}
	return false
}

// Adds a symbol to heap
func (m *T) HeapAdd(sym symbol.T, tk *token.T) {
	m.Heap[sym] = tk
}

// Gets a heap token
func (m *T) HeapGet(sym symbol.T) (tk *token.T, ok bool) {
	for i := len(m.Pmachines) - 1; i >= 0; i-- {
		tk, ok = m.Pmachines[i].Heap[sym]
		if ok {
			return
		}
	}
	return
}

// Returns a position (symbol, 0) for created tokens.
func (m *T) MkPos() *token.PosT {
	return token.NewPos(m.prg[0].Pos.Source, 0)
}

// Panic with a formatted error message.
func (m *T) Fail(t, template string, values ...interface{}) {
	panic(&Error{Machine: m, Type: t, Message: fmt.Sprintf(template, values...)})
}

// "Type error" panic with a formatted error message.
func (m *T) Failt(template string, values ...interface{}) {
	panic(&Error{
		Machine: m, Type: "Type error", Message: fmt.Sprintf(template, values...),
	})
}

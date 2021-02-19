// Copyright 06-Jan-2021 ºDeme
// GNU General Public License - V3 <http://www.gnu.org/licenses/>

// Virtual machine
package machine

import (
	"fmt"
	"github.com/dedeme/dmstack/cts"
	"github.com/dedeme/dmstack/symbol"
	"github.com/dedeme/dmstack/token"
)

type TError struct {
	e string
}

func EType() *TError {
	return &TError{"Type error"}
}
func EImport() *TError {
	return &TError{"Import error"}
}
func EStack() *TError {
	return &TError{"Stack error"}
}
func EMachine() *TError {
	return &TError{"Machine error"}
}
func EAssert() *TError {
	return &TError{"Assert error"}
}
func EExpect() *TError {
	return &TError{"Expect error"}
}
func ERuntime() *TError {
	return &TError{"Runtime error"}
}
func ERange() *TError {
	return &TError{"Index out of range error"}
}
func ENotFound() *TError {
	return &TError{"Not found error"}
}
func ESys() *TError {
	return &TError{"'sys' module error"}
}
func EMath() *TError {
	return &TError{"Math error"}
}
func EJs() *TError {
	return &TError{"Json error"}
}
func EFile() *TError {
	return &TError{"File error"}
}
func ECustom(t string) *TError {
	return &TError{t}
}
func (e *TError) String() string {
	return e.e
}

// Dmstack error
type Error struct {
	Machine *T
	Type    *TError
	Message string
}

// Machine structure
type T struct {
	Source    symbol.T //file path without extension.
	Pmachines []*T     // Parent machines including this.
	Stack     *[]*token.T
	Heap      map[symbol.T]*token.T
	prg       []*token.T
	ix        int
}

// Creates a new virtual machine to use with procedures.
//    module   : Symbol of file path without extension.
//    pmachines: Parent machines.
//    prg      : Token type Procedure.
//    isThread : 'true' if 'proc' is a thread.
func new(
	module symbol.T, pmachines []*T, proc *token.T, isThread bool,
) *T {
	l := len(pmachines) - 1

	stack := &[]*token.T{}
	if l >= 0 && !isThread {
		m2 := pmachines[l]
		stack = m2.Stack
	}

	prg, ok := proc.P()
	if !ok {
		panic("Expected a procedure")
	}
	heap := proc.GetHeap()

	m := &T{
		module,
		nil,
		stack,
		heap,
		prg,
		0,
	}
	m.Pmachines = append(pmachines, m)

	return m
}

// Creates a new virtual machine to use with procedures.
//    module   : Symbol of file path without extension.
//    pmachines: Parent machines.
//    prg      : Token type Procedure.
func New(
	module symbol.T, pmachines []*T, proc *token.T,
) *T {
	return new(module, pmachines, proc, false)
}

// Creates a new virtual machine to use with procedures.
//    module   : Symbol of file path without extension.
//    pmachines: Parent machines.
//    prg      : Token type Procedure.
func NewThread(
	module symbol.T, pmachines []*T, proc *token.T,
) *T {
	return new(module, pmachines, proc, true)
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
		m.Fail(EStack(), "Stack is empty")
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

// Returns a position (symbol, 0) for created tokens.
func (m *T) MkPos() *token.PosT {
	return token.NewPos(m.prg[0].Pos.Source, 0)
}

// Panic with a formatted error message.
func (m *T) Fail(t *TError, template string, values ...interface{}) {
	panic(&Error{Machine: m, Type: t, Message: fmt.Sprintf(template, values...)})
}

// "Type error" panic with a formatted error message.
func (m *T) Failt(template string, values ...interface{}) {
	panic(&Error{
		Machine: m, Type: EType(), Message: fmt.Sprintf(template, values...),
	})
}

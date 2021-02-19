// Copyright 20-Jan-2021 ºDeme
// GNU General Public License - V3 <http://www.gnu.org/licenses/>

// Runner machine.
package machine

import (
	"fmt"
	"github.com/dedeme/dmcoffee/cts"
	"github.com/dedeme/dmcoffee/heap"
	"github.com/dedeme/dmcoffee/statement"
	"github.com/dedeme/dmcoffee/token"
	"github.com/dedeme/dmcoffee/value"
	//	"runtime/debug"
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
func EUnexpected() *TError {
	return &TError{"Unexpected error"}
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

// Dmcoffee error
type Error struct {
	Machine *T
	Type    *TError
	Message string
}

// Machine structure
type T struct {
	pmachines []*T // Parent machines including this.
	hp        heap.T
	code      []*statement.T
	depth     int
	stIx      int
	tkIx      int
}

// Creates a new virtual machine.
//    pmachines: Parent machines.
//    hp       : Heap with value symbols.
//    code     : Code to run.
func New(
	pmachines []*T, hp heap.T, code []*statement.T, depth int,
) *T {
	r := &T{pmachines, hp, code, depth, -1, -1}
	r.pmachines = append(pmachines, r)
	return r
}

// Creates an initial virtual machine.
//    code     : Code to run.
func NewInit(code []*statement.T) *T {
	return New([]*T{}, heap.New(), code, 0)
}

// Create a virtual machine from a parent machine and a token type Procedure.
//    m     : Parent machine.
//    proc  : Value type procedure.
//    values: Values for proc.
func FromProc(m *T, proc value.T, values []value.T) *T {
	php, pars, code, ok := value.P(proc)
	if !ok {
		panic("proc is not a Procedure")
	}
	if len(pars) != len(values) {
		m.Fail(
			ERuntime(), "Expected %v values in procedure call. Found $v",
			len(pars), len(values),
		)
	}
	hp := heap.From(php).Copy()
	for i, par := range pars {
		if hp.Put(par, values[i]) {
			continue
		}
		m.Fail(
			ERuntime(), "Parameter identifier '%v' is duplicate (in heap)", par,
		)
	}
	return New(m.pmachines, hp, code, m.depth+1)
}

// Returns the depth of 'm'
func (m *T) Depth() int {
	return m.depth
}

// Returns the next statement or ok=false if there is no more.
// This function advances the statement counter.
func (m *T) NextStatement() (st *statement.T, ok bool) {
	if m.stIx < len(m.code)-1 {
		ok = true
		m.stIx++
		st = m.code[m.stIx]
	}
	return
}

// Returns the next statement or ok=false if there is no more.
// This function does not advance the statement counter.
func (m *T) PeekStatement() (st *statement.T, ok bool) {
	if m.stIx < len(m.code)-1 {
		ok = true
		st = m.code[m.stIx+1]
	}
	return
}

// Returns the next token or ok=false if there is no more.
// This function advances the token counter.
func (m *T) NextToken() (tk *token.T, ok bool) {
	if st, ok2 := m.PeekStatement(); ok2 {
		tks := st.Tokens
		if m.tkIx < len(tks)-1 {
			ok = true
			m.tkIx++
			tk = tks[m.tkIx]
		}
	}
	return
}

// Returns the next token or ok=false if there is no more.
// This function does not advance the token counter.
func (m *T) PeekToken() (tk *token.T, ok bool) {
	if st, ok2 := m.PeekStatement(); ok2 {
		tks := st.Tokens
		if m.tkIx < len(tks)-1 {
			ok = true
			tk = tks[m.tkIx+1]
		}
	}
	return
}

// Returns the stack trace of 'm'
func (m *T) StackTrace() []string {
	var r []string
	l := len(m.pmachines) - 1
	fmt.Println(l)
	for i := range m.pmachines {
		mch := m.pmachines[l-i]
		stIx := mch.stIx
		tkIx := mch.tkIx
		if tkIx < 0 {
			tkIx = 0
		}
		if len(mch.code) <= stIx || stIx < 0 {
			continue
		}
		tks := mch.code[stIx].Tokens
		if len(tks) <= tkIx {
			continue
		}
		tk := tks[tkIx]
		pos := tk.Pos
		r = append(r, fmt.Sprintf("%v.%v:%v: %v",
			pos.Source, cts.SourceExtension, pos.Nline, tk.StringDraft(),
		))
	}

	return r
}

// Panic with a formatted error message.
func (m *T) Fail(t *TError, template string, values ...interface{}) {
	panic(&Error{m, t, fmt.Sprintf(template, values...)})
}

// "Type error" panic with a formatted error message.
func (m *T) Failt(template string, values ...interface{}) {
	panic(&Error{m, EType(), fmt.Sprintf(template, values...)})
}

// "Unexpected error" panic with a formatted error message.
func (m *T) Faile(expectedValue, actualValue interface{}) {
	panic(&Error{m, EUnexpected(), fmt.Sprintf(
		"Expected: %v\nActual: %v", expectedValue, actualValue)})
}

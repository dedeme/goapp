// Copyright 15-Jan-2021 ºDeme
// GNU General Public License - V3 <http://www.gnu.org/licenses/>

// Error of Reader, Runner and Evaluator
package fail

import (
	"github.com/dedeme/dmcoffee/token"
)

type Error struct {
	e string
}

func EType() *Error {
	return &Error{"Type error"}
}
func EImport() *Error {
	return &Error{"Import error"}
}
func EStack() *Error {
	return &Error{"Stack error"}
}
func EMachine() *Error {
	return &Error{"Machine error"}
}
func EAssert() *Error {
	return &Error{"Assert error"}
}
func EExpect() *Error {
	return &Error{"Expect error"}
}
func ERuntime() *Error {
	return &Error{"Runtime error"}
}
func ERange() *Error {
	return &Error{"Index out of range error"}
}
func ENotFound() *Error {
	return &Error{"Not found error"}
}
func ESys() *Error {
	return &Error{"'sys' module error"}
}
func EMath() *Error {
	return &Error{"Math error"}
}
func EJs() *Error {
	return &Error{"Json error"}
}
func EFile() *Error {
	return &Error{"File error"}
}
func ESyntax() *Error {
	return &Error{"Syntax error"}
}
func ECustom(t string) *Error {
	return &Error{t}
}
func (e *Error) String() string {
	return e.e
}

// Dmstack error
type T struct {
	Pos     *token.PosT
	Type    *Error
	Message string
}

// Constructor
func New(pos *token.PosT, errorType *Error, message string) *T {
	return &T{pos, errorType, message}
}

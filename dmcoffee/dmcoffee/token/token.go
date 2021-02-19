// Copyright 14-Jan-2021 ºDeme
// GNU General Public License - V3 <http://www.gnu.org/licenses/>

// Token data.
package token

import (
	"encoding/json"
	"fmt"
	"github.com/dedeme/dmcoffee/operator"
	"github.com/dedeme/dmcoffee/symbol"
	"strings"
)

// Token type.
type T struct {
	tk  t
	Pos *PosT
}

// Token position.
type PosT struct {
	Source symbol.T
	Nline  int
}

// Creates a new PosT
func NewPos(source symbol.T, nLine int) *PosT {
	return &PosT{source, nLine}
}

// Token type.
type TypeT int

// Token types enumeration.
const (
	Bool TypeT = iota
	Int
	Float
	String
	Symbol
	Operator
)

type t interface {
	tp() TypeT
}

type bo bool

func (tk bo) tp() TypeT {
	return Bool
}

type in int64

func (tk in) tp() TypeT {
	return Int
}

type fl float64

func (tk fl) tp() TypeT {
	return Float
}

type st string

func (tk st) tp() TypeT {
	return String
}

type sy symbol.T

func (tk sy) tp() TypeT {
	return Symbol
}

type op operator.T

func (tk op) tp() TypeT {
	return Operator
}

// Returns one of next values:
//   "Bool"
//   "Int"
//   "Float"
//   "String"
//   "Symbol"
//   "Operator"
func (t TypeT) String() string {
	switch t {
	case Bool:
		return "Bool"
	case Int:
		return "Int"
	case Float:
		return "Float"
	case String:
		return "String"
	case Symbol:
		return "Symbol"
	case Operator:
		return "Operator"
	default:
		panic(fmt.Errorf("Token type not valid (%v)", float32(t)))
	}
}

// Creates a token of type Bool.
func NewB(value bool, p *PosT) *T {
	return &T{bo(value), p}
}

// Creates a token of type Int.
func NewI(value int64, p *PosT) *T {
	return &T{in(value), p}
}

// Creates a token of type Float.
func NewF(value float64, p *PosT) *T {
	return &T{fl(value), p}
}

// Creates a token of type String.
func NewS(value string, p *PosT) *T {
	return &T{st(value), p}
}

// Creates a token of type Symbol.
func NewSy(value symbol.T, p *PosT) *T {
	return &T{sy(value), p}
}

// Creates a token of type Operator.
func NewO(value operator.T, p *PosT) *T {
	return &T{op(value), p}
}

// Returns the value of a token of type Bool.
//    If 'tk' is not of the spected type, it returns ok = false.
func (tk *T) B() (value bool, ok bool) {
	if tk.tk.tp() == Bool {
		value = bool(tk.tk.(bo))
		ok = true
	}
	return
}

// Returns the value of a token of type Int.
//    If 'tk' is not of the spected type, it returns ok = false.
func (tk *T) I() (value int64, ok bool) {
	if tk.tk.tp() == Int {
		value = int64(tk.tk.(in))
		ok = true
	}
	return
}

// Returns the value of a token of type Float.
//    If 'tk' is not of the spected type, it returns ok = false.
func (tk *T) F() (value float64, ok bool) {
	if tk.tk.tp() == Float {
		value = float64(tk.tk.(fl))
		ok = true
	}
	return
}

// Returns the value of a token of type String.
//    If 'tk' is not of the spected type, it returns ok = false.
func (tk *T) S() (value string, ok bool) {
	if tk.tk.tp() == String {
		value = string(tk.tk.(st))
		ok = true
	}
	return
}

// Returns the value of a token of type Symbol.
//    If 'tk' is not of the spected type, it returns ok = false.
func (tk *T) Sy() (value symbol.T, ok bool) {
	if tk.tk.tp() == Symbol {
		value = symbol.T(tk.tk.(sy))
		ok = true
	}
	return
}

// Returns the value of a token of type Operator.
//    If 'tk' is not of the spected type, it returns ok = false.
func (tk *T) O() (value operator.T, ok bool) {
	if tk.tk.tp() == Operator {
		value = operator.T(tk.tk.(op))
		ok = true
	}
	return
}

// Returns the type of 'tk'.
func (tk *T) Type() TypeT {
	return tk.tk.tp()
}

// Returns the following format:
//    Bool, Int, Float and String -> Its JSON respresentation.
//    Symbol -> sym.String() JSON representation
//    Operator -> sym.String() JSON representation
// If its result has more than 50 characters, it is truncated.
func (tk *T) StringDraft() string {
	s := tk.String()
	l := len(s)
	r := strings.NewReader(s)
	var b strings.Builder
	ix := 0
	irn := 0
	for {
		if ix >= l {
			return b.String()
		}
		if irn >= 50 {
			return b.String() + "..."
		}
		rn, n, _ := r.ReadRune()
		b.WriteRune(rn)
		ix += n
		irn++
	}
}

// Returns the following format:
//    Bool, Int and Float -> Its JSON respresentation.
//    String -> Its value
//    Procedure -> (element,element,...)
//    List -> [element,element,...]
//    Object -> {key:element,key:element,...}
//    Map -> [key:element,key:element,...]
//    Symbol -> "<" + sym.String() JSON representation + ">"
//    Native -> "<" + sym.String() JSON representation + "|" +
//                    PointerValue JSON representation + ">"
// Proceudure, List, Object and Map show strings in JSON format.
func (tk *T) String() string {
	switch tk.tk.tp() {
	case Bool:
		v, _ := json.Marshal(bool(tk.tk.(bo)))
		return string(v)
	case Int:
		v, _ := json.Marshal(int64(tk.tk.(in)))
		return string(v)
	case Float:
		v, _ := json.Marshal(float64(tk.tk.(fl)))
		return string(v)
	case String:
		v, _ := json.Marshal(string(tk.tk.(st)))
		return string(v)
	case Symbol:
		return symbol.T(tk.tk.(sy)).String()
	case Operator:
		return operator.T(tk.tk.(op)).String()
	default:
		panic(fmt.Errorf("Token type not valid (%v)", tk.tk.tp()))
	}
}

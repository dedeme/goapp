// Copyright 04-Jan-2021 ºDeme
// GNU General Public License - V3 <http://www.gnu.org/licenses/>

// Token data.
package token

import (
	"encoding/json"
	"fmt"
	"github.com/dedeme/dmstack/operator"
	"github.com/dedeme/dmstack/symbol"
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
	Procedure
	Array
	Map
	Symbol
	Operator
	Native
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

type pr struct {
	heap map[symbol.T]*T
	code []*T
}

func (tk pr) tp() TypeT {
	return Procedure
}

type ar []*T

func (tk ar) tp() TypeT {
	return Array
}

type ma map[string]*T

func (tk ma) tp() TypeT {
	return Map
}

type sy symbol.T

func (tk sy) tp() TypeT {
	return Symbol
}

type op operator.T

func (tk op) tp() TypeT {
	return Operator
}

type na struct {
	o     operator.T
	value interface{}
}

func (tk na) tp() TypeT {
	return Native
}

// Returns one of next values:
//   "Bool"
//   "Int"
//   "Float"
//   "String"
//   "Procedure
//   "Array"
//   "Map"
//   "Symbol"
//   "Operator"
//   "Native"
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
	case Procedure:
		return "Procedure"
	case Array:
		return "Array"
	case Map:
		return "Map"
	case Symbol:
		return "Symbol"
	case Operator:
		return "Operator"
	case Native:
		return "Native"
	default:
		panic(fmt.Errorf("Token type not valid (%v)", float32(t)))
	}
}

// Returns one of next values:
//   "Bool" -> b
//   "Int" -> i
//   "Float" -> f
//   "String" -> s
//   "Procedure -> p
//   "Array" -> a
//   "Map" -> m
//   "Native" -> <Identifier>
func (tk T) TypeCode() string {
	t := tk.Type()
	switch t {
	case Bool:
		return "b"
	case Int:
		return "i"
	case Float:
		return "f"
	case String:
		return "s"
	case Procedure:
		return "p"
	case Array:
		return "a"
	case Map:
		return "m"
	case Symbol:
		return "y"
	case Operator:
		return "o"
	case Native:
		v := tk.tk.(na)
		o := v.o
		return "<" + o.String() + ">"
	default:
		panic(fmt.Errorf("Token type not valid (%v)", t))
	}
}

// Returns the number of types codified in typeCodes (the string returned by
// 'TypeCode()'. If typeCodes is not a valid typeCode string, an error is
// returned.
func CountTypes(typeCodes string) (n int, err error) {
	rd := strings.NewReader(typeCodes)
	isLess := false
	for {
		if rd.Len() == 0 {
			return
		}
		rn, _, errr := rd.ReadRune()
		if errr != nil {
			err = fmt.Errorf("'%v' is a not valid string.", typeCodes)
			return
		}
		if isLess {
			if string(rn) == ">" {
				isLess = false
				n++
			}
			continue
		}
		if string(rn) == "<" {
			isLess = true
			continue
		}
		n++
	}
	if isLess {
		err = fmt.Errorf("'%v': '<' not closed.", typeCodes)
		return
	}
	return
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

// Creates a token of type Procedure.
func NewP(value []*T, p *PosT) *T {
	return &T{pr{map[symbol.T]*T{}, value}, p}
}

// Creates a token of type Array.
func NewA(value []*T, p *PosT) *T {
	return &T{ar(value), p}
}

// Creates a token of type Map.
func NewM(value map[string]*T, p *PosT) *T {
	return &T{ma(value), p}
}

// Creates a token of type Symbol.
func NewSy(value symbol.T, p *PosT) *T {
	return &T{sy(value), p}
}

// Creates a token of type Operator.
func NewO(value operator.T, p *PosT) *T {
	return &T{op(value), p}
}

// Creates a token of type Native.
//   NOTE: 'value' must be a pointer.
func NewN(o operator.T, value interface{}, p *PosT) *T {
	return &T{na{o, value}, p}
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

// Returns the value of a token of type Procedure.
//    If 'tk' is not of the spected type, it returns ok = false.
func (tk *T) P() (value []*T, ok bool) {
	if tk.tk.tp() == Procedure {
		value = []*T(tk.tk.(pr).code)
		ok = true
	}
	return
}

// Returns the value of a token of type Array.
//    If 'tk' is not of the spected type, it returns ok = false.
func (tk *T) A() (value []*T, ok bool) {
	if tk.tk.tp() == Array {
		value = []*T(tk.tk.(ar))
		ok = true
	}
	return
}

// Returns the value of a token of type Map.
//    If 'tk' is not of the spected type, it returns ok = false.
func (tk *T) M() (value map[string]*T, ok bool) {
	if tk.tk.tp() == Map {
		value = map[string]*T(tk.tk.(ma))
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

// Returns the value of a token of type Native.
//    If 'tk' is not of the spected type, it returns ok = false.
func (tk *T) N() (o operator.T, value interface{}, ok bool) {
	if tk.tk.tp() == Native {
		v := tk.tk.(na)
		o = v.o
		value = v.value
		ok = true
	}
	return
}

// Change an array value.
func (tk *T) SetA(value []*T) (ok bool) {
	if tk.tk.tp() == Array {
		tk.tk = ar(value)
		ok = true
	}
	return
}

// Sets heap of 'tk' with a copy of 'heap'
func (tk *T) SetHeap(heap map[symbol.T]*T) {
	if tk.tk.tp() == Procedure {
		h := tk.tk.(pr).heap
		for k := range h {
			delete(h, k)
		}
		for k, v := range heap {
			h[k] = v
		}

		return
	}
	panic("'value' is not a procedure")
}

// Returns a copy of heap
func (tk *T) GetHeap() map[symbol.T]*T {
	if tk.tk.tp() == Procedure {
		r := map[symbol.T]*T{}
		proc := tk.tk.(pr)
		for k, v := range proc.heap {
			r[k] = v
		}
		return r
	}
	panic("'value' is not a procedure")
}

// Returns the type of 'tk'.
func (tk *T) Type() TypeT {
	return tk.tk.tp()
}

// Return a new Token which is a deep copy of 'tk'.
//    If token is of type Native, the native object is sallowly copied.
func (tk *T) Clone() *T {
	switch tk.tk.tp() {
	case Bool:
		return NewB(bool(tk.tk.(bo)), tk.Pos)
	case Int:
		return NewI(int64(tk.tk.(in)), tk.Pos)
	case Float:
		return NewF(float64(tk.tk.(fl)), tk.Pos)
	case String:
		return NewS(string(tk.tk.(st)), tk.Pos)
	case Procedure:
		v := tk.tk.(pr).code
		var v2 []*T
		for _, e := range v {
			v2 = append(v2, e.Clone())
		}
		r := NewP(v2, tk.Pos)
		h := tk.tk.(pr).heap
		h2 := make(map[symbol.T]*T)
		for k, val := range h {
			h2[k] = val.Clone()
		}
		r.SetHeap(h2)
		return r
	case Array:
		v := []*T(tk.tk.(ar))
		var v2 []*T
		for _, e := range v {
			v2 = append(v2, e.Clone())
		}
		return NewA(v2, tk.Pos)
	case Map:
		v := map[string]*T(tk.tk.(ma))
		v2 := make(map[string]*T)
		for k, val := range v {
			v2[k] = val.Clone()
		}
		return NewM(v2, tk.Pos)
	case Symbol:
		return NewSy(symbol.T(tk.tk.(sy)), tk.Pos)
	case Operator:
		return NewO(operator.T(tk.tk.(op)), tk.Pos)
	case Native:
		v := tk.tk.(na)
		o := v.o
		value := v.value
		return NewN(o, value, tk.Pos)
	default:
		panic(fmt.Errorf("Token type not valid (%v)", tk.tk.tp()))
	}
}

// Returns 'true' if 'tk' is equals to 'tk2' without taking into account
// its position.
// Native types are compared by pointer.
func (tk *T) Eq(tk2 *T) bool {
	if tk.tk.tp() != tk2.tk.tp() {
		return false
	}

	switch tk.tk.tp() {
	case Bool:
		return bool(tk.tk.(bo)) == bool(tk2.tk.(bo))
	case Int:
		return int64(tk.tk.(in)) == int64(tk2.tk.(in))
	case Float:
		return float64(tk.tk.(fl)) == float64(tk2.tk.(fl))
	case String:
		return string(tk.tk.(st)) == string(tk2.tk.(st))
	case Procedure:
		v := tk.tk.(pr).code
		v2 := tk2.tk.(pr).code
		if len(v) != len(v2) {
			return false
		}
		for i := 0; i < len(v); i++ {
			if !v[i].Eq(v2[i]) {
				return false
			}
		}
		h := tk.tk.(pr).heap
		h2 := tk2.tk.(pr).heap
		if len(h) != len(h2) {
			return false
		}
		for k, val := range h {
			if val2, ok := h2[k]; !ok || !val.Eq(val2) {
				return false
			}
		}
		return true
	case Array:
		v := []*T(tk.tk.(ar))
		v2 := []*T(tk2.tk.(ar))
		if len(v) != len(v2) {
			return false
		}
		for i := 0; i < len(v); i++ {
			if !v[i].Eq(v2[i]) {
				return false
			}
		}
		return true
	case Map:
		v := map[string]*T(tk.tk.(ma))
		v2 := map[string]*T(tk2.tk.(ma))
		if len(v) != len(v2) {
			return false
		}
		for k, val := range v {
			if val2, ok := v2[k]; !ok || !val.Eq(val2) {
				return false
			}
		}
		return true
	case Symbol:
		return symbol.T(tk.tk.(sy)) == symbol.T(tk2.tk.(sy))
	case Operator:
		return operator.T(tk.tk.(op)) == operator.T(tk2.tk.(op))
	case Native:
		v := tk.tk.(na)
		o := v.o
		value := v.value
		v2 := tk2.tk.(na)
		o2 := v2.o
		value2 := v2.value
		return o == o2 && value == value2
	default:
		panic(fmt.Errorf("Token type not valid (%v)", tk.tk.tp()))
	}
}

func (tk *T) str() string {
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
	case Procedure:
		var v []string
		for _, e := range []*T(tk.tk.(pr).code) {
			v = append(v, e.str())
		}
		return "(" + strings.Join(v, " ") + ")"
	case Array:
		var v []string
		for _, e := range []*T(tk.tk.(ar)) {
			v = append(v, e.str())
		}
		return "[" + strings.Join(v, ",") + "]"
	case Map:
		var v []string
		for k, val := range map[string]*T(tk.tk.(ma)) {
			v = append(v, k+":"+val.str())
		}
		return "{" + strings.Join(v, ",") + "}"
	case Symbol:
		return "<." + symbol.T(tk.tk.(sy)).String() + ".>"
	case Operator:
		return "<." + operator.T(tk.tk.(op)).String() + ".>"
	case Native:
		v := tk.tk.(na)
		return fmt.Sprintf("<.%v|%p.>", v.o.String(), v.value)
	default:
		panic(fmt.Errorf("Token type not valid (%v)", tk.tk.tp()))
	}
}

// Returns the following format:
//    Bool, Int, Float and String -> Its JSON respresentation.
//    Procedure -> (element,element,...)
//    List -> [element,element,...]
//    Object -> {key:element,key:element,...}
//    Map -> [key:element,key:element,...]
//    Symbol -> "<" + sym.String() JSON representation + ">"
//    Native -> "<" + sym.String() JSON representation + "|" +
//                    PointerValue JSON representation + ">"
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
	if tk.Type() == String {
		r, _ := tk.S()
		return r
	}
	return tk.str()
}

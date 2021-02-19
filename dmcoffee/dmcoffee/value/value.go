// Copyright 20-Jan-2021 ºDeme
// GNU General Public License - V3 <http://www.gnu.org/licenses/>

// Value data.
package value

import (
	//	"encoding/json"
	"encoding/json"
	"fmt"
	"github.com/dedeme/dmcoffee/operator"
	"github.com/dedeme/dmcoffee/statement"
	"github.com/dedeme/dmcoffee/symbol"
	"strings"
)

// Value type.
type TypeT int

// Value types enumeration.
const (
	Bool TypeT = iota
	Int
	Float
	String
	Procedure
	Array
	Map
	Native
)

type T interface {
	tp() TypeT
}

type bo bool

func (val bo) tp() TypeT {
	return Bool
}

type in int64

func (val in) tp() TypeT {
	return Int
}

type fl float64

func (val fl) tp() TypeT {
	return Float
}

type st string

func (val st) tp() TypeT {
	return String
}

type pr struct {
	hp   map[symbol.T]T
	pars []symbol.T
	code []*statement.T
}

func (val pr) tp() TypeT {
	return Procedure
}

type ar []T

func (val ar) tp() TypeT {
	return Array
}

type ma map[string]T

func (val ma) tp() TypeT {
	return Map
}

type na struct {
	o   operator.T
	val interface{}
}

func (val na) tp() TypeT {
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
	case Native:
		return "Native"
	default:
		panic(fmt.Errorf("Value type not valid (%v)", float32(t)))
	}
}

// Creates a value of type Bool.
func NewB(val bool) T {
	return bo(val)
}

// Creates a value of type Int.
func NewI(val int64) T {
	return in(val)
}

// Creates a value of type Float.
func NewF(val float64) T {
	return fl(val)
}

// Creates a value of type String.
func NewS(val string) T {
	return st(val)
}

// Creates a value of type Procedure.
//    hp  : Parent machine heap.
//    pars: Parameter names.
//    code: Array of statements.
func NewP(hp map[symbol.T]T, pars []symbol.T, code []*statement.T) T {
	return pr{hp, pars, code}
}

// Creates a value of type Array.
func NewA(val []T) T {
	return ar(val)
}

// Creates a value of type Map.
func NewM(val map[string]T) T {
	return ma(val)
}

// Creates a value of type Native.
//   NOTE: 'val' must be a pointer.
func NewN(o operator.T, val interface{}) T {
	return na{o, val}
}

// Returns a value of type Bool.
//    If 'val' is not of the spected type, it returns ok = false.
func B(val T) (v bool, ok bool) {
	if val.tp() == Bool {
		v = bool(val.(bo))
		ok = true
	}
	return
}

// Returns a value of type Int.
//    If 'val' is not of the spected type, it returns ok = false.
func I(val T) (v int64, ok bool) {
	if val.tp() == Int {
		v = int64(val.(in))
		ok = true
	}
	return
}

// Returns the val of a value  of type Float.
//    If 'val' is not of the spected type, it returns ok = false.
func F(val T) (v float64, ok bool) {
	if val.tp() == Float {
		v = float64(val.(fl))
		ok = true
	}
	return
}

// Returns a value of type String.
//    If 'val' is not of the spected type, it returns ok = false.
func S(val T) (v string, ok bool) {
	if val.tp() == String {
		v = string(val.(st))
		ok = true
	}
	return
}

// Returns vals of a value of type Procedure.
//    If 'val' is not of the spected type, it returns ok = false.
//    Return is:
//      hp  : Heap inherited from procedure.
//      pars: Parameter names.
//      code: Array of statements.
//      ok  : 'true' if val is a Procedure.
func P(val T) (hp map[symbol.T]T, pars []symbol.T, code []*statement.T, ok bool) {
	if val.tp() == Procedure {
		v := val.(pr)
		hp = v.hp
		pars = v.pars
		code = v.code
		ok = true
	}
	return
}

// Returns the val of a value of type Array.
//    If 'val' is not of the spected type, it returns ok = false.
func A(val T) (v []T, ok bool) {
	if val.tp() == Array {
		v = []T(val.(ar))
		ok = true
	}
	return
}

// Returns the val of a value of type Map.
//    If 'val' is not of the spected type, it returns ok = false.
func M(val T) (v map[string]T, ok bool) {
	if val.tp() == Map {
		v = map[string]T(val.(ma))
		ok = true
	}
	return
}

// Returns the val of a value of type Native.
//    If 'val' is not of the spected type, it returns ok = false.
//    Return is:
//      o : Operator identifier of type.
//      v : Pointer to native value.
//      ok: 'true' if val is a Native.
func N(val T) (o operator.T, v interface{}, ok bool) {
	if val.tp() == Native {
		n := val.(na)
		o = n.o
		v = n.val
		ok = true
	}
	return
}

// Returns the type of 'tk'.
func Type(val T) TypeT {
	return val.tp()
}

// Return a new Value which is a deep copy of 'tk'.
//    If 'val' is type Native, the native object is sallowly copied.
//    if 'val' is type Procedure, hp and code are sallowly copied.
func Clone(val T) T {
	switch val.tp() {
	case Bool:
		return NewB(bool(val.(bo)))
	case Int:
		return NewI(int64(val.(in)))
	case Float:
		return NewF(float64(val.(fl)))
	case String:
		return NewS(string(val.(st)))
	case Procedure:
		v := val.(pr)
		var pars []symbol.T
		for _, e := range v.pars {
			pars = append(pars, e)
		}
		var code []*statement.T
		for _, st := range code {
			code = append(code, st)
		}
		hp := map[symbol.T]T{}
		for k, v := range v.hp {
			hp[k] = v
		}
		r := NewP(hp, pars, code)
		return r
	case Array:
		v := []T(val.(ar))
		var v2 []T
		for _, e := range v {
			v2 = append(v2, Clone(e))
		}
		return NewA(v2)
	case Map:
		v := map[string]T(val.(ma))
		v2 := make(map[string]T)
		for k, vv := range v {
			v2[k] = Clone(vv)
		}
		return NewM(v2)
	case Native:
		v := val.(na)
		o := v.o
		vv := v.val
		return NewN(o, vv)
	default:
		panic(fmt.Errorf("Value type not valid (%v)", val.tp()))
	}
}

// Returns 'true' if 'v1' is equals to 'v2'
// Native types are compared by pointer.
// Element code of Procedure types are compared by pointer.
func Eq(v1, v2 T) bool {
	if v1.tp() != v2.tp() {
		return false
	}

	switch v1.tp() {
	case Bool:
		return bool(v1.(bo)) == bool(v2.(bo))
	case Int:
		return int64(v1.(in)) == int64(v2.(in))
	case Float:
		return float64(v1.(fl)) == float64(v2.(fl))
	case String:
		return string(v1.(st)) == string(v2.(st))
	case Procedure:
		v := v1.(pr)
		vv := v2.(pr)

		pars := v.pars
		pars2 := vv.pars
		for i := 0; i < len(pars); i++ {
			if pars[i] != pars2[i] {
				return false
			}
		}

		code := v.code
		code2 := vv.code
		if len(code) != len(code2) {
			return false
		}
		for i := 0; i < len(code); i++ {
			st1 := code[i]
			st2 := code2[i]
			if st2 != st1 {
				return false
			}
		}
		h := v.hp
		h2 := vv.hp
		if len(h) != len(h2) {
			return false
		}
		for k, val := range h {
			if val2, ok := h2[k]; !ok || !Eq(val, val2) {
				return false
			}
		}
		return true
	case Array:
		v := []T(v1.(ar))
		vv := []T(v2.(ar))
		if len(v) != len(vv) {
			return false
		}
		for i := 0; i < len(v); i++ {
			if !Eq(v[i], vv[i]) {
				return false
			}
		}
		return true
	case Map:
		v := map[string]T(v1.(ma))
		vv := map[string]T(v2.(ma))
		if len(v) != len(vv) {
			return false
		}
		for k, val := range v {
			if val2, ok := vv[k]; !ok || !Eq(val, val2) {
				return false
			}
		}
		return true
	case Native:
		v := v1.(na)
		o := v.o
		val1 := v.val
		vv := v2.(na)
		o2 := vv.o
		val2 := vv.val
		return o == o2 && val1 == val2
	default:
		panic(fmt.Errorf("Value type not valid (%v)", v1.tp()))
	}
}

func str(val T) string {
	switch val.tp() {
	case Bool:
		v, _ := json.Marshal(bool(val.(bo)))
		return string(v)
	case Int:
		v, _ := json.Marshal(int64(val.(in)))
		return string(v)
	case Float:
		v, _ := json.Marshal(float64(val.(fl)))
		return string(v)
	case String:
		v, _ := json.Marshal(string(val.(st)))
		return string(v)
	case Procedure:
		v := val.(pr)
		var ps []string
		for _, e := range v.pars {
			ps = append(ps, e.String())
		}
		var sts []string
		for _, st := range v.code {
			var tks []string
			for _, e := range st.Tokens {
				tks = append(tks, e.String())
			}
			sts = append(sts, strings.Join(tks, " "))
		}
		return "(" + strings.Join(ps, " ") + ") -> {\n" +
			strings.Join(sts, "\n") + "\n}"
	case Array:
		var v []string
		for _, e := range []T(val.(ar)) {
			v = append(v, str(e))
		}
		return "[" + strings.Join(v, ",") + "]"
	case Map:
		var es []string
		for k, v := range map[string]T(val.(ma)) {
			es = append(es, k+":"+str(v))
		}
		return "{" + strings.Join(es, ",") + "}"
	case Native:
		v := val.(na)
		return fmt.Sprintf("<.%v|%p.>", v.o.String(), v.val)
	default:
		panic(fmt.Errorf("Token type not valid (%v)", val.tp()))
	}
}

// Returns the following format:
//    Bool, Int, Float and String -> Its JSON respresentation.
//    Procedure -> (param1, param2...) -> { tokens }
//    Array -> [element,element,...]
//    Map -> {key:element,key:element,...}
//    Native -> "<" + sym.String() JSON representation + "|" +
//                    PointerValue JSON representation + ">"
// If its result has more than 50 characters, it is truncated.
func StrDraft(val T) string {
	s := str(val)
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
//    Procedure -> (param1, param2...) -> { tokens }
//    Array -> [element,element,...]
//    Map -> {key:element,key:element,...}
//    Native -> "<" + sym.String() JSON representation + "|" +
//                    PointerValue JSON representation + ">"
// Proceudure, Array, and Map show strings in JSON format.
func Str(val T) string {
	if Type(val) == String {
		r, _ := S(val)
		return r
	}
	return str(val)
}

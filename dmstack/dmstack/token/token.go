// Copyright 27-Apr-2020 ºDeme
// GNU General Public License - V3 <http://www.gnu.org/licenses/>

// Token data.
package token

import (
	"encoding/json"
	"fmt"
	"github.com/dedeme/dmstack/symbol"
  "strings"
)

// Token type.
type T struct {
	tk t
	Pos *PosT
}

// Token position.
type PosT struct {
	Source symbol.T
	Nline  int
}

// Creates a new PosT
func NewPos (source symbol.T, nLine int) *PosT {
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
	List
	Map
	Symbol
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

type pr []*T

func (tk pr) tp() TypeT {
	return Procedure
}

type li []*T

func (tk li) tp() TypeT {
	return List
}

type ma map[*T]*T

func (tk ma) tp() TypeT {
	return Map
}

type sy symbol.T

func (tk sy) tp() TypeT {
	return Symbol
}

type na struct {
	sym   symbol.T
	value interface{}
}

func (tk na) tp() TypeT {
	return Native
}

// Returns one of next values:
//	 "Bool"
//	 "Int"
//	 "Float"
//	 "String"
//   "Procedure
//	 "List"
//	 "Map"
//	 "Symbol"
//	 "Native"
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
	case List:
		return "List"
	case Map:
		return "Map"
	case Symbol:
		return "Symbol"
	case Native:
		return "Native"
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

// Creates a token of type Procedure.
func NewP(value []*T, p *PosT) *T {
	return &T{pr(value), p}
}

// Creates a token of type List.
func NewL(value []*T, p *PosT) *T {
	return &T{li(value), p}
}

// Creates a token of type Map.
func NewM(value map[*T]*T, p *PosT) *T {
	return &T{ma(value), p}
}

// Creates a token of type Symbol.
func NewSy(value symbol.T, p *PosT) *T {
	return &T{sy(value), p}
}

// Creates a token of type Native.
//   NOTE: 'value' must be a pointer.
func NewN(sym symbol.T, value interface{}, p *PosT) *T {
	return &T{na{sym, value}, p}
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
		value = []*T(tk.tk.(pr))
		ok = true
	}
	return
}

// Returns the value of a token of type List.
//    If 'tk' is not of the spected type, it returns ok = false.
func (tk *T) L() (value []*T, ok bool) {
	if tk.tk.tp() == List {
		value = []*T(tk.tk.(li))
		ok = true
	}
	return
}

// Returns the value of a token of type Map.
//    If 'tk' is not of the spected type, it returns ok = false.
func (tk *T) M() (value map[*T]*T, ok bool) {
	if tk.tk.tp() == Map {
		value = map[*T]*T(tk.tk.(ma))
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

// Returns the value of a token of type Native.
//    If 'tk' is not of the spected type, it returns ok = false.
func (tk *T) N() (sym symbol.T, value interface{}, ok bool) {
	if tk.tk.tp() == Native {
		v := tk.tk.(na)
		sym = v.sym
		value = v.value
		ok = true
	}
	return
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
		v := []*T(tk.tk.(pr))
		var v2 []*T
		for _, e := range v {
			v2 = append(v2, e.Clone())
		}
		return NewP(v2, tk.Pos)
	case List:
		v := []*T(tk.tk.(li))
		var v2 []*T
		for _, e := range v {
			v2 = append(v2, e.Clone())
		}
		return NewL(v2, tk.Pos)
	case Map:
		v := map[*T]*T(tk.tk.(ma))
		v2 := make(map[*T]*T)
		for k, val := range v {
			v2[k] = val.Clone()
		}
		return NewM(v2, tk.Pos)
	case Symbol:
		return NewSy(symbol.T(tk.tk.(sy)), tk.Pos)
	case Native:
 		v := tk.tk.(na)
		sym := v.sym
		value := v.value
		return NewN(sym, value, tk.Pos)
	default:
		panic(fmt.Errorf("Token type not valid (%v)", tk.tk.tp()))

	}
}

// Returns 'true' if 'tk' is equals to 'tk2' without taking into account
// its position.
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
		v := []*T(tk.tk.(pr))
		v2 := []*T(tk2.tk.(pr))
		if len(v) != len(v2) {
			return false
		}
		for i := 0; i < len(v); i++ {
			if !v[i].Eq(v2[i]) {
				return false
			}
		}
		return true
	case List:
		v := []*T(tk.tk.(li))
		v2 := []*T(tk2.tk.(li))
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
		v := map[*T]*T(tk.tk.(ma))
		v2 := map[*T]*T(tk2.tk.(ma))
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
	case Native:
 		v := tk.tk.(na)
		sym := v.sym
		value := v.value
 		v2 := tk2.tk.(na)
		sym2 := v2.sym
		value2 := v2.value
		return sym == sym2 && value == value2
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
		for _, e := range []*T(tk.tk.(pr)) {
			v = append(v, e.str())
		}
		return "(" + strings.Join(v, ",") + ")"
	case List:
		var v []string
		for _, e := range []*T(tk.tk.(li)) {
			v = append(v, e.str())
		}
		return "[" + strings.Join(v, ",") + "]"
	case Map:
    var v []string
    for k, val := range map[*T]*T(tk.tk.(ma)) {
      v = append(v, k.str() + ":" + val.str())
    }
		return "{" + strings.Join(v, ",") + "}"
	case Symbol:
		return "<." + symbol.T(tk.tk.(sy)).String() + ".>"
	case Native:
		v := tk.tk.(na)
		return fmt.Sprintf("<.%v|%p.>", v.sym.String(), v.value)
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
  r := tk.str()
  if len(r) > 50 {
    r = ""
    for i, rn := range r {
      if i >= 47 {
        break
      }
      r += string(rn)
    }
    r += "..."
  }
  return r
}

// Returns the following format:
//    Bool, Int and Float -> Its JSON respresentation.
//    String -> Its token value.
//    Procedure -> (element,element,...)
//    List -> [element,element,...]
//    Object -> {key:element,key:element,...}
//    Map -> [key:element,key:element,...]
//    Symbol -> "<" + sym.String() JSON representation + ">"
//    Native -> "<" + sym.String() JSON representation + "|" +
//                    PointerValue JSON representation + ">"
// Proceudure, List, Object and Map show strings in JSON format.
func (tk *T) String() string {
  if tk.tk.tp() == String  {
    return string(tk.tk.(st))
  } else {
    return tk.str()
  }
}

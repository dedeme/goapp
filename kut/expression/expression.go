// Copyright 02-Feb-2022 ºDeme
// GNU General Public License - V3 <http://www.gnu.org/licenses/>

// Expression data.
package expression

import (
	"fmt"
	"strings"
)

type emptyT struct{}

func (e *emptyT) String() string {
	return "<empty expression>"
}

type Type int

const (
	Final Type = iota
	Arr
	Map
	Func
	Sym

	Range

	ExPt
	ExSq
	ExPr

	Switch

	Not
	Minus

	Add
	Sub
	Mul
	Div
	Mod

	And
	Or

	Greater
	GreaterEq
	Less
	LessEq
	Eq
	Neq

	Ternary
)

type T struct {
	Type  Type
	Value interface{}
}

func New(t Type, v interface{}) *T {
	return &T{t, v}
}

func MkFinal(v interface{}) *T {
	return &T{Final, v}
}

func MkEmpty() *T {
	return &T{Final, &emptyT{}}
}

func (ex *T) IsEmpty() (ok bool) {
	if ex.Type == Final {
		switch ex.Value.(type) {
		case *emptyT:
			ok = true
		}
	}
	return
}

func (ex *T) IsBinary() (ok bool) {
	switch ex.Type {
	case Add, Sub, Mul, Div, Mod, And, Or,
		Greater, GreaterEq, Less, LessEq, Eq, Neq:
		ok = true
	}
	return
}

func (ex *T) IsString() (ok bool) {
	if ex.Type == Final {
		_, ok = ex.Value.(string)
	}
	return
}

func (ex *T) IsFunctionCall() (ok bool) {
	switch ex.Type {
	case ExPr:
		ok = true
	case ExPt:
		ok = ex.Value.([]*T)[1].IsFunctionCall()
	}
	return
}

func (e *T) String() (s string) {
	bin := func(tx string) (s string) {
		ps := e.Value.([]*T)
		if ps[0].IsBinary() {
			s += "(" + ps[0].String() + ")"
		} else {
			s += ps[0].String()
		}
		s += tx
		if ps[1].IsBinary() {
			s += "(" + ps[1].String() + ")"
		} else {
			s += ps[1].String()
		}
		return
	}
	list := func(l []*T) (s string) {
		var ss []string
		for _, e := range l {
			ss = append(ss, e.String())
		}
		return strings.Join(ss, ",")
	}

	switch e.Type {
	case Final:
		v := e.Value
		if vs, ok := v.(string); ok {
			s = "\"" + strings.ReplaceAll(
				strings.ReplaceAll(vs, "\\", "\\\\"), "\"", "\\\"") +
				"\""
		} else if vs, ok := v.([]*T); ok {
			var ss []string
			for _, ex := range vs {
				ss = append(ss, ex.String())
			}
			s = "[" + strings.Join(ss, ", ") + "]"
		} else if vs, ok := v.(map[string]*T); ok {
			var ss []string
			for k, v := range vs {
				s = "\"" + strings.ReplaceAll(
					strings.ReplaceAll(k, "\\", "\\\\"), "\"", "\\\"") +
					"\""
				ss = append(ss, s+": "+v.String())
			}
			s = "{" + strings.Join(ss, ", ") + "}"
		} else {
			s = fmt.Sprint(e.Value)
		}
	case Arr:
		var ss []string
		for _, ex := range e.Value.([]*T) {
			ss = append(ss, ex.String())
		}
		s = "[" + strings.Join(ss, ", ") + "]"
	case Map:
		var ss []string
		for k, v := range e.Value.(map[string]*T) {
			s = "\"" + strings.ReplaceAll(
				strings.ReplaceAll(k, "\\", "\\\\"), "\"", "\\\"") +
				"\""
			ss = append(ss, s+": "+v.String())
		}
		s = "{" + strings.Join(ss, ", ") + "}"
	case Sym:
		s = e.Value.(string)
	case Func:
		s = fmt.Sprint(e.Value)
	case ExPt:
		ps := e.Value.([]*T)
		s = ps[0].String() + "." + ps[1].String()
	case ExSq:
		ps := e.Value.([]*T)
		s = ps[0].String() + "[" + ps[1].String() + "]"
	case ExPr:
		ps := e.Value.([]interface{})
		s = (ps[0].(*T)).String() + "(" + list(ps[1].([]*T)) + ")"
	case Switch:
		ps := e.Value.([]interface{})
		ex0 := ps[0].(*T)
		s = "switch (" + ex0.String() + ") {\n"
		for _, conds := range ps[1].([][]*T) {
			s += "  " + conds[0].String() + ": " + conds[1].String() + ";\n"
		}
		s += "}"
	case Range:
		ps := e.Value.([]*T)
		s = "[" + ps[0].String() + ":" + ps[1].String() + "]"
	case Not:
		s = "!" + e.Value.(*T).String()
	case Minus:
		s = "-" + e.Value.(*T).String()
	case Add:
		s = bin("+")
	case Sub:
		s = bin("-")
	case Mul:
		s = bin("*")
	case Div:
		s = bin("/")
	case Mod:
		s = bin("%")
	case And:
		s = bin("&")
	case Or:
		s = bin("|")
	case Greater:
		s = bin(">")
	case GreaterEq:
		s = bin(">=")
	case Less:
		s = bin("<")
	case LessEq:
		s = bin("<=")
	case Eq:
		s = bin("==")
	case Neq:
		s = bin("!=")
	case Ternary:
		ps := e.Value.([]*T)
		s = ps[0].String() + " ? " + ps[1].String() + " : " + ps[2].String()
	default:
		panic("Unknown expresion type")
	}

	if len(s) > 70 {
		s = s[:67] + "..."
	}
	return
}

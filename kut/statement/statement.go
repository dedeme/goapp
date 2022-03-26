// Copyright 02-Feb-2022 ºDeme
// GNU General Public License - V3 <http://www.gnu.org/licenses/>

// Statement data.
package statement

import (
	"github.com/dedeme/kut/expression"
	"github.com/dedeme/kut/fileix"
)

type Type int

const (
	Empty Type = iota
	Assign
	AddAs
	SubAs
	MulAs
	DivAs
	AndAs
	OrAs
	FunctionCalling
	Block
	CloseBlock

	Break
	Continue
	Trace
	Assert
	Return

	While
	If
	Else
	For
	ForIx
	ForR
	Switch

	Import
)

type T struct {
	File  int
	Nline int
	Type  Type
	Value interface{}
}

func New(file, nline int, t Type, v interface{}) *T {
	return &T{file, nline, t, v}
}

func NewBuilt(t Type, v interface{}) *T {
	return &T{-1, 0, t, v}
}

func (st *T) String() (s string) {
	switch st.Type {
	case Empty:
		s = ";"
	case Assign:
		parts := st.Value.([]*expression.T)
		s = parts[0].String() + " = " + parts[1].String() + ";"
	case AddAs:
		parts := st.Value.([]*expression.T)
		s = parts[0].String() + " += " + parts[1].String() + ";"
	case SubAs:
		parts := st.Value.([]*expression.T)
		s = parts[0].String() + " -= " + parts[1].String() + ";"
	case MulAs:
		parts := st.Value.([]*expression.T)
		s = parts[0].String() + " *= " + parts[1].String() + ";"
	case DivAs:
		parts := st.Value.([]*expression.T)
		s = parts[0].String() + " /= " + parts[1].String() + ";"
	case AndAs:
		parts := st.Value.([]*expression.T)
		s = parts[0].String() + " &= " + parts[1].String() + ";"
	case OrAs:
		parts := st.Value.([]*expression.T)
		s = parts[0].String() + " |= " + parts[1].String() + ";"
	case FunctionCalling:
		s = (st.Value.(*expression.T)).String() + ";"
	case Block:
		s = "{"
		for _, stat := range st.Value.([]*T) {
			s += stat.String()
		}
		s += "}"
	case CloseBlock:
		s = "}"
	case Break:
		s = "break;"
	case Continue:
		s = "continue;"
	case Trace:
		s = "trace " + (st.Value.(*expression.T)).String() + ";"
	case Assert:
		s = "assert " + (st.Value.(*expression.T)).String() + ";"
	case Return:
		if st.Value == nil {
			s = "return;"
		} else {
			s = "return " + (st.Value.(*expression.T)).String() + ";"
		}
	case While:
		ps := st.Value.([]interface{})
		cond := ""
		if ps[0].(*expression.T) != nil {
			cond = (ps[0].(*expression.T)).String()
		}
		s = "while (" + cond + ") " + (ps[1].(*T)).String()
	case If:
		ps := st.Value.([]interface{})
		s = "if (" + ps[0].(*expression.T).String() + ") " + (ps[1].(*T)).String()
		if ps[2] != nil {
			s += "\nelse " + (ps[2].(*T)).String()
		}
	case Else:
		s = "else " + (st.Value.(*T)).String()
	case For:
		ps := st.Value.([]interface{})
		s = "for (" + ps[0].(string) + " : " +
			(ps[1].(*expression.T)).String() + ") " +
			(ps[2].(*T)).String()
	case ForIx:
		ps := st.Value.([]interface{})
		s = "for (" + ps[0].(string) + ", " + ps[1].(string) + " : " +
			(ps[2].(*expression.T)).String() + ") " +
			(ps[3].(*T)).String()
	case ForR:
		ps := st.Value.([]interface{})
		s = "for (" + ps[0].(string) + " : " +
			(ps[1].(*expression.T)).String() + ":" +
			(ps[2].(*expression.T)).String() + ") " +
			(ps[3].(*T)).String()
	case Switch:
		ps := st.Value.([]interface{})
		s = "switch (" + (ps[0].(*expression.T)).String() + ") {\n"
		for _, cond := range ps[1].([][]interface{}) {
			s += "  " + cond[0].(*expression.T).String() + ": " +
				cond[1].(*T).String() + "\n"
		}
		s += "}"
	case Import:
		ps := st.Value.([]interface{})
		s = "import " + fileix.Get(ps[0].(int)) + " : " + ps[1].(string) + ";"
	default:
		panic("Unknown statement type")
	}

	if len(s) > 70 {
		s = s[:67] + "..."
	}
	return
}

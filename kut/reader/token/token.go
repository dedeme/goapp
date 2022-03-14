// Copyright 01-Feb-2022 ºDeme
// GNU General Public License - V3 <http://www.gnu.org/licenses/>

// Token data
package token

import "fmt"

type Type int

// Token types
const (
	Bool Type = iota
	Int
	Float
	String
	LineComment
	Comment
	Symbol
	Operator
)

func (t Type) String() string {
	switch t {
	case Bool:
		return "Bool"
	case Int:
		return "Int"
	case Float:
		return "Float"
	case String:
		return "String"
	case LineComment:
		return "LineComment"
	case Comment:
		return "Comment"
	case Symbol:
		return "Symbol"
	case Operator:
		return "Operator"
	}
	return "Unknown"
}

// Token data
type T struct {
	Type  Type
	Value interface{}
}

// Constructor.
func New(tokenType Type, value interface{}) *T {
	return &T{tokenType, value}
}

func (tk *T) IsUnary() (ok bool) {
	if tk.Type == Operator {
		switch tk.Value.(string) {
		case "!", "-":
			ok = true
		}
	}
	return
}

func (tk *T) IsBinary() (ok bool) {
	if tk.Type == Operator {
		switch tk.Value.(string) {
		case "*", "/", "%", "+", "-",
			"==", "!=", ">", ">=", "<", "<=", "&", "|":
			ok = true
		}
	}
	return
}

func (tk *T) IsBinary1() (ok bool) {
	if tk.Type == Operator {
		switch tk.Value.(string) {
		case "*", "/", "%":
			ok = true
		}
	}
	return
}

func (tk *T) IsBinary2() (ok bool) {
	if tk.Type == Operator {
		switch tk.Value.(string) {
		case "+", "-":
			ok = true
		}
	}
	return
}

func (tk *T) IsBinary3() (ok bool) {
	if tk.Type == Operator {
		switch tk.Value.(string) {
		case "==", "!=", ">", ">=", "<", "<=":
			ok = true
		}
	}
	return
}

func (tk *T) IsBinary4() (ok bool) {
	if tk.Type == Operator {
		switch tk.Value.(string) {
		case "&", "|":
			ok = true
		}
	}
	return
}

func (tk *T) IsTernary() bool {
	if tk.Type == Operator && tk.Value.(string) == "?" {
		return true
	}
	return false
}

func (tk *T) IsAsign() (ok bool) {
	if tk.Type == Operator {
		switch tk.Value.(string) {
		case "=", "+=", "->", "*=", "/=", "%=", "|=", "&=":
			ok = true
		}
	}
	return
}

func (tk *T) IsEquals() bool {
	if tk.Type == Operator && tk.Value.(string) == "=" {
		return true
	}
	return false
}

func (tk *T) IsComma() bool {
	if tk.Type == Operator && tk.Value.(string) == "," {
		return true
	}
	return false
}

func (tk *T) IsColon() bool {
	if tk.Type == Operator && tk.Value.(string) == ":" {
		return true
	}
	return false
}

func (tk *T) IsSemicolon() bool {
	if tk.Type == Operator && tk.Value.(string) == ";" {
		return true
	}
	return false
}

func (tk *T) IsBackSlash() bool {
	if tk.Type == Operator && tk.Value.(string) == "\\" {
		return true
	}
	return false
}

func (tk *T) IsArrow() bool {
	if tk.Type == Operator && tk.Value.(string) == "->" {
		return true
	}
	return false
}

func (tk *T) IsOpenPar() bool {
	if tk.Type == Operator && tk.Value.(string) == "(" {
		return true
	}
	return false
}

func (tk *T) IsClosePar() bool {
	if tk.Type == Operator && tk.Value.(string) == ")" {
		return true
	}
	return false
}

func (tk *T) IsOpenSquare() bool {
	if tk.Type == Operator && tk.Value.(string) == "[" {
		return true
	}
	return false
}

func (tk *T) IsCloseSquare() bool {
	if tk.Type == Operator && tk.Value.(string) == "]" {
		return true
	}
	return false
}

func (tk *T) IsOpenBracket() bool {
	if tk.Type == Operator && tk.Value.(string) == "{" {
		return true
	}
	return false
}

func (tk *T) IsCloseBracket() bool {
	if tk.Type == Operator && tk.Value.(string) == "}" {
		return true
	}
	return false
}

func (tk *T) IsElse() bool {
	if tk.Type == Symbol && tk.Value.(string) == "else" {
		return true
	}
	return false
}

func (t *T) String() string {
	return fmt.Sprintf("%v: %v", t.Type, t.Value)
}

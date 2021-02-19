// Copyright 14-Jan-2021 ºDeme
// GNU General Public License - V3 <http://www.gnu.org/licenses/>

// operator indexer.
package operator

type T int32

const (
	And          = T(iota) // &&
	Assign                 // =
	Ccurly                 // }
	Colon                  // :
	Comma                  // ,
	Cparenthesis           // )
	Csquare                // ]
	Div                    // /
	Eq                     // ==
	Function               // ->
	Greater                // >
	GreaterEq              // >=
	Less                   // <
	LessEq                 // <=
	Minus                  // -
	MinusMinus             // --
	Mod                    // %
	Mul                    // *
	Neq                    // !=
	Not                    // !
	Ocurly                 // {
	Oparenthesis           // (
	Osquare                // [
	Or                     // ||
	Plus                   // +
	PlusPlus               // ++
	Point                  // .
	Power                  // ^
	Question               // ?
	RefDiv                 // </
	RefGet                 // >>
	RefMinus               // <-
	RefMod                 // <%
	RefMul                 // <*
	RefPlus                // <+
	RefPower               // <^
	RefSet                 // <<
	RefUp                  // ^^
	Semicolon              // ;
	Stack                  // @

	Blob
	Thread
	Iterator
	Js
	Date
	File
	Iserver
	IserverRq
	Exc

	SystemCount
)

var ops []string

func Initialize() {
	var tmp [SystemCount]string
	tmp[And] = "&&"
	tmp[Assign] = "="
	tmp[Ccurly] = "}"
	tmp[Colon] = ":"
	tmp[Comma] = ","
	tmp[Cparenthesis] = ")"
	tmp[Csquare] = "]"
	tmp[Div] = "/"
	tmp[Eq] = "=="
	tmp[Function] = "->"
	tmp[Greater] = ">"
	tmp[GreaterEq] = ">="
	tmp[Less] = "<"
	tmp[LessEq] = "<="
	tmp[Minus] = "-"
	tmp[MinusMinus] = "--"
	tmp[Mod] = "%"
	tmp[Mul] = "*"
	tmp[Neq] = "!="
	tmp[Not] = "!"
	tmp[Ocurly] = "{"
	tmp[Oparenthesis] = "("
	tmp[Osquare] = "["
	tmp[Or] = "||"
	tmp[Plus] = "+"
	tmp[PlusPlus] = "++"
	tmp[Point] = "."
	tmp[Power] = "^"
	tmp[Question] = "?"
	tmp[RefDiv] = "</"
	tmp[RefGet] = ">>"
	tmp[RefMinus] = "<-"
	tmp[RefMod] = "<%"
	tmp[RefMul] = "<*"
	tmp[RefPlus] = "<+"
	tmp[RefPower] = "<^"
	tmp[RefSet] = "<<"
	tmp[RefUp] = "^^"
	tmp[Semicolon] = ";"
	tmp[Stack] = "@"

	tmp[Blob] = "Blob"
	tmp[Thread] = "Thread"
	tmp[Iterator] = "Iterator"
	tmp[Js] = "Js"
	tmp[Date] = "Date"
	tmp[File] = "File"
	tmp[Iserver] = "Iserver"
	tmp[IserverRq] = "IserverRq"
	tmp[Exc] = "Exc"

	ops = tmp[0 : SystemCount-1]
}

// Creates a new Operator for 's'. If 's' has already an operator, it returns
// it.
//    s : String to make the operator.
func New(s string) T {
	i := -1
	for ix, op := range ops {
		if op == s {
			i = ix
			break
		}
	}
	if i == -1 {
		i = len(ops)
		ops = append(ops, s)
	}
	return T(i)
}

// Returns the string pointed by 'sym'.
func (op T) String() string {
	return ops[op]
}

// Returns true if 'op' is a binary operator and its precedence.
func (op T) Binary() (prec int, ok bool) {
	ok = true
	switch op {
	case Or:
		prec = 0
	case And:
		prec = 1
	case Eq, Neq:
		prec = 2
	case Greater, GreaterEq, Less, LessEq:
		prec = 3
	case Plus, Minus:
		prec = 4
	case Mul, Div, Mod:
		prec = 5
	case Power:
		prec = 6
	default:
		ok = false
	}
	return
}

// Returns the operator matching ch or 'ok=false' if none matchs.
//    ch: Byte to compare.
func GetO1(ch byte) (op T, ok bool) {
	ok = true
	switch ch {
	case '=':
		op = Assign
	case '!':
		op = Not
	case '>':
		op = Greater
	case '<':
		op = Less
	case '^':
		op = Power
	case '+':
		op = Plus
	case '-':
		op = Minus
	case '(':
		op = Oparenthesis
	case ')':
		op = Cparenthesis
	case '[':
		op = Osquare
	case ']':
		op = Csquare
	case '{':
		op = Ocurly
	case '}':
		op = Ccurly
	case '.':
		op = Point
	case ',':
		op = Comma
	case ':':
		op = Colon
	case ';':
		op = Semicolon
	case '/':
		op = Div
	case '%':
		op = Mod
	case '*':
		op = Mul
	case '@':
		op = Stack
	case '?':
		op = Question
	default:
		ok = false
	}
	return
}

// Returns the operator matching ch1 and ch2 or 'ok=false' if none matchs.
//    ch1: First byte to compare.
//    ch2: Second byte or 0 if the operator is of only one byte.
func GetO2(ch1, ch2 byte) (op T, ok bool) {
	ok = true
	switch ch1 {
	case '&':
		switch ch2 {
		case '&':
			op = And
		default:
			ok = false
		}
	case '|':
		switch ch2 {
		case '|':
			op = Or
		default:
			ok = false
		}
	case '=':
		switch ch2 {
		case '=':
			op = Eq
		default:
			ok = false
		}
	case '!':
		switch ch2 {
		case '=':
			op = Neq
		default:
			ok = false
		}
	case '>':
		switch ch2 {
		case '=':
			op = GreaterEq
		case '>':
			op = RefGet
		default:
			ok = false
		}
	case '<':
		switch ch2 {
		case '=':
			op = LessEq
		case '<':
			op = RefSet
		case '+':
			op = RefPlus
		case '-':
			op = RefMinus
		case '*':
			op = RefMul
		case '/':
			op = RefDiv
		case '%':
			op = RefMod
		case '^':
			op = RefPower
		default:
			ok = false
		}
	case '^':
		switch ch2 {
		case '^':
			op = RefUp
		default:
			ok = false
		}
	case '+':
		switch ch2 {
		case '+':
			op = PlusPlus
		default:
			ok = false
		}
	case '-':
		switch ch2 {
		case '-':
			op = MinusMinus
		case '>':
			op = Function
		default:
			ok = false
		}
	default:
		ok = false
	}
	return
}

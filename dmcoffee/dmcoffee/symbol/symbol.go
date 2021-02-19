// Copyright 14-Jan-2021 ºDeme
// GNU General Public License - V3 <http://www.gnu.org/licenses/>

// Symbol indexer.
package symbol

type T int32

const (
	Assert = T(iota)
	Break
	Clone
	Continue
	Dup
	Elif
	Else
	Expect
	Eval
	Fail
	For
	If
	Import
	Loop
	Mrun
	Pop
	Puts
	Return
	Run
	Swap
	Sync
	This
	Throw
	ToString
	Try
	While

	systemReserved

	Arr
	B64
	Blob
	Cryp
	Date
	Exc
	File
	Float
	Int
	It
	Js
	Map
	Math
	Path
	Str
	Sys

	CodeEvaluation
	FromJs
	From
	Get
	Newp
	Push
	Put
	ToJs

	systemCount
)

var sharp int32
var syms []string

func Initialize() {
	var tmp [systemCount]string
	tmp[Assert] = "assert"
	tmp[Break] = "break"
	tmp[Clone] = "clone"
	tmp[Continue] = "continue"
	tmp[Dup] = "dup"
	tmp[Elif] = "elif"
	tmp[Else] = "else"
	tmp[Eval] = "eval"
	tmp[Expect] = "expect"
	tmp[Fail] = "fail"
	tmp[For] = "for"
	tmp[If] = "if"
	tmp[Import] = "import"
	tmp[Loop] = "loop"
	tmp[Mrun] = "mrun"
	tmp[Puts] = "puts"
	tmp[Pop] = "pop"
	tmp[Return] = "return"
	tmp[Run] = "run"
	tmp[Swap] = "swap"
	tmp[Sync] = "sync"
	tmp[This] = "this"
	tmp[Throw] = "throw"
	tmp[ToString] = "toString"
	tmp[Try] = "try"
	tmp[While] = "while"

	tmp[Arr] = "arr"
	tmp[B64] = "b64"
	tmp[Blob] = "blob"
	tmp[Cryp] = "cryp"
	tmp[Date] = "date"
	tmp[Exc] = "exc"
	tmp[File] = "file"
	tmp[Float] = "float"
	tmp[Int] = "int"
	tmp[It] = "it"
	tmp[Js] = "js"
	tmp[Map] = "map"
	tmp[Math] = "math"
	tmp[Path] = "path"
	tmp[Str] = "str"
	tmp[Sys] = "sys"

	tmp[CodeEvaluation] = "Code evaluation"
	tmp[FromJs] = "fromJs"
	tmp[From] = "from"
	tmp[Get] = "get"
	tmp[Newp] = "new"
	tmp[Push] = "push"
	tmp[Put] = "put"
	tmp[ToJs] = "toJs"

	syms = tmp[0 : systemCount-1]
}

// Creates a new Symbol for 's'. If 's' has already a symbol, it returns it.
//    s : String to make symbol.
func New(s string) T {
	i := -1
	for ix, sym := range syms {
		if sym == s {
			i = ix
			break
		}
	}
	if i == -1 {
		i = len(syms)
		syms = append(syms, s)
	}
	return T(i)
}

// Returns the string pointed by 'sym'.
func (sym T) String() string {
	return syms[sym]
}

// Returns 'true' if 'sym' is reserved.
func (sym T) Reserved() bool {
	return sym < systemReserved
}

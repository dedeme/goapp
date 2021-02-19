// Copyright 03-Jan-2021 ºDeme
// GNU General Public License - V3 <http://www.gnu.org/licenses/>

// operator indexer.
package operator

type T int32

type Kv struct {
	Key   T
	Value T
}

const (
	Ampersand  = T(iota) // &
	And                  // &&
	Assign               // :
	Div                  // /
	Eq                   // ==
	Equals               // =
	Greater              // >
	GreaterEq            // >=
	Less                 // <
	LessEq               // <=
	Minus                // -
	MinusMinus           // --
	Mod                  // %
	Mult                 // *
	Neq                  // !=
	Not                  // !
	Or                   // ||
	Plus                 // +
	PlusPlus             // ++
	Point                // .
	ProcHeap             // heap
	RefGet               // >>
	RefSet               // >>
	RefUp                // ^^
	Stack                // @
	StackCheck           // @?
	StackClose           // @-
	StackOpen            // @+
	StackStop            // @!

	Blob_
	Thread_
	Iterator_
	Js_
	Date_
	File_
	Iserver_
	IserverRq_
	Exc_

	SystemCount
)

var syms []string

func Initialize() {
	var tmp [SystemCount]string
	tmp[Ampersand] = "&"
	tmp[And] = "&&"
	tmp[Assign] = ":"
	tmp[Div] = "/"
	tmp[Eq] = "=="
	tmp[Equals] = "="
	tmp[Greater] = ">"
	tmp[GreaterEq] = ">="
	tmp[Less] = "<"
	tmp[LessEq] = "<="
	tmp[Minus] = "-"
	tmp[MinusMinus] = "--"
	tmp[Mod] = "%"
	tmp[Mult] = "*"
	tmp[Neq] = "!="
	tmp[Not] = "!"
	tmp[Or] = "||"
	tmp[Plus] = "+"
	tmp[PlusPlus] = "++"
	tmp[Point] = "."
	tmp[ProcHeap] = "= heap"
	tmp[RefGet] = ">>"
	tmp[RefSet] = "<<"
	tmp[RefUp] = "^^"
	tmp[Stack] = "= @"
	tmp[StackCheck] = "= @?"
	tmp[StackOpen] = "= @+"
	tmp[StackClose] = "= @-"
	tmp[StackStop] = "= @!"

	tmp[Blob_] = "=Blob"
	tmp[Thread_] = "=Thread"
	tmp[Iterator_] = "=Iterator"
	tmp[Js_] = "=Js"
	tmp[Date_] = "=Date"
	tmp[File_] = "=File"
	tmp[Iserver_] = "=Iserver"
	tmp[IserverRq_] = "=IserverRq"
	tmp[Exc_] = "=Exc"

	syms = tmp[0 : SystemCount-1]
}

// Creates a new Symbol for 's'. If 's' has already a symbol, it returns this.
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

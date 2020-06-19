// Copyright 26-Apr-2020 ºDeme
// GNU General Public License - V3 <http://www.gnu.org/licenses/>

// Symbol indexer.
package symbol

import (
"fmt"
)

type T int32

type Kv struct {
  Key T
  Value T
}

const (
  Ampersand = T(iota) // &
  Assert
  Break
  Continue
  Data
  Div // /
  Elif
  Else
  Eq // ==
  Equals // =
  Expect
  Eval
  For
  Function // =>
  Greater // >
  GreaterEq // >=
  If
  Import
  Less // <
  LessEq // <=
  Loop
  Minus // -
  Mrun
  Mult // *
  Nop
  Plus // +
  Puts
  Recursive
  Run
  Stack // @
  StackCheck // @?
  StackClose // @-
  StackOpen // @+
  StackStop // @!
  Sync
  This
  ToStr
  While

  B64
  Blob
  Clock
  Cryp
  Exc
  File
  Float
  Int
  Iserver
  It
  Js
  Lst
  Map
  Math
  Path
  Stk
  Str
  Sys
  Time
  Wrap

  FromJs
  FromList
  Get
  ToJs

  Ref_
  Option_
  Either_
  Map_
  Blob_
  Thread_
  Iterator_
  File_
  Iserver_
  IserverRq_
  Exc_

  SystemCount
)

var sharp int32
var syms []string

func Initialize () {
  var tmp [SystemCount]string
  tmp[Ampersand] = "&"
  tmp[Assert] = "assert"
  tmp[Break] = "break"
  tmp[Continue] = "continue"
  tmp[Data] = "data"
  tmp[Div] = "/"
  tmp[Elif] = "elif"
  tmp[Else] = "else"
  tmp[Eq] = "=="
  tmp[Equals] = "="
  tmp[Eval] = "eval"
  tmp[Expect] = "expect"
  tmp[For] = "for"
  tmp[Function] = "=>"
  tmp[Greater] = ">"
  tmp[GreaterEq] = ">="
  tmp[If] = "if"
  tmp[Import] = "import"
  tmp[Less] = "<"
  tmp[LessEq] = "<="
  tmp[Loop] = "loop"
  tmp[Minus] = "-"
  tmp[Mrun] = "mrun"
  tmp[Mult] = "*"
  tmp[Nop] = "nop"
  tmp[Plus] = "+"
  tmp[Puts] = "puts"
  tmp[Recursive] = "recursive"
  tmp[Run] = "run"
  tmp[Stack] = "= @"
  tmp[StackCheck] = "= @?"
  tmp[StackOpen] = "= @+"
  tmp[StackClose] = "= @-"
  tmp[StackStop] = "= @!"
  tmp[Sync] = "sync"
  tmp[This] = "this"
  tmp[ToStr] = "toStr"
  tmp[While] = "while"

  tmp[B64] = "b64"
  tmp[Blob] = "blob"
  tmp[Clock] = "clock"
  tmp[Cryp] = "cryp"
  tmp[Exc] = "exc"
  tmp[File] = "file"
  tmp[Float] = "float"
  tmp[Int] = "int"
  tmp[Iserver] = "iserver"
  tmp[It] = "it"
  tmp[Js] = "js"
  tmp[Lst] = "lst"
  tmp[Map] = "map"
  tmp[Math] = "math"
  tmp[Path] = "path"
  tmp[Stk] = "stk"
  tmp[Str] = "str"
  tmp[Sys] = "sys"
  tmp[Time] = "time"
  tmp[Wrap] = "wrap"

  tmp[FromJs] = "fromJs"
  tmp[FromList] = "fromList"
  tmp[Get] = "get"
  tmp[ToJs] = "toJs"

  tmp[Ref_] = "= Ref"
  tmp[Option_] = "= Option"
  tmp[Either_] = "= Either"
  tmp[Map_] = "= Map"
  tmp[Blob_] = "= Blob"
  tmp[Thread_] = "= Thread"
  tmp[Iterator_] = "= Iterator"
  tmp[File_] = "= File"
  tmp[Iserver_] = "= Iserver"
  tmp[IserverRq_] = "= IserverRq"
  tmp[Exc_] = "= Exc"

  syms = tmp[0:SystemCount-1]
}

// Creates a new Symbol for 's'. If 's' has already a symbol, it returns this.
//    s : String to make symbol.
func New(s string) T {
  if s == "#" {
    sharp++
    return Nop
  }

  if s[len(s) - 1] == '#' {
    s = fmt.Sprintf("%v·%v", s, sharp)
  }

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

// Copyright 29-Sep-2020 ºDeme
// GNU General Public License - V3 <http://www.gnu.org/licenses/>

// B64 management.
package b64

import (
	"encoding/base64"
	"github.com/dedeme/dmstack/machine"
	"github.com/dedeme/dmstack/symbol"
	"github.com/dedeme/dmstack/token"
)

var b64 = base64.StdEncoding

// Returns a B64 string from a normal string.
//    m: Virtual machine.
func prEncode(m *machine.T) {
	tk := m.PopT(token.String)
	s, _ := tk.S()
	m.Push(token.NewS(b64.EncodeToString([]byte(s)), m.MkPos()))
}

// Returns a B64 string from a Blob.
//    m: Virtual machine.
func prEncodeBytes(m *machine.T) {
	tk := m.PopT(token.Native)
	sym, b, _ := tk.N()
	if sym != symbol.Blob_ {
		m.Failt("\n  Expected: Blob object.\n  Actual  : '%v'.", sym)
	}
	m.Push(token.NewS(b64.EncodeToString(b.([]byte)), m.MkPos()))
}

// Returns a string from a B64 string. If decoding fails, it raises a
// "B64 error".
//    m: Virtual machine.
func prDecode(m *machine.T) {
	tk := m.PopT(token.String)
	s, _ := tk.S()
	b, err := b64.DecodeString(s)
	if err != nil {
		m.Fail("B64 error", "Wrong B64 string:\n%v", s)
	}
	m.Push(token.NewS(string(b), m.MkPos()))
}

// Returns a Blob from a B64 string. If decoding fails raises an "B64 error".
//    m: Virtual machine.
func prDecodeBytes(m *machine.T) {
	tk := m.PopT(token.String)
	s, _ := tk.S()
	b, err := b64.DecodeString(s)
	if err != nil {
		m.Fail("B64 error", "Wrong B64 string:\n%v", s)
	}
	m.Push(token.NewN(symbol.Blob_, b, m.MkPos()))
}

// Processes date procedures.
//    m: Virtual machine.
//    run: Function which running a machine.
func Proc(m *machine.T, run func(m *machine.T)) {
	tk, ok := m.PrgNext()
	if !ok {
		m.Failt("'b64' procedure is missing")
	}
	sy, ok := tk.Sy()
	if !ok {
		m.Failt(
			"\n  Expected: 'b64' procedure.\n  Actual  : '%v'.", tk.StringDraft(),
		)
	}
	switch sy {
	case symbol.New("encode"):
		prEncode(m)
	case symbol.New("encodeBytes"):
		prEncodeBytes(m)
	case symbol.New("decode"):
		prDecode(m)
	case symbol.New("decodeBytes"):
		prDecodeBytes(m)
	default:
		m.Failt("'b64' does not contains the procedure '%v'.", sy.String())
	}
}

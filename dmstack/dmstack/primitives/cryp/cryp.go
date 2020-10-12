// Copyright 29-Sep-2020 ºDeme
// GNU General Public License - V3 <http://www.gnu.org/licenses/>

// Encrypt and decript procedures.
package cryp

import (
	"crypto/rand"
	"encoding/base64"
	"github.com/dedeme/dmstack/machine"
	"github.com/dedeme/dmstack/symbol"
	"github.com/dedeme/dmstack/token"
)

var b64 = base64.StdEncoding

// Auxiliar function.
func popStr(m *machine.T) string {
	tk := m.PopT(token.String)
	s, _ := tk.S()
	return s
}

// Auxiliar function.
func pushStr(m *machine.T, s string) {
	m.Push(token.NewS(s, m.MkPos()))
}

// Auxiliar function.
func keyf(key string, lg int64) string {
	k, err := b64.DecodeString(b64.EncodeToString([]byte(
		key + "codified in irreversibleDeme is good, very good!\n\r8@@",
	)))

	if err != nil {
		panic(err)
	}

	lenk := int64(len(k))
	sum := int64(0)
	for i := int64(0); i < lenk; i++ {
		sum += int64(k[i])
	}

	lg2 := lg + lenk
	r := make([]byte, lg2)
	r1 := make([]byte, lg2)
	r2 := make([]byte, lg2)
	ik := int64(0)
	for i := int64(0); i < lg2; i++ {
		v1 := int64(k[ik])
		v2 := v1 + int64(k[v1%lenk])
		v3 := v2 + int64(k[v2%lenk])
		v4 := v3 + int64(k[v3%lenk])
		sum = (sum + i + v4) & 255
		r1[i] = byte(sum)
		r2[i] = byte(sum)
		ik++
		if ik == lenk {
			ik = 0
		}
	}

	for i := int64(0); i < lg2; i++ {
		v1 := int64(r2[i])
		v2 := v1 + int64(r2[v1%lg2])
		v3 := v2 + int64(r2[v2%lg2])
		v4 := v3 + int64(r2[v3%lg2])
		sum = (sum + v4) & 255
		r2[i] = byte(sum)
		r[i] = byte((sum + int64(r1[i])) & 255)
	}

	return b64.EncodeToString(r)[:lg]
}

// Returns a random B64 string of length 'l'.
//    m: Virtual machine.
func prGenk(m *machine.T) {
	tk := m.PopT(token.Int)
	l, _ := tk.I()

	arr := make([]byte, l)
	_, err := rand.Read(arr)
	if err != nil {
		panic(err)
	}

	pushStr(m, b64.EncodeToString(arr)[:l])
}

// Returns 'key' codified in irreversible way, using 'lg' B64 characters.
//    m: Virtual machine.
func prKey(m *machine.T) {
	tk := m.PopT(token.Int)
	lg, _ := tk.I()
	key := popStr(m)

	pushStr(m, keyf(key, lg))
}

// Encodes 'msg' with key 'key'.
//    m: Virtual machine.
func prEncode(m *machine.T) {
	msg := popStr(m)
	key := popStr(m)

	ma := b64.EncodeToString([]byte(msg))
	lg := int64(len(ma))
	k := keyf(key, lg)
	mb := []byte(ma)
	kb := []byte(k)
	r := make([]byte, lg)
	for i := int64(0); i < lg; i++ {
		r[i] = mb[i] + kb[i]
	}

	pushStr(m, b64.EncodeToString(r))
}

// Decodes 'c' using key 'key'. 'c' was codified with 'cryp'.
//    m: Virtual machine.
func prDecode(m *machine.T) {
	c := popStr(m)
	key := popStr(m)

	mb, err := b64.DecodeString(c)
	if err != nil {
		m.Fail("B64 error", "Wrong B64 string:\n%v", c)
	}
	lg := int64(len(mb))
	k := keyf(key, lg)
	kb := []byte(k)
	r := make([]byte, lg)
	for i := int64(0); i < lg; i++ {
		r[i] = mb[i] - kb[i]
	}
	mb, err = b64.DecodeString(string(r))
	if err != nil {
		m.Fail("Cryp error", "Code can not be decrypted:\n%v", c)
	}

	pushStr(m, string(mb))
}

// Processes date procedures.
//    m: Virtual machine.
//    run: Function which running a machine.
func Proc(m *machine.T, run func(m *machine.T)) {
	tk, ok := m.PrgNext()
	if !ok {
		m.Failt("'cryp' procedure is missing")
	}
	sy, ok := tk.Sy()
	if !ok {
		m.Failt(
			"\n  Expected: 'cryp' procedure.\n  Actual  : '%v'.", tk.StringDraft(),
		)
	}
	switch sy {
	case symbol.New("genk"):
		prGenk(m)
	case symbol.New("key"):
		prKey(m)
	case symbol.New("encode"):
		prEncode(m)
	case symbol.New("decode"):
		prDecode(m)
	default:
		m.Failt("'cryp' does not contains the procedure '%v'.", sy.String())
	}
}

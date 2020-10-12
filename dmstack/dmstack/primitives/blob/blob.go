// Copyright 29-Sep-2020 ºDeme
// GNU General Public License - V3 <http://www.gnu.org/licenses/>

// Blob procedures.
package blob

import (
	"github.com/dedeme/dmstack/machine"
	"github.com/dedeme/dmstack/symbol"
	"github.com/dedeme/dmstack/token"
)

// Auxiliar function
func popBlob(m *machine.T) []byte {
	tk := m.PopT(token.Native)
	sym, i, _ := tk.N()
	if sym != symbol.Blob_ {
		m.Failt("\n  Expected: Blob object.\n  Actual  : '%v'.", sym)
	}
	return i.([]byte)
}

// Auxiliar function
func pushBlob(m *machine.T, d []byte) {
	m.Push(token.NewN(symbol.Blob_, d, m.MkPos()))
}

// Returns a '0' blob with length 'l'. If 'l' < 0, it raises a "Blob error".
//    m: Virtual machine.
func prNew(m *machine.T) {
	tk := m.PopT(token.Int)
	l, _ := tk.I()
	if l < 0 {
		m.Fail("Blob error", "Blob length < 0 (%v)", l)
	}
	pushBlob(m, make([]byte, l))
}

// Fill a blob with a value.
//    m: Virtual machine.
func prFill(m *machine.T) {
	tk := m.PopT(token.Int)
	l, _ := tk.I()
	b := popBlob(m)
	by := byte(l & 255)
	for i := range b {
		b[i] = by
	}
	pushBlob(m, b)
}

// Returns a blob from a list of Int.
//    m: Virtual machine.
func prFrom(m *machine.T) {
	tk1 := m.PopT(token.List)
	l, _ := tk1.L()
	b := make([]byte, len(l))
	for i, tk := range l {
		v, ok := tk.I()
		if !ok {
			m.Failt("\n  Expected: Int object.\n  Actual  : '%v'.", tk)
		}
		b[i] = byte(v & 255)
	}
	pushBlob(m, b)
}

// Returns a blob from a String.
//    m: Virtual machine.
func prFromStr(m *machine.T) {
	tk := m.PopT(token.String)
	s, _ := tk.S()
	pushBlob(m, []byte(s))
}

// Returns a byte of blob. If i is out of range return an
// "Index out of range error".
//    m: Virtual machine.
func prGet(m *machine.T) {
	tk1 := m.PopT(token.Int)
	i, _ := tk1.I()
	b := popBlob(m)
	if i < 0 || i > int64(len(b)) {
		m.Fail("Index out of range error", "%v", i)
	}
	m.Push(token.NewI(int64(b[i]), m.MkPos()))
}

// Sets a byte of blob. If i is out of range return an
// "Index out of range error".
//    m: Virtual machine.
func prSet(m *machine.T) {
	tk2 := m.PopT(token.Int)
	v, _ := tk2.I()
	tk1 := m.PopT(token.Int)
	i, _ := tk1.I()
	b := popBlob(m)
	if i < 0 || i > int64(len(b)) {
		m.Fail("Index out of range error", "%v", i)
	}
	b[i] = byte(v & 255)
}

// Sets a byte of blob applying 'p' to the previous value. If i is out of range
// return an "Index out of range error".
//    m: Virtual machine.
func prUp(m *machine.T, run func(m *machine.T)) {
	tk2 := m.PopT(token.Procedure)
	p, _ := tk2.P()
	tk1 := m.PopT(token.Int)
	i, _ := tk1.I()
	b := popBlob(m)
	if i < 0 || i > int64(len(b)) {
		m.Fail("Index out of range error", "%v", i)
	}
	m2 := machine.NewIsolate(m.SourceDir, m.Pmachines, p)
	m2.Push(token.NewI(int64(b[i]), m.MkPos()))
	run(m2)
	tk3 := m2.PopT(token.Int)
	v, _ := tk3.I()

	b[i] = byte(v & 255)
}

// Returns the size of 'b'
//    m: Virtual machine.
func prSize(m *machine.T) {
	b := popBlob(m)
	m.Push(token.NewI(int64(len(b)), m.MkPos()))
}

// Returns 'true' if b1 == b2
//    m: Virtual machine.
func prEq(m *machine.T) {
	b2 := popBlob(m)
	b1 := popBlob(m)
	var r bool
	if len(b1) == len(b2) {
		r = true
		for i, e := range b1 {
			if e != b2[i] {
				r = false
				break
			}
		}
	}
	m.Push(token.NewB(r, m.MkPos()))
}

// Returns 'true' if b1 != b2
//    m: Virtual machine.
func prNeq(m *machine.T) {
	b2 := popBlob(m)
	b1 := popBlob(m)
	r := true
	if len(b1) == len(b2) {
		r = false
		for i, e := range b1 {
			if e != b2[i] {
				r = true
				break
			}
		}
	}
	m.Push(token.NewB(r, m.MkPos()))
}

// Auxiliar function
func sub(m *machine.T, b []byte, begin, end int64) []byte {
	l := int64(len(b))
	if begin < 0 {
		begin = l + begin
	}
	if end < 0 {
		end = l + end
	}
	if begin < 0 || begin > l {
		m.Fail("Index out of range error", "%v", begin)
	}
	if end > l {
		m.Fail("Index out of range error", "%v", end)
	}
	if end < begin {
		end = begin
	}

	newL := end - begin
	r := make([]byte, newL)
	for i := int64(0); i < newL; i++ {
		r[i] = b[begin]
		begin++
	}
	return r
}

// Returns a sub-blob from begin (inclusive) to end (exclusive)
//    m: Virtual machine.
func prSub(m *machine.T) {
	tk2 := m.PopT(token.Int)
	end, _ := tk2.I()
	tk1 := m.PopT(token.Int)
	begin, _ := tk1.I()
	b := popBlob(m)
	pushBlob(m, sub(m, b, begin, end))
}

// Equals to sub(0, end)
//    m: Virtual machine.
func prLeft(m *machine.T) {
	tk1 := m.PopT(token.Int)
	end, _ := tk1.I()
	b := popBlob(m)
	pushBlob(m, sub(m, b, 0, end))
}

// Equals to sub(begin, len(b))
//    m: Virtual machine.
func prRight(m *machine.T) {
	tk1 := m.PopT(token.Int)
	begin, _ := tk1.I()
	b := popBlob(m)
	pushBlob(m, sub(m, b, begin, int64(len(b))))
}

// Concatenates two blobs.
//    m: Virtual machine.
func prPlus2(m *machine.T) {
	b2 := popBlob(m)
	b1 := popBlob(m)
	var r []byte
	r = append(append(r, b1...), b2...)
	pushBlob(m, r)
}

// Returns b1 + b2
//    m: Virtual machine.
func prPlus(m *machine.T) {
	b2 := popBlob(m)
	b1 := popBlob(m)
	l := len(b1)
	l2 := len(b2)
	if l2 < l {
		l = l2
	}
	r := make([]byte, l)
	for i := 0; i < l; i++ {
		r[i] = b1[i] + b2[i]
	}
	pushBlob(m, r)
}

// Returns b1 - b2
//    m: Virtual machine.
func prMinus(m *machine.T) {
	b2 := popBlob(m)
	b1 := popBlob(m)
	l := len(b1)
	l2 := len(b2)
	if l2 < l {
		l = l2
	}
	r := make([]byte, l)
	for i := 0; i < l; i++ {
		r[i] = b1[i] - b2[i]
	}
	pushBlob(m, r)
}

// Returns b1 & b2
//    m: Virtual machine.
func prAnd(m *machine.T) {
	b2 := popBlob(m)
	b1 := popBlob(m)
	l := len(b1)
	l2 := len(b2)
	if l2 < l {
		l = l2
	}
	r := make([]byte, l)
	for i := 0; i < l; i++ {
		r[i] = b1[i] & b2[i]
	}
	pushBlob(m, r)
}

// Returns b1 | b2
//    m: Virtual machine.
func prOr(m *machine.T) {
	b2 := popBlob(m)
	b1 := popBlob(m)
	l := len(b1)
	l2 := len(b2)
	if l2 < l {
		l = l2
	}
	r := make([]byte, l)
	for i := 0; i < l; i++ {
		r[i] = b1[i] | b2[i]
	}
	pushBlob(m, r)
}

// Returns b1 ^ b2
//    m: Virtual machine.
func prXor(m *machine.T) {
	b2 := popBlob(m)
	b1 := popBlob(m)
	l := len(b1)
	l2 := len(b2)
	if l2 < l {
		l = l2
	}
	r := make([]byte, l)
	for i := 0; i < l; i++ {
		r[i] = b1[i] ^ b2[i]
	}
	pushBlob(m, r)
}

// Returns a list of Int from a blob.
//    m: Virtual machine.
func prTo(m *machine.T) {
	pos := m.MkPos()
	b := popBlob(m)
	var l []*token.T
	for _, v := range b {
		l = append(l, token.NewI(int64(v), pos))
	}
	m.Push(token.NewL(l, pos))
}

// Returns a String from a blob.
//    m: Virtual machine.
func prToStr(m *machine.T) {
	b := popBlob(m)
	m.Push(token.NewS(string(b), m.MkPos()))
}

// Processes date procedures.
//    m: Virtual machine.
//    run: Function which running a machine.
func Proc(m *machine.T, run func(m *machine.T)) {
	tk, ok := m.PrgNext()
	if !ok {
		m.Failt("'blob' procedure is missing")
	}
	sy, ok := tk.Sy()
	if !ok {
		m.Failt(
			"\n  Expected: 'blob' procedure.\n  Actual  : '%v'.", tk.StringDraft(),
		)
	}
	switch sy {
	case symbol.New("new"):
		prNew(m)
	case symbol.New("fill"):
		prFill(m)
	case symbol.From:
		prFrom(m)
	case symbol.New("fromStr"):
		prFromStr(m)
	case symbol.New("get"):
		prGet(m)
	case symbol.New("set"):
		prSet(m)
	case symbol.New("up"):
		prUp(m, run)
	case symbol.New("size"):
		prSize(m)
	case symbol.New("=="):
		prEq(m)
	case symbol.New("!="):
		prNeq(m)
	case symbol.New("sub"):
		prSub(m)
	case symbol.New("left"):
		prLeft(m)
	case symbol.New("right"):
		prRight(m)
	case symbol.New("++"):
		prPlus2(m)
	case symbol.New("+"):
		prPlus(m)
	case symbol.New("-"):
		prMinus(m)
	case symbol.New("&"):
		prAnd(m)
	case symbol.New("|"):
		prOr(m)
	case symbol.New("^"):
		prXor(m)
	case symbol.New("to"):
		prTo(m)
	case symbol.New("toStr"):
		prToStr(m)
	default:
		m.Failt("'blob' does not contains the procedure '%v'.", sy.String())
	}
}

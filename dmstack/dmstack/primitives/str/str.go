// Copyright 24-Jul-2020 ºDeme
// GNU General Public License - V3 <http://www.gnu.org/licenses/>

// String procedures.
package str

import (
	"github.com/dedeme/dmstack/machine"
	"github.com/dedeme/dmstack/primitives/sys"
	"github.com/dedeme/dmstack/symbol"
	"github.com/dedeme/dmstack/token"
	"io"
	"strings"
  "strconv"
)

// Auxiliar function
func popStr(m *machine.T) (s string) {
	tk := m.PopT(token.String)
	s, _ = tk.S()
	return
}

// Auxiliar function
func pushStr(m *machine.T, s string) {
	m.Push(token.NewS(s, m.MkPos()))
}

// Auxiliar function
func runesLength(m *machine.T, tk *token.T, s string) int64 {
	rd := strings.NewReader(s)
	c := int64(0)
	for {
		_, _, e := rd.ReadRune()
		if e == nil {
			c++
			continue
		}
		if e == io.EOF {
			break
		}
		m.Failf("Wrong UTF-8 string '%s'", tk.StringDraft())
	}
	return c
}

// Creeates a UTF-8 string for a ISO string.
//    m : Virtual machine.
func prFromIso(m *machine.T) {
	s := popStr(m)
	var bf strings.Builder
	for i := 0; i < len(s); i++ {
		ch := s[i]
		if ch < 0x80 {
			bf.WriteByte(ch)
			continue
		}
		bf.WriteByte(0xc0 | (ch&0xc0)>>6)
		bf.WriteByte(0x80 | (ch & 0x3f))
	}
	pushStr(m, bf.String())
}

// Compares two strings.
//    m : Virtual machine.
func prCmp(m *machine.T) {
	s2 := popStr(m)
	s1 := popStr(m)
	m.Push(token.NewI(int64(sys.Collator()(s1, s2)), m.MkPos()))
}

// Returns true if s1 == s2.
//    m : Virtual machine.
func prEq(m *machine.T) {
	s2 := popStr(m)
	s1 := popStr(m)
	m.Push(token.NewB(sys.Collator()(s1, s2) == 0, m.MkPos()))
}

// Returns true if s1 != s2.
//    m : Virtual machine.
func prNeq(m *machine.T) {
	s2 := popStr(m)
	s1 := popStr(m)
	m.Push(token.NewB(sys.Collator()(s1, s2) != 0, m.MkPos()))
}

// Returns true if s1 > s2.
//    m : Virtual machine.
func prGreater(m *machine.T) {
	s2 := popStr(m)
	s1 := popStr(m)
	m.Push(token.NewB(sys.Collator()(s1, s2) > 0, m.MkPos()))
}

// Returns true if s1 >= s2.
//    m : Virtual machine.
func prGreaterEq(m *machine.T) {
	s2 := popStr(m)
	s1 := popStr(m)
	m.Push(token.NewB(sys.Collator()(s1, s2) >= 0, m.MkPos()))
}

// Returns true if s1 < s2.
//    m : Virtual machine.
func prLess(m *machine.T) {
	s2 := popStr(m)
	s1 := popStr(m)
	m.Push(token.NewB(sys.Collator()(s1, s2) < 0, m.MkPos()))
}

// Returns true if s1 <= s2.
//    m : Virtual machine.
func prLessEq(m *machine.T) {
	s2 := popStr(m)
	s1 := popStr(m)
	m.Push(token.NewB(sys.Collator()(s1, s2) <= 0, m.MkPos()))
}

// Returns string bytes length.
//    m : Virtual machine.
func prEnds(m *machine.T) {
	sub := popStr(m)
	s := popStr(m)
	m.Push(token.NewB(strings.HasSuffix(s, sub), m.MkPos()))
}

// Returns string bytes length.
//    m : Virtual machine.
func prStarts(m *machine.T) {
	sub := popStr(m)
	s := popStr(m)
	m.Push(token.NewB(strings.HasPrefix(s, sub), m.MkPos()))
}

// Returns index of the first occurrence of 'sub' in 's' or -1 if 'sub' is not
// found.
//    m : Virtual machine.
func prIndex(m *machine.T) {
	sub := popStr(m)
	s := popStr(m)
	m.Push(token.NewI(int64(strings.Index(s, sub)), m.MkPos()))
}

// Returns index of the first occurrence of 'sub' in 's' from 'i' inclusive,
// or -1 if 'sub' is not found (e.g. "abcdc" "c" 3 str.indexFrom -> 4).
//    m : Virtual machine.
func prIndexFrom(m *machine.T) {
  tk := m.PopT(token.Int)
  i, _ := tk.I()
	sub := popStr(m)
	tk2 := m.PopT(token.String)
	s, _ := tk2.S()
  if i > int64(len(s)) {
    m.Failf("Index out of range (%v in %v)", i, tk2.StringDraft())
  }
  r := int64(strings.Index(s[i:], sub))
  if r != -1 {
    r += i
  }
	m.Push(token.NewI(r, m.MkPos()))
}

// Returns index of the last occurrence of 'sub' in 's' or -1 if 'sub' is not
// found.
//    m : Virtual machine.
func prLastIndex(m *machine.T) {
	sub := popStr(m)
	s := popStr(m)
	m.Push(token.NewI(int64(strings.LastIndex(s, sub)), m.MkPos()))
}

// Replaces 'old' by 'new' in 's'
//    m : Virtual machine.
func prReplace(m *machine.T) {
  new := popStr(m)
  old := popStr(m)
  s := popStr(m)
  pushStr(m, strings.ReplaceAll(s, old, new))
}

// Joins a string list with separator 'sep'.
//    m : Virtual machine.
func prJoin(m *machine.T) {
  sep := popStr(m)
  tk := m.PopT(token.List)
  ss, _ := tk.L()
  var as []string
  for _, sTk := range ss {
    s, ok := sTk.S()
    if !ok {
      m.Failf("Expected: String.\nActual  : %v.", sTk.StringDraft())
    }
    as = append(as, s)
  }
  pushStr(m, strings.Join(as, sep))
}

// Splits 's' by separator 'sep' in a List.
//    m : Virtual machine.
func prSplit(m *machine.T) {
  sep := popStr(m)
  s := popStr(m)
  var tks []*token.T
  for _, elem := range strings.Split(s, sep) {
    tks = append(tks, token.NewS(elem, m.MkPos()))
  }
  m.Push(token.NewL(tks, m.MkPos()))
}

// Splits 's' by separator 'sep' in a List, 'trimming' each element of it.
//    m : Virtual machine.
func prSplitTrim(m *machine.T) {
  sep := popStr(m)
  s := popStr(m)
  var tks []*token.T
  for _, elem := range strings.Split(s, sep) {
    tks = append(tks, token.NewS(strings.TrimSpace(elem), m.MkPos()))
  }
  m.Push(token.NewL(tks, m.MkPos()))
}

// Removes spaces at the beginning and at the end of 's'
//    m : Virtual machine.
func prTrim(m *machine.T) {
  s := popStr(m)
  pushStr(m, strings.TrimSpace(s))
}

// Removes spaces at the beginning of 's'
//    m : Virtual machine.
func prLtrim(m *machine.T) {
  s := popStr(m)
  s = strings.TrimSpace(s + "*")
  pushStr(m, s[:len(s) - 1])
}

// Removes spaces at the end of 's'
//    m : Virtual machine.
func prRtrim(m *machine.T) {
  s := popStr(m)
  s = strings.TrimSpace("*" + s)
  pushStr(m, s[1:])
}

// Sets 's' to uppercase
//    m : Virtual machine.
func prToUpper(m *machine.T) {
  pushStr(m, strings.ToUpper(popStr(m)))
}

// Sets 's' to lowercase
//    m : Virtual machine.
func prToLower(m *machine.T) {
  pushStr(m, strings.ToLower(popStr(m)))
}

// Returns true if every rune of 's' is a digit (0-9).
//    m : Virtual machine.
func prDigits(m *machine.T) {
  s := popStr(m)
  ok := true
  for i := 0; i < len(s); i++ {
    ch := s[i]
    if ch < 48 || ch > 57 {
      ok = false
      break
    }
  }
  m.Push(token.NewB(ok, m.MkPos()))
}

// Returns true if 's' is a number (float).
//    m : Virtual machine.
func prNumber(m *machine.T) {
  s := popStr(m)
  _, err := strconv.ParseFloat(s, 64)
  ok := false
  if err == nil {
    ok = true
  }
  m.Push(token.NewB(ok, m.MkPos()))
}

// Returns 's' replacing '.' by '' and ',' by '.'.
//    m : Virtual machine.
func prRegularizeIso(m *machine.T) {
  s := popStr(m)
  pushStr(m, strings.ReplaceAll(strings.ReplaceAll(s, ".", ""), ",", "."))
}

// Returns 's' replacing ',' by ''.
//    m : Virtual machine.
func prRegularizeEn(m *machine.T) {
  s := popStr(m)
  pushStr(m, strings.ReplaceAll(s, ",", ""))
}

// Returns string bytes length.
//    m : Virtual machine.
func prLen(m *machine.T) {
	m.Push(token.NewI(int64(len(popStr(m))), m.MkPos()))
}

// Returns string runes length.
//    m : Virtual machine.
func prRunesLen(m *machine.T) {
	tk := m.PopT(token.String)
	s, _ := tk.S()
	m.Push(token.NewI(runesLength(m, tk, s), m.MkPos()))
}

// Returns the byte in position i.
//    m : Virtual machine.
func prGet(m *machine.T) {
	tk1 := m.PopT(token.Int)
	i, _ := tk1.I()
	tk2 := m.PopT(token.String)
	s, _ := tk2.S()
	if i < 0 || i >= int64(len(s)) {
		m.Failf("Index out of range (%v in %v)", i, tk2.StringDraft())
	}
	pushStr(m, string(s[i]))
}

// Returns the rune in position i.
//    m : Virtual machine.
func prGetRune(m *machine.T) {
	tk1 := m.PopT(token.Int)
	i, _ := tk1.I()
	tk2 := m.PopT(token.String)
	s, _ := tk2.S()
	if i < 0 || i >= runesLength(m, tk2, s) {
		m.Failf("Index out of range (%v in %v)", i, tk2.StringDraft())
	}
	rd := strings.NewReader(s)
	ix := int64(0)
	for {
		if ix == i {
			rn, _, _ := rd.ReadRune()
			pushStr(m, string(rn))
			break
		}
		_, _, _ = rd.ReadRune()
		ix++
	}
}

// Returns byte code of the first byte,
func prRune(m *machine.T) {
	rd := strings.NewReader(popStr(m))
	rn, _, err := rd.ReadRune()
	if err != nil {
		m.Fail("String is empty or invalid")
	}
	m.Push(token.NewI(int64(rn), m.MkPos()))
}

// Returns the character-byte with value n.
func prFromRune(m *machine.T) {
	tk := m.PopT(token.Int)
	i, _ := tk.I()
	r := []rune{rune(i)}
	pushStr(m, string(r))
}

// Executes p with each rune.
func prEach(m *machine.T, run func(m *machine.T)) {
	tk := m.PopT(token.Procedure)
	p, _ := tk.P()
	rd := strings.NewReader(popStr(m))
	for {
		rn, _, err := rd.ReadRune()
		if err != nil {
			return
		}
		m2 := machine.New(m.SourceDir, m.Pmachines, p)
		m2.Push(token.NewS(string(rn), m2.MkPos()))
		run(m2)
	}
}

func subaux(
	m *machine.T, tk *token.T, s string, begin int64, end int64,
) string {
	l := int64(len(s))
	if begin < 0 {
		begin = l + begin
	}
	if end < 0 {
		end = l + end
	}
	if begin < 0 {
		m.Failf("Index out of range (%v in %v)", begin, tk.StringDraft())
	}
	if begin > l {
		m.Failf("Index out of range (%v in %v)", begin, tk.StringDraft())
	}
	if end > l {
		m.Failf("Index out of range (%v in %v)", end, tk.StringDraft())
	}
	if end < begin {
		end = begin
	}
	return s[begin:end]
}

// Return the substring between 'begin' (inclusive) and 'end' (exclusive).
// If 'begin' or 'end' is less than 0, its value is substracted to len(s).
//    m : Virtual machine.
func prSub(m *machine.T) {
	tk1 := m.PopT(token.Int)
	end, _ := tk1.I()
	tk2 := m.PopT(token.Int)
	begin, _ := tk2.I()
	tk3 := m.PopT(token.String)
	s, _ := tk3.S()

	pushStr(m, subaux(m, tk3, s, begin, end))
}

// Return the substring between 0 (inclusive) and 'end' (exclusive).
// If 'end' is less than 0, its value is substracted to len(s).
//    m  : Virtual machine.
//    run: Function which running a machine.
func prLeft(m *machine.T) {
	tk1 := m.PopT(token.Int)
	end, _ := tk1.I()
	tk2 := m.PopT(token.String)
	s, _ := tk2.S()

	pushStr(m, subaux(m, tk2, s, 0, end))
}

// Return the substring between 'begin' (inclusive) and 'len(s)' (exclusive).
// If 'begin' is less than 0, its value is substracted to len(s).
//    m : Virtual machine.
func prRight(m *machine.T) {
	tk1 := m.PopT(token.Int)
	begin, _ := tk1.I()
	tk2 := m.PopT(token.String)
	s, _ := tk2.S()

	pushStr(m, subaux(m, tk2, s, begin, int64(len(s))))
}

// Processes string procedures.
//    m  : Virtual machine.
//    run: Function which running a machine.
func Proc(m *machine.T, run func(m *machine.T)) {
	tk, ok := m.PrgNext()
	if !ok {
		m.Fail("'str' procedure is missing")
	}
	sy, ok := tk.Sy()
	if !ok {
		m.Failf("Expected: 'str' procedure.\nActual  : %v.", tk.StringDraft())
	}
	switch sy {
	case symbol.New("fromIso"):
		prFromIso(m)

	case symbol.New("cmp"):
		prCmp(m)
	case symbol.New(">"):
		prGreater(m)
	case symbol.New("=="):
		prEq(m)
	case symbol.New("!="):
		prNeq(m)
	case symbol.New(">="):
		prGreaterEq(m)
	case symbol.New("<"):
		prLess(m)
	case symbol.New("<="):
		prLessEq(m)

	case symbol.New("ends?"):
		prEnds(m)
	case symbol.New("starts?"):
		prStarts(m)
	case symbol.New("index"):
		prIndex(m)
	case symbol.New("indexFrom"):
		prIndexFrom(m)
	case symbol.New("lastIndex"):
		prLastIndex(m)
	case symbol.New("replace"):
		prReplace(m)
	case symbol.New("join"):
		prJoin(m)
	case symbol.New("split"):
		prSplit(m)
	case symbol.New("splitTrim"):
		prSplitTrim(m)
	case symbol.New("trim"):
		prTrim(m)
	case symbol.New("ltrim"):
		prLtrim(m)
	case symbol.New("rtrim"):
		prRtrim(m)
	case symbol.New("toUpper"):
		prToUpper(m)
	case symbol.New("toLower"):
		prToLower(m)

	case symbol.New("digits?"):
		prDigits(m)
	case symbol.New("number?"):
		prNumber(m)
	case symbol.New("regularizeIso"):
		prRegularizeIso(m)
	case symbol.New("regularizeEn"):
		prRegularizeEn(m)

	case symbol.New("len"):
		prLen(m)
	case symbol.New("runesLen"):
		prRunesLen(m)
	case symbol.New("get"):
		prGet(m)
	case symbol.New("getRune"):
		prGetRune(m)
	case symbol.New("rune"):
		prRune(m)
	case symbol.New("fromRune"):
		prFromRune(m)
	case symbol.New("each"):
		prEach(m, run)

	case symbol.New("sub"):
		prSub(m)
	case symbol.New("left"):
		prLeft(m)
	case symbol.New("right"):
		prRight(m)
	default:
		m.Failf("'str' does not contains the procedure '%v'.", sy.String())
	}
}

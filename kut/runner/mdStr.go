// Copyright 03-Mar-2022 ºDeme
// GNU General Public License - V3 <http://www.gnu.org/licenses/>

package runner

import (
	"fmt"
	"github.com/dedeme/kut/builtin/bfail"
	"github.com/dedeme/kut/builtin/bfunction"
	"github.com/dedeme/kut/expression"
	"strings"
	"unicode/utf16"
)

/// Operates with bytes.
/// \s, i -> s
func strAt(args []*expression.T) (ex *expression.T, err error) {
	switch s := (args[0].Value).(type) {
	case string:
		switch i := (args[1].Value).(type) {
		case int64:
			ex = expression.MkFinal(string(s[i]))
		default:
			err = bfail.Type(args[1], "int")
		}
	default:
		err = bfail.Type(args[0], "string")
	}
	return
}

/// Operates with bytes.
/// \s, i -> s
func strDrop(args []*expression.T) (ex *expression.T, err error) {
	switch s := (args[0].Value).(type) {
	case string:
		switch i := (args[1].Value).(type) {
		case int64:
			if i <= 0 {
				ex = expression.MkFinal(s)
			} else if i >= int64(len(s)) {
				ex = expression.MkFinal("")
			} else {
				ex = expression.MkFinal(s[i:])
			}
		default:
			err = bfail.Type(args[1], "int")
		}
	default:
		err = bfail.Type(args[0], "string")
	}
	return
}

/// \s, s -> b
func strEnds(args []*expression.T) (ex *expression.T, err error) {
	switch s := (args[0].Value).(type) {
	case string:
		switch sub := (args[1].Value).(type) {
		case string:
			ex = expression.MkFinal(strings.HasSuffix(s, sub))
		default:
			err = bfail.Type(args[1], "string")
		}
	default:
		err = bfail.Type(args[0], "string")
	}
	return
}

/// \s -> s
func strFromIso(args []*expression.T) (ex *expression.T, err error) {
	switch s := (args[0].Value).(type) {
	case string:
		var r []byte
		for _, b := range []byte(s) {
			if b < 0x80 {
				r = append(r, b)
			} else {
				r = append(r, 0xc0|(b&0xc0)>>6, 0x80|(b&0x3f))
			}
		}
		ex = expression.MkFinal(string(r))
	default:
		err = bfail.Type(args[0], "string")
	}
	return
}

/// \[s...] -> s
func strFromRunes(args []*expression.T) (ex *expression.T, err error) {
	switch ss := (args[0].Value).(type) {
	case []*expression.T:
		var runes []string
		for _, sex := range ss {
			switch s := (sex.Value).(type) {
			case string:
				runes = append(runes, s)
			default:
				err = bfail.Type(sex, "string")
			}
		}
		if err == nil {
			ex = expression.MkFinal(strings.Join(runes, ""))
		}
	default:
		err = bfail.Type(args[0], "array")
	}
	return
}

/// \[i...] -> s
func strFromUtf16(args []*expression.T) (ex *expression.T, err error) {
	switch ss := (args[0].Value).(type) {
	case []*expression.T:
		var runes []uint16
		for _, sex := range ss {
			switch s := (sex.Value).(type) {
			case int64:
				runes = append(runes, uint16(s))
			default:
				err = bfail.Type(sex, "int")
			}
		}
		if err == nil {
			ex = expression.MkFinal(string(utf16.Decode(runes)))
		}
	default:
		err = bfail.Type(args[0], "array")
	}
	return
}

/// Returns the result of fill 'f' with values of 'vs'
/// \s, a -> s
func strFmt(args []*expression.T) (ex *expression.T, err error) {
	switch f := (args[0].Value).(type) {
	case string:
		switch vs := (args[1].Value).(type) {
		case []*expression.T:
			var vvs []interface{}
			for _, e := range vs {
				vvs = append(vvs, e.Value)
			}
			ex = expression.MkFinal(fmt.Sprintf(f, vvs...))
		default:
			err = bfail.Type(args[1], "array")
		}
	default:
		err = bfail.Type(args[0], "string")
	}
	return
}

/// Returns the index (counting by byte) of 'subs' in 's'.
/// \s, s -> b
func strIndex(args []*expression.T) (ex *expression.T, err error) {
	switch s := (args[0].Value).(type) {
	case string:
		switch subs := (args[1].Value).(type) {
		case string:
			ex = expression.MkFinal(int64(strings.Index(s, subs)))
		default:
			err = bfail.Type(args[1], "string")
		}
	default:
		err = bfail.Type(args[0], "string")
	}
	return
}

/// Returns the index (counting by byte) of 'subs' in 's', searching from ix,
/// inclusive. 'ix' is a byte (not rune) position.
/// \s, s, i -> b
func strIndexFrom(args []*expression.T) (ex *expression.T, err error) {
	switch s := (args[0].Value).(type) {
	case string:
		switch subs := (args[1].Value).(type) {
		case string:
			switch ix := (args[2].Value).(type) {
			case int64:
				r := int64(-1)
				if ix >= int64(0) && ix < int64(len(s)) {
					r = int64(strings.Index(s[ix:], subs))
					if r != -1 {
						r += ix
					}
				}
				ex = expression.MkFinal(r)
			default:
				err = bfail.Type(args[1], "int")
			}
		default:
			err = bfail.Type(args[1], "string")
		}
	default:
		err = bfail.Type(args[0], "string")
	}
	return
}

/// Returns the index (counting by byte) of 'subs' in 's'.
/// \s, s -> b
func strLastIndex(args []*expression.T) (ex *expression.T, err error) {
	switch s := (args[0].Value).(type) {
	case string:
		switch subs := (args[1].Value).(type) {
		case string:
			ex = expression.MkFinal(int64(strings.LastIndex(s, subs)))
		default:
			err = bfail.Type(args[1], "string")
		}
	default:
		err = bfail.Type(args[0], "string")
	}
	return
}

/// Operates with bytes.
/// \s, i -> s
func strLeft(args []*expression.T) (ex *expression.T, err error) {
	switch s := (args[0].Value).(type) {
	case string:
		switch i := (args[1].Value).(type) {
		case int64:
			if i < 0 {
				i = int64(len(s)) + i
			}
			ex = expression.MkFinal(s[:i])
		default:
			err = bfail.Type(args[1], "int")
		}
	default:
		err = bfail.Type(args[0], "string")
	}
	return
}

/// Returns the length (counting by byte).
/// \s -> i
func strLen(args []*expression.T) (ex *expression.T, err error) {
	switch s := (args[0].Value).(type) {
	case string:
		ex = expression.MkFinal(int64(len(s)))
	default:
		err = bfail.Type(args[0], "string")
	}
	return
}

/// \s -> s
func strLtrim(args []*expression.T) (ex *expression.T, err error) {
	switch s := (args[0].Value).(type) {
	case string:
		ex = expression.MkFinal(strings.TrimLeft(s, " \n\t\r"))
	default:
		err = bfail.Type(args[0], "string")
	}
	return
}

/// \s, s, s -> s
func strReplace(args []*expression.T) (ex *expression.T, err error) {
	switch s := (args[0].Value).(type) {
	case string:
		switch old := (args[1].Value).(type) {
		case string:
			switch new := (args[2].Value).(type) {
			case string:
				ex = expression.MkFinal(strings.ReplaceAll(s, old, new))
			default:
				err = bfail.Type(args[2], "string")
			}
		default:
			err = bfail.Type(args[1], "string")
		}
	default:
		err = bfail.Type(args[0], "string")
	}
	return
}

/// Operates with bytes.
/// \s, i -> s
func strRight(args []*expression.T) (ex *expression.T, err error) {
	switch s := (args[0].Value).(type) {
	case string:
		switch i := (args[1].Value).(type) {
		case int64:
			if i < 0 {
				i = int64(len(s)) + i
			}
			ex = expression.MkFinal(s[i:])
		default:
			err = bfail.Type(args[1], "int")
		}
	default:
		err = bfail.Type(args[0], "string")
	}
	return
}

/// \s -> s
func strRtrim(args []*expression.T) (ex *expression.T, err error) {
	switch s := (args[0].Value).(type) {
	case string:
		ex = expression.MkFinal(strings.TrimRight(s, " \n\t\r"))
	default:
		err = bfail.Type(args[0], "string")
	}
	return
}

/// \s, s -> [s...]
func strSplit(args []*expression.T) (ex *expression.T, err error) {
	switch s := (args[0].Value).(type) {
	case string:
		switch sep := (args[1].Value).(type) {
		case string:
			var exs []*expression.T
			for _, e := range strings.Split(s, sep) {
				exs = append(exs, expression.MkFinal(e))
			}
			ex = expression.MkFinal(exs)
		default:
			err = bfail.Type(args[1], "string")
		}
	default:
		err = bfail.Type(args[0], "string")
	}
	return
}

/// \s, s -> [s...]
func strSplitTrim(args []*expression.T) (ex *expression.T, err error) {
	switch s := (args[0].Value).(type) {
	case string:
		switch sep := (args[1].Value).(type) {
		case string:
			var exs []*expression.T
			for _, e := range strings.Split(s, sep) {
				exs = append(exs, expression.MkFinal(strings.TrimSpace(e)))
			}
			ex = expression.MkFinal(exs)
		default:
			err = bfail.Type(args[1], "string")
		}
	default:
		err = bfail.Type(args[0], "string")
	}
	return
}

/// \s, s -> b
func strStarts(args []*expression.T) (ex *expression.T, err error) {
	switch s := (args[0].Value).(type) {
	case string:
		switch sub := (args[1].Value).(type) {
		case string:
			ex = expression.MkFinal(strings.HasPrefix(s, sub))
		default:
			err = bfail.Type(args[1], "string")
		}
	default:
		err = bfail.Type(args[0], "string")
	}
	return
}

/// Operates with bytes.
/// \s, i, i -> s
func strSub(args []*expression.T) (ex *expression.T, err error) {
	switch s := (args[0].Value).(type) {
	case string:
		switch start := (args[1].Value).(type) {
		case int64:
			switch end := (args[2].Value).(type) {
			case int64:
				if start < 0 {
					start = int64(len(s)) + start
				}
				if end < 0 {
					end = int64(len(s)) + end
				}
				ex = expression.MkFinal(s[start:end])
			default:
				err = bfail.Type(args[1], "int")
			}
		default:
			err = bfail.Type(args[1], "int")
		}
	default:
		err = bfail.Type(args[0], "string")
	}
	return
}

/// Operates with bytes.
/// \s, i -> s
func strTake(args []*expression.T) (ex *expression.T, err error) {
	switch s := (args[0].Value).(type) {
	case string:
		switch i := (args[1].Value).(type) {
		case int64:
			if i <= 0 {
				ex = expression.MkFinal("")
			} else if i >= int64(len(s)) {
				ex = expression.MkFinal(s)
			} else {
				ex = expression.MkFinal(s[:i])
			}
		default:
			err = bfail.Type(args[1], "int")
		}
	default:
		err = bfail.Type(args[0], "string")
	}
	return
}

/// \s -> s
func strToLower(args []*expression.T) (ex *expression.T, err error) {
	switch s := (args[0].Value).(type) {
	case string:
		ex = expression.MkFinal(strings.ToLower(s))
	default:
		err = bfail.Type(args[0], "string")
	}
	return
}

/// \s -> [s...]
func strToRunes(args []*expression.T) (ex *expression.T, err error) {
	switch s := (args[0].Value).(type) {
	case string:
		var runes []*expression.T
		for _, e := range []rune(s) {
			runes = append(runes, expression.MkFinal(string(e)))
		}
		ex = expression.MkFinal(runes)
	default:
		err = bfail.Type(args[0], "string")
	}
	return
}

/// \s -> s
func strToUpper(args []*expression.T) (ex *expression.T, err error) {
	switch s := (args[0].Value).(type) {
	case string:
		ex = expression.MkFinal(strings.ToUpper(s))
	default:
		err = bfail.Type(args[0], "string")
	}
	return
}

/// \s -> [i...]
func strToUtf16(args []*expression.T) (ex *expression.T, err error) {
	switch s := (args[0].Value).(type) {
	case string:
		var codepoints []*expression.T
		for _, e := range utf16.Encode([]rune(s)) {
			codepoints = append(codepoints, expression.MkFinal(int64(e)))
		}
		ex = expression.MkFinal(codepoints)
	default:
		err = bfail.Type(args[0], "string")
	}
	return
}

/// \s -> s
func strTrim(args []*expression.T) (ex *expression.T, err error) {
	switch s := (args[0].Value).(type) {
	case string:
		ex = expression.MkFinal(strings.TrimSpace(s))
	default:
		err = bfail.Type(args[0], "string")
	}
	return
}

func strGet(fname string) (fn *bfunction.T, ok bool) {
	ok = true
	switch fname {
	case "drop":
		fn = bfunction.New(2, strDrop)
	case "ends":
		fn = bfunction.New(2, strEnds)
	case "fmt":
		fn = bfunction.New(2, strFmt)
	case "fromIso":
		fn = bfunction.New(1, strFromIso)
	case "fromRunes":
		fn = bfunction.New(1, strFromRunes)
	case "fromUtf16":
		fn = bfunction.New(1, strFromUtf16)
	case "index":
		fn = bfunction.New(2, strIndex)
	case "indexFrom":
		fn = bfunction.New(3, strIndexFrom)
	case "lastIndex":
		fn = bfunction.New(2, strLastIndex)
	case "len":
		fn = bfunction.New(1, strLen)
	case "ltrim":
		fn = bfunction.New(1, strLtrim)
	case "replace":
		fn = bfunction.New(3, strReplace)
	case "rtrim":
		fn = bfunction.New(1, strRtrim)
	case "split":
		fn = bfunction.New(2, strSplit)
	case "splitTrim":
		fn = bfunction.New(2, strSplitTrim)
	case "starts":
		fn = bfunction.New(2, strStarts)
	case "take":
		fn = bfunction.New(2, strTake)
	case "toLower":
		fn = bfunction.New(1, strToLower)
	case "toRunes":
		fn = bfunction.New(1, strToRunes)
	case "toUpper":
		fn = bfunction.New(1, strToUpper)
	case "toUtf16":
		fn = bfunction.New(1, strToUtf16)
	case "trim":
		fn = bfunction.New(1, strTrim)
	default:
		ok = false
	}

	return
}

// Copyright 03-Mar-2022 ºDeme
// GNU General Public License - V3 <http://www.gnu.org/licenses/>

package runner

import (
	"encoding/json"
	"fmt"
	"github.com/dedeme/kut/builtin/bfail"
	"github.com/dedeme/kut/builtin/bfunction"
	"github.com/dedeme/kut/expression"
	"strconv"
	"strings"
)

func nextByte(s string, ch byte, ix int) (pos int, ok bool) {
	pos = ix
	l := len(s)
	quotes := false
	bar := false
	squarn := 0
	bracketn := 0
	ok = true
	var c byte
	for {
		if pos == l {
			ok = false
			break
		}
		c = s[pos]
		if quotes {
			if bar {
				bar = false
			} else {
				if c == '\\' {
					bar = true
				} else if c == '"' {
					quotes = false
				}
			}
		} else {
			if c == ch &&
				((c == ']' && squarn == 1 && bracketn == 0) ||
					(c == '}' && squarn == 0 && bracketn == 1) ||
					(squarn == 0 && bracketn == 0)) {
				break
			} else if c == '"' {
				quotes = true
			} else if c == '[' {
				squarn++
			} else if c == ']' {
				squarn--
			} else if c == '{' {
				bracketn++
			} else if c == '}' {
				bracketn--
			}
		}
		pos++
	}
	return
}

/// \s -> b
func jsIsNull(args []*expression.T) (ex *expression.T, err error) {
	switch s := (args[0].Value).(type) {
	case string:
		ex = expression.MkFinal(strings.TrimSpace(s) == "null")
	default:
		err = bfail.Type(args[0], "string")
	}
	return
}

/// \s -> b
func jsRb(args []*expression.T) (ex *expression.T, err error) {
	switch s := (args[0].Value).(type) {
	case string:
		var v bool
		err = json.Unmarshal([]byte(s), &v)
		if err == nil {
			ex = expression.MkFinal(v)
		} else {
			err = bfail.Mk(fmt.Sprintf("%v in\n'%v'", err.Error(), s))
		}
	default:
		err = bfail.Type(args[0], "string")
	}
	return
}

/// \s -> i
func jsRi(args []*expression.T) (ex *expression.T, err error) {
	switch s := (args[0].Value).(type) {
	case string:
		var v int64
		err = json.Unmarshal([]byte(s), &v)
		if err == nil {
			ex = expression.MkFinal(v)
		} else {
			err = bfail.Mk(fmt.Sprintf("%v in\n'%v'", err.Error(), s))
		}
	default:
		err = bfail.Type(args[0], "string")
	}
	return
}

/// \s -> f
func jsRf(args []*expression.T) (ex *expression.T, err error) {
	switch s := (args[0].Value).(type) {
	case string:
		var v float64
		err = json.Unmarshal([]byte(s), &v)
		if err == nil {
			ex = expression.MkFinal(v)
		} else {
			err = bfail.Mk(fmt.Sprintf("%v in\n'%v'", err.Error(), s))
		}
	default:
		err = bfail.Type(args[0], "string")
	}
	return
}

/// \s -> s
func jsRs(args []*expression.T) (ex *expression.T, err error) {
	switch s := (args[0].Value).(type) {
	case string:
		var v string
		err = json.Unmarshal([]byte(s), &v)
		if err == nil {
			ex = expression.MkFinal(v)
		} else {
			err = bfail.Mk(fmt.Sprintf("%v in\n'%v'", err.Error(), s))
		}
	default:
		err = bfail.Type(args[0], "string")
	}
	return
}

/// \s -> [s...]
func jsRa(args []*expression.T) (ex *expression.T, err error) {
	switch s := (args[0].Value).(type) {
	case string:
		s = strings.TrimSpace(s)
		var rs []*expression.T
		if !strings.HasPrefix(s, "[") {
			err = bfail.Mk(fmt.Sprintf("Array does not start with '[' in\n'%v'", s))
			return
		}
		if !strings.HasSuffix(s, "]") {
			err = bfail.Mk(fmt.Sprintf("Array does not end with ']' in\n'%v'", s))
			return
		}

		s2 := strings.TrimSpace(s[1 : len(s)-1])
		l := len(s2)
		if l == 0 {
			ex = expression.MkFinal(rs)
			return
		}

		i := 0
		var e string
		for {
			if i2, ok := nextByte(s2, ',', i); ok {
				e = strings.TrimSpace(s2[i:i2])
				if e == "" {
					err = bfail.Mk(fmt.Sprintf("Missing elements in\n'%v'", s))
					return
				}
				rs = append(rs, expression.MkFinal(e))
				i = i2 + 1
				continue
			}
			e = strings.TrimSpace(s2[i:l])
			if e == "" {
				err = bfail.Mk(fmt.Sprintf("Missing elements in\n'%v'", s))
				return
			}
			rs = append(rs, expression.MkFinal(e))
			break
		}
		ex = expression.MkFinal(rs)
	default:
		err = bfail.Type(args[0], "string")
	}
	return
}

/// \s -> m[s...]
func jsRo(args []*expression.T) (ex *expression.T, err error) {
	switch s := (args[0].Value).(type) {
	case string:
		s = strings.TrimSpace(s)
		v := make(map[string]*expression.T)
		if !strings.HasPrefix(s, "{") {
			err = bfail.Mk(fmt.Sprintf("Object does not start with '{' in\n'%v'", s))
			return
		}
		if !strings.HasSuffix(s, "}") {
			err = bfail.Mk(fmt.Sprintf("Object does not end with '}' in\n'%v'", s))
			return
		}
		s2 := strings.TrimSpace(s[1 : len(s)-1])
		l := len(s2)
		if l == 0 {
			ex = expression.MkFinal(v)
			return
		}
		i := 0
		var kjs string
		var k string
		var val string
		for {
			i2, ok := nextByte(s2, ':', i)
			if !ok {
				err = bfail.Mk(fmt.Sprintf("Expected ':' in\n'%v'", s2))
				return
			}
			kjs = strings.TrimSpace(s2[i:i2])
			if kjs == "" {
				err = bfail.Mk(fmt.Sprintf("Key missing in\n'%v'", s))
				return
			}
			err = json.Unmarshal([]byte(kjs), &k)
			if err != nil {
				return
			}

			i = i2 + 1

			if i2, ok := nextByte(s2, ',', i); ok {
				val = strings.TrimSpace(s2[i:i2])
				if val == "" {
					err = bfail.Mk(fmt.Sprintf("Value missing in\n'%v'", s))
					return
				}
				v[k] = expression.MkFinal(val)
				i = i2 + 1
				continue
			}
			val = strings.TrimSpace(s2[i:l])
			if val == "" {
				panic(fmt.Sprintf("Value missing in\n'%v'", s))
			}
			v[k] = expression.MkFinal(val)
			break
		}
		ex = expression.MkFinal(v)
	default:
		err = bfail.Type(args[0], "string")
	}
	return
}

/// \-> s
func jsWn(args []*expression.T) (ex *expression.T, err error) {
	ex = expression.MkFinal("null")
	return
}

/// \b -> s
func jsWb(args []*expression.T) (ex *expression.T, err error) {
	switch v := (args[0].Value).(type) {
	case bool:
		j, _ := json.Marshal(v)
		ex = expression.MkFinal(string(j))
	default:
		err = bfail.Type(args[0], "bool")
	}
	return
}

/// \i -> s
func jsWi(args []*expression.T) (ex *expression.T, err error) {
	switch v := (args[0].Value).(type) {
	case int64:
		j, _ := json.Marshal(v)
		ex = expression.MkFinal(string(j))
	default:
		err = bfail.Type(args[0], "int")
	}
	return
}

/// \f -> s
func jsWf(args []*expression.T) (ex *expression.T, err error) {
	switch v := (args[0].Value).(type) {
	case float64:
		j, _ := json.Marshal(v)
		ex = expression.MkFinal(string(j))
	default:
		err = bfail.Type(args[0], "float")
	}
	return
}

/// \f, i -> s
func jsWf2(args []*expression.T) (ex *expression.T, err error) {
	switch v := (args[0].Value).(type) {
	case float64:
		switch dec := (args[1].Value).(type) {
		case int64:
			d := int(dec)
			if d < 0 {
				d = 0
			} else if d > 9 {
				d = 9
			}
			fm := "%." + strconv.Itoa(d) + "f"

			ex = expression.MkFinal(fmt.Sprintf(fm, mathFRound(v, dec)))
		default:
			err = bfail.Type(args[0], "int")
		}
	default:
		err = bfail.Type(args[0], "float")
	}
	return
}

/// \s -> s
func jsWs(args []*expression.T) (ex *expression.T, err error) {
	switch s := (args[0].Value).(type) {
	case string:
		j, _ := json.Marshal(s)
		ex = expression.MkFinal(string(j))
	default:
		err = bfail.Type(args[0], "string")
	}
	return
}

/// \a -> s
func jsWa(args []*expression.T) (ex *expression.T, err error) {
	switch a := (args[0].Value).(type) {
	case []*expression.T:
		var b strings.Builder
		b.WriteByte('[')
		for i, js := range a {
			if i > 0 {
				b.WriteByte(',')
			}
			switch v := (js.Value).(type) {
			case string:
				b.WriteString(v)
			default:
				err = bfail.Type(js, "js-string")
			}
		}
		b.WriteByte(']')
		if err == nil {
			ex = expression.MkFinal(b.String())
		}
	default:
		err = bfail.Type(args[0], "array")
	}
	return
}

/// \m -> s
func jsWo(args []*expression.T) (ex *expression.T, err error) {
	switch m := (args[0].Value).(type) {
	case map[string]*expression.T:
		var b strings.Builder
		b.WriteByte('{')
		more := false
		for k, js := range m {
			if more {
				b.WriteByte(',')
			} else {
				more = true
			}
			jk, _ := json.Marshal(k)
			b.WriteString(string(jk))
			b.WriteByte(':')
			switch v := (js.Value).(type) {
			case string:
				b.WriteString(v)
			default:
				err = bfail.Type(js, "js-string")
			}
		}
		b.WriteByte('}')
		if err == nil {
			ex = expression.MkFinal(b.String())
		}
	default:
		err = bfail.Type(args[0], "map")
	}
	return
}

func jsGet(fname string) (fn *bfunction.T, ok bool) {
	ok = true
	switch fname {
	case "isNull":
		fn = bfunction.New(1, jsIsNull)
	case "rb":
		fn = bfunction.New(1, jsRb)
	case "ri":
		fn = bfunction.New(1, jsRi)
	case "rf":
		fn = bfunction.New(1, jsRf)
	case "rs":
		fn = bfunction.New(1, jsRs)
	case "ra":
		fn = bfunction.New(1, jsRa)
	case "ro":
		fn = bfunction.New(1, jsRo)
	case "wn":
		fn = bfunction.New(0, jsWn)
	case "wb":
		fn = bfunction.New(1, jsWb)
	case "wi":
		fn = bfunction.New(1, jsWi)
	case "wf":
		fn = bfunction.New(1, jsWf)
	case "wf2":
		fn = bfunction.New(2, jsWf2)
	case "ws":
		fn = bfunction.New(1, jsWs)
	case "wa":
		fn = bfunction.New(1, jsWa)
	case "wo":
		fn = bfunction.New(1, jsWo)
	default:
		ok = false
	}

	return
}

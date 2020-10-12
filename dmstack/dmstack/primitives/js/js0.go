// Copyright 27-Sep-2020 ºDeme
// GNU General Public License - V3 <http://www.gnu.org/licenses/>

// Json procedures.
package js

import (
	"encoding/json"
	"fmt"
	"github.com/dedeme/dmstack/machine"
	"github.com/dedeme/dmstack/symbol"
	"github.com/dedeme/dmstack/token"
	"strings"
)

// Auxiliar function
func excf(m *machine.T, template string, values ...interface{}) {
	panic(&machine.Error{
		Machine: m, Type: "Js error", Message: fmt.Sprintf(template, values...),
	})
}

// Creates a 'js' from a string.
func prFrom(m *machine.T) {
	tk := m.PopT(token.String)
	s, _ := tk.S()
	pushJs(m, s)
}

// Returns 'js' as a string.
func prTo(m *machine.T) {
	s := popJs(m)
	m.Push(token.NewS(s, m.MkPos()))
}

// Returns 'true' if s is null.
func prIsNull(m *machine.T) {
	s := popJs(m)
	m.Push(token.NewB(s == "null", m.MkPos()))
}

// Reads Bool
func prRb(m *machine.T) {
	s := popJs(m)
	var v bool
	err := json.Unmarshal([]byte(s), &v)
	if err != nil {
		excf(m, "\n  Expected: Bool value\n  Actual  : %v", s)
	}
	m.Push(token.NewB(v, m.MkPos()))
}

// Reads Int
func prRi(m *machine.T) {
	s := popJs(m)
	var v int64
	err := json.Unmarshal([]byte(s), &v)
	if err != nil {
		excf(m, "\n  Expected: Int value\n  Actual  : %v", s)
	}
	m.Push(token.NewI(v, m.MkPos()))
}

// Reads Float
func prRf(m *machine.T) {
	s := popJs(m)
	var v float64
	err := json.Unmarshal([]byte(s), &v)
	if err != nil {
		excf(m, "\n  Expected: Float value\n  Actual  : %v", s)
	}
	m.Push(token.NewF(v, m.MkPos()))
}

// Reads String
func prRs(m *machine.T) {
	s := popJs(m)
	var v string
	err := json.Unmarshal([]byte(s), &v)
	if err != nil {
		excf(m, "\n  Expected: String value\n  Actual  : %v", s)
	}
	m.Push(token.NewS(v, m.MkPos()))
}

// Reads Array
func prRa(m *machine.T) {
	s := popJs(m)
	if !strings.HasPrefix(s, "[") {
		excf(m, "Array does not start with '[' in\n'%v'", s)
	}
	if !strings.HasSuffix(s, "]") {
		excf(m, "Array does not end with ']' in\n'%v'", s)
	}

	pos := m.MkPos()
	var v []*token.T
	s2 := strings.TrimSpace(s[1 : len(s)-1])
	l := len(s2)
	if l != 0 {
		i := 0
		var e string
		for {
			if i2, ok := nextByte(s2, ',', i); ok {
				e = strings.TrimSpace(s2[i:i2])
				if e == "" {
					excf(m, "Missing elements in\n'%v'", s)
				}
				v = append(v, token.NewN(symbol.Js_, e, pos))
				i = i2 + 1
				continue
			}

			e = strings.TrimSpace(s2[i:l])
			if e == "" {
				excf(m, "Missing elements in\n'%v'", s)
			}
			v = append(v, token.NewN(symbol.Js_, e, pos))
			break
		}
	}
	m.Push(token.NewL(v, pos))
}

// Reads Object
func prRo(m *machine.T) {
	s := popJs(m)

	if !strings.HasPrefix(s, "{") {
		excf(m, "Object does not start with '{' in\n'%v'", s)
	}
	if !strings.HasSuffix(s, "}") {
		excf(m, "Object does not end with '}' in\n'%v'", s)
	}

	pos := m.MkPos()
	v := map[string]*token.T{}
	s2 := strings.TrimSpace(s[1 : len(s)-1])
	l := len(s2)

	if l != 0 {
		i := 0
		var kjs string
		var k string
		var val string
		for {
			i2, ok := nextByte(s2, ':', i)

			if !ok {
				excf(m, "Expected ':' in\n'%v'", s2)
			}
			kjs = strings.TrimSpace(s2[i:i2])
			if kjs == "" {
				excf(m, "Key missing in\n'%v'", s)
			}

			err := json.Unmarshal([]byte(kjs), &k)
			if err != nil {
				excf(m, "\n  Expected: String value\n  Actual  : %v", kjs)
			}

			i = i2 + 1
			if i2, ok := nextByte(s2, ',', i); ok {
				val = strings.TrimSpace(s2[i:i2])
				if val == "" {
					excf(m, "Value missing in\n'%v'", s)
				}
				v[k] = token.NewN(symbol.Js_, val, pos)

				i = i2 + 1
				continue
			}

			val = strings.TrimSpace(s2[i:l])
			if val == "" {
				excf(m, "Value missing in\n'%v'", s)
			}
			v[k] = token.NewN(symbol.Js_, val, pos)
			break
		}
	}

	m.Push(token.NewM(v, pos))
}

// Writes 'null'
func prWn(m *machine.T) {
	pushJs(m, "null")
}

// Writes Bool
func prWb(m *machine.T) {
	tk := m.PopT(token.Bool)
	v, _ := tk.B()
	s, _ := json.Marshal(v)
	pushJs(m, string(s))
}

// Writes Int
func prWi(m *machine.T) {
	tk := m.PopT(token.Int)
	v, _ := tk.I()
	s, _ := json.Marshal(v)
	pushJs(m, string(s))
}

// Writes Float
func prWf(m *machine.T) {
	tk := m.PopT(token.Float)
	v, _ := tk.F()
	s, _ := json.Marshal(v)
	pushJs(m, string(s))
}

// Writes String
func prWs(m *machine.T) {
	tk := m.PopT(token.String)
	v, _ := tk.S()
	s, _ := json.Marshal(v)
	pushJs(m, string(s))
}

// Writes Array
func prWa(m *machine.T) {
	tk := m.PopT(token.List)
	v, _ := tk.L()

	var b strings.Builder
	b.WriteByte('[')
	for i, tk := range v {
		if i > 0 {
			b.WriteByte(',')
		}
		sym, s, _ := tk.N()
		if sym != symbol.Js_ {
			m.Failt("\n  Expected: Json object.\n  Actual  : '%v'.", sym)
		}
		b.WriteString(s.(string))
	}
	b.WriteByte(']')

	pushJs(m, b.String())
}

// Writes Object
func prWo(m *machine.T) {
	tk := m.PopT(token.Map)
	v, _ := tk.M()

	var b strings.Builder
	b.WriteByte('{')
	more := false
	for k, tk := range v {
		if more {
			b.WriteByte(',')
		} else {
			more = true
		}

		kjs, _ := json.Marshal(k)
		b.WriteString(string(kjs))

		b.WriteByte(':')

		sym, s, _ := tk.N()
		if sym != symbol.Js_ {
			m.Failt("\n  Expected: Json object.\n  Actual  : '%v'.", sym)
		}
		b.WriteString(s.(string))
	}

	b.WriteByte('}')

	pushJs(m, b.String())
}

// Copyright 02-Mar-2022 ºDeme
// GNU General Public License - V3 <http://www.gnu.org/licenses/>

package runner

import (
	"crypto/rand"
	"encoding/base64"
	"github.com/dedeme/kut/builtin/bfail"
	"github.com/dedeme/kut/builtin/bfunction"
	"github.com/dedeme/kut/expression"
)

func mkKey(key string, lg int) string {
	k := []byte(
		key + "codified in irreversibleDeme is good, very good!\n\r8@@")

	lenk := len(k)
	sum := 0
	for i := 0; i < lenk; i++ {
		sum += int(k[i])
	}

	lg2 := lg + lenk
	r := make([]byte, lg2)
	r1 := make([]byte, lg2)
	r2 := make([]byte, lg2)
	ik := 0
	for i := 0; i < lg2; i++ {
		v1 := int(k[ik])
		v2 := v1 + int(k[v1%lenk])
		v3 := v2 + int(k[v2%lenk])
		v4 := v3 + int(k[v3%lenk])
		sum = (sum + i + v4) & 255
		r1[i] = byte(sum)
		r2[i] = byte(sum)
		ik++
		if ik == lenk {
			ik = 0
		}
	}

	for i := 0; i < lg2; i++ {
		v1 := int(r2[i])
		v2 := v1 + int(r2[v1%lg2])
		v3 := v2 + int(r2[v2%lg2])
		v4 := v3 + int(r2[v3%lg2])
		sum = (sum + v4) & 255
		r2[i] = byte(sum)
		r[i] = byte((sum + int(r1[i])) & 255)
	}

	return base64.StdEncoding.EncodeToString(r)[:lg]
}

// \s, s -> s
func crypEncode(args []*expression.T) (ex *expression.T, err error) {
	switch k := (args[0].Value).(type) {
	case string:
		switch s := (args[1].Value).(type) {
		case string:
			m := base64.StdEncoding.EncodeToString([]byte(s))
			lg := len(m)
			k := mkKey(k, lg)
			mb := []byte(m)
			kb := []byte(k)
			r := make([]byte, lg)
			for i := 0; i < lg; i++ {
				r[i] = mb[i] + kb[i]
			}
			ex = expression.MkFinal(base64.StdEncoding.EncodeToString(r))
		default:
			err = bfail.Type(args[1], "string")
		}
	default:
		err = bfail.Type(args[0], "string")
	}
	return
}

// \s, s -> s
func crypDecode(args []*expression.T) (ex *expression.T, err error) {
	switch k := (args[0].Value).(type) {
	case string:
		switch s := (args[1].Value).(type) {
		case string:
			var mb []byte
			mb, err = base64.StdEncoding.DecodeString(s)
			if err == nil {
				lg := len(mb)
				k := mkKey(k, lg)
				kb := []byte(k)
				r := make([]byte, lg)
				for i := 0; i < lg; i++ {
					r[i] = mb[i] - kb[i]
				}
				mb, err = base64.StdEncoding.DecodeString(string(r))
				if err == nil {
					ex = expression.MkFinal(string(mb))
				}
			}
		default:
			err = bfail.Type(args[1], "string")
		}
	default:
		err = bfail.Type(args[0], "string")
	}
	return
}

// \i -> s
func crypGenK(args []*expression.T) (ex *expression.T, err error) {
	switch lg := (args[0].Value).(type) {
	case int64:
		arr := make([]byte, lg)
		_, err = rand.Read(arr)
		if err == nil {
			ex = expression.MkFinal(base64.StdEncoding.EncodeToString(arr)[:lg])
		}
	default:
		err = bfail.Type(args[0], "int")
	}
	return
}

// \s, i -> s
func crypKey(args []*expression.T) (ex *expression.T, err error) {
	switch k := (args[0].Value).(type) {
	case string:
		switch lg := (args[1].Value).(type) {
		case int64:
			ex = expression.MkFinal(mkKey(k, int(lg)))
		default:
			err = bfail.Type(args[1], "int")
		}
	default:
		err = bfail.Type(args[0], "string")
	}
	return
}

func crypGet(fname string) (fn *bfunction.T, ok bool) {
	ok = true
	switch fname {
	case "encode":
		fn = bfunction.New(2, crypEncode)
	case "decode":
		fn = bfunction.New(2, crypDecode)
	case "genK":
		fn = bfunction.New(1, crypGenK)
	case "key":
		fn = bfunction.New(2, crypKey)
	default:
		ok = false
	}

	return
}

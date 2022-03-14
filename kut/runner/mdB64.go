// Copyright 02-Mar-2022 ºDeme
// GNU General Public License - V3 <http://www.gnu.org/licenses/>

package runner

import (
	"encoding/base64"
	"github.com/dedeme/kut/builtin/bfail"
	"github.com/dedeme/kut/builtin/bfunction"
	"github.com/dedeme/kut/expression"
)

func b64Decode(args []*expression.T) (ex *expression.T, err error) {
	switch s := (args[0].Value).(type) {
	case string:
		var bs []byte
		bs, err = base64.StdEncoding.DecodeString(s)
		if err == nil {
			ex = expression.MkFinal(string(bs))
		}
	default:
		err = bfail.Type(args[0], "string")
	}
	return
}

func b64DecodeBytes(args []*expression.T) (ex *expression.T, err error) {
	switch s := (args[0].Value).(type) {
	case string:
		var bs []byte
		bs, err = base64.StdEncoding.DecodeString(s)
		if err == nil {
			ex = expression.MkFinal(bs)
		}
	default:
		err = bfail.Type(args[0], "string")
	}
	return
}

func b64Encode(args []*expression.T) (ex *expression.T, err error) {
	switch s := (args[0].Value).(type) {
	case string:
		r := base64.StdEncoding.EncodeToString([]byte(s))
		if err == nil {
			ex = expression.MkFinal(r)
		}
	default:
		err = bfail.Type(args[0], "string")
	}
	return
}

func b64EncodeBytes(args []*expression.T) (ex *expression.T, err error) {
	switch bs := (args[0].Value).(type) {
	case []byte:
		r := base64.StdEncoding.EncodeToString(bs)
		if err == nil {
			ex = expression.MkFinal(r)
		}
	default:
		err = bfail.Type(args[0], "<bytes>")
	}
	return
}

func b64Get(fname string) (fn *bfunction.T, ok bool) {
	ok = true
	switch fname {
	case "decode":
		fn = bfunction.New(1, b64Decode)
	case "decodeBytes":
		fn = bfunction.New(1, b64DecodeBytes)
	case "encode":
		fn = bfunction.New(1, b64Encode)
	case "encodeBytes":
		fn = bfunction.New(1, b64EncodeBytes)
	default:
		ok = false
	}

	return
}

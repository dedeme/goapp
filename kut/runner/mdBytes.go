// Copyright 02-Mar-2022 ºDeme
// GNU General Public License - V3 <http://www.gnu.org/licenses/>

package runner

import (
	"github.com/dedeme/kut/builtin/bfail"
	"github.com/dedeme/kut/builtin/bfunction"
	"github.com/dedeme/kut/expression"
)

func bytesNew(args []*expression.T) (ex *expression.T, err error) {
	switch ix := (args[0].Value).(type) {
	case int64:
		ex = expression.MkFinal(make([]byte, ix))
	default:
		err = bfail.Type(args[0], "int")
	}
	return
}

func bytesCheckType(args []*expression.T) (ex *expression.T, err error) {
	switch (args[0].Value).(type) {
	case []byte:
		ex = expression.MkFinal(true)
	default:
		ex = expression.MkFinal(false)
	}
	return
}

func bytesFromArr(args []*expression.T) (ex *expression.T, err error) {
	switch a := (args[0].Value).(type) {
	case []*expression.T:
		var r []byte
		for _, e := range a {
			switch b := (e.Value).(type) {
			case int64:
				r = append(r, byte(b))
			default:
				err = bfail.Type(e, "int")
				break
			}
		}
		if err == nil {
			ex = expression.MkFinal(r)
		}
	default:
		err = bfail.Type(args[0], "array")
	}
	return
}

func bytesToArr(args []*expression.T) (ex *expression.T, err error) {
	switch a := (args[0].Value).(type) {
	case []byte:
		var r []*expression.T
		for _, e := range a {
			r = append(r, expression.MkFinal(int64(e)))
		}
		ex = expression.MkFinal(r)
	default:
		err = bfail.Type(args[0], "bytes")
	}
	return
}

func bytesToStr(args []*expression.T) (ex *expression.T, err error) {
	switch a := (args[0].Value).(type) {
	case []byte:
		ex = expression.MkFinal(string(a))
	default:
		err = bfail.Type(args[0], "bytes")
	}
	return
}

func bytesFromStr(args []*expression.T) (ex *expression.T, err error) {
	switch s := (args[0].Value).(type) {
	case string:
		ex = expression.MkFinal([]byte(s))
	default:
		err = bfail.Type(args[0], "string")
	}
	return
}

func bytesFget(args []*expression.T) (ex *expression.T, err error) {
	switch bs := (args[0].Value).(type) {
	case []byte:
		switch ix := (args[1].Value).(type) {
		case int64:
			ex = expression.MkFinal(int64(bs[ix]))
		default:
			err = bfail.Type(args[1], "int")
		}
	default:
		err = bfail.Type(args[0], "bytes")
	}
	return
}

func bytesSet(args []*expression.T) (ex *expression.T, err error) {
	switch bs := (args[0].Value).(type) {
	case []byte:
		switch ix := (args[1].Value).(type) {
		case int64:
			switch b := (args[2].Value).(type) {
			case int64:
				bs[ix] = byte(b)
			default:
				err = bfail.Type(args[2], "int")
			}
		default:
			err = bfail.Type(args[1], "int")
		}
	default:
		err = bfail.Type(args[0], "bytes")
	}
	return
}

func bytesSize(args []*expression.T) (ex *expression.T, err error) {
	switch a := (args[0].Value).(type) {
	case []byte:
		ex = expression.MkFinal(int64(len(a)))
	default:
		err = bfail.Type(args[0], "bytes")
	}
	return
}

func bytesGet(fname string) (fn *bfunction.T, ok bool) {
	ok = true
	switch fname {
	case "new":
		fn = bfunction.New(1, bytesNew)
	case "fromArr":
		fn = bfunction.New(1, bytesFromArr)
	case "fromStr":
		fn = bfunction.New(1, bytesFromStr)
	case "checkType":
		fn = bfunction.New(1, bytesCheckType)
	case "toArr":
		fn = bfunction.New(1, bytesToArr)
	case "toStr":
		fn = bfunction.New(1, bytesToStr)
	case "get":
		fn = bfunction.New(2, bytesFget)
	case "set":
		fn = bfunction.New(3, bytesSet)
	case "size":
		fn = bfunction.New(1, bytesSize)
	default:
		ok = false
	}

	return
}

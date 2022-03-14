// Copyright 23-Feb-2022 ºDeme
// GNU General Public License - V3 <http://www.gnu.org/licenses/>

package runner

import (
	"fmt"
	"github.com/dedeme/kut/builtin/bfail"
	"github.com/dedeme/kut/builtin/bfunction"
	"github.com/dedeme/kut/expression"
	"github.com/dedeme/kut/function"
	"github.com/dedeme/kut/runner/fail"
	"io/ioutil"
	"math/rand"
	"os"
	"os/exec"
	"time"
)

var sysFffail *function.T
var rndV int64

// \-> [s...]
func sysArgs(args []*expression.T) (ex *expression.T, err error) {
	var r []*expression.T
	for i, a := range os.Args {
		if i > 1 {
			r = append(r, expression.MkFinal(a))
		}
	}
	ex = expression.MkFinal(r)
	return
}

// \s, [s...] -> [s, s]
func sysCmd(args []*expression.T) (ex *expression.T, err error) {
	switch c := (args[0].Value).(type) {
	case string:
		switch ps := (args[1].Value).(type) {
		case []*expression.T:
			var ps2 []string
			for _, e := range ps {
				switch v := (e.Value).(type) {
				case string:
					ps2 = append(ps2, v)
				default:
					err = bfail.Type(args[0], "string")
					return
				}
			}
			cmd := exec.Command(c, ps2...)
			stdout, er := cmd.StdoutPipe()
			if er != nil {
				err = er
				return
			}
			stderr, er := cmd.StderrPipe()
			if er != nil {
				err = er
				return
			}

			err = cmd.Start()
			if err != nil {
				return
			}

			rpOut, er := ioutil.ReadAll(stdout)
			if er != nil {
				err = er
				return
			}
			rpOutS := string(rpOut)

			rpError, er := ioutil.ReadAll(stderr)
			if er != nil {
				err = er
				return
			}
			rpErrorS := string(rpError)

			er = cmd.Wait()
			if er != nil && len(rpError) == 0 {
				rpErrorS = err.Error()
			}
			ex = expression.MkFinal([]*expression.T{
				expression.MkFinal(rpOutS),
				expression.MkFinal(rpErrorS)})
		default:
			err = bfail.Type(args[1], "array")
		}
	default:
		err = bfail.Type(args[0], "string")
	}
	return
}

// \function -> ()
func sysFfail(args []*expression.T) (ex *expression.T, err error) {
	switch f := (args[0].Value).(type) {
	case *function.T:
		sysFffail = f
	default:
		err = bfail.Type(args[0], "function")
	}
	return
}

// \s->()
func sysFail(args []*expression.T) (ex *expression.T, err error) {
	switch s := (args[0].Value).(type) {
	case string:
		if sysFffail == nil {
			err = bfail.Mk(s)
		} else {
			err = fail.MkSysError(s, sysFffail)
		}
	default:
		err = bfail.Type(args[0], "string")
	}
	return
}

// \s->()
func sysPrint(args []*expression.T) (ex *expression.T, err error) {
	switch v := (args[0].Value).(type) {
	case string:
		fmt.Print(v)
	default:
		err = bfail.Type(args[0], "string")
	}
	return
}

// \s->()
func sysPrintln(args []*expression.T) (ex *expression.T, err error) {
	switch v := (args[0].Value).(type) {
	case string:
		fmt.Println(v)
	default:
		err = bfail.Type(args[0], "string")
	}
	return
}

// \->()
func sysRand(args []*expression.T) (ex *expression.T, err error) {
	rndV += time.Now().UnixMilli()
	rand.Seed(rndV)
	return
}

// \i->()
func sysSleep(args []*expression.T) (ex *expression.T, err error) {
	switch v := (args[0].Value).(type) {
	case int64:
		time.Sleep(time.Duration(v) * time.Millisecond)
	default:
		err = bfail.Type(args[0], "string")
	}
	return
}

func sysGet(fname string) (fn *bfunction.T, ok bool) {
	ok = true
	switch fname {
	case "args":
		fn = bfunction.New(0, sysArgs)
	case "cmd":
		fn = bfunction.New(2, sysCmd)
	case "fail":
		fn = bfunction.New(1, sysFail)
	case "ffail":
		fn = bfunction.New(1, sysFfail)
	case "print":
		fn = bfunction.New(1, sysPrint)
	case "println":
		fn = bfunction.New(1, sysPrintln)
	case "rand":
		fn = bfunction.New(0, sysRand)
	case "sleep":
		fn = bfunction.New(1, sysSleep)
	default:
		ok = false
	}

	return
}

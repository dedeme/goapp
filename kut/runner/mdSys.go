// Copyright 23-Feb-2022 ºDeme
// GNU General Public License - V3 <http://www.gnu.org/licenses/>

package runner

import (
	"bufio"
	"fmt"
	"github.com/dedeme/kut/builtin/bfail"
	"github.com/dedeme/kut/builtin/bfunction"
	"github.com/dedeme/kut/expression"
	"github.com/dedeme/kut/function"
	"github.com/dedeme/kut/module"
	"github.com/dedeme/kut/runner/fail"
	"io"
	"io/ioutil"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

var SysFffail *function.T
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

			rpOutS := ""
			rpOut, er := ioutil.ReadAll(stdout)
			if er != nil {
				err = er
				return
			}
			rpOutS = string(rpOut)

			rpErrorS := ""
			rpError, er := ioutil.ReadAll(stderr)
			if er != nil {
				err = er
				return
			}
			rpErrorS = string(rpError)

			er = cmd.Wait()
			if er != nil && len(rpErrorS) == 0 {
				rpErrorS = er.Error()
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

// \-> m[s]
func sysEnviron(args []*expression.T) (ex *expression.T, err error) {
	env := map[string]*expression.T{}
	for _, e := range os.Environ() {
		ix := strings.IndexByte(e, '=')
		env[e[:ix]] = expression.MkFinal(e[ix+1:])
	}
	ex = expression.MkFinal(env)
	return
}

// \i->()
func sysExit(args []*expression.T) (ex *expression.T, err error) {
	switch s := (args[0].Value).(type) {
	case int64:
		os.Exit(int(s))
	default:
		err = bfail.Type(args[0], "int")
	}
	return
}

// \s->()
func sysFail(args []*expression.T) (ex *expression.T, err error) {
	switch s := (args[0].Value).(type) {
	case string:
		if SysFffail == nil {
			err = bfail.Mk(s)
		} else {
			err = fail.MkSysError(s, SysFffail)
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
		SysFffail = f
	default:
		err = bfail.Type(args[0], "function")
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

// \-> s
func sysReadLine(args []*expression.T) (ex *expression.T, err error) {
	rd := bufio.NewReader(os.Stdin)
	s, er := rd.ReadString('\n')
	if er != nil && er != io.EOF {
		err = er
		return
	}
	ex = expression.MkFinal(s[:len(s)-1])
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

// \* -> s
func sysType(args []*expression.T) (ex *expression.T, err error) {
	var t string
	switch (args[0].Value).(type) {
	case bool:
		t = "bool"
	case int64:
		t = "int"
	case float64:
		t = "float"
	case string:
		t = "string"
	case []*expression.T:
		t = "array"
	case map[string]*expression.T:
		t = "map"
	case *function.T, *bfunction.T:
		t = "function"
	case *module.T, *BModuleT:
		t = "module"
	default:
		t = "object"
	}
	ex = expression.MkFinal(t)
	return
}

func sysGet(fname string) (fn *bfunction.T, ok bool) {
	ok = true
	switch fname {
	case "args":
		fn = bfunction.New(0, sysArgs)
	case "cmd":
		fn = bfunction.New(2, sysCmd)
	case "environ":
		fn = bfunction.New(0, sysEnviron)
	case "exit":
		fn = bfunction.New(1, sysExit)
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
	case "readLine":
		fn = bfunction.New(0, sysReadLine)
	case "sleep":
		fn = bfunction.New(1, sysSleep)
	case "type":
		fn = bfunction.New(1, sysType)
	default:
		ok = false
	}

	return
}

// Copyright 09-Jan-2021 ºDeme
// GNU General Public License - V3 <http://www.gnu.org/licenses/>

// System procedures.
package sys

import (
	"bufio"
	"fmt"
	"github.com/dedeme/dmstack/machine"
	"github.com/dedeme/dmstack/operator"
	"github.com/dedeme/dmstack/symbol"
	"github.com/dedeme/dmstack/token"
	"io/ioutil"
	"math/rand"
	"os"
	"os/exec"
	"path"
	"strings"
	"syscall"
	"time"
)

var home, udir, uname, locale string

// Auxliar function
func regularizeEs(rn rune) rune {
	switch rn {
	case 255: // á
		return 97
	case 233: // é
		return 101
	case 237: // í
		return 105
	case 243: // ó
		return 111
	case 250: // ú
		return 117
	case 252: // ü
		return 117
	case 193: // Á
		return 65
	case 201: // É
		return 69
	case 205: // Í
		return 73
	case 211: // Ó
		return 79
	case 218: // Ú
		return 85
	case 220: // Ü
		return 85
	default:
		return rn
	}
}

// Auxliar function
func regularizeEs2(rn rune) rune {
	switch rn {
	case 241: // ñ
		return 110 // n
	case 209: // Ñ
		return 78 // N
	default:
		return rn
	}
}

// Creates a locale collator which returns 1, 0 or -1, according to s1 be
// >, == or < to s2. Only available "en" and "es".
func Collator() func(s1, s2 string) int {
	if strings.HasPrefix(locale, "es_") {
		return func(s1, s2 string) int {
			if s1 == s2 {
				return 0
			}
			r1 := strings.NewReader(s1)
			r2 := strings.NewReader(s2)
			for {
				rn1, _, err1 := r1.ReadRune()
				if err1 != nil {
					_, _, err2 := r2.ReadRune()
					if err2 == nil {
						return -1
					}
					return 0
				}
				rn2, _, err2 := r2.ReadRune()
				if err2 != nil {
					return 1
				}
				rn1 = regularizeEs(rn1)
				rn2 = regularizeEs(rn2)
				if rn1 == rn2 {
					continue
				}
				rn1b := regularizeEs2(rn1)
				rn2b := regularizeEs2(rn2)
				if rn1b == rn2b {
					if rn1 == 241 || rn1 == 209 { // ñ, Ñ
						return 1
					}
					return -1
				}
				if rn1b > rn2b {
					return 1
				}
				return -1
			}
		}
	}
	return strings.Compare
}

// Initializes program.
//    m: Virtual machine.
func prInit(m *machine.T) {
	tk := m.PopT(token.String)
	s, _ := tk.S()
	if len(s) == 0 {
		m.Fail(machine.ESys(), "Empty string is not allowed for 'home' directory")
	}

	for _, e := range os.Environ() {
		if strings.HasPrefix(e, "HOME=") {
			udir = e[5:]
			continue
		}
		if strings.HasPrefix(e, "USERNAME=") {
			uname = e[9:]
			continue
		}
		if strings.HasPrefix(e, "LANG=") {
			locale = e[5:]
			continue
		}
	}

	if udir == "" {
		m.Fail(machine.ESys(), "User directory can not be set")
	}

	home = path.Join(udir, ".dmStackApp", s)
	os.MkdirAll(home, 0755)

	rand.Seed(time.Now().UnixNano())
}

// Returns 'true' if 'sys' was initialized.
//    m: Virtual machine.
func prIsInit(m *machine.T) {
	m.Push(token.NewB(home != "", m.MkPos()))
}

// Returns user name.
//    m: Virtual machine.
func prUname(m *machine.T) {
	m.Push(token.NewS(uname, m.MkPos()))
}

// Returns user dir.
//    m: Virtual machine.
func prUdir(m *machine.T) {
	m.Push(token.NewS(udir, m.MkPos()))
}

// Returns system locale.
//    m: Virtual machine.
func prLocale(m *machine.T) {
	m.Push(token.NewS(locale, m.MkPos()))
}

//  Returns program home. Its value is set with 'sys.init'.
//    m: Virtual machine.
func prHome(m *machine.T) {
	m.Push(token.NewS(home, m.MkPos()))
}

//  Returns program arguments. args[0] is the program name.
//    m: Virtual machine.
func prArgs(m *machine.T) {
	pos := m.MkPos()
	var args []*token.T
	for _, e := range os.Args {
		args = append(args, token.NewS(e, pos))
	}
	m.Push(token.NewA(args, pos))
}

// Executes an external program.
//    m: Virtual machine.
func prCmd(m *machine.T) {
	tk := m.PopT(token.Array)
	a, _ := tk.A()
	var args []string
	for _, e := range a {
		s, ok := e.S()
		if !ok {
			m.Failt(
				"\n  Expected: String value.\n  Actual:  '%v'.", e.StringDraft(),
			)
		}
		args = append(args, s)
	}
	if len(args) == 0 {
		m.Fail(machine.ESys(), "'cmd' argument list is empty.")
	}

	cmd := exec.Command(args[0], args[1:]...)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		m.Fail(machine.ESys(), "Pipe of stdout failed in cmd '%v'.", args[0])

	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		m.Fail(machine.ESys(), "Pipe of stderror failed in cmd '%v'.", args[0])
	}

	if err := cmd.Start(); err != nil {
		m.Fail(machine.ESys(), "Fail starting cmd '%v'.", args[0])
	}

	rpOut, err := ioutil.ReadAll(stdout)
	if err != nil {
		m.Fail(machine.ESys(), "Fail reading stdout of cmd '%v'.", args[0])
	}

	rpError, err := ioutil.ReadAll(stderr)
	if err != nil {
		m.Fail(machine.ESys(), "Fail reading sterror of cmd '%v'.", args[0])
	}

	if err := cmd.Wait(); err != nil {
		if rpError == nil || len(rpError) == 0 {
			m.Push(token.NewS(string(rpOut), m.MkPos()))
			m.Push(token.NewS(err.Error(), m.MkPos()))
			return
		}
	}

	m.Push(token.NewS(string(rpOut), m.MkPos()))
	m.Push(token.NewS(string(rpError), m.MkPos()))
}

// Stop program 'i' milliseconds.
//    m: Virtual machine.
func prSleep(m *machine.T) {
	tk := m.PopT(token.Int)
	i, _ := tk.I()
	time.Sleep(time.Duration(i) * time.Millisecond)
}

//  Run a free tread.
//    m: Virtual machine.
func prFreeThread(m *machine.T, run func(m *machine.T)) {
	tk := m.PopT(token.Procedure)
	go run(machine.NewThread(m.Source, m.Pmachines, tk))
}

// Run a joinable thread.
//    m: Virtual machine.
func prThread(m *machine.T, run func(m *machine.T)) {
	tk := m.PopT(token.Procedure)
	ch := make(chan int)
	m.Push(token.NewN(operator.Thread_, ch, m.MkPos()))
	go func() {
		run(machine.NewThread(m.Source, m.Pmachines, tk))
		ch <- 1
	}()
}

// Stops the progran until the thread-channel 'ch' returns a value.
//    m: Virtual machine.
func prJoin(m *machine.T) {
	tk := m.PopT(token.Native)
	sy, ch, _ := tk.N()
	if sy != operator.Thread_ {
		m.Failt(
			"Expected: Value of type <=Thread>.\nActual  : '%v'.",
			tk.StringDraft(),
		)
	}
	<-ch.(chan int)
}

// Print 's' in stdout
//    m: Virtual machine.
func prPrint(m *machine.T) {
	tk := m.PopT(token.String)
	s, _ := tk.S()
	fmt.Print(s)
}

// Print 's' in stdout
//    m: Virtual machine.
func prPrintln(m *machine.T) {
	tk := m.PopT(token.String)
	s, _ := tk.S()
	fmt.Println(s)
}

// Print 's' in stderr
//    m: Virtual machine.
func prError(m *machine.T) {
	tk := m.PopT(token.String)
	s, _ := tk.S()
	fmt.Fprintln(os.Stderr, s)
}

// Read stdin until a end of line is written. If an error happens, it
// returns "".
//    m: Virtual machine.
func prGetLine(m *machine.T) {
	reader := bufio.NewReader(os.Stdin)
	s, err := reader.ReadString('\n')
	if err != nil {
		s = ""
	}
	m.Push(token.NewS(s[:len(s)-1], m.MkPos()))
}

// Read stdin until a text equal to 'tx' is written. If an error happens, it
// returns "".
//    m: Virtual machine.
func prGetText(m *machine.T) {
	tk := m.PopT(token.String)
	tx, _ := tk.S()
	if tx == "" {
		m.Fail(machine.ESys(), "Unexpected empty string")
	}
	tx += "\n"
	var r strings.Builder
	for {
		reader := bufio.NewReader(os.Stdin)
		s, err := reader.ReadString('\n')
		if err != nil {
			r.Reset()
			break
		}
		if strings.HasSuffix(s, tx) {
			r.WriteString(s[:len(s)-len(tx)])
			break
		}
		r.WriteString(s)
	}
	m.Push(token.NewS(r.String(), m.MkPos()))
}

// Read stdin until a end of line is written. If an error happens, it
// returns "".
//    m: Virtual machine.
func prGetPass(m *machine.T) {
	attrs := syscall.ProcAttr{
		Dir:   "",
		Env:   []string{},
		Files: []uintptr{os.Stdin.Fd(), os.Stdout.Fd(), os.Stderr.Fd()},
		Sys:   nil}
	var ws syscall.WaitStatus

	pid, err := syscall.ForkExec(
		"/bin/stty",
		[]string{"stty", "-echo"},
		&attrs)
	if err != nil {
		m.Push(token.NewS("", m.MkPos()))
	}

	_, err = syscall.Wait4(pid, &ws, 0, nil)
	if err != nil {
		m.Push(token.NewS("", m.MkPos()))
	}

	reader := bufio.NewReader(os.Stdin)
	text, err := reader.ReadString('\n')
	if err != nil {
		m.Push(token.NewS("", m.MkPos()))
	}

	pid, err = syscall.ForkExec(
		"/bin/stty",
		[]string{"stty", "echo"},
		&attrs)
	if err != nil {
		panic("Sys error\nFail re-enabling echo.")
	}

	_, err = syscall.Wait4(pid, &ws, 0, nil)
	if err != nil {
		panic("Sys error\nFail re-enabling echo.")
	}

	m.Push(token.NewS(text[:len(text)-1], m.MkPos()))
}

// Processes system procedures.
//    m   : Virtual machine.
//    proc: Procedure
//    run : Function which running a machine.
func Proc(m *machine.T, proc symbol.T, run func(m *machine.T)) {
	switch proc {
	case symbol.New("init"):
		prInit(m)
	case symbol.New("initialized"):
		prIsInit(m)
	case symbol.New("uname"):
		prUname(m)
	case symbol.New("udir"):
		prUdir(m)
	case symbol.New("locale"):
		prLocale(m)
	case symbol.New("home"):
		prHome(m)
	case symbol.New("args"):
		prArgs(m)
	case symbol.New("cmd"):
		prCmd(m)
	case symbol.New("sleep"):
		prSleep(m)
	case symbol.New("freeThread"):
		prFreeThread(m, run)
	case symbol.New("thread"):
		prThread(m, run)
	case symbol.New("join"):
		prJoin(m)
	case symbol.New("print"):
		prPrint(m)
	case symbol.New("println"):
		prPrintln(m)
	case symbol.New("error"):
		prError(m)
	case symbol.New("getLine"):
		prGetLine(m)
	case symbol.New("getText"):
		prGetText(m)
	case symbol.New("getPass"):
		prGetPass(m)
	default:
		m.Failt("'sys' does not contains the procedure '%v'.", proc.String())
	}
}

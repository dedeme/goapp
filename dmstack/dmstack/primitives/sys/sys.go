// Copyright 24-Jul-2020 ºDeme
// GNU General Public License - V3 <http://www.gnu.org/licenses/>

// System procedures.
package sys

import (
	"github.com/dedeme/dmstack/machine"
	"github.com/dedeme/dmstack/symbol"
	"github.com/dedeme/dmstack/token"
	"os"
	"path"
	"strings"
	"time"
  "math/rand"
)

var home, udir, uname, locale string

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

func prInit(m *machine.T) {
	tk := m.PopT(token.String)
	s, _ := tk.S()
	if len(s) == 0 {
		m.Fail("Empty string is not allowed for 'home' directory")
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
		m.Fail("User directory can not be set")
	}

	home = path.Join(udir, ".dmStackApp", s)
  os.MkdirAll(home, 0755)

  rand.Seed(time.Now().UnixNano())
}


func prUname(m *machine.T) {
  m.Push(token.NewS(uname, m.MkPos()))
}

func prUdir(m *machine.T) {
  m.Push(token.NewS(udir, m.MkPos()))
}

func prLocale(m *machine.T) {
  m.Push(token.NewS(locale, m.MkPos()))
}

func prHome(m *machine.T) {
  m.Push(token.NewS(home, m.MkPos()))
}

func prSleep(m *machine.T) {
	tk := m.PopT(token.Int)
	i, _ := tk.I()
	time.Sleep(time.Duration(i) * time.Millisecond)
}

func prFreeThread(m *machine.T, run func(m *machine.T)) {
	tk := m.PopT(token.Procedure)
	p, _ := tk.P()
	go run(machine.NewIsolate(m.SourceDir, m.Pmachines, p))
}

func prThread(m *machine.T, run func(m *machine.T)) {
	tk := m.PopT(token.Procedure)
	p, _ := tk.P()
	ch := make(chan int)
	m.Push(token.NewN(symbol.Thread_, ch, m.MkPos()))
	go func() {
		run(machine.NewIsolate(m.SourceDir, m.Pmachines, p))
		ch <- 1
	}()
}

func prJoin(m *machine.T) {
	tk := m.PopT(token.Native)
	sy, ch, _ := tk.N()
	if sy != symbol.Thread_ {
		m.Failf(
			"Expected: Value of type <= Thread>.\nActual  : %v.",
			tk.StringDraft(),
		)
	}
	<-ch.(chan int)
}

// Processes system procedures.
func Proc(m *machine.T, run func(m *machine.T)) {
	tk, ok := m.PrgNext()
	if !ok {
		m.Fail("'sys' procedure is missing")
	}
	sy, ok := tk.Sy()
	if !ok {
		m.Failf("Expected: 'sys' procedure.\nActual  : %v.", tk.StringDraft())
	}
	switch sy {
	case symbol.New("init"):
		prInit(m)
	case symbol.New("uname"):
		prUname(m)
	case symbol.New("udir"):
		prUdir(m)
	case symbol.New("locale"):
		prLocale(m)
	case symbol.New("home"):
		prHome(m)
	case symbol.New("sleep"):
		prSleep(m)
	case symbol.New("freeThread"):
		prFreeThread(m, run)
	case symbol.New("thread"):
		prThread(m, run)
	case symbol.New("join"):
		prJoin(m)
	default:
		m.Failf("'sys' does not contains the procedure '%v'.", sy.String())
	}
}

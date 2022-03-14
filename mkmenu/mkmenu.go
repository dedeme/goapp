// Copyright 13-Jun-2021 ºDeme
// GNU General Public License - V3 <http://www.gnu.org/licenses/>

// Start program.
package main

import (
	"fmt"
	"github.com/dedeme/golib/file"
	"os"
	"strings"
)

const (
	FMENU = "menu.txt"
	FFLUX = /*"./flmenu.txt" //*/ "/home/deme/.fluxbox/menu"
)

type tab struct {
	v string
}

func processLine(out *os.File, nline int, l string, tb *tab) {
	l = strings.TrimSpace(l)
	if l == "" || l[0] == '#' {
		return
	}

	switch l[0] {
	case '-':
		out.Write([]byte(fmt.Sprintf("%v[nop]\n", tb.v)))
	case '/':
		out.Write([]byte(fmt.Sprintf(
			"%v[submenu] (%v)\n", tb.v, l[1:])))
		tb.v += "  "
	case '<':
		if len(tb.v) > 2 {
			tb.v = tb.v[2:]
		}
		out.Write([]byte(fmt.Sprintf("%v[end]\n", tb.v)))
	case '*':
		fields := strings.Split(l, "||")
		lfields := len(fields)
		switch lfields {
		case 2, 3:
			if lfields == 3 {
				fields[1] = fields[2]
			}
			out.Write([]byte(fmt.Sprintf("%v[exec] (%v) {%v}\n",
				tb.v, fields[0][1:], fields[1])))
		default:
			fmt.Printf(
				"Bad number of fields (%v).\n%v: %v\n",
				lfields, nline, l)
		}
	default:
		fmt.Printf(
			"Invalid directive. A line must start with '/', '-', '>' or '*'.\n%v: %v\n",
			nline, l)
	}
}

// Program entry.
func main() {
  if !file.Exists(FMENU) {
    panic(FMENU + " not found\nTo run tests change name of menuXX.txt")
  }

	out := file.OpenWrite(FFLUX)
	out.Write([]byte("[begin] (Principal)\n"))

	nline := 1
	tb := &tab{"  "}
	file.Lines(FMENU, func(ln string) bool {
		processLine(out, nline, ln, tb)
		nline++
		return false
	})

	out.Write([]byte("" +
		"  [nop]\n" +
		"  [submenu] (fluxbox)\n" +
		"    [include] (/etc/X11/fluxbox/fluxbox-menu)\n" +
		"  [end]\n" +
		"[end]\n"))

	out.Close()

	fmt.Println("Menú terminado\n")
}

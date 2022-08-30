// Copyright 13-Jun-2021 ºDeme
// GNU General Public License - V3 <http://www.gnu.org/licenses/>

// Start program.
package main

import (
	"fmt"
	"github.com/dedeme/ktlib/file"
	"github.com/dedeme/ktlib/arr"
	"github.com/dedeme/ktlib/str"
	"strings"
)

const (
	FMENU = "menu.txt"
	FFLUX = /*"./flmenu.txt" //*/ "/home/deme/.fluxbox/menu"
)

type tab struct {
	v string
}

func processLine(out *file.T, nline int, l string, tb *tab) {
	l = strings.TrimSpace(l)
	if l == "" || l[0] == '#' {
		return
	}

	switch l[0] {
	case '-':
		file.WriteText(out, fmt.Sprintf("%v[nop]\n", tb.v))
	case '/':
		file.WriteText(out, fmt.Sprintf(
			"%v[submenu] (%v)\n", tb.v, l[1:]))
		tb.v += "  "
	case '<':
		if len(tb.v) > 2 {
			tb.v = tb.v[2:]
		}
		file.WriteText(out, fmt.Sprintf("%v[end]\n", tb.v))
	case '*':
		fields := strings.Split(l, "||")
		lfields := len(fields)
		switch lfields {
		case 2, 3:
			if lfields == 3 {
				fields[1] = fields[2]
			}
			file.WriteText(out, fmt.Sprintf("%v[exec] (%v) {%v}\n",
				tb.v, fields[0][1:], fields[1]))
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

	out := file.Wopen(FFLUX)
	file.WriteText(out, "[begin] (Principal)\n")

	tb := &tab{"  "}
	arr.EachIx(str.Split(file.Read(FMENU), "\n"), func(ln string, nline int) {
		processLine(out, nline, ln, tb)
		nline++
	})

	file.WriteText(out, "" +
		"  [nop]\n" +
		"  [submenu] (fluxbox)\n" +
		"    [include] (/etc/X11/fluxbox/fluxbox-menu)\n" +
		"  [end]\n" +
		"[end]\n")

	file.Close(out)

	fmt.Println("Menú terminado\n")
}

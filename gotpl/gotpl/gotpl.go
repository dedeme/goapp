// Copyright 08-Apr-2020 ºDeme
// GNU General Public License - V3 <http://www.gnu.org/licenses/>

package main

import (
	"fmt"
	"github.com/dedeme/golib/file"
	"github.com/dedeme/golib/sys"
	"github.com/dedeme/gotpl/tpls"
	"os"
	"strings"
)

func process(l string) (r string, err string) {
	key, value := tpls.Kv(l)

	switch key {
	case "for":
		r, err = tpls.Pfor(value)
	case "each":
		r, err = tpls.Peach(value)
	case "map":
		r, err = tpls.Pmap(value)
	case "mapTo":
		r, err = tpls.PmapTo(value)
	case "mapFrom":
		r, err = tpls.PmapFrom(value)
	case "filter":
		r, err = tpls.Pfilter(value)
	case "sort":
		r, err = tpls.Psort(value)
	case "props":
		r, err = tpls.Pprops(value)
	case "toJs":
		r, err = tpls.PtoJs(value)
	case "fromJs":
		r, err = tpls.PfromJs(value)
	case "copyright":
		r, err = tpls.Pcopyright(value)
	default:
		err = fmt.Sprintf("Key '%v' is unknown", key)
	}
	//fmt.Println(r)
	return
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("File to format is missing.\nUse:\ngotpl <goFile>")
		return
	}

	f := os.Args[1]
	if !strings.HasSuffix(f, ".go") {
		fmt.Println("File must end in '.go'")
		os.Exit(1)
	}

	code := strings.Split(file.ReadAll(f), "\n")
	changed := false
	var newCode []string
	for _, l := range code {
		l2 := strings.TrimSpace(l)
		if strings.HasPrefix(l2, "·") {
			tx, emsg := process(l2[2:])
			if emsg != "" {
				fmt.Println(emsg)
				os.Exit(1)
			}
			newCode = append(newCode, tx)
			changed = true
		} else {
			newCode = append(newCode, l)
		}
	}

	if changed {
		file.WriteAll(f, strings.Join(newCode, "\n"))

		_, err := sys.Cmd("go", "fmt", f)
		if len(err) > 0 {
			fmt.Println("go format error:")
			fmt.Println(string(err))
			os.Exit(1)
		}
	}
}

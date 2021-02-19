// Copyright 14-Jan-2021 ºDeme
// GNU General Public License - V3 <http://www.gnu.org/licenses/>

// Program arguments management.
package args

import (
	"fmt"
	"github.com/dedeme/dmcoffee/cts"
	"os"
)

// If program must be ran in production mode.
var Production = true

// Arguments of dmcoffee. Currently only:
//    -d : Running in debug mode.
//    dms : Path of file .dms to run.
var CoffeeArgs = make(map[string]string)

// Arguments of program .dms.
var Args = []string{}

func help() {
	fmt.Println(
		"dmcoffee v. " + cts.Version + ".\n" +
			"Copyright 14-Jun-2021 ºDeme\n" +
			"GNU General Public License - V3 <http://www.gnu.org/licenses/>\n" +
			"\n" +
			"Use:\n" +
			"  dmcoffee [Options] dmcProgram [-- dmcProgramOptions]\n" +
			"or to show this message:\n" +
			"  dmcoffee\n" +
			"-----------------------\n" +
			"Options:\n" +
			"  -d : Execute in debug mode. \n" +
			"dmcProgram: Path to .dmc file, with or without extension .dmc.\n" +
			"dmcProgramOptions: Options to be passed to .dmc program.\n",
	)
}

// Extracts arguments.
func Initialize() (ok bool) {
	argc := len(os.Args)
	if argc < 2 {
		help()
		return
	}

	c := 1
	for ; c < argc; c++ {
		a := os.Args[c]

		if a[0] == '-' {
			if a == "--" {
				break
			} else if a == "-d" {
				Production = false
				CoffeeArgs["-d"] = ""
				continue
			}

			fmt.Printf("Unkown option '%v'\n\n=======================\n", a)
			help()
			return
		} else {
			CoffeeArgs[cts.SourceExtension] = a
			c++
			break
		}
	}

	if _, ok := CoffeeArgs[cts.SourceExtension]; !ok {
		fmt.Println(
			"Name of ." + cts.SourceExtension + " file is missing.\n\n" +
				"=======================",
		)
		help()
		return false
	}

	if c == argc {
		return true
	}

	if os.Args[c] != "--" {
		fmt.Println(
			"Expected '--'.\n\n" +
				"=======================",
		)
		help()
		return
	}

	c++
	for ; c < argc; c++ {
		Args = append(Args, os.Args[c])
	}

	return true
}

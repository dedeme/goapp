// Copyright 25-Apr-2020 ºDeme
// GNU General Public License - V3 <http://www.gnu.org/licenses/>

// Program arguments management.
package args

import (
	"fmt"
	"os"
  "github.com/dedeme/dmstack/cts"
)

// If program must be ran in production mode.
var Production = true

// Arguments of dmstack. Currently only:
//    -d : Running in debug mode.
//    dms : Path of file .dms to run.
var Stkargs = make(map[string]string)

// Arguments of program .dms.
var Args = []string{}

func help() {
	fmt.Println(
		"dmstack v. " + cts.Version + ".\n" +
			"Copyright 25-Apr-2020 ºDeme\n" +
			"GNU General Public License - V3 <http://www.gnu.org/licenses/>\n" +
			"\n" +
			"Use:\n" +
			"  dmstack [Options] dmsProgram [-- dmsProgramOptions]\n" +
			"or to show this message:\n" +
			"  dmstack\n" +
			"-----------------------\n" +
			"Options: They are\n" +
			"  -d : Execute in debug mode. \n" +
			"dmsProgram: Path to .dms file, with or without extension .dms.\n" +
			"dmsProgramOptions: Options to be passed to .dms program.\n",
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
        Stkargs["-d"] = ""
        continue
      }

      fmt.Printf ("Unkown option '%v'\n\n=======================\n", a)
      help()
      return
    } else {
      Stkargs["dms"] = a
      c++
      break
    }
  }

  if _, ok := Stkargs["dms"]; !ok {
    fmt.Println(
      "Name of .dms file is missing.\n\n"+
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
      "Expected '--'.\n\n"+
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

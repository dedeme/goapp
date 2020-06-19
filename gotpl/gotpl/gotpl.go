// Copyright 08-Apr-2020 ºDeme
// GNU General Public License - V3 <http://www.gnu.org/licenses/>

package main

import (
	"fmt"
	"github.com/dedeme/golib/file"
	"github.com/dedeme/golib/sys"
	"github.com/dedeme/gotpl/tpls/directs"
	"github.com/dedeme/gotpl/tpls/structs"
	"os"
	"strings"
)

type state int
const (
  normalCode state = iota
  definition1
  definition2
)

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
  startStruct := -1
  var structCode []string
	var newCode []string
  st := normalCode
	for i, l := range code {
    l2 := strings.TrimSpace(l)
    if st == definition1 {
      newCode = append(newCode, l)
      structCode = append(structCode, l2)
      if l2 == "*/" {
        st = definition2
      }
    } else if st == definition2 {
      if l2 == "//===" {
        st = normalCode
        newCd, err := structs.Process(startStruct, structCode)
        startStruct = -1
        if err != nil {
          fmt.Println(err.Error())
          os.Exit(1)
        }
        newCode = append(newCode, newCd...)
        newCode = append(newCode, l)
      }
    } else {
      if strings.HasPrefix(l2, "·") {
        tx, emsg := directs.Process(l2[2:])
        if emsg != "" {
          fmt.Println(emsg)
          os.Exit(1)
        }
        newCode = append(newCode, tx)
        changed = true
      } else if l2 == "/*·" {
        newCode = append(newCode, "/* ·")
        structCode = []string{}
        st = definition1
        startStruct = i + 1
        changed = true
      } else {
        newCode = append(newCode, l)
      }
    }
	}
  if startStruct != -1 {
    fmt.Printf("Line %v: End of file reached", startStruct)
    os.Exit(1)
  }

	if changed {
		file.WriteAll(f, strings.Join(newCode, "\n"))

		_, err := sys.Cmd("go", "fmt", f)
		if len(err) > 0 {
			fmt.Println("go format error:")
			fmt.Println(string(err))
			os.Exit(1)
		}
	} else {
    fmt.Println("No template fund")
    os.Exit(1)
  }
}

// Copyright 07-Nov-2020 ºDeme
// GNU General Public License - V3 <http://www.gnu.org/licenses/>

package zenity

import (
	"github.com/dedeme/golib/sys"
	"os"
	"strings"
)

func cwd () string {
  if wd, err := os.Getwd(); err == nil {
    return wd
  }
  msg := "Working directory not found"
  Msgbox(msg)
  panic(msg)
}

func Menu() string {
	path := cwd()
	out, _ := sys.Cmd(
		"zenity", "--title", "Hxtpl", "--list",
		"--text", "Path:\n"+path,
		"--radiolist",
		"--column", "", "--column", "Option",
		"FALSE", "New Semicomponent",
	)
	return strings.TrimSpace(string(out))
}

func Msgbox(tx string) {
	sys.Cmd(
		"zenity", "--title", "Hxtpl", "--info",
		"--text", tx,
	)
}

func Input (title, tx string) string {
		path :=  cwd()
    out, _ := sys.Cmd(
		"zenity", "--title", "Hxtpl: " + title, "--entry",
		"--text", "Path:\n"+path+tx,
    )
  return strings.TrimSpace(string(out))
}

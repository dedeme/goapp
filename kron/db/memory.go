// Copyright 07-Jul-2020 ºDeme
// GNU General Public License - V3 <http://www.gnu.org/licenses/>

package db

import (
	"github.com/dedeme/golib/file"
	"github.com/dedeme/golib/json"
	"github.com/dedeme/golib/sys"
	"path"
)

var memPath string

func MemInit() {
	memPath = path.Join(sys.Home(), "Memory.tb")
	if !file.Exists(memPath) {
		MemReset()
	}
}

func MemRead () string {
  return json.FromString(file.ReadAll(memPath)).Rs()
}

func MemWrite (command string) {
  file.WriteAll(memPath, json.Ws(command).String())
}

func MemReset () {
  MemWrite("")
}

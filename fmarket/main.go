// Copyright 01-Dic-2022 ºDeme
// GNU General Public License - V3 <http://www.gnu.org/licenses/>

// Main file.
package main

import (
	"github.com/dedeme/fmarket/data/cts"
	"github.com/dedeme/fmarket/data/model"
	"github.com/dedeme/fmarket/db"
	"github.com/dedeme/fmarket/start"
	"github.com/dedeme/fmarket/tests"
	"github.com/dedeme/ktlib/arr"
	"github.com/dedeme/ktlib/js"
	"github.com/dedeme/ktlib/log"
	"github.com/dedeme/ktlib/sys"
)

func help() {
	sys.Println(
		"Use:\n" +
			"  fmarket [parameter]\n" +
			"Where parameter can be:\n" +
			"  start : Starts daily calculations.\n" +
			"  models: Returns a JSON array with model identifiers\n" +
			"  test  : Makes program tests.\n" +
			"  help  : Shows this message.\n",
	)
}

func main() {
	args := sys.Args()
	if len(args) != 2 {
		help()
		return
	}

	sys.Rand()
	db.Initialize()
	log.Initialize(cts.LogPath)

	switch args[1] {
	case "start":
		start.Run()
	case "models":
		sys.Println(js.Wa(arr.Map(model.List(), func(md *model.T) string {
			return js.Ws(md.Id())
		})))
	case "test":
		tests.All()
	case "help":
		help()
	default:
		panic("Unexpected argument '" + args[1] + "'")
	}
}

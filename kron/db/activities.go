// Copyright 07-Jul-2020 ºDeme
// GNU General Public License - V3 <http://www.gnu.org/licenses/>

package db

import (
	"github.com/dedeme/golib/file"
	"github.com/dedeme/golib/json"
	"github.com/dedeme/golib/sys"
	"github.com/dedeme/kron/data"
	"path"
)

var actPath string

func ActInit() {
	actPath = path.Join(sys.Home(), "Activities.tb")
	if !file.Exists(actPath) {
		file.WriteAll(actPath, "[]")
	}
}

func ActRead() (acts []*data.ActEntry) {
	for _, e := range json.FromString(file.ReadAll(actPath)).Ra() {
		acts = append(acts, data.ActFromJs(e))
	}
	return
}

func ActWrite(acts []*data.ActEntry) {
	var a []json.T
	for _, e := range acts {
		a = append(a, e.ToJs())
	}
	file.WriteAll(actPath, json.Wa(a).String())
}

func ActAdd(e *data.ActEntry) {
	acts := ActRead()
	if len(acts) == 0 {
		e.Id = 0
	} else {
		e.Id = acts[len(acts)-1].Id + 1
	}
	acts = append(acts, e)
	ActWrite(acts)
}

func ActDel(id int) (ok bool) {
	var acts []*data.ActEntry
	for _, e := range ActRead() {
		if e.Id != id {
			acts = append(acts, e)
		} else {
      ok = true
    }
	}
  if ok {
    ActWrite(acts)
  }
  return
}

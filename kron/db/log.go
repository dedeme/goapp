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

var logPath string

func LogInit() {
	logPath = path.Join(sys.Home(), "Log.tb")
	if !file.Exists(logPath) {
		file.WriteAll(logPath, "[]")
	}
}

func LogRead() (log []*data.LogEntry) {
	logJs := json.FromString(file.ReadAll(logPath)).Ra()
	for _, e := range logJs {
		log = append(log, data.LogFromJs(e))
	}
	return
}

func logAdd(entry *data.LogEntry) {
	old := LogRead()
	l := len(old)
	var new []json.T
	for i, e := range old {
		if l-i < 100 {
			new = append(new, e.ToJs())
		}
	}
	new = append(new, entry.ToJs())
	file.WriteAll(logPath, json.Wa(new).String())
}

func LogInfo(msg string) {
	entry := data.NewLogInfo(msg)
	logAdd(entry)
}

func LogError(msg string) {
	entry := data.NewLogError(msg)
	logAdd(entry)
}

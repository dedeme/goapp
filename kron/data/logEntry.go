// Copyright 07-Jul-2020 ºDeme
// GNU General Public License - V3 <http://www.gnu.org/licenses/>

package data

import (
	"github.com/dedeme/golib/date"
	"github.com/dedeme/golib/json"
	"runtime/debug"
	"strings"
)

type LogEntry struct {
	isError bool
	time    string
	msg     string
}

func time() string {
	return date.Now().Format("%D/%M/%Y(%t)")
}

func NewLogInfo(msg string) *LogEntry {
	return &LogEntry{isError: false, time: time(), msg: msg}
}

func NewLogError(msg string) *LogEntry {
	ls1 := strings.Split(string(debug.Stack()), "\n")
	var ls []string
	for _, e := range ls1 {
		ls = append(ls, "    "+e)
	}
	stack := strings.Join(ls, "\n")
	return &LogEntry{isError: true, time: time(), msg: msg + "\n" + stack}
}

func (e *LogEntry) IsError() bool {
	return e.isError
}

func (e *LogEntry) String() string {
	tp := " - "
	if e.isError {
		tp = " = "
	}
	return e.time + tp + e.msg
}

func (e *LogEntry) ToJs() json.T {
	return json.Wa([]json.T{
		json.Wb(e.isError),
		json.Ws(e.time),
		json.Ws(e.msg),
	})
}

func LogFromJs(js json.T) *LogEntry {
	a := js.Ra()
	return &LogEntry{
		isError: a[0].Rb(),
		time:    a[1].Rs(),
		msg:     a[2].Rs(),
	}
}

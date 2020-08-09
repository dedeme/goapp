// Copyright 07-Jul-2020 ºDeme
// GNU General Public License - V3 <http://www.gnu.org/licenses/>

package db

import (
	"github.com/dedeme/golib/file"
	"github.com/dedeme/golib/sys"
	"github.com/dedeme/golib/date"
	"path"
)

// Return "" if 'start.txt' already exists.
func StartOn() (timeStamp string) {
	p := path.Join(sys.Home(), "start.txt")
	if file.Exists(p) {
		return
	}
  timeStamp = date.Now().Format("%Y%M%D-%T")
  file.WriteAll(p, timeStamp);
  return
}

func StartTimeStamp () string {
  p := path.Join(sys.Home(), "start.txt")
  if !file.Exists(p) {
    return ""
  }
  return file.ReadAll(p)
}

func StartOf() {
  file.Remove(path.Join(sys.Home(), "start.txt"))
}


// Copyright 24-Jul-2021 ºDeme
// GNU General Public License - V3 <http://www.gnu.org/licenses/>

// Filesystem operations for restoring.
package restore

import (
	"github.com/dedeme/FastBackups/data/cts"
	"github.com/dedeme/FastBackups/db/control"
	"github.com/dedeme/golib/cryp"
	"github.com/dedeme/golib/file"
	"github.com/dedeme/golib/sys"
	"os"
	"path"
	"strings"
)

// Restores backups
//    key       : Key to decode tar files.
//    failReport: Function to inform of fails.
func Run(key string, failReport func(string)) {
	d := cts.DRestoreSource()
	var tars []string
	for _, info := range file.List(d) {
		name := info.Name()
		if strings.HasSuffix(name, ".tgz") {
			tars = append(tars, name)
		}
	}

	if len(tars) == 0 {
		failReport("No '.tgz' found in\n'" + d + "'")
		return
	}

	dTarget := cts.DRestoreTarget()
	file.Remove(dTarget)
	file.Mkdir(dTarget)

	os.Chdir(dTarget)
	for _, tar := range tars {
		ftarget := file.OpenWrite(tar)
		defer func() {
			ftarget.Close()
		}()

		first := true
		file.Lines(path.Join("..", "source", tar), func(l string) bool {
			if first {
				first = false
				return false
			}
			file.Write(ftarget, cryp.Decryp(key, l))
			return false
		})

		cout, cerr := sys.Cmd("tar", "-xf", tar)
		if cerr == nil || len(cerr) > 0 {
			failReport("Fail restoring " + tar + "\n" + string(cerr))
			return
		}
		if cout == nil || len(cout) > 0 {
			failReport("Fail restoring " + tar + "\n" + string(cout))
			return
		}

		file.Remove(tar)
	}
}

// Returns key control and file name of one .tgz or 'ok = false' if there
// is no one.
func KeyControl() (key, controlFile string, ok bool) {
	d := cts.DRestoreSource()
	for _, info := range file.List(d) {
		name := info.Name()
		if strings.HasSuffix(name, ".tgz") {
			controlFile = path.Join("restore", "source", name)
			f := path.Join(d, name)
			key = control.FileKey(f)
			ok = true
			break
		}
	}
	return
}

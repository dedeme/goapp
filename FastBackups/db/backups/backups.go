// Copyright 24-Jul-2021 ºDeme
// GNU General Public License - V3 <http://www.gnu.org/licenses/>

// Filesystem operations for backups.
package backups

import (
	"github.com/dedeme/FastBackups/data/cts"
	"github.com/dedeme/FastBackups/db/control"
	"github.com/dedeme/golib/cryp"
	"github.com/dedeme/golib/file"
	"github.com/dedeme/golib/sys"
	"io"
	"os"
	"path"
	"strings"
)

func compress(key, source, name, target string, failReport func(string)) {
	os.Chdir(source)
	tar := name + ".tgz"
	cout, cerr := sys.Cmd("tar", "-czf", tar, name)
	if cerr == nil || len(cerr) > 0 {
		failReport("Fail compressing " + tar + "\n" + string(cerr))
	}
	if cout == nil || len(cout) > 0 {
		failReport("Fail compressing " + tar + "\n" + string(cout))
	}

	fsource := file.OpenRead(tar)
	bs := make([]byte, 500)
	ftarget := file.OpenWrite(path.Join(target, tar))

	defer func() {
		fsource.Close()
		ftarget.Close()
		file.Remove(tar)
	}()

	n, err := fsource.Read(bs)
	if n == 0 || err == io.EOF {
		return
	}

	code := cryp.Cryp(key, string(bs[:n]))
	keyCode := cryp.Cryp(code, cryp.Key(key, 450))

	file.Write(ftarget, keyCode+"\n")
	file.Write(ftarget, code+"\n")

	for {
		n, err := fsource.Read(bs)
		if n == 0 || err == io.EOF {
			return
		}
		file.Write(ftarget, cryp.Cryp(key, string(bs[:n]))+"\n")
	}
}

// Makes backups
//    key       : Key to encode tar files.
//    failReport: Function to inform of fails.
func Run(key string, failReport func(string)) {
	d := cts.DBackupsSource()
	var dirs []string
	for _, info := range file.List(d) {
		dpath := path.Join(d, info.Name())
		if file.IsDirectory(dpath) {
			dirs = append(dirs, dpath)
		}
	}

	if len(dirs) == 0 {
		failReport("No directories found in\n'" + d + "'")
		return
	}

	dTarget := cts.DBackupsTarget()
	for _, info := range file.List(dTarget) {
		file.Remove(path.Join(dTarget, info.Name()))
	}

	for _, d := range dirs {
		base := path.Base(d)
		d2 := path.Join(d, base)
		if file.Exists(d2) {
			compress(key, d, base, dTarget, failReport)
		} else {
			failReport(d2 + "\nnot found.")
		}
	}
}

// Returns key control and file name of one .tgz or 'ok = false' if there
// is no one.
func KeyControl() (key, controlFile string, ok bool) {
	d := cts.DBackupsTarget()
	for _, info := range file.List(d) {
		name := info.Name()
		if strings.HasSuffix(name, ".tgz") {
			controlFile = path.Join("backups", "target", name)
			f := path.Join(d, name)
			key = control.FileKey(f)
			ok = true
			break
		}
	}
	return
}

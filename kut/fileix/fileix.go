// Copyright 02-Feb-2022 ºDeme
// GNU General Public License - V3 <http://www.gnu.org/licenses/>

// Index of code files.
package fileix

import (
	"errors"
	"fmt"
	"github.com/dedeme/golib/file"
	"os"
	"path"
)

var Root = "./"
var paths []string

// Add path to code file paths and returns its index.
// If path already has been added, it do nothing and also returns its index.
//    fix: Index of file which imports 'ipath' or -1 if 'ipath' is the main file.
//    ipath: Import file path. It can be:
//             -Absolute          : In any case.
//             -Relative to 'fix' : fix must be != -1 and 'path.parent(fix)/ipath'
//                                  must be found in the file system.
//             -Relative to 'Root': In any case.
func Add(fix int, ipath string) int {
	if fix != -1 && !path.IsAbs(ipath) {
		newPath := path.Join(path.Dir(paths[fix]), ipath)
		if file.Exists(path.Join(Root, newPath+".kut")) {
			ipath = newPath
		}
	}
	for ix, p := range paths {
		if p == ipath {
			return ix
		}
	}
	paths = append(paths, ipath)
	return len(paths) - 1
}

// Returns the path with index 'ix' "shorted" to 50 bytes or raise 'panic'.
func Get(ix int) (s string) {
	if ix < 0 {
		s = "Built-in"
		return
	}
	s = paths[ix]
	if !path.IsAbs(s) {
		s = path.Join(Root, s)
	}
	if len(s) > 50 {
		s = "..." + s[len(s)-47:]
	}
	return
}

func GetFail(ix int) (s string) {
	if ix < 0 {
		return "Built-in"
	}
	if !path.IsAbs(s) {
		s = path.Clean(path.Join(Root, paths[ix]+".kut"))
	}
	if !path.IsAbs(s) {
		dir, err := os.Getwd()
		if err == nil {
			s = path.Clean(path.Join(dir, s))
		}
	}
	return
}

// Read the file with index 'ix'
func Read(ix int) (code string, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = errors.New(fmt.Sprintf("%v", e))
		}
	}()

	p := paths[ix]
	if !path.IsAbs(p) {
		p = path.Join(Root, p)
	}

	code = file.ReadAll(p + ".kut")
	return
}

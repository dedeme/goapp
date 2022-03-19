// Copyright 02-Feb-2022 ºDeme
// GNU General Public License - V3 <http://www.gnu.org/licenses/>

// Index of code files.
package fileix

import (
	"errors"
	"fmt"
	"github.com/dedeme/golib/file"
  "path"
)

var Root = "./";
var paths []string


// Add path to code file paths and returns its index.
// If path already has been added, it do nothing and also returns its index.
func Add(fpath string) int {
	for ix, p := range paths {
		if p == fpath {
			return ix
		}
	}
	paths = append(paths, fpath)
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

// Returns the complete path with index 'ix' or raise 'panic'.
func GetComplete(ix int) (s string) {
	if ix < 0 {
		return "Built-in"
	}
  s = paths[ix]
  if !path.IsAbs(s) {
    s = path.Join(Root, s)
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

	code = file.ReadAll(GetComplete(ix) + ".kut")
	return
}

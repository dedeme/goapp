// Copyright 24-Jul-2021 ºDeme
// GNU General Public License - V3 <http://www.gnu.org/licenses/>

// Algorithm to control password from files.
package control

import (
	"github.com/dedeme/golib/cryp"
	"github.com/dedeme/golib/file"
)

// Returns the key of a tar file.
func FileKey(fpath string) string {
	var lines []string
	file.Lines(fpath, func(l string) bool {
		lines = append(lines, l)
		if len(lines) > 1 {
			return true
		}
		return false
	})
	if len(lines) == 2 {
		return cryp.Decryp(lines[1], lines[0])
	}
	return "%"
}

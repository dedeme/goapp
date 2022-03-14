// Copyright 24-Jul-2021 ºDeme
// GNU General Public License - V3 <http://www.gnu.org/licenses/>

// Data base initialization
package db

import (
	"github.com/dedeme/FastBackups/data/cts"
	"github.com/dedeme/golib/file"
)

// Initializes data base
func Initialize() {
	if !file.Exists(cts.Root()) {
		file.Mkdirs(cts.DBackupsSource())
		file.Mkdirs(cts.DBackupsTarget())
		file.Mkdirs(cts.DRestoreSource())
		file.Mkdirs(cts.DRestoreTarget())
	}
}

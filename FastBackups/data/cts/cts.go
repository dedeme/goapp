// Copyright 24-Jul-2021 ºDeme
// GNU General Public License - V3 <http://www.gnu.org/licenses/>

// Program constants
package cts

import (
	"github.com/dedeme/golib/file"
	"path"
)

func Root() string {
	return path.Join(file.HomeDir(), "FastBackups")
}

// Backups directory
func DBackups() string {
	return path.Join(Root(), "backups")
}

// Backups source directory
// This directory contains links to the directory parent of that which are
// going to be backed (orginal directory).
// Links must have the same name that the original directory.
func DBackupsSource() string {
	return path.Join(DBackups(), "source")
}

// Backups target directory.
// Directory containing encoded .tar's comming from 'dBackupSource'
func DBackupsTarget() string {
	return path.Join(DBackups(), "target")
}

// Restore directory
func DRestore() string {
	return path.Join(Root(), "restore")
}

// Restore source directory
// Directory containing encoded .tar's, placed manually.
// It can have several files, but they should not be fron the same orignal
// directory.
func DRestoreSource() string {
	return path.Join(DRestore(), "source")
}

// Restore target directory
// After click 'Restore' it will contain orginal directories extracted from
// .tar's in 'source'.
func DRestoreTarget() string {
	return path.Join(DRestore(), "target")
}

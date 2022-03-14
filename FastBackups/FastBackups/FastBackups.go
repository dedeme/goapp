// Copyright 22-Jan-2021 ºDeme
// GNU General Public License - V3 <http://www.gnu.org/licenses/>

// Start program.
package main

import (
	"github.com/dedeme/FastBackups/db"
	"github.com/dedeme/FastBackups/db/backups"
	"github.com/dedeme/FastBackups/db/restore"
	"github.com/dedeme/golib/cryp"
	"github.com/dedeme/golib/file"
	"github.com/dedeme/golib/sys"
	"path"
	"strings"
)

func dir() string {
	return path.Join(file.HomeDir(), "FastBackups")
}

func getControlKey() (key, controlFile string, ok bool) {
	key, controlFile, ok = backups.KeyControl()
	if !ok {
		key, controlFile, ok = restore.KeyControl()
	}
	return
}

func mkBackups(key string) {
	backups.Run(key, func(info string) {
		sys.Cmd("zenity",
			"--no-wrap",
			"--title", "Fast Backups",
			"--warning", "--text="+info,
		)
	})
	sys.Cmd("zenity",
		"--no-wrap",
		"--title", "Fast Backups",
		"--info", "--text=Backups finished.",
	)
}

func mkRestore(key string) {
	restore.Run(key, func(info string) {
		sys.Cmd("zenity",
			"--no-wrap",
			"--title", "Fast Backups",
			"--warning", "--text="+info,
		)
	})
	sys.Cmd("zenity",
		"--no-wrap",
		"--title", "Fast Backups",
		"--info", "--text=Restoring finished.",
	)
}

// Program entry.
func main() {
	db.Initialize()

	cout, cerr := sys.Cmd("zenity",
		"--title", "Fast Backups",
		"--password",
	)
	if len(cerr) > 0 {
		return
	}
	if len(cout) == 0 || cout[0] == 10 {
		sys.Cmd("zenity",
			"--no-wrap",
			"--title", "Fast Backups",
			"--error", "--text=No password entered.",
		)
		return
	}

	key := cryp.Key(strings.TrimSpace(string(cout)), 450)
	key2, controlFile, ok2 := getControlKey()
	if ok2 {
		if key2 != cryp.Key(key, 450) {
			sys.Cmd("zenity",
				"--no-wrap",
				"--title", "Fast Backups",
				"--error", "--text=Key of\n'"+controlFile+
					"'\ndoes not match password",
			)
			return
		}
	} else {
		cout, cerr := sys.Cmd("zenity",
			"--title", "Fast Backups (Retype Pass)",
			"--password",
		)
		if len(cerr) > 0 {
			return
		}
		if len(cout) == 0 || cout[0] == 10 {
			sys.Cmd("zenity",
				"--no-wrap",
				"--title", "Fast Backups",
				"--error", "--text=First and second passwords do not mach.",
			)
			return
		}

		if key != cryp.Key(strings.TrimSpace(string(cout)), 450) {
			sys.Cmd("zenity",
				"--no-wrap",
				"--title", "Fast Backups",
				"--error", "--text=First and second passwords do not mach.",
			)
			return
		}
	}

	cout, cerr = sys.Cmd("zenity",
		"--title", "Fast Backups",
		"--list", "--radiolist",
		"--column=", "--column=",
		"FALSE", "Backups",
		"FALSE", "Restore",
	)
	if len(cerr) > 0 {
		return
	}
	option := strings.TrimSpace(string(cout))
	if len(option) == 0 {
		return
	}

	if option == "Backups" {
		mkBackups(key)
	} else if option == "Restore" {
		mkRestore(key)
	} else {
		sys.Cmd("zenity",
			"--no-wrap",
			"--title", "Fast Backups",
			"--error", "--text=Unexpected selection '"+option+"'",
		)
		return
	}
}

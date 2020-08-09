// Copyright 07-Jul-2020 ºDeme
// GNU General Public License - V3 <http://www.gnu.org/licenses/>

package main

import (
	"github.com/dedeme/golib/sys"
	"github.com/dedeme/kron/db"
	"os"
)

func main() {
	sys.Initialize("kron")
	db.LogInit()
	db.ActInit()
	db.MemInit()

	if len(os.Args) == 2 {
		if os.Args[1] == "help" {
			help() // help.go
			return
		}
		if os.Args[1] == "list" {
			list() // actions.go
			return
		}
		if os.Args[1] == "start" {
			start() // actions.go
			return
		}
		if os.Args[1] == "stop" {
			stop() // actions.go
			return
		}
		if os.Args[1] == "test" {
			test("") // actions.go
			return
		}
		if os.Args[1] == "log" {
			printLog(false) // actions.go
			return
		}
	} else if len(os.Args) > 2 {
		if os.Args[1] == "help" && len(os.Args) == 3 {
			if os.Args[2] == "add" {
				addHelp() // help.go
				return
			}
			if os.Args[2] == "mem" {
				memHelp() // help.go
				return
			}
			if os.Args[2] == "list" {
				listHelp() // help.go
				return
			}
			if os.Args[2] == "del" {
				delHelp() // help.go
				return
			}
			if os.Args[2] == "log" {
				logHelp() // help.go
				return
			}
		}

		if os.Args[1] == "log" {
			if os.Args[2] == "all" && len(os.Args) == 3 {
				printLog(true)
				return
			}
		}

		if os.Args[1] == "mem" {
			if len(os.Args) == 3 {
				mem(os.Args[2])
				return
			}
		}

		if os.Args[1] == "add" {
			if os.Args[2] == "d" {
				if len(os.Args) == 4 {
					addD(os.Args[3], "")
					return
				} else if len(os.Args) == 5 {
					addD(os.Args[3], os.Args[4])
					return
				}
			} else if os.Args[2] == "m" {
				if len(os.Args) == 5 {
					addM(os.Args[3], os.Args[4], "")
					return
				} else if len(os.Args) == 6 {
					addM(os.Args[3], os.Args[4], os.Args[5])
					return
				}
			}
		}

		if os.Args[1] == "del" {
			if len(os.Args) == 3 {
				del(os.Args[2])
				return
			}
		}

		if os.Args[1] == "test" {
			if len(os.Args) == 3 {
				test(os.Args[2])
				return
			}
		}
	}

	help() // help.go
}

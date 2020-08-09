// Copyright 07-Jul-2020 ºDeme
// GNU General Public License - V3 <http://www.gnu.org/licenses/>

package main

import (
	"fmt"
)

func help() {
	fmt.Println(`
Use: kron action options

action:
  help
    options: [action]
    Shows program information.
    examples:
      kron help
      kron help add
  mem
    options: command
    Saves command in memory.
    examples:
      kron mem prg
      kron mem "prg -o ammount.txt"
  add
    options: [d moment [exclude] | m day moment [exclude]]
    Adds the current memorized activity to the task list.
    examples:
      kron add d 17:5 SD
      kron add m 5 7 7,8,12
  list
    Shows in memory and saved activities.
    examples:
      kron list
  del
    Removes a saved activity.
    examples:
      kron del 234
  test
    Tries command in memory or saved.
    examples:
      kron test
      kron test 32
  start
    Starts server.
    examples:
      kron start
  stop
    Stops server.
    examples:
      kron stop
  log
    Shows program log (errors or all the entries).
    examples:
      kron log
      kron log a
  `)
}

func addHelp() {
	fmt.Println(`
Use: kron add d moment [exclude]
 or: kron add m day moment [exclude]

Adds the current memorized activity to the task list.

Case 1: kron add d moment [exclude]
  moment: hour | hour:minute
    Moment to do the activity.
  exclude: String with L, M, X, J, V, S, D
    Week days excluded for activity.
  examples:
    kron add d 17
    kron add d 7:30
    kron add d 7:5 SD

Case 2 kron add m day moment [exclude]
  day: Day number from 1 to 31
    Month day.
  moment: hour | hour:minute
    Moment to do the activity.
  exclude: Month numbers separate by comma and without spaces
    Month excluded for activity
  examples:
    kron add m 5 17:5
    kron add m 12 7:30
    kron add m 5 7 7,8,12
  `)
}

func memHelp() {
	fmt.Println(`
Use: kron mem command

Saves command in memory.

  command: System command (surrounded by quotes if necessary)
  examples:
    kron mem prg
    kron mem "prg -o ammount.txt"
`)
}

func listHelp() {
	fmt.Println(`
Use: kron list

Shows in memory and saved activities.

  examples:
    kron list
`)
}

func delHelp() {
	fmt.Println(`
Use: kron del identifier

Removes a saved activity.

  identifier: Activity identifier conforming is showed by 'kron list'
  examples:
    kron del 234
`)
}

func testHelp() {
	fmt.Println(`
Use: kron test [id]

Tries command in memory or saved.

  id: Activity identifier. If it is missing, it tries memory activity.
  examples:
    kron test
    kron test 32
`)
}

func startHelp() {
	fmt.Println(`
Use: kron start

Starts server.

  examples:
    kron start
`)
}

func stopHelp() {
	fmt.Println(`
Use: kron start

Stops server.

  examples:
    kron stop
`)
}

func logHelp() {
	fmt.Println(`
Use: kron log [all]

Shows program log (errors or all the entries).

  all: Shows every entry of log.
  examples:
    kron log
    kron log a
`)
}

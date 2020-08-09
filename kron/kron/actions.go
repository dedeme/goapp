// Copyright 07-Jul-2020 ºDeme
// GNU General Public License - V3 <http://www.gnu.org/licenses/>

package main

import (
	"fmt"
	"github.com/dedeme/golib/date"
	"github.com/dedeme/golib/sys"
	"github.com/dedeme/kron/data"
	"github.com/dedeme/kron/db"
	"sort"
	"strconv"
	"strings"
)

func printLog(all bool) {
	printed := false
	for _, e := range db.LogRead() { // log.go
		if all || e.IsError() {
			fmt.Println(e)
			printed = true
		}
	}
	if !printed {
		fmt.Println("Log is empty.")
	}
}

func list() {
	fmt.Printf("Memory:\n  %v\n\n", db.MemRead())
	fmt.Println("Activities:")
	acts := db.ActRead()
	sort.Slice(acts, func(i, j int) bool {
		if acts[i].IsMonthly == acts[j].IsMonthly {
			return acts[i].Command < acts[j].Command
		} else {
			return !acts[i].IsMonthly
		}
	})
	for _, e := range acts {
		fmt.Println(e)
	}
}

func mem(command string) {
	db.MemWrite(command)
	fmt.Printf("Memorized %v\n", command)
}

func addD(moment, excludes string) {
	comm := db.MemRead()
	if comm == "" {
		fmt.Println("There is no command in memory.")
		return
	}
	entry, err := data.NewDaily(moment, excludes, comm)
	if err != nil {
		fmt.Println(err)
		return
	}
	db.ActAdd(entry)
	db.MemReset()
	fmt.Println("Added activity:\n" + entry.String())
}

func addM(day, moment, excludes string) {
	comm := db.MemRead()
	if comm == "" {
		fmt.Println("There is no command in memory.")
		return
	}
	entry, err := data.NewMonthly(day, moment, excludes, comm)
	if err != nil {
		fmt.Println(err)
		return
	}
	db.ActAdd(entry)
	db.MemReset()
	fmt.Println("Added activity:\n" + entry.String())
}

func del(id string) {
	n, err := strconv.Atoi(id)
	if err != nil {
		fmt.Println("Activity identifier must be a number.")
		return
	}
	if db.ActDel(n) {
		fmt.Println("Deleted activity", id)
		return
	}
	fmt.Printf("Activity %v not found.", id)
}

func start() {
	timeStamp := db.StartOn()
	if timeStamp == "" {
		fmt.Println("Kron is already started.")
		return
	}
  msg := "Server started: " + timeStamp
	fmt.Println(msg)
  db.LogInfo(msg)
	for {
		ts := db.StartTimeStamp()
		if ts != timeStamp {
			break
		}

		dt := date.Now()
    dtS := dt.String()
		entries := db.ActRead()
    toWrite := false
		for _, e := range entries {
			edt := date.NewTime(dt.Day(), dt.Month(), dt.Year(), e.Hour, e.Minute, 0)
			toDo := false
			if e.IsMonthly {
				mOk := true
				month := dt.Month()
				for _, m := range e.Mexceptions {
					if m == month {
						mOk = false
						break
					}
				}
				toDo = mOk && e.Day == dt.Day() && dt.DfTime(edt) > 0 && dtS != e.Done
			} else {
        wDay := "DLMXJVS"[dt.Weekday()]
        dOk := strings.IndexByte(e.Dexceptions, wDay) == -1
        toDo = dOk && dt.DfTime(edt) > 0 && dtS != e.Done
			}
			if toDo {
        comm := e.Command
				go func() {
					ps := strings.Split(comm, " ")
					o, e := sys.Cmd(ps[0], ps[1:]...)
					msg := "Running " + comm
					if len(e) > 0 {
						msg += fmt.Sprintf("\nError:\n%v\n", string(e))
					} else {
						if len(o) > 0 {
							msg += fmt.Sprintf("\n%v", string(o))
						}
					}
          db.LogInfo(msg)
				}()
        e.Done = dtS
        toWrite = true;
			}
		}

    if toWrite {
      db.ActWrite(entries)
    }

		sys.Sleep(5000)
	}
  msg = "Server stoped: " + timeStamp
	fmt.Println(msg)
  db.LogInfo(msg)
}

func stop() {
	db.StartOf()
	fmt.Println("Sever will be stopped at most in 1 minute...")
}

func test(id string) {
	var comm string
	if id == "" {
		comm = db.MemRead()
		if comm == "" {
			fmt.Println("Memory command is not defined.")
			return
		}
	} else {
		n, err := strconv.Atoi(id)
		if err != nil {
			fmt.Println("Activity identifier must be a number.")
			return
		}
		for _, e := range db.ActRead() {
			if e.Id == n {
				comm = e.Command
				break
			}
		}
		if comm == "" {
			fmt.Printf("Activity %v not found\n", id)
		}
	}

	fmt.Println("Testing " + comm)
	ps := strings.Split(comm, " ")
	o, e := sys.Cmd(ps[0], ps[1:]...)

	if len(e) > 0 {
		fmt.Printf("Error:\n%v\n", string(e))
	} else {
		fmt.Print("Ok")
		if len(o) > 0 {
			fmt.Printf(":\n%v", string(o))
		}
		fmt.Println("")
	}
}

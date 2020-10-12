// Copyright 24-Apr-2020 ºDeme
// GNU General Public License - V3 <http://www.gnu.org/licenses/>

// Start program.
package main

import (
	"fmt"
	"github.com/dedeme/dmstack/args"
	"github.com/dedeme/dmstack/imports"
	"github.com/dedeme/dmstack/machine"
	"github.com/dedeme/dmstack/reader"
	"github.com/dedeme/dmstack/runner"
	"github.com/dedeme/dmstack/symbol"
	"io/ioutil"
	"os"
	"path"
	"runtime/debug"
	"strings"
)

// Run main process and recover machine panics
func run(m *machine.T) {
	defer func() {
		if err := recover(); err != nil {
			e, ok := err.(*machine.Error)
			if ok {
				fmt.Printf("%v: %v\n", e.Type, e.Message)
				fmt.Println(strings.Join(e.Machine.StackTrace(), "\n"))
				os.Exit(0)
			} else {
				panic(err)
			}
		}
	}()

	runner.Run(m)
}

// Process main .dms file.
func process(file, module string) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
			debug.PrintStack()
			os.Exit(1)
		}
	}()

	bs, err := ioutil.ReadFile(file)
	if err != nil {
		panic(err)
	}

	rd := reader.New(module, string(bs))
	prg, ok := rd.Process().P()
	if !ok {
		panic("Reader process does not return a Program")
	}

	m := machine.NewIsolate(path.Dir(file), []*machine.T{}, prg)
	run(m)
}

// Program entry.
func main() {
	symbol.Initialize()
	if ok := args.Initialize(); !ok {
		return
	}

	file := args.Stkargs["dms"]
	if !strings.HasSuffix(file, ".dms") {
		file = file + ".dms"
	}
	file = path.Clean(file)
	if _, err := os.Stat(file); err != nil {
		fmt.Printf("File '%v' not found.", file)
		return
	}

	module := file[0 : len(file)-4]
	source := symbol.New(module)
	imports.PutOnWay(source)

	process(file, module)
}

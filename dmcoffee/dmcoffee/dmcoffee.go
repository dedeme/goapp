// Copyright 14-Jan-2021 ºDeme
// GNU General Public License - V3 <http://www.gnu.org/licenses/>

// Start program.
package main

import (
	"fmt"
	"github.com/dedeme/dmcoffee/args"
	"github.com/dedeme/dmcoffee/cts"
	"github.com/dedeme/dmcoffee/fail"
	"github.com/dedeme/dmcoffee/imports"
	"github.com/dedeme/dmcoffee/machine"
	"github.com/dedeme/dmcoffee/operator"
	"github.com/dedeme/dmcoffee/reader"
	"github.com/dedeme/dmcoffee/runner"
	"github.com/dedeme/dmcoffee/symbol"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

// Process main .dmc file.
func process(file string, module symbol.T) {
	defer func() {
		if err := recover(); err != nil {
			e, ok := err.(*fail.T)
			if ok {
				fmt.Printf(
					"%v.%v:%v: %v:\n%v\n",
					e.Pos.Source, cts.SourceExtension, e.Pos.Nline,
					e.Type, e.Message,
				)
				os.Exit(0)
			} else {
				panic(err)
			}
		}
	}()

	bs, err := ioutil.ReadFile(file)
	if err != nil {
		panic(err)
	}

	rd := reader.New(module, string(bs))
	prg := rd.Process()
	runner.Run(machine.NewInit(prg))
}

// Program entry.
func main() {
	operator.Initialize()
	symbol.Initialize()

	if ok := args.Initialize(); !ok {
		return
	}

	file := args.CoffeeArgs[cts.SourceExtension]
	if !strings.HasSuffix(file, "."+cts.SourceExtension) {
		file = file + "." + cts.SourceExtension
	}
	file = path.Clean(file)
	if _, err := os.Stat(file); err != nil {
		fmt.Printf("File '%v' not found.", file)
		return
	}

	module := symbol.New(file[0 : len(file)-4])
	imports.Add(module)
	process(file, module)
}

// Copyright 30-Mar-2022 ºDeme
// GNU General Public License - V3 <http://www.gnu.org/licenses/>

// Module checker.
package checker

import (
	"errors"
	"fmt"
	"github.com/dedeme/kut/checker/cksym"
	"github.com/dedeme/kut/fileix"
	"github.com/dedeme/kut/modules"
	"github.com/dedeme/kut/reader/token"
	"github.com/dedeme/kut/reader/txReader"
)

func checkCode(isTop bool, isFor bool, layers [][]*cksym.T, tx *txReader.T) (
	errs []error,
) {
	if !isFor {
		lastLayer := []*cksym.T{}
		if isTop {
			md := modules.GetOk(tx.File)
			for k, v := range md.Heap0 {
				lastLayer = append(lastLayer, cksym.New(k, tx.File, v.Nline))
			}
		}
		layers = append(layers, lastLayer)
	}
	var nextTk *token.T
	for {
		var end bool
		var ers []error
		end, nextTk, ers = checkStatement(isTop, false, nextTk, layers, tx) // stChecker.go
		errs = append(errs, ers...)
		if end {
			break
		}
	}

	for _, s := range layers[len(layers)-1] {
		if !s.Used {
			errs = append(errs, errors.New(fmt.Sprintf(
				"%v:%v: Symbol not used (%v)",
				fileix.Get(s.File), s.Nline, s.Name,
			)))
		}
	}

	return
}

func Run() {
	ne := 0
	for i, m := range modules.List() {
		if m == nil {
			continue
		}

		kutCode, err := fileix.Read(i)
		if err != nil {
			fmt.Println(err)
			return
		}
		errs := checkCode(true, false, [][]*cksym.T{}, txReader.New(i, kutCode))
		for _, e := range errs {
			if ne >= 15 {
				fmt.Println("... and more.")
				break
			}
			fmt.Println(e)
		}

		if ne >= 15 {
			break
		}
	}
}

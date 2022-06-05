// Copyright 30-Mar-2022 ºDeme
// GNU General Public License - V3 <http://www.gnu.org/licenses/>

// Checker module.
package cksym

import (
	"errors"
	"fmt"
	"github.com/dedeme/kut/fileix"
)

type T struct {
	Name  string
	File  int
	Nline int
	Used  bool
}

// Creates a new check symbol
//    name: Symbol name.
//    fileIx: Index of file where 'sym' apears.
//    nline : File number in 'fileIx'.
// NOTE: Field User is set to 'false'
func New(name string, fileIx, nline int) *T {
	return &T{name, fileIx, nline, false}
}

// Find symbol 'sym' in layers and set its field Used to 'true' returning
// an error if such symbol was NOT found or 'nil' otherwise.
func ErrIfNotFound(layers [][]*T, sym *T) error {
	for i := len(layers) - 1; i >= 0; i-- {
		a := layers[i]
		for _, s := range a {
			if s.Name == sym.Name {
				s.Used = true
				return nil
			}
		}
	}
	return errors.New(fmt.Sprintf(
		"%v:%v: Symbol not declared (%v)",
		fileix.GetFail(sym.File), sym.Nline, sym.Name,
	))
}

// Find symbol 'sym' in the last layer of 'layers' returning an error
// if such symbol was found or an error otherwise.
func ErrIfFound(layers [][]*T, sym *T) error {
	for _, s := range layers[len(layers)-1] {
		if s.Name == sym.Name {
			return errors.New(fmt.Sprintf(
				"%v:%v: Symbol '%v' already declared in line %v",
				fileix.GetFail(sym.File), sym.Nline, sym.Name, s.Nline,
			))
		}
	}
	return nil
}

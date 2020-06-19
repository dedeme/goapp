// Copyright 26-Apr-2020 ºDeme
// GNU General Public License - V3 <http://www.gnu.org/licenses/>

package tests

import (
	"github.com/dedeme/dmstack/imports"
	"github.com/dedeme/dmstack/symbol"
	"testing"
)

func TestOnWay(t *testing.T) {
	s := symbol.New("abc")
	imports.PutOnWay(s)
	if !(imports.IsOnWay(s)) {
		t.Fatal(fail)
	}
	imports.QuitOnWay(s)
	if imports.IsOnWay(s) {
		t.Fatal(failNot)
	}
}

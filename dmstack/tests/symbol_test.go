// Copyright 26-Apr-2020 ºDeme
// GNU General Public License - V3 <http://www.gnu.org/licenses/>

package tests

import (
	"github.com/dedeme/dmstack/symbol"
	"testing"
)

func TestConsts(t *testing.T) {
  symbol.Initialize()

	if r := eqs(
		symbol.If.String(),
		"if",
	); r != "" {
		t.Fatal(r)
	}
	if r := eqs(
		symbol.Nop.String(),
		"nop",
	); r != "" {
		t.Fatal(r)
	}

}

func TestNew(t *testing.T) {
	if r := eqs(
		symbol.New("abc").String(),
		"abc",
	); r != "" {
		t.Fatal(r)
	}
}

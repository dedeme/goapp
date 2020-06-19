// Copyright 24-Apr-2020 ºDeme
// GNU General Public License - V3 <http://www.gnu.org/licenses/>

package tests

import (
	"testing"
)

func TestFirst(t *testing.T) {

	if r := eqi(
		1+1,
		2,
	); r != "" {
		t.Fatal(r)
	}

}

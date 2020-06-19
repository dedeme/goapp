// Copyright 24-Apr-2020 ºDeme
// GNU General Public License - V3 <http://www.gnu.org/licenses/>

package tests

import (
	"fmt"
)

const fail = "\nActual  : false\nExpected: true\n"

const failNot = "\nActual  : true\nExpected: false\n"

func eqs(actual, expected string) string {
	if actual != expected {
		return fmt.Sprintf("\nActual  : %v\nExpected: %v\n", actual, expected)
	}
	return ""
}

func eqi(actual, expected int) string {
	if actual != expected {
		return fmt.Sprintf("\nActual  : %v\nExpected: %v\n", actual, expected)
	}
	return ""
}

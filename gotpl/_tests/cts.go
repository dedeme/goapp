// Copyright 31-Aug-2017 ºDeme
// GNU General Public License - V3 <http://www.gnu.org/licenses/>

package _tests

import (
	"fmt"
)

const fail = "\nActual  : false\nExpected: true\n"

func eq(actual, expected string) string {
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

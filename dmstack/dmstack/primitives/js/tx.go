// Copyright 28-Sep-2020 ºDeme
// GNU General Public License - V3 <http://www.gnu.org/licenses/>

// Text management.
package js

func nextByte(s string, ch byte, ix int) (pos int, ok bool) {
	pos = ix
	l := len(s)
	quotes := false
	bar := false
	squarn := 0
	bracketn := 0
	ok = true
	var c byte
	for {
		if pos == l {
			ok = false
			break
		}
		c = s[pos]
		if quotes {
			if bar {
				bar = false
			} else {
				if c == '\\' {
					bar = true
				} else if c == '"' {
					quotes = false
				}
			}
		} else {
			if c == ch &&
				((c == ']' && squarn == 1 && bracketn == 0) ||
					(c == '}' && squarn == 0 && bracketn == 1) ||
					(squarn == 0 && bracketn == 0)) {
				break
			} else if c == '"' {
				quotes = true
			} else if c == '[' {
				squarn++
			} else if c == ']' {
				squarn--
			} else if c == '{' {
				bracketn++
			} else if c == '}' {
				bracketn--
			}
		}
		pos++
	}
	return
}

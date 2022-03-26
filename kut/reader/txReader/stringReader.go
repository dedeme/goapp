// Copyright 01-Feb-2022 ºDeme
// GNU General Public License - V3 <http://www.gnu.org/licenses/>

package txReader

import (
	"fmt"
	"github.com/dedeme/kut/reader/token"
	"strings"
)

func (r *T) readUnicode() (tx string, err error) {
	isHex := func(b byte) bool {
		return (b >= '0' && b <= '9') ||
			(b >= 'a' && b <= 'f') || (b >= 'A' && b <= 'F')
	}
	var b byte
	var eof bool
	for i := 0; i < 4; i++ {
		b, eof, err = r.readByte()
		if err != nil {
			return
		}
		if eof {
			err = r.Fail(fmt.Sprintf("Bad unicode sequence (%v).", tx))
			return
		}
		if !isHex(b) {
			err = r.Fail(fmt.Sprintf("Bad unicode sequence (%v%v).", tx, string(b)))
			return
		}
		tx += string(b)
	}

	return
}

func (r *T) readEscape(close byte) (tx string, err error) {
	var b byte
	var eof bool
	b, eof, err = r.readByte()
	if err != nil {
		return
	}
	if eof {
		err = r.Fail("Unclosed quotes.")
		return
	}
	switch b {
	case close:
		tx = string(b)
	case '\\':
		tx = "\\"
	case 'n':
		tx = "\n"
	case 't':
		tx = "\t"
	case 'b':
		tx = "\b"
	case 'f':
		tx = "\f"
	case 'r':
		tx = "\r"
	case 'u':
		tx, err = r.readUnicode()
	default:
		err = r.Fail(fmt.Sprintf("Bad escape squence '\\%v'.", string(b)))
	}
	return
}

func (r *T) readSimpleString(close byte) (tk *token.T, err error) {
	rnClose := rune(close)
	var rn rune
	var eof bool
	var bf strings.Builder
	for {
		rn, eof, err = r.readRune()
		if err != nil {
			break
		}
		if eof || rn == '\n' {
			err = r.Fail("Unclosed quotes.")
			break
		}
		if rn == '\\' {
			var tx string
			tx, err = r.readEscape(close)
			if err != nil {
				break
			}
			bf.WriteString(tx)
			continue
		}
		if rn == rnClose {
			tk = token.New(token.String, bf.String())
			break
		}
		bf.WriteRune(rn)
	}

	return
}

func (r *T) readMultilineString() (tk *token.T, err error) {
	format := func(s string) string {
		ss := strings.Split(s, "\n")
		cut := -1
		for _, e := range ss {
			if strings.TrimSpace(e) == "" {
				continue
			}
			n := 0
			for i := 0; i < len(e); i++ {
				if e[i] == ' ' {
					n++
				} else {
					break
				}
			}
			if cut < 0 || n < cut {
				cut = n
			}
		}
		if cut > 0 {
			var newSs []string
			for _, e := range ss {
				if strings.TrimSpace(e) == "" {
					newSs = append(newSs, "")
				} else {
					newSs = append(newSs, e[cut:])
				}
			}
			ss = newSs
		}
		return strings.Join(ss, "\n")
	}

	nline := r.Nline
	var b byte
	var eof bool
	b, eof, err = r.readByte()
	if err != nil {
		return
	}
	if eof {
		err = r.FailLine("Unclosed quotes.", nline)
	}
	if b != '\n' {
		err = r.FailLine("Expected end of line after \"\"\".", nline)
		return
	}

	var bf strings.Builder
	var rn rune
	nqs := 0
	for {
		rn, eof, err = r.readRune()
		if err != nil {
			break
		}
		if eof {
			err = r.FailLine("Unclosed quotes.", nline)
			break
		}
		if rn == '"' {
			nqs++
			if nqs == 3 {
				tk = token.New(token.String, format(bf.String()))
				break
			}
			continue
		}
		if nqs > 0 {
			bf.WriteByte('"')
			if nqs > 1 {
				bf.WriteByte('"')
			}
			nqs = 0
		}
		bf.WriteRune(rn)
	}
	return
}

func (r *T) readString(b byte) (tk *token.T, err error) {
	if b == '\'' {
		tk, err = r.readSimpleString('\'')
		return
	}

	var eof bool
	var rn rune
	rn, eof, err = r.readRune()
	if err == nil {
		if eof {
			err = r.Fail("Unclosed quotes.")
		} else if rn == '"' {
			rn, eof, err = r.readRune()
			if err == nil {
				if eof {
					tk = token.New(token.String, "")
				} else {
					if rn == '"' {
						tk, err = r.readMultilineString()
					} else {
						err = r.unreadRune(rn)
						if err == nil {
							tk = token.New(token.String, "")
						}
					}
				}
			}
		} else {
			err = r.unreadRune(rn)
			if err == nil {
				tk, err = r.readSimpleString('"')
			}
		}
	}

	return
}

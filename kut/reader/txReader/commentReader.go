// Copyright 01-Feb-2022 ºDeme
// GNU General Public License - V3 <http://www.gnu.org/licenses/>

package txReader

import (
	"github.com/dedeme/kut/reader/token"
	"strings"
)

func (r *T) readLineComment() (tk *token.T, err error) {
	var rn rune
	var eof bool
	var bf strings.Builder
	for {
		rn, eof, err = r.readRune() // in txReader.go
		if err != nil {
			break
		}
		if eof || rn == '\n' {
			tk = token.New(token.LineComment, bf.String())
			break
		}
		bf.WriteRune(rn)
	}
	return
}

func (r *T) readLongComment() (tk *token.T, err error) {
	nline := r.Nline
	var rn rune
	var eof bool
	var bf strings.Builder
	for {
		rn, eof, err = r.readRune() // in txReader.go
		if err != nil {
			break
		}
		if eof {
			err = r.FailLine("Unclosed comment.", nline)
			break
		}
		if rn == '*' {
			var b byte
			var eof bool
			b, eof, err = r.readByte() // in txReader.go
			if err != nil {
				break
			}
			if eof {
				err = r.FailLine("Unclosed comment.", nline)
				break
			}
			if b == '/' {
				tk = token.New(token.Comment, bf.String())
				break
			}
			bf.WriteByte('*')
		}
		bf.WriteRune(rn)
	}
	return
}

func (r *T) readComment() (ok bool, err error) {
	var b byte
	var eof bool
	b, eof, err = r.readByte() // in txReader.go
	if err != nil {
		return
	}
	if eof {
		err = r.Fail("Unexpected end of file.")
		return
	}

	switch b {
	case '/':
		_, err = r.readLineComment()
		ok = true
	case '*':
		_, err = r.readLongComment()
		ok = true
	default:
		err = r.unreadByte(b)
	}

	return
}

// Returns 'tk' with commentary text or 'nil'.
func (r *T) readCommentText() (tk *token.T, err error) {
	var b byte
	var eof bool
	b, eof, err = r.readByte() // in txReader.go
	if err != nil {
		return
	}
	if eof {
		err = r.Fail("Unexpected end of file.")
		return
	}

	switch b {
	case '/':
		tk, err = r.readLineComment()
	case '*':
		tk, err = r.readLongComment()
	default:
		err = r.unreadByte(b)
	}

	return
}

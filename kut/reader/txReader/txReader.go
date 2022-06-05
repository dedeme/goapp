// Copyright 01-Feb-2022 ºDeme
// GNU General Public License - V3 <http://www.gnu.org/licenses/>

// Text to tokens reader.
package txReader

import (
	"errors"
	"fmt"
	"github.com/dedeme/kut/fileix"
	"github.com/dedeme/kut/reader/token"
	"strings"
)

const (
	opChs = ";=,.()+-*/!<>[]{}%&|?:\\"
)

type T struct {
	tx    *strings.Reader
	File  int
	nline int
	Nline int
}

func New(file int, tx string) *T {
	return &T{strings.NewReader(tx), file, 1, 1}
}

func (r *T) Fail(msg string) error {
	return errors.New(fmt.Sprintf(
		"%v\n  %v:%v:", msg, fileix.GetFail(r.File), r.Nline))
}

func (r *T) FailLine(msg string, nline int) error {
	return errors.New(fmt.Sprintf(
		"%v\n  %v:%v:", msg, fileix.GetFail(r.File), nline))
}

func (r *T) FailExpect(expected string, found string, nline int) error {
	return errors.New(fmt.Sprintf(
		"Expected: %v\nFound: %v\n  %v:%v:",
		expected, found, fileix.GetFail(r.File), nline))
}

func (r *T) readByte() (b byte, eof bool, err error) {
	if r.Nline != r.nline {
		r.Nline = r.nline
	}

	var rn rune
	var size int
	rn, size, err = r.tx.ReadRune()
	if err != nil && err.Error() != "EOF" {
		err = r.Fail(err.Error())
	} else {
		if size == 0 || err != nil {
			err = nil
			eof = true
		} else if size == 1 {
			b = byte(rn)
			if b == '\n' {
				r.nline++
			}
		} else {
			err = r.Fail(fmt.Sprintf("Unexpected character '%v'", string(rn)))
		}
	}

	return
}

func (r *T) readRune() (rn rune, eof bool, err error) {
	if r.Nline != r.nline {
		r.Nline = r.nline
	}

	var size int
	rn, size, err = r.tx.ReadRune()
	if err != nil && err.Error() != "EOF" {
		err = r.Fail(err.Error())
	} else {
		if size == 0 || err != nil {
			eof = true
		} else {
			if rn == '\n' {
				r.nline++
			}
		}
	}
	return
}

func (r *T) unreadByte(b byte) (err error) {
	if err = r.tx.UnreadByte(); err != nil {
		err = r.Fail(err.Error())
	}
	if b == '\n' {
		r.nline--
	}
	return
}

func (r *T) unreadRune(rn rune) (err error) {
	if err = r.tx.UnreadRune(); err != nil {
		err = r.Fail(err.Error())
	}
	if rn == '\n' {
		r.nline--
	}
	return
}

func (r *T) readSymbol(b byte) (tk *token.T, err error) {
	var eof bool
	var bf strings.Builder
	s := ""
	bf.WriteByte(b)
	for {
		b, eof, err = r.readByte()
		if err != nil {
			break
		}
		if eof {
			s = bf.String()
			break
		}
		if (b >= 'a' && b <= 'z') || (b >= 'A' && b <= 'Z') ||
			(b >= '0' && b <= '9') {
			bf.WriteByte(b)
		} else {
			r.unreadByte(b)
			s = bf.String()
			break
		}
	}

	switch s {
	case "":
	case "true":
		tk = token.New(token.Bool, true)
	case "false":
		tk = token.New(token.Bool, false)
	default:
		tk = token.New(token.Symbol, bf.String())
	}
	return
}

func (r *T) readOperator(b byte) (tk *token.T, err error) {
	switch b {
	case '=', '!', '>', '<', '+', '*', '/', '%', '|', '&':
		var eof bool
		var b2 byte
		b2, eof, err = r.readByte()
		if err == nil {
			if !eof && b2 == '=' {
				tk = token.New(token.Operator, string([]byte{b, b2}))
			} else {
				r.unreadByte(b2)
				tk = token.New(token.Operator, string(b))
			}
		}
	case '-':
		var eof bool
		var b2 byte
		b2, eof, err = r.readByte()
		if err == nil {
			if !eof && (b2 == '=' || b2 == '>') {
				tk = token.New(token.Operator, string([]byte{b, b2}))
			} else {
				r.unreadByte(b2)
				tk = token.New(token.Operator, string(b))
			}
		}
	default:
		tk = token.New(token.Operator, string(b))
	}
	return
}

func (r *T) ReadToken() (tk *token.T, eof bool, err error) {
	var b byte
	b, eof, err = r.readByte()

	if err != nil || eof {
		return
	}

	switch {
	case b <= ' ':
		tk, eof, err = r.ReadToken()
	case b == '/':
		var ok bool
		ok, err = r.readComment() // in commentReader.go
		if ok {
			tk, eof, err = r.ReadToken()
		} else {
			tk, err = r.readOperator(b)
		}
	case b == '"', b == '\'':
		tk, err = r.readString(b) // stringReader.go
	case (b >= 'a' && b <= 'z') || (b >= 'A' && b <= 'Z'):
		tk, err = r.readSymbol(b)
	case strings.IndexByte(opChs, b) != -1:
		tk, err = r.readOperator(b)
	case b >= '0' && b <= '9':
		tk, err = r.readNumber(b) // in numberReader.go
	default:
		err = r.Fail(fmt.Sprintf("Unexpected character '%v'", string(b)))
	}

	return
}

// Copyright 01-Feb-2022 ºDeme
// GNU General Public License - V3 <http://www.gnu.org/licenses/>

package txReader

import (
	"github.com/dedeme/kut/reader/token"
	"strconv"
	"strings"
)

func (r *T) readExponent(bf *strings.Builder) (tk *token.T, err error) {
	var b byte
	var eof bool
	b, eof, err = r.readByte()
	if err != nil {
		return
	}
	if eof {
		err = r.Fail("Exponent has no digits.")
		return
	}
	if b == '+' || b == '-' {
		bf.WriteByte(b)
		b, eof, err = r.readByte()
		if err != nil {
			return
		}
		if eof {
			err = r.Fail("Exponent has no digits.")
			return
		}
	}

	var bf2 strings.Builder
	for {
		if b >= '0' && b <= '9' {
			bf2.WriteByte(b)
			b, eof, err = r.readByte()
			if err != nil {
				break
			}
			if eof {
				break
			}
			continue
		}
		r.unreadByte(b)
		break
	}
	if err != nil {
		return
	}
	if bf.Len() == 0 {
		err = r.Fail("Exponent has no digits.")
		return
	}
	n, e := strconv.ParseFloat(bf.String()+bf2.String(), 64)
	if e != nil {
		err = r.Fail(e.Error())
		return
	}
	tk = token.New(token.Float, n)
	return
}

func (r *T) readFloat(bf *strings.Builder) (tk *token.T, err error) {
	var b byte
	var eof bool
	for {
		b, eof, err = r.readByte()
		if err != nil {
			break
		}
		if eof {
			n, _ := strconv.ParseFloat(bf.String(), 64)
			tk = token.New(token.Float, n)
			break
		}
		if b >= '0' && b <= '9' {
			bf.WriteByte(b)
		} else if b == 'e' || b == 'E' {
			bf.WriteByte(b)
			tk, err = r.readExponent(bf)
			break
		} else {
			r.unreadByte(b)
			n, e := strconv.ParseFloat(bf.String(), 64)
			if e != nil {
				err = r.Fail(err.Error())
				break
			}
			tk = token.New(token.Float, n)
			break
		}
	}
	return
}

func (r *T) readNumber(b byte) (tk *token.T, err error) {
	var eof bool
	var bf strings.Builder
	bf.WriteByte(b)
	for {
		b, eof, err = r.readByte()
		if err != nil {
			break
		}
		if eof {
			n, _ := strconv.ParseInt(bf.String(), 10, 64)
			tk = token.New(token.Int, n)
			break
		}
		if b >= '0' && b <= '9' {
			bf.WriteByte(b)
		} else if b == '_' {
			continue
		} else if b == '.' {
			bf.WriteByte(b)
			tk, err = r.readFloat(&bf)
			break
		} else {
			r.unreadByte(b)
			n, _ := strconv.ParseInt(bf.String(), 10, 64)
			tk = token.New(token.Int, n)
			break
		}
	}
	return
}

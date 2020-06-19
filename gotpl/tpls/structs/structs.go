// Copyright 13-Apr-2020 ºDeme
// GNU General Public License - V3 <http://www.gnu.org/licenses/>

package structs

import (
	format "fmt"
	"github.com/dedeme/gotpl/tpls/model"
	"strings"
)

type state int

const (
	identifier state = iota
	parameters
	variables
	functions
)

func splitSpaces(s string) []string {
	var r []string
	for _, e := range strings.Split(s, " ") {
		if e != "" {
			r = append(r, e)
		}
	}
	return r
}

func cutSignum(s string) string {
	if s[0] == '-' || s[0] == '+' {
		return s[1:]
	}
	return s
}

func fmt(tpl string, ps ...interface{}) string {
	return format.Sprintf(tpl, ps...)
}

func println(s string) {
	format.Println(s)
}

func emsg(ix int, msg string) error {
	return format.Errorf("Line %v: %v", ix, msg)
}

func reservedWord(word string) error {
	if strings.Index(" this _tmp __tmp ", " "+word+" ") != -1 {
		return format.Errorf("Field '%v' is not allowed", word)
	}
	return nil
}

func checkType(tp string) (rs string, err error) {
	if tp == "*" {
		err = format.Errorf("'*': Wrong type")
		return
	}
	rs = tp
	if strings.HasPrefix(tp, "map[]") {
		rs = "map[string]" + tp[5:]
	}
	return
}

// 'ix' is the number of first line of 'source' in code.
func Process(ix int, source []string) (result []string, err error) {

	var md *model.T
	st := identifier
	doc := []string{}

	for i, l := range source {
		i++
		l2 := strings.TrimSpace(l)
		if l2 == "" {
			source = append(source, l)
			continue
		} else if strings.HasPrefix(l2, "#") {
			doc = append(doc, l2)
			continue
		} else if l2 == "*/" {
			if md == nil {
				err = emsg(ix, "Struct identifier is missing")
				return
			}
			result, err = write(md)
			return
		}

		switch st {
		case identifier:
			if l2 == "-" || l2 == "---" {
				err = emsg(ix, "Struct identifier is missing")
				return
			}
			md = model.New(l2, doc)
			doc = []string{}
			st = parameters
		case parameters:
			switch l2 {
			case "-":
				doc = []string{}
				st = variables
			case "--":
				st = functions
			default:
				ps := splitSpaces(l2)
				if len(ps) != 2 {
					err = emsg(ix+i, "Expected 'name type'")
					return
				}
				if err = reservedWord(cutSignum(ps[0])); err != nil {
					err = emsg(ix+i, err.Error())
					return
				}
				ps[1], err = checkType(ps[1])
				if err != nil {
					err = emsg(ix+i, err.Error())
					return
				}
				md.AddParam(ps[0], ps[1], doc)
				doc = []string{}
			}
		case variables:
			switch l2 {
			case "--":
				st = functions
			default:
				ps := splitSpaces(l2)
				if len(ps) == 2 {
					if err = reservedWord(cutSignum(ps[0])); err != nil {
						err = emsg(ix+i, err.Error())
						return
					}
					ps[1], err = checkType(ps[1])
					if err != nil {
						err = emsg(ix+i, err.Error())
						return
					}
					md.AddVar(ps[0], ps[1], doc)
					doc = []string{}
				} else if len(ps) > 3 {
					if err = reservedWord(cutSignum(ps[0])); err != nil {
						err = emsg(ix+i, err.Error())
						return
					}
					ps[1], err = checkType(ps[1])
					if err != nil {
						err = emsg(ix+i, err.Error())
						return
					}
					md.AddVarValue(ps[0], ps[1], strings.Join(ps[2:], " "), doc)
					doc = []string{}
				} else {
					err = emsg(ix+i, "Expected 'name type [value]'")
					return
				}
			}
		default: // functions
			ps := strings.Split(l2, "··")
			for j, e := range ps {
				j++
				ps2 := strings.Split(strings.TrimSpace(e), "·")
        var fnPs [] string
        for _, e := range ps2 {
          fnPs = append(fnPs, strings.TrimSpace(e))
        }
				if len(fnPs) == 1 {
					md.AddFunc(fmt("%v(%v)", ix+i, j), fnPs[0])
				} else {
					md.AddFuncParams(fmt("%v(%v)", ix+i, j), fnPs[0], fnPs[1:])
				}
			}
		}
	}

	err = emsg(ix, "End of file reached")
	return
}

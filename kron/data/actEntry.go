// Copyright 07-Jul-2020 ºDeme
// GNU General Public License - V3 <http://www.gnu.org/licenses/>

package data

import (
	"fmt"
	"github.com/dedeme/golib/json"
	"github.com/dedeme/golib/date"
	"strconv"
	"strings"
)

type ActEntry struct {
	Id          int
	IsMonthly   bool
	Day         int
	Hour        int
	Minute      int
	Dexceptions string
	Mexceptions []int
	Command     string
	Done        string
}

func processMoment(mm string) (h, m int, err error) {
	errv := fmt.Errorf("Moment value not valid: '%v'", mm)
	ps := strings.Split(mm, ":")
	if len(ps) == 1 {
		h, err = strconv.Atoi(mm)
		if err == nil {
			if h < 0 || h > 23 {
				err = errv
			} else {
				m = 0
			}
		} else {
			err = errv
		}
	} else if len(ps) == 2 {
		h, err = strconv.Atoi(ps[0])
		if err == nil {
			m, err = strconv.Atoi(ps[1])
			if err != nil ||
				h < 0 || h > 23 ||
				m < 0 || m > 59 {
				err = errv
			}
		} else {
			err = errv
		}
	} else {
		err = errv
	}
	return
}

func processDExceptions(excs string) (es string, err error) {
	for _, b := range []byte(excs) {
		if strings.IndexByte("LMXJVSD", b) != -1 && strings.IndexByte(es, b) == -1 {
			es += string(b)
		} else {
			err = fmt.Errorf("Exceptions value not valid: '%v'", excs)
			break
		}
	}
	return
}

func processDay(day string) (d int, err error) {
	d, err = strconv.Atoi(day)
	if err != nil || d < 1 || d > 31 {
		err = fmt.Errorf("Day value not valid: '%v'", day)
	}
	return
}

func processMExceptions(excs string) (e []int, err error) {
  if excs == "" {
    return
  }
	errv := fmt.Errorf("Exceptions value not valid: '%v'", excs)
	ps := strings.Split(excs, ",")
	if len(ps) > 0 {
		for _, n := range ps {
			d, err2 := strconv.Atoi(n)
			if err2 != nil || d < 1 || d > 12 {
				err = errv
				break
			}
			for _, d2 := range e {
				if d == d2 {
					err = errv
				}
			}
			if err != nil {
				break
			}
			e = append(e, d)
		}
	} else {
		err = errv
	}
	return
}

func NewDaily(moment, exceptions, command string) (e *ActEntry, err error) {
	h, m, err := processMoment(moment)
	if err == nil {
		excs, err2 := processDExceptions(exceptions)
    err = err2
		if err == nil {
			e = &ActEntry{
				Id:          -1,
				IsMonthly:   false,
				Day:         -1,
				Hour:        h,
				Minute:      m,
				Dexceptions: excs,
				Mexceptions: []int{},
				Command:     command,
				Done:        "",
			}
		}
	}
	return
}

func NewMonthly(
	day, moment, exceptions, command string,
) (e *ActEntry, err error) {
	d, err := processDay(day)
	if err == nil {
		h, m, err2 := processMoment(moment)
    err = err2
		if err == nil {
			excs, err2 := processMExceptions(exceptions)
      err = err2
			if err == nil {
				e = &ActEntry{
					Id:          -1,
					IsMonthly:   true,
					Day:         d,
					Hour:        h,
					Minute:      m,
					Dexceptions: "",
					Mexceptions: excs,
					Command:     command,
					Done:        "",
				}
			}
		}
	}
	return
}

func (e *ActEntry) String () string {
  dt := date.NewTime(e.Day, 12, 2020, e.Hour, e.Minute, 0)
  if e.IsMonthly {
    return fmt.Sprintf(
      "(%v) m %v%v %v <%v>\n    %v",
      e.Id, dt.Format("%D("), dt.Format("%t")[:5] + ")", e.Mexceptions,
      e.Done, e.Command,
    )
  } else {
    return fmt.Sprintf(
      "(%v) d %v %v <%v>\n    %v",
      e.Id, dt.Format("%t")[:5], e.Dexceptions, e.Done, e.Command,
    )
  }
}

func (e *ActEntry) ToJs() json.T {
	var excs []json.T
	for _, e := range e.Mexceptions {
		excs = append(excs, json.Wi(e))
	}
	return json.Wa([]json.T{
		json.Wi(e.Id),
		json.Wb(e.IsMonthly),
		json.Wi(e.Day),
		json.Wi(e.Hour),
		json.Wi(e.Minute),
		json.Ws(e.Dexceptions),
		json.Wa(excs),
		json.Ws(e.Command),
		json.Ws(e.Done),
	})
}

func ActFromJs(js json.T) *ActEntry {
	a := js.Ra()
	var excs []int
	for _, e := range a[6].Ra() {
		excs = append(excs, e.Ri())
	}
	return &ActEntry{
		Id:          a[0].Ri(),
		IsMonthly:   a[1].Rb(),
		Day:         a[2].Ri(),
		Hour:        a[3].Ri(),
		Minute:      a[4].Ri(),
		Dexceptions: a[5].Rs(),
		Mexceptions: excs,
		Command:     a[7].Rs(),
		Done:        a[8].Rs(),
	}
}

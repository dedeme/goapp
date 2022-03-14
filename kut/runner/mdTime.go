// Copyright 03-Mar-2022 ºDeme
// GNU General Public License - V3 <http://www.gnu.org/licenses/>

package runner

import (
	"github.com/dedeme/kut/builtin/bfail"
	"github.com/dedeme/kut/builtin/bfunction"
	"github.com/dedeme/kut/expression"
	"strconv"
	"strings"
	"time"
)

func timeAllDigits(s string) bool {
	for i := 0; i < len(s); i++ {
		ch := s[i]
		if ch < '0' || ch > '9' {
			return false
		}
	}
	return true
}

func timeToInt(s string) (n int) {
	n, _ = strconv.Atoi(s)
	return
}

func timeFformat(tmp string, tm int64) string {
	d := time.UnixMilli(tm)
	tmp = strings.ReplaceAll(tmp, "%d", time.Time(d).Format("2"))
	tmp = strings.ReplaceAll(tmp, "%D", time.Time(d).Format("02"))
	tmp = strings.ReplaceAll(tmp, "%m", time.Time(d).Format("1"))
	tmp = strings.ReplaceAll(tmp, "%M", time.Time(d).Format("01"))
	tmp = strings.ReplaceAll(tmp, "%y", time.Time(d).Format("06"))
	tmp = strings.ReplaceAll(tmp, "%Y", time.Time(d).Format("2006"))
	tmp = strings.ReplaceAll(tmp, "%t", time.Time(d).Format("15:04:05"))
	tmp = strings.ReplaceAll(tmp, "%T", time.Time(d).Format("15:04:05.000"))
	tmp = strings.ReplaceAll(tmp, "%%", "%")
	return tmp
}

// \i, i, i, i, i, i -> i
func timeNew(args []*expression.T) (ex *expression.T, err error) {
	switch day := (args[0].Value).(type) {
	case int64:
		switch month := (args[1].Value).(type) {
		case int64:
			switch year := (args[2].Value).(type) {
			case int64:
				switch hour := (args[3].Value).(type) {
				case int64:
					switch min := (args[4].Value).(type) {
					case int64:
						switch sec := (args[5].Value).(type) {
						case int64:
							var lc *time.Location
							lc, err = time.LoadLocation("Local")
							if err == nil {
								ex = expression.MkFinal(
									time.Date(int(year), time.Month(int(month)), int(day),
										int(hour), int(min), int(sec), 0, lc).UnixMilli())
							}
						default:
							err = bfail.Type(args[5], "int")
						}
					default:
						err = bfail.Type(args[4], "int")
					}
				default:
					err = bfail.Type(args[3], "int")
				}
			default:
				err = bfail.Type(args[2], "int")
			}
		default:
			err = bfail.Type(args[1], "int")
		}
	default:
		err = bfail.Type(args[0], "int")
	}
	return
}

// \i, i, i -> i
func timeNewDate(args []*expression.T) (ex *expression.T, err error) {
	switch day := (args[0].Value).(type) {
	case int64:
		switch month := (args[1].Value).(type) {
		case int64:
			switch year := (args[2].Value).(type) {
			case int64:
				var lc *time.Location
				lc, err = time.LoadLocation("Local")
				if err == nil {
					ex = expression.MkFinal(
						time.Date(int(year), time.Month(int(month)), int(day),
							12, 0, 0, 0, lc).UnixMilli())
				}
			default:
				err = bfail.Type(args[2], "int")
			}
		default:
			err = bfail.Type(args[1], "int")
		}
	default:
		err = bfail.Type(args[0], "int")
	}
	return
}

// \-> i
func timeNow(args []*expression.T) (ex *expression.T, err error) {
	ex = expression.MkFinal(time.Now().UnixMilli())
	return
}

// \s, s -> i
func timeFromEn(args []*expression.T) (ex *expression.T, err error) {
	switch d := (args[0].Value).(type) {
	case string:
		switch sep := (args[1].Value).(type) {
		case string:
			ps := strings.Split(d, sep)
			if len(ps) != 3 || !timeAllDigits(ps[0]) ||
				!timeAllDigits(ps[1]) || !timeAllDigits(ps[2]) {
				err = bfail.Mk("'" + d + "' bad english date.")
			} else {
				var lc *time.Location
				lc, err = time.LoadLocation("Local")
				if err == nil {
					ex = expression.MkFinal(
						time.Date(timeToInt(ps[2]), time.Month(timeToInt(ps[0])),
							timeToInt(ps[1]), 12, 0, 0, 0, lc).UnixMilli())
				}
			}
		default:
			err = bfail.Type(args[1], "string")
		}
	default:
		err = bfail.Type(args[0], "string")
	}
	return
}

// \s, s -> i
func timeFromIso(args []*expression.T) (ex *expression.T, err error) {
	switch d := (args[0].Value).(type) {
	case string:
		switch sep := (args[1].Value).(type) {
		case string:
			ps := strings.Split(d, sep)
			if len(ps) != 3 || !timeAllDigits(ps[0]) ||
				!timeAllDigits(ps[1]) || !timeAllDigits(ps[2]) {
				err = bfail.Mk("'" + d + "' bad ISO date.")
			} else {
				var lc *time.Location
				lc, err = time.LoadLocation("Local")
				if err == nil {
					ex = expression.MkFinal(
						time.Date(timeToInt(ps[2]), time.Month(timeToInt(ps[1])),
							timeToInt(ps[0]), 12, 0, 0, 0, lc).UnixMilli())
				}
			}
		default:
			err = bfail.Type(args[1], "string")
		}
	default:
		err = bfail.Type(args[0], "string")
	}
	return
}

// \s -> i
func timeFromStr(args []*expression.T) (ex *expression.T, err error) {
	switch d := (args[0].Value).(type) {
	case string:
		if len(d) != 8 || !timeAllDigits(d) {
			err = bfail.Mk("'" + d + "' bad date.")
		} else {
			var lc *time.Location
			lc, err = time.LoadLocation("Local")
			if err == nil {
				ex = expression.MkFinal(
					time.Date(timeToInt(d[:4]), time.Month(timeToInt(d[4:6])),
						timeToInt(d[6:]), 12, 0, 0, 0, lc).UnixMilli())
			}
		}
	default:
		err = bfail.Type(args[0], "string")
	}
	return
}

// \i, s -> i
func timeFromClock(args []*expression.T) (ex *expression.T, err error) {
	switch d := (args[0].Value).(type) {
	case int64:
		switch cl := (args[1].Value).(type) {
		case string:
			ps := strings.Split(cl, ":")
			if len(ps) != 3 || !timeAllDigits(ps[0]) ||
				!timeAllDigits(ps[1]) || !timeAllDigits(ps[2]) ||
				len(ps[0]) != 2 || len(ps[1]) != 2 || len(ps[2]) != 2 ||
				timeToInt(ps[0]) < 0 || timeToInt(ps[0]) > 23 ||
				timeToInt(ps[1]) < 0 || timeToInt(ps[1]) > 59 ||
				timeToInt(ps[2]) < 0 || timeToInt(ps[1]) > 59 {
				err = bfail.Mk("'" + cl + "' bad clock.")
			} else {
				tm := time.UnixMilli(d)
				var lc *time.Location
				lc, err = time.LoadLocation("Local")
				if err == nil {
					ex = expression.MkFinal(
						time.Date(tm.Year(), tm.Month(), tm.Day(),
							timeToInt(ps[0]), timeToInt(ps[1]), timeToInt(ps[2]),
							0, lc).UnixMilli())
				}
			}
		default:
			err = bfail.Type(args[1], "string")
		}
	default:
		err = bfail.Type(args[0], "int")
	}
	return
}

// \i -> i
func timeDay(args []*expression.T) (ex *expression.T, err error) {
	switch d := (args[0].Value).(type) {
	case int64:
		ex = expression.MkFinal(int64(time.UnixMilli(d).Day()))
	default:
		err = bfail.Type(args[0], "int")
	}
	return
}

// \i -> i
func timeMonth(args []*expression.T) (ex *expression.T, err error) {
	switch d := (args[0].Value).(type) {
	case int64:
		ex = expression.MkFinal(int64(time.UnixMilli(d).Month()))
	default:
		err = bfail.Type(args[0], "int")
	}
	return
}

// \i -> i
func timeYear(args []*expression.T) (ex *expression.T, err error) {
	switch d := (args[0].Value).(type) {
	case int64:
		ex = expression.MkFinal(int64(time.UnixMilli(d).Year()))
	default:
		err = bfail.Type(args[0], "int")
	}
	return
}

// \i -> i
func timeWeekday(args []*expression.T) (ex *expression.T, err error) {
	switch d := (args[0].Value).(type) {
	case int64:
		ex = expression.MkFinal(int64(time.UnixMilli(d).Weekday()))
	default:
		err = bfail.Type(args[0], "int")
	}
	return
}

// \i -> i
func timeYearDay(args []*expression.T) (ex *expression.T, err error) {
	switch d := (args[0].Value).(type) {
	case int64:
		ex = expression.MkFinal(int64(time.UnixMilli(d).YearDay()))
	default:
		err = bfail.Type(args[0], "int")
	}
	return
}

// \i -> i
func timeHour(args []*expression.T) (ex *expression.T, err error) {
	switch d := (args[0].Value).(type) {
	case int64:
		ex = expression.MkFinal(int64(time.UnixMilli(d).Hour()))
	default:
		err = bfail.Type(args[0], "int")
	}
	return
}

// \i -> i
func timeMinute(args []*expression.T) (ex *expression.T, err error) {
	switch d := (args[0].Value).(type) {
	case int64:
		ex = expression.MkFinal(int64(time.UnixMilli(d).Minute()))
	default:
		err = bfail.Type(args[0], "int")
	}
	return
}

// \i -> i
func timeSecond(args []*expression.T) (ex *expression.T, err error) {
	switch d := (args[0].Value).(type) {
	case int64:
		ex = expression.MkFinal(int64(time.UnixMilli(d).Second()))
	default:
		err = bfail.Type(args[0], "int")
	}
	return
}

// \i, i -> i
func timeAddDays(args []*expression.T) (ex *expression.T, err error) {
	switch d := (args[0].Value).(type) {
	case int64:
		switch days := (args[1].Value).(type) {
		case int64:
			tm := time.UnixMilli(d).AddDate(0, 0, int(days))
			ex = expression.MkFinal(tm.UnixMilli())
		default:
			err = bfail.Type(args[1], "int")
		}
	default:
		err = bfail.Type(args[0], "int")
	}
	return
}

// \i, i -> i
func timeDfDays(args []*expression.T) (ex *expression.T, err error) {
	switch d1 := (args[0].Value).(type) {
	case int64:
		switch d2 := (args[1].Value).(type) {
		case int64:
			d1m := d1 / 86400000
			if d1 < 0 {
				d1m--
			}
			d2m := d2 / 86400000
			if d2 < 0 {
				d2m--
			}
			ex = expression.MkFinal(d1m - d2m)
		default:
			err = bfail.Type(args[1], "int")
		}
	default:
		err = bfail.Type(args[0], "int")
	}
	return
}

// \i, i -> i
func timeEqDay(args []*expression.T) (ex *expression.T, err error) {
	switch d1 := (args[0].Value).(type) {
	case int64:
		switch d2 := (args[1].Value).(type) {
		case int64:
			d1m := d1 / 86400000
			if d1 < 0 {
				d1m--
			}
			d2m := d2 / 86400000
			if d2 < 0 {
				d2m--
			}
			ex = expression.MkFinal(d1m == d2m)
		default:
			err = bfail.Type(args[1], "int")
		}
	default:
		err = bfail.Type(args[0], "int")
	}
	return
}

// \s, i -> i
func timeFormat(args []*expression.T) (ex *expression.T, err error) {
	switch tmp := (args[0].Value).(type) {
	case string:
		switch tm := (args[1].Value).(type) {
		case int64:
			ex = expression.MkFinal(timeFformat(tmp, tm))
		default:
			err = bfail.Type(args[1], "int")
		}
	default:
		err = bfail.Type(args[0], "string")
	}
	return
}

// \s, i -> i
func timeToStr(args []*expression.T) (ex *expression.T, err error) {
	switch tm := (args[0].Value).(type) {
	case int64:
		ex = expression.MkFinal(timeFformat("%Y%M%D", tm))
	default:
		err = bfail.Type(args[0], "int")
	}
	return
}

// \s, i -> i
func timeToEn(args []*expression.T) (ex *expression.T, err error) {
	switch tm := (args[0].Value).(type) {
	case int64:
		ex = expression.MkFinal(timeFformat("%M-%D-%Y", tm))
	default:
		err = bfail.Type(args[0], "int")
	}
	return
}

// \s, i -> i
func timeToIso(args []*expression.T) (ex *expression.T, err error) {
	switch tm := (args[0].Value).(type) {
	case int64:
		ex = expression.MkFinal(timeFformat("%D/%M/%Y", tm))
	default:
		err = bfail.Type(args[0], "int")
	}
	return
}

func timeGet(fname string) (fn *bfunction.T, ok bool) {
	ok = true
	switch fname {
	case "new":
		fn = bfunction.New(6, timeNew)
	case "newDate":
		fn = bfunction.New(3, timeNewDate)
	case "now":
		fn = bfunction.New(0, timeNow)
	case "fromEn":
		fn = bfunction.New(2, timeFromEn)
	case "fromIso":
		fn = bfunction.New(2, timeFromIso)
	case "fromStr":
		fn = bfunction.New(1, timeFromStr)
	case "fromClock":
		fn = bfunction.New(2, timeFromClock)
	case "day":
		fn = bfunction.New(1, timeDay)
	case "month":
		fn = bfunction.New(1, timeMonth)
	case "year":
		fn = bfunction.New(1, timeYear)
	case "weekday":
		fn = bfunction.New(1, timeWeekday)
	case "yearDay":
		fn = bfunction.New(1, timeYearDay)
	case "hour":
		fn = bfunction.New(1, timeHour)
	case "minute":
		fn = bfunction.New(1, timeMinute)
	case "second":
		fn = bfunction.New(1, timeSecond)
	case "addDays":
		fn = bfunction.New(2, timeAddDays)
	case "dfDays":
		fn = bfunction.New(2, timeDfDays)
	case "eqDay":
		fn = bfunction.New(2, timeEqDay)
	case "format":
		fn = bfunction.New(2, timeFormat)
	case "toStr":
		fn = bfunction.New(1, timeToStr)
	case "toEn":
		fn = bfunction.New(1, timeToEn)
	case "toIso":
		fn = bfunction.New(1, timeToIso)

	default:
		ok = false
	}

	return
}

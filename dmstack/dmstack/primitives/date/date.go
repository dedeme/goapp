// Copyright 10-Jan-2021 ºDeme
// GNU General Public License - V3 <http://www.gnu.org/licenses/>

// Date procedures.
package date

import (
	"github.com/dedeme/dmstack/machine"
	"github.com/dedeme/dmstack/operator"
	"github.com/dedeme/dmstack/symbol"
	"github.com/dedeme/dmstack/token"
	"strings"
	"time"
)

// Auxiliar function
func template(t string) string {
	t = strings.ReplaceAll(t, "%d", "2")
	t = strings.ReplaceAll(t, "%D", "02")
	t = strings.ReplaceAll(t, "%m", "1")
	t = strings.ReplaceAll(t, "%M", "01")
	t = strings.ReplaceAll(t, "%y", "06")
	t = strings.ReplaceAll(t, "%Y", "2006")
	t = strings.ReplaceAll(t, "%t", "15:04:05")
	t = strings.ReplaceAll(t, "%T", "15:04:05.000")
	t = strings.ReplaceAll(t, "%%", "%")
	return t
}

// Auxiliar function
func popInt(m *machine.T) int {
	tk := m.PopT(token.Int)
	i, _ := tk.I()
	return int(i)
}

// Auxiliar function
func pushInt(m *machine.T, i int) {
	m.Push(token.NewI(int64(i), m.MkPos()))
}

// Auxiliar function
func popDate(m *machine.T) time.Time {
	tk := m.PopT(token.Native)
	o, i, _ := tk.N()
	if o != operator.Date_ {
		m.Failt("\n Expected: Date object.\n  Actual  : '%v'.", o)
	}
	return i.(time.Time)
}

// Auxiliar function
func pushDate(m *machine.T, d time.Time) {
	m.Push(token.NewN(operator.Date_, d, m.MkPos()))
}

// Returns a date from year, month and day values.
//    m: Virtual machine.
func prNew(m *machine.T) {
	pushInt(m, 12)
	pushInt(m, 0)
	pushInt(m, 0)
	pushInt(m, 0)
	prNewTime(m)
}

// Returns a date from all its values.
//    m: Virtual machine.
func prNewTime(m *machine.T) {
	millis := popInt(m)
	s := popInt(m)
	mi := popInt(m)
	h := popInt(m)
	d := popInt(m)
	mt := popInt(m)
	y := popInt(m)

	lc, err := time.LoadLocation("Local")
	if err != nil {
		panic(err)
	}
	pushDate(m, time.Date(y, time.Month(mt), d, h, mi, s, millis*1000000, lc))
}

// Returns the current date.
//    m: Virtual machine.
//    m: Virtual machine.
func prNow(m *machine.T) {
	pushDate(m, time.Now())
}

// Auxiliar function
func from(m *machine.T, template string) {
	tk := m.PopT(token.String)
	s, _ := tk.S()
	d, err := time.Parse(template, s)

	if err != nil {
		m.Failt("'%v' is not a valid date", tk.StringDraft())
	}
	pushDate(m, d)
}

// Returns a Date from a string with format "yyyymmdd". If it fails, throws
// a "Date error".
//    m: Virtual machine.
func prFrom(m *machine.T) {
	from(m, "20060102")
}

// Returns a Date from a string with format "dd/mm/yyyy". If it fails, throws
// a "Date error".
//    m: Virtual machine.
func prFromIso(m *machine.T) {
	from(m, "02/01/2006")
}

// Returns a Date from a string with format "mm/dd/yyyy". If it fails, throws
// a "Date error".
//    m: Virtual machine.
func prFromEn(m *machine.T) {
	from(m, "01/02/2006")
}

// Returns a Date from a string with format.. If it fails, throws
// a "Date error".
//    m: Virtual machine.
func prFromFormat(m *machine.T) {
	tk := m.PopT(token.String)
	t, _ := tk.S()
	from(m, template(t))
}

// Auxiliar function
func to(m *machine.T, template string) {
	d := popDate(m)
	m.Push(token.NewS(d.Format(template), m.MkPos()))
}

// Returns a string with format.
//    m: Virtual machine.
func prFormat(m *machine.T) {
	tk := m.PopT(token.String)
	t, _ := tk.S()
	to(m, template(t))
}

// Returns a string with format "yyyymmdd"
//    m: Virtual machine.
func prTo(m *machine.T) {
	to(m, "20060102")
}

// Returns a string with format "dd/mm/yyyy"
//    m: Virtual machine.
func prToIso(m *machine.T) {
	to(m, "02/01/2006")
}

// Returns a string with format "mm/dd/yyyy"
//    m: Virtual machine.
func prToEn(m *machine.T) {
	to(m, "01/02/2006")
}

// Returns the year of d
//    m: Virtual machine.
func prYear(m *machine.T) {
	pushInt(m, time.Time(popDate(m)).Year())
}

// Returns the month of d
//    m: Virtual machine.
func prMonth(m *machine.T) {
	pushInt(m, int(time.Time(popDate(m)).Month()))
}

// Returns the day of d
//    m: Virtual machine.
func prDay(m *machine.T) {
	pushInt(m, time.Time(popDate(m)).Day())
}

// Returns the hour of d
//    m: Virtual machine.
func prHour(m *machine.T) {
	pushInt(m, time.Time(popDate(m)).Hour())
}

// Returns the minutes of d
//    m: Virtual machine.
func prMin(m *machine.T) {
	pushInt(m, time.Time(popDate(m)).Minute())
}

// Returns the seconds of d
//    m: Virtual machine.
func prSec(m *machine.T) {
	pushInt(m, time.Time(popDate(m)).Second())
}

// Returns the millisecond of d
//    m: Virtual machine.
func prMillis(m *machine.T) {
	pushInt(m, time.Time(popDate(m)).Nanosecond()/1000000)
}

// Returns the week number an year of d.
//    m: Virtual machine.
func prWeek(m *machine.T) {
	y, d := time.Time(popDate(m)).ISOWeek()
	pushInt(m, d)
	pushInt(m, y)
}

// Returns the weekday of d
//    m: Virtual machine.
func prWeekday(m *machine.T) {
	pushInt(m, int(time.Time(popDate(m)).Weekday()))
}

// Adds a number of days.
//    m: Virtual machine.
func prAdd(m *machine.T) {
	days := popInt(m)
	d := popDate(m)
	pushDate(m, time.Time(d).Add(time.Hour*24*time.Duration(days)))
}

// Adds a number of seconds.
//    m: Virtual machine.
func prAddSec(m *machine.T) {
	seconds := popInt(m)
	d := popDate(m)
	pushDate(m, time.Time(d).Add(time.Second*time.Duration(seconds)))
}

// Adds a number of milliseconds.
//    m: Virtual machine.
func prAddMillis(m *machine.T) {
	millis := popInt(m)
	d := popDate(m)
	pushDate(m, time.Time(d).Add(time.Millisecond*time.Duration(millis)))
}

func dfDays(d1, d2 time.Time) int64 {
	df := time.Time(d1).Sub(time.Time(d2)).Milliseconds()
	dv := df / 86400000
	if df >= 0 && df%86400000 >= 43200000 {
		dv++
	} else if df%86400000 <= -43200000 {
		dv--
	}
	return dv
}

// Returns d1 - d2 in days
//    m: Virtual machine.
func prDf(m *machine.T) {
	d2 := popDate(m)
	d1 := popDate(m)
	m.Push(token.NewI(dfDays(d1, d2), m.MkPos()))
}

// returns d1 - d2 in milliseconds
//    m: Virtual machine.
func prDfTime(m *machine.T) {
	d2 := popDate(m)
	d1 := popDate(m)
	df := time.Time(d1).Sub(time.Time(d2)).Milliseconds()
	m.Push(token.NewI(df, m.MkPos()))
}

// Returns 'true' if d1 == d2 comparing only the date.
//    m: Virtual machine.
func prEq(m *machine.T) {
	d2 := popDate(m)
	d1 := popDate(m)
	df := dfDays(d1, d2)
	m.Push(token.NewB(df == 0, m.MkPos()))
}

// Returns 'true' if d1 == d2.
//    m: Virtual machine.
func prEqTime(m *machine.T) {
	d2 := popDate(m)
	d1 := popDate(m)
	df := time.Time(d1).Sub(time.Time(d2)).Milliseconds()
	m.Push(token.NewB(df == 0, m.MkPos()))
}

// Processes date procedures.
//    m   : Virtual machine.
//    proc: Procedure
func Proc(m *machine.T, proc symbol.T) {
	switch proc {
	case symbol.New("new"):
		prNew(m)
	case symbol.New("newTime"):
		prNewTime(m)
	case symbol.New("now"):
		prNow(m)
	case symbol.From:
		prFrom(m)
	case symbol.New("fromIso"):
		prFromIso(m)
	case symbol.New("fromEn"):
		prFromEn(m)
	case symbol.New("fromFormat"):
		prFromFormat(m)
	case symbol.New("format"):
		prFormat(m)
	case symbol.New("to"):
		prTo(m)
	case symbol.New("toIso"):
		prToIso(m)
	case symbol.New("toEn"):
		prToEn(m)
	case symbol.New("year"):
		prYear(m)
	case symbol.New("month"):
		prMonth(m)
	case symbol.New("day"):
		prDay(m)
	case symbol.New("hour"):
		prHour(m)
	case symbol.New("min"):
		prMin(m)
	case symbol.New("sec"):
		prSec(m)
	case symbol.New("millis"):
		prMillis(m)
	case symbol.New("week"):
		prWeek(m)
	case symbol.New("weekday"):
		prWeekday(m)
	case symbol.New("add"):
		prAdd(m)
	case symbol.New("addSec"):
		prAddSec(m)
	case symbol.New("addMillis"):
		prAddMillis(m)
	case symbol.New("df"):
		prDf(m)
	case symbol.New("dfTime"):
		prDfTime(m)
	case symbol.New("cmp"):
		prDf(m)
	case symbol.New("cmpTime"):
		prDfTime(m)
	case symbol.New("eq"):
		prEq(m)
	case symbol.New("eqTime"):
		prEqTime(m)
	default:
		m.Failt("'date' does not contains the procedure '%v'.", proc.String())
	}
}

// Copyright 01-May-2020 ºDeme
// GNU General Public License - V3 <http://www.gnu.org/licenses/>

package tests

import (
	"github.com/dedeme/dmstack/reader"
	"github.com/dedeme/dmstack/token"
	"github.com/dedeme/dmstack/symbol"
	"testing"
)

func TestBlank(t *testing.T) {

	code := `a1
    //ds fsf
    a /* ccc */ b
    c // sdfef
    /**/d /*
    dsf
    */e
    f`

	rd := reader.New("[TEST]", code)
	tk := rd.Process()
	if !(tk.Type() == token.Procedure) {
		t.Fatal(fail)
	}
	prg, ok := tk.P()
	if !(ok) {
		t.Fatal(fail)
	}
	if r := eqi(
		len(prg),
		7,
	); r != "" {
		t.Fatal(r)
	}

	s, _ := prg[0].Sy()
	if r := eqs(
		s.String(),
		"a1",
	); r != "" {
		t.Fatal(r)
	}

	s, _ = prg[1].Sy()
	if r := eqs(
		s.String(),
		"a",
	); r != "" {
		t.Fatal(r)
	}

	s, _ = prg[2].Sy()
	if r := eqs(
		s.String(),
		"b",
	); r != "" {
		t.Fatal(r)
	}

	s, _ = prg[3].Sy()
	if r := eqs(
		s.String(),
		"c",
	); r != "" {
		t.Fatal(r)
	}

	s, _ = prg[4].Sy()
	if r := eqs(
		s.String(),
		"d",
	); r != "" {
		t.Fatal(r)
	}

	s, _ = prg[5].Sy()
	if r := eqs(
		s.String(),
		"e",
	); r != "" {
		t.Fatal(r)
	}

	s, _ = prg[6].Sy()
	if r := eqs(
		s.String(),
		"f",
	); r != "" {
		t.Fatal(r)
	}
}

func TestBool(t *testing.T) {
	code := "true a"
	rd := reader.New("[TEST]", code)
	tk := rd.Process()
	prg, _ := tk.P()
	if r := eqi(
		len(prg),
		2,
	); r != "" {
		t.Fatal(r)
	}
	s, _ := prg[0].B()
	if !(s == true) {
		t.Fatal(fail)
	}

	code = "  false : a  ,;"
	rd = reader.New("[TEST]", code)
	tk = rd.Process()
	prg, _ = tk.P()
	if r := eqi(
		len(prg),
		2,
	); r != "" {
		t.Fatal(r)
	}
	s, _ = prg[0].B()
	if !(s == false) {
		t.Fatal(fail)
	}

	code = "b  false : a  ,;"
	rd = reader.New("[TEST]", code)
	tk = rd.Process()
	prg, _ = tk.P()
	if r := eqi(
		len(prg),
		3,
	); r != "" {
		t.Fatal(r)
	}
	s, _ = prg[1].B()
	if !(s == false) {
		t.Fatal(fail)
	}

	code = ":b:  true   ,;"
	rd = reader.New("[TEST]", code)
	tk = rd.Process()
	prg, _ = tk.P()
	if r := eqi(
		len(prg),
		2,
	); r != "" {
		t.Fatal(r)
	}
	s, _ = prg[1].B()
	if !(s == true) {
		t.Fatal(fail)
	}
}

func TestNumber(t *testing.T) {
	code := "0 a"
	rd := reader.New("[TEST]", code)
	tk := rd.Process()
	prg, _ := tk.P()
	if r := eqi(
		len(prg),
		2,
	); r != "" {
		t.Fatal(r)
	}
	s, _ := prg[0].I()
	if !(s == 0) {
		t.Fatal(fail)
	}

	code = "  -0 : a  ,;"
	rd = reader.New("[TEST]", code)
	tk = rd.Process()
	prg, _ = tk.P()
	if r := eqi(
		len(prg),
		2,
	); r != "" {
		t.Fatal(r)
	}
	s, _ = prg[0].I()
	if !(s == 0) {
		t.Fatal(fail)
	}

	code = "b  1234 : a  ,;"
	rd = reader.New("[TEST]", code)
	tk = rd.Process()
	prg, _ = tk.P()
	if r := eqi(
		len(prg),
		3,
	); r != "" {
		t.Fatal(r)
	}
	s, _ = prg[1].I()
	if !(s == 1234) {
		t.Fatal(fail)
	}

	code = "b  -1234 : a  ,;"
	rd = reader.New("[TEST]", code)
	tk = rd.Process()
	prg, _ = tk.P()
	if r := eqi(
		len(prg),
		3,
	); r != "" {
		t.Fatal(r)
	}
	s, _ = prg[1].I()
	if !(s == -1234) {
		t.Fatal(fail)
	}

	code = ":b:  0xff   ,;"
	rd = reader.New("[TEST]", code)
	tk = rd.Process()
	prg, _ = tk.P()
	if r := eqi(
		len(prg),
		2,
	); r != "" {
		t.Fatal(r)
	}
	s, _ = prg[1].I()
	if !(s == 255) {
		t.Fatal(fail)
	}

	code = ":b:  0x-ff   ,;"
	rd = reader.New("[TEST]", code)
	tk = rd.Process()
	prg, _ = tk.P()
	if r := eqi(
		len(prg),
		2,
	); r != "" {
		t.Fatal(r)
	}
	s, _ = prg[1].I()
	if !(s == -255) {
		t.Fatal(fail)
	}

	code = "0.0 a"
	rd = reader.New("[TEST]", code)
	tk = rd.Process()
	prg, _ = tk.P()
	if r := eqi(
		len(prg),
		2,
	); r != "" {
		t.Fatal(r)
	}
	n, _ := prg[0].F()
	if !(n == 0.0) {
		t.Fatal(fail)
	}

	code = "  -0. : a  ,;"
	rd = reader.New("[TEST]", code)
	tk = rd.Process()
	prg, _ = tk.P()
	if r := eqi(
		len(prg),
		2,
	); r != "" {
		t.Fatal(r)
	}
	n, _ = prg[0].F()
	if !(n == 0.0) {
		t.Fatal(fail)
	}

	code = "b  -1234.22 : a  ,;"
	rd = reader.New("[TEST]", code)
	tk = rd.Process()
	prg, _ = tk.P()
	if r := eqi(
		len(prg),
		3,
	); r != "" {
		t.Fatal(r)
	}
	n, _ = prg[1].F()
	if !(n == -1234.22) {
		t.Fatal(fail)
	}

	code = "b  -01234.001 : a  ,;"
	rd = reader.New("[TEST]", code)
	tk = rd.Process()
	prg, _ = tk.P()
	if r := eqi(
		len(prg),
		3,
	); r != "" {
		t.Fatal(r)
	}
	n, _ = prg[1].F()
	if !(n == -1234.001) {
		t.Fatal(fail)
	}

	code = ":b:  1.e1   ,;"
	rd = reader.New("[TEST]", code)
	tk = rd.Process()
	prg, _ = tk.P()
	if r := eqi(
		len(prg),
		2,
	); r != "" {
		t.Fatal(r)
	}
	n, _ = prg[1].F()
	if !(n == 10.0) {
		t.Fatal(fail)
	}

	code = ":b:  1.E-1   ,;"
	rd = reader.New("[TEST]", code)
	tk = rd.Process()
	prg, _ = tk.P()
	if r := eqi(
		len(prg),
		2,
	); r != "" {
		t.Fatal(r)
	}
	n, _ = prg[1].F()
	if !(n == 0.1) {
		t.Fatal(fail)
	}

}

func TestString(t *testing.T) {
  symbol.Initialize ()
	code := "\"\""
	rd := reader.New("[TEST]", code)
	tk := rd.Process()
	prg, _ := tk.P()
	if r := eqi(
		len(prg),
		1,
	); r != "" {
		t.Fatal(r)
	}
	n, _ := prg[0].S()
	if !(n == "") {
		t.Fatal(fail)
	}

	code = "\"\" abc "
	rd = reader.New("[TEST]", code)
	tk = rd.Process()
	prg, _ = tk.P()
	if r := eqi(
		len(prg),
		2,
	); r != "" {
		t.Fatal(r)
	}
	n, _ = prg[0].S()
	if !(n == "") {
		t.Fatal(fail)
	}

	code = "bc \"\" abc "
	rd = reader.New("[TEST]", code)
	tk = rd.Process()
	prg, _ = tk.P()
	if r := eqi(
		len(prg),
		3,
	); r != "" {
		t.Fatal(r)
	}
	n, _ = prg[1].S()
	if !(n == "") {
		t.Fatal(fail)
	}

	code = "bc \"\"  "
	rd = reader.New("[TEST]", code)
	tk = rd.Process()
	prg, _ = tk.P()
	if r := eqi(
		len(prg),
		2,
	); r != "" {
		t.Fatal(r)
	}
	n, _ = prg[1].S()
	if !(n == "") {
		t.Fatal(fail)
	}

	code = "\"año\\\"B\\\":\\n \\\\A \\\\n \\u263a\""
	rd = reader.New("[TEST]", code)
	tk = rd.Process()
	prg, _ = tk.P()
	if r := eqi(
		len(prg),
		1,
	); r != "" {
		t.Fatal(r)
	}
	n, _ = prg[0].S()
	if r := eqs(
		n,
		"año\"B\":\n \\A \\n \u263a",
	); r != "" {
		t.Fatal(r)
	}

	code = "\"año\\\"B\\\":\\n \\\\A \\\\n \\u263a\" .:. abc "
	rd = reader.New("[TEST]", code)
	tk = rd.Process()
	prg, _ = tk.P()
	if r := eqi(
		len(prg),
		2,
	); r != "" {
		t.Fatal(r)
	}
	n, _ = prg[0].S()
	if r := eqs(
		n,
		"año\"B\":\n \\A \\n \u263a",
	); r != "" {
		t.Fatal(r)
	}

	code = " bc  \"año\\\"B\\\":\\n \\\\A \\\\n \\u263a\" abc "
	rd = reader.New("[TEST]", code)
	tk = rd.Process()
	prg, _ = tk.P()
	if r := eqi(
		len(prg),
		3,
	); r != "" {
		t.Fatal(r)
	}
	n, _ = prg[1].S()
	if r := eqs(
		n,
		"año\"B\":\n \\A \\n \u263a",
	); r != "" {
		t.Fatal(r)
	}

	code = " bc  \"año\\\"B\\\":\\n \\\\A \\\\n \\u263a\" "
	rd = reader.New("[TEST]", code)
	tk = rd.Process()
	prg, _ = tk.P()
	if r := eqi(
		len(prg),
		2,
	); r != "" {
		t.Fatal(r)
	}
	n, _ = prg[1].S()
	if r := eqs(
		n,
		"año\"B\":\n \\A \\n \u263a",
	); r != "" {
		t.Fatal(r)
	}
}

func TestString2(t *testing.T) {
  symbol.Initialize ()

	code := "`\nThis\nis \"so\" 3245€`"
	rd := reader.New("[TEST]", code)
	tk := rd.Process()
	prg, _ := tk.P()
	if r := eqi(
		len(prg),
		1,
	); r != "" {
		t.Fatal(r)
	}
	n, _ := prg[0].S()
	if r := eqs(
		n,
		"This\nis \"so\" 3245€",
	); r != "" {
		t.Fatal(r)
	}

	code = "`TX\nThis\nis \"so\" 32`45€TX`"
	rd = reader.New("[TEST]", code)
	tk = rd.Process()
	prg, _ = tk.P()
	if r := eqi(
		len(prg),
		1,
	); r != "" {
		t.Fatal(r)
	}
	n, _ = prg[0].S()
	if r := eqs(
		n,
		"This\nis \"so\" 32`45€",
	); r != "" {
		t.Fatal(r)
	}

	code = "a `TX\nThis\nis \"so\" 32`45€TX` b"
	rd = reader.New("[TEST]", code)
	tk = rd.Process()
	prg, _ = tk.P()
	if r := eqi(
		len(prg),
		3,
	); r != "" {
		t.Fatal(r)
	}
	n, _ = prg[1].S()
	if r := eqs(
		n,
		"This\nis \"so\" 32`45€",
	); r != "" {
		t.Fatal(r)
	}
}

func TestProcedure(t *testing.T) {
  symbol.Initialize()
	code := "()"
	rd := reader.New("[TEST]", code)
	tk := rd.Process()
	prg, _ := tk.P()
	if r := eqi(
		len(prg),
		1,
	); r != "" {
		t.Fatal(r)
	}
	p, _ := prg[0].P()
	if r := eqi(
		len(p),
		0,
	); r != "" {
		t.Fatal(r)
	}

	code = " a ()"
	rd = reader.New("[TEST]", code)
	tk = rd.Process()
	prg, _ = tk.P()
	if r := eqi(
		len(prg),
		2,
	); r != "" {
		t.Fatal(r)
	}
	p, _ = prg[1].P()
	if r := eqi(
		len(p),
		0,
	); r != "" {
		t.Fatal(r)
	}

	code = "a ()  b"
	rd = reader.New("[TEST]", code)
	tk = rd.Process()
	prg, _ = tk.P()
	if r := eqi(
		len(prg),
		3,
	); r != "" {
		t.Fatal(r)
	}
	p, _ = prg[1].P()
	if r := eqi(
		len(p),
		0,
	); r != "" {
		t.Fatal(r)
	}

	code = " a ( 101 )"
	rd = reader.New("[TEST]", code)
	tk = rd.Process()
	prg, _ = tk.P()
	if r := eqi(
		len(prg),
		2,
	); r != "" {
		t.Fatal(r)
	}
	p, _ = prg[1].P()
	if r := eqi(
		len(p),
		1,
	); r != "" {
		t.Fatal(r)
	}
	if r := eqs(
		p[0].String(),
		"101",
	); r != "" {
		t.Fatal(r)
	}

	code = "a(101 : b =)c d"
	rd = reader.New("[TEST]", code)
	tk = rd.Process()
	prg, _ = tk.P()
	if r := eqi(
		len(prg),
		4,
	); r != "" {
		t.Fatal(r)
	}
	p, _ = prg[1].P()
	if r := eqi(
		len(p),
		3,
	); r != "" {
		t.Fatal(r)
	}
	if r := eqs(
		p[0].String(),
		"101",
	); r != "" {
		t.Fatal(r)
	}
	if r := eqs(
		p[1].String(),
		"<.b.>",
	); r != "" {
		t.Fatal(r)
	}
	if r := eqs(
		p[2].String(),
		"<.=.>",
	); r != "" {
		t.Fatal(r)
	}


  dataTk := token.NewSy(symbol.Data, nil).String()
  mapTk := token.NewSy(symbol.Map, nil).String()
  fromListTk := token.NewSy(symbol.FromList, nil).String()
	code = " a (101 (true 2) [\"a\"] {\"a\": 1} ) c"
	rd = reader.New("[TEST]", code)
	tk = rd.Process()
	prg, _ = tk.P()
	if r := eqi(
		len(prg),
		3,
	); r != "" {
		t.Fatal(r)
	}
	p, _ = prg[1].P()
	if r := eqi(
		len(p),
		8,
	); r != "" {
		t.Fatal(r)
	}
	if r := eqs(
		p[0].String(),
		"101",
	); r != "" {
		t.Fatal(r)
	}
	if r := eqs(
		p[1].String(),
		"(true,2)",
	); r != "" {
		t.Fatal(r)
	}
	if r := eqs(
		p[2].String(),
		"(\"a\")",
	); r != "" {
		t.Fatal(r)
	}
	if r := eqs(
		p[3].String(),
		dataTk,
	); r != "" {
		t.Fatal(r)
	}
	if r := eqs(
		p[4].String(),
		"(\"a\",1)",
	); r != "" {
		t.Fatal(r)
	}
	if r := eqs(
		p[5].String(),
		dataTk,
	); r != "" {
		t.Fatal(r)
	}
	if r := eqs(
		p[6].String(),
		mapTk,
	); r != "" {
		t.Fatal(r)
	}
	if r := eqs(
		p[7].String(),
		fromListTk,
	); r != "" {
		t.Fatal(r)
	}

	code = "[102, (33 a =), [\"bc\", \"\"], {\n  \"a\": 1\n  \"b\": 2\n}]"
	rd = reader.New("[TEST]", code)
	tk = rd.Process()
	if r := eqs(
		tk.String(),
		"((" +
    "102,"+
    "(33,<.a.>,<.=.>),"+
    "(\"bc\",\"\")," + dataTk + "," +
    "(\"a\",1,\"b\",2)," + dataTk + "," + mapTk + "," + fromListTk +
    "))",
	); r != "" {
		t.Fatal(r)
	}

}


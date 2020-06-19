// Copyright 28-Apr-2020 ºDeme
// GNU General Public License - V3 <http://www.gnu.org/licenses/>

package tests

import (
	"github.com/dedeme/dmstack/symbol"
	"github.com/dedeme/dmstack/token"
	"testing"
)

func TestCons(t *testing.T) {
  symbol.Initialize ()

  pos := token.NewPos(symbol.New("a"), 0)
	vB, ok := token.NewB(true, pos).B()
	if !(vB && ok) {
		t.Fatal(fail)
	}

	vI, ok := token.NewI(23, pos).I()
	if !(vI == 23 && ok) {
		t.Fatal(fail)
	}

	vF, ok := token.NewF(23.56, pos).F()
	if !(vF == 23.56 && ok) {
		t.Fatal(fail)
	}

	vP, ok := token.NewP([]*token.T{}, pos).P()
	if !(vP != nil && ok) {
		t.Fatal(fail)
	}

	vL, ok := token.NewL([]*token.T{}, pos).L()
	if !(vL != nil && ok) {
		t.Fatal(fail)
	}

	vM, ok := token.NewM(make(map[*token.T]*token.T), pos).M()
	if !(vM != nil && ok) {
		t.Fatal(fail)
	}

	vSy, ok := token.NewSy(symbol.New("abc"), pos).Sy()
	if !(vSy.String() == "abc" && ok) {
		t.Fatal(fail)
	}

  s := "xyz"
	vNs, vNv, ok := token.NewN(symbol.New("abz"), &s, pos).N()
	if !(vNs.String() == "abz" && *(vNv.(*string)) == "xyz" && ok) {
		t.Fatal(fail)
	}

}

func TestType(t *testing.T) {
  pos := token.NewPos(symbol.New("a"), 0)

	if r := eqs(
		token.NewB(true, pos).Type().String(),
		"Bool",
	); r != "" {
		t.Fatal(r)
	}

	if r := eqs(
		token.NewF(35.32, pos).Type().String(),
		"Float",
	); r != "" {
		t.Fatal(r)
	}

	if r := eqs(
		token.NewS("35", pos).Type().String(),
		"String",
	); r != "" {
		t.Fatal(r)
	}

	if r := eqs(
		token.NewP([]*token.T{}, pos).Type().String(),
		"Procedure",
	); r != "" {
		t.Fatal(r)
	}

	if r := eqs(
		token.NewL([]*token.T{}, pos).Type().String(),
		"List",
	); r != "" {
		t.Fatal(r)
	}

	if r := eqs(
		token.NewM(make(map[*token.T]*token.T), pos).Type().String(),
		"Map",
	); r != "" {
		t.Fatal(r)
	}

	if r := eqs(
		token.NewSy(symbol.New("abc"), pos).Type().String(),
		"Symbol",
	); r != "" {
		t.Fatal(r)
	}

  s := "xyz"
	if r := eqs(
		token.NewN(symbol.New("abz"), &s, pos).Type().String(),
		"Native",
	); r != "" {
		t.Fatal(r)
	}

}

func TestClone(t *testing.T) {
  pos := token.NewPos(symbol.New("a"), 0)

	tk := token.NewB(false, pos)
	if !(tk.Eq(tk.Clone())) {
		t.Fatal(fail)
	}
	tk = token.NewI(33, pos)
	if !(tk.Eq(tk.Clone())) {
		t.Fatal(fail)
	}
	tk = token.NewF(12.55, pos)
	if !(tk.Eq(tk.Clone())) {
		t.Fatal(fail)
	}
	tk = token.NewS("abcd", pos)
	if !(tk.Eq(tk.Clone())) {
		t.Fatal(fail)
	}

	tk0 := token.NewL([]*token.T{token.NewI(12, pos), token.NewS("bc", pos)}, pos)
	m1 := make(map[*token.T]*token.T)
	m1[token.NewS("a", pos)] = token.NewI(54, pos)
	m1[token.NewS("b", pos)] = tk0
	m1[token.NewS("C", pos)] = token.NewS("ZZZ", pos)
	tk1 := token.NewM(m1, pos)

	tk = token.NewP([]*token.T{tk0, token.NewS("xxx", pos), tk1}, pos)
	if !(tk.Eq(tk.Clone())) {
		t.Fatal(fail)
	}

	tk = token.NewL([]*token.T{tk0, token.NewS("xxx", pos), tk1}, pos)
	if !(tk.Eq(tk.Clone())) {
		t.Fatal(fail)
	}

	tk = token.NewM(m1, pos)
	if !(tk.Eq(tk.Clone())) {
		t.Fatal(fail)
	}

	tk = token.NewSy(symbol.New("xxx"), pos)
	if !(tk.Eq(tk.Clone())) {
		t.Fatal(fail)
	}

  s := "a test"
	tk = token.NewN(symbol.New("ss"), &s, pos)
	if !(tk.Eq(tk.Clone())) {
		t.Fatal(fail)
	}
}

func TestTkToString(t *testing.T) {
  pos := token.NewPos(symbol.New("a"), 0)

	sum := func(s string) int {
		r := 0
		bs := []byte(s)
		for i := 0; i < len(bs); i++ {
      r += int(bs[i])
		}
    return r
	}

	tk := token.NewB(false, pos)
	if r := eqs(
		tk.String(),
		"false",
	); r != "" {
		t.Fatal(r)
	}

	tk = token.NewI(33, pos)
	if r := eqs(
		tk.String(),
		"33",
	); r != "" {
		t.Fatal(r)
	}

	tk = token.NewF(12.55, pos)
	if r := eqs(
		tk.String(),
		"12.55",
	); r != "" {
		t.Fatal(r)
	}

	tk = token.NewS("ab\"€\"cd", pos)
	if r := eqs(
		tk.String(),
		"ab\"€\"cd",
	); r != "" {
		t.Fatal(r)
	}

	tk = token.NewS("ab\"€\"cd", pos)
	if r := eqs(
		tk.StringDraft(),
		"\"ab\\\"€\\\"cd\"",
	); r != "" {
		t.Fatal(r)
	}

	tk0 := token.NewL([]*token.T{token.NewI(12, pos), token.NewS("bc", pos)}, pos)
	m1 := make(map[*token.T]*token.T)
	m1[token.NewS("a", pos)] = token.NewI(54, pos)
	m1[token.NewS("b", pos)] = tk0
	m1[token.NewS("C", pos)] = token.NewS("ZZZ", pos)
	tk1 := token.NewM(m1, pos)

	tk = token.NewP([]*token.T{tk0, token.NewS("xxx", pos), tk1}, pos)
	if r := eqi(
		sum(tk.String()),
		sum("([12,\"bc\"],\"xxx\",{\"b\":[12,\"bc\"],\"C\":\"ZZZ\",\"a\":54})"),
	); r != "" {
		t.Fatal(r)
	}

	tk = token.NewL([]*token.T{tk0, token.NewS("xxx", pos), tk1}, pos)
	if r := eqi(
		sum(tk.String()),
		sum("[[12,\"bc\"],\"xxx\",{\"b\":[12,\"bc\"],\"C\":\"ZZZ\",\"a\":54}]"),
	); r != "" {
		t.Fatal(r)
	}

	tk = token.NewM(m1, pos)
	if r := eqi(
		sum(tk.String()),
		sum("{\"a\":54,\"b\":[12,\"bc\"],\"C\":\"ZZZ\"}"),
	); r != "" {
		t.Fatal(r)
	}

	tk = token.NewSy(symbol.New("xxx"), pos)
	if r := eqs(
		tk.String(),
		"<.xxx.>",
	); r != "" {
		t.Fatal(r)
	}

  s := "a test"
	tk = token.NewN(symbol.New("ss"), &s, pos)
	if r := eqs(
		tk.String()[0:7],
		"<.ss|0x",
	); r != "" {
		t.Fatal(r)
	}
	if r := eqi(
		len(tk.String()),
		19,
	); r != "" {
		t.Fatal(r)
	}
}

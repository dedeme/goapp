// Copyright 13-Apr-2020 ºDeme
// GNU General Public License - V3 <http://www.gnu.org/licenses/>

package funcs

import (
	"github.com/dedeme/gotpl/tpls/model"
)

func mkType(m *model.T) string {
  if m.IsBasicType() {
    return m.Id
  }
  return "*" + m.Id
}

func mkName(m *model.T, name string) string {
  if m.IsMainType() {
    return name
  }
  return name + capitalize(m.Id)
}

func wFilter(m *model.T) (result []string) {
  tp := mkType(m)
	result = append(result, fmt(
		"func %v(sl []%v, fn func (%v) bool) (rs []%v) {",
		mkName(m, "Filter"), tp, tp, tp,
	))
	result = append(result, "  for _, e := range sl {")
	result = append(result, "    if fn(e) {")
	result = append(result, "      rs = append(rs, e)")
	result = append(result, "    }")
	result = append(result, "  }")
	result = append(result, "  return")
	result = append(result, "}")
	return
}

func wMap(m *model.T, fnName, tpTarget string) (result []string) {
  tp := mkType(m)
	result = append(result, fmt(
		"func %v(sl []%v, fn func (%v) %v) (rs []%v) {",
		fnName, tp, tp, tpTarget, tpTarget,
	))
	result = append(result, "  for _, e := range sl {")
	result = append(result, "    rs = append(rs, fn(e))")
	result = append(result, "  }")
	result = append(result, "  return")
	result = append(result, "}")
	return
}

func wTake(m *model.T) (result []string) {
  tp := mkType(m)
	result = append(result, fmt(
		"func %v(sl []%v, n int) (rs []%v) {",
		mkName(m, "Take"), tp, tp,
	))
	result = append(result, "  if n < 1 {")
	result = append(result, "    return")
	result = append(result, "  }")
	result = append(result, "  if n > len(sl) {")
	result = append(result, "    n = len(sl)")
	result = append(result, "  }")
	result = append(result, "  return sl[0:n]")
	result = append(result, "}")
	return
}

func wDrop(m *model.T) (result []string) {
  tp := mkType(m)
	result = append(result, fmt(
		"func %v(sl []%v, n int) (rs []%v) {",
		mkName(m, "Drop"), tp, tp,
	))
	result = append(result, "  if n < 1 {")
	result = append(result, "    return sl")
	result = append(result, "  }")
	result = append(result, "  if n >= len(sl) {")
	result = append(result, "    return")
	result = append(result, "  }")
	result = append(result, "  return sl[n:]")
	result = append(result, "}")
	return
}


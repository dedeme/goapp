// Copyright 08-Apr-2020 ºDeme
// GNU General Public License - V3 <http://www.gnu.org/licenses/>

package tpls

// map template.
//
// Forms:
//    ·sort·iname·type
//      ·sort·sortInt·string
//      -----------
//      var tg []string
//      for _, e := range src {
//        tg = append(tg, fcopy(e+4))
//      }
func Psort(pars string) (r string, err string) {
	p1, rs := Kv(pars)
	if p1 == "" {
		err = fmt("Parameters missings in '%v'", pars)
		return
	}
	p2, rs := Kv(rs)
	if p2 == "" {
		err = fmt("Parameters missings in '%v'", pars)
		return
	}
	r = fmt(
		"type %v []%v\n"+
			"func (s %v) Less(i, j int) bool {\n"+
			"    return s[i] < s[j]\n"+
			"}\n"+
			"func (s %v) Len() int {\n"+
			"    return len(s)\n"+
			"}\n"+
			"func (s %v) Swap(i, j int) {\n"+
			"    s[i], s[j] = s[j], s[i]\n"+
			"}\n",
		p1, p2, p1, p1, p1,
	)
	return
}

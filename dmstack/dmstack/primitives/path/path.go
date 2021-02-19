// Copyright 11-Jan-2021 ºDeme
// GNU General Public License - V3 <http://www.gnu.org/licenses/>

// File path management.
package path

import (
	"github.com/dedeme/dmstack/machine"
	"github.com/dedeme/dmstack/symbol"
	"github.com/dedeme/dmstack/token"
	"path"
	"strings"
)

// Auxiliar function.
func popStr(m *machine.T) string {
	tk := m.PopT(token.String)
	s, _ := tk.S()
	return s
}

// Auxiliar function.
func pushStr(m *machine.T, s string) {
	m.Push(token.NewS(s, m.MkPos()))
}

// Returns a canonical representation of 'p'.
//    m: Virtual machine.
func prCanonical(m *machine.T) {
	pushStr(m, path.Clean(popStr(m)))
}

// Returns 'p1' joined to 'p2'.
//    m: Virtual machine.
func prAdd(m *machine.T) {
	p2 := popStr(m)
	p1 := popStr(m)
	pushStr(m, path.Join(p1, p2))
}

// Returns paths of 'a' joined.
//    m: Virtual machine.
func prJoin(m *machine.T) {
	tk := m.PopT(token.Array)
	a, _ := tk.A()
	var ps []string
	for _, e := range a {
		p, ok := e.S()
		if !ok {
			m.Failt("\n  Expected: String.\n  Actual:  '%v'.", e.StringDraft())
		}
		ps = append(ps, p)
	}
	pushStr(m, path.Join(ps...))
}

// Returns the extension of 'p'.
//    m: Virtual machine.
func prExtension(m *machine.T) {
	p := popStr(m)
	i := strings.LastIndexByte(p, '.')
	if i == -1 {
		pushStr(m, "")
		return
	}
	pushStr(m, p[i:])
}

// Auxiliar function
func name(p string) string {
	if p == "" || strings.HasSuffix(p, "/") {
		return ""
	}
	i := strings.LastIndexByte(p, '/')
	if i == -1 {
		return p
	}
	return p[i+1:]
}

// Returns the name with extension of 'p'.
//    m: Virtual machine.
func prName(m *machine.T) {
	pushStr(m, name(popStr(m)))
}

// Returns the name without extension of 'p'.
//    m: Virtual machine.
func prOnlyName(m *machine.T) {
	n := name(popStr(m))
	i := strings.LastIndexByte(n, '.')
	if i == -1 {
		pushStr(m, n)
		return
	}
	pushStr(m, n[:i])
}

// Returns the parent of 'p'
//    m: Virtual machine.
func prParent(m *machine.T) {
	p := popStr(m)
	i := strings.LastIndexByte(p, '/')
	if i == -1 {
		i = 0
	}
	pushStr(m, p[:i])
}

// Processes path procedures.
//    m   : Virtual machine.
//    proc: Procedure
func Proc(m *machine.T, proc symbol.T) {
	switch proc {
	case symbol.New("canonical"):
		prCanonical(m)
	case symbol.New("add"):
		prAdd(m)
	case symbol.New("join"):
		prJoin(m)
	case symbol.New("extension"):
		prExtension(m)
	case symbol.New("name"):
		prName(m)
	case symbol.New("onlyName"):
		prOnlyName(m)
	case symbol.New("parent"):
		prParent(m)
	default:
		m.Failt("'path' does not contains the procedure '%v'.", proc.String())
	}
}

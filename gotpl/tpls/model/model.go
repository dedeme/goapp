// Copyright 13-Apr-2020 ºDeme
// GNU General Public License - V3 <http://www.gnu.org/licenses/>

// Model for struct templates
package model

import (
	"fmt"
	"strings"
)

type Param struct {
	Id   string
	Type string
	Doc  []string
}

type Var struct {
	Id    string
	Type  string
	Value string // If it is not defined, its value is ""
	Doc   []string
}

type Func struct {
	Nline  string
	Id     string
	Params []string
}

type T struct {
	Id     string
	Params []*Param
	Vars   []*Var
	Funcs  []*Func
	Doc    []string
}

func New(id string, doc []string) *T {
	return &T{id, []*Param{}, []*Var{}, []*Func{}, doc}
}

func (m *T) AddParam(id, tp string, doc []string) error {
	for _, e := range m.Params {
		if e.Id == id {
			return fmt.Errorf("Parameter '%v' is duplicated", id)
		}
	}
	m.Params = append(m.Params, &Param{id, tp, doc})
	return nil
}

func (m *T) AddVar(id, tp string, doc []string) error {
	for _, e := range m.Vars {
		if e.Id == id {
			return fmt.Errorf("Variable '%v' is duplicated", id)
		}
	}
	m.Vars = append(m.Vars, &Var{id, tp, "", doc})
	return nil
}

func (m *T) AddVarValue(id, tp, value string, doc []string) error {
	for _, e := range m.Vars {
		if e.Id == id {
			return fmt.Errorf("Variable '%v' is duplicated", id)
		}
	}
	m.Vars = append(m.Vars, &Var{id, tp, value, doc})
	return nil
}

func (m *T) AddFunc(nline, id string) error {
	for _, e := range m.Funcs {
		if e.Id == id {
			return fmt.Errorf("Function '%v' is duplicated", id)
		}
	}
	m.Funcs = append(m.Funcs, &Func{nline, id, []string{}})
	return nil
}

func (m *T) AddFuncParams(nline, id string, params []string) error {
	for _, e := range m.Funcs {
		if e.Id == id {
			return fmt.Errorf("Function '%v' is duplicated", id)
		}
	}
	m.Funcs = append(m.Funcs, &Func{nline, id, params})
	return nil
}

func (m *T) IsBasicType() bool {
	return strings.Index(
		" bool int int32 int64 float32 float64 string ",
		" "+m.Id+" ",
	) != -1
}

func (m *T) IsMainType() bool {
	return m.Id == "T" || m.Id == "t"
}

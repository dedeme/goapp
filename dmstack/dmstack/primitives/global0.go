// Copyright 16-May-2020 ºDeme
// GNU General Public License - V3 <http://www.gnu.org/licenses/>

package primitives

import (
	"github.com/dedeme/dmstack/machine"
	"github.com/dedeme/dmstack/token"
	"github.com/dedeme/dmstack/imports"
	"github.com/dedeme/dmstack/reader"
  "fmt"
  "os"
  "path"
  "io/ioutil"
)

// Shows and consume a draft of the last element of stack.
//    m: Virtual machine.
func prPuts (m *machine.T) {
  tk := m.Pop()
  fmt.Println(tk.StringDraft())
}

// Raise a fail if the last element of stack is not a Bool or it is 'false'.
//    m: Virtual machine.
func prAssert (m *machine.T) {
  ok, _ := m.Pop1(token.Bool).B()
  if ok {
    return
  }
  m.Fail("Assert error")
}

// Raise a fail if the two last elements of stack are not equals.
//    m: Virtual machine.
func prExpect (m *machine.T) {
  expected := m.Pop()
  actual := m.Pop()
  if expected.Eq(actual) {
    return
  }
  m.Failf(
    "Expect error:\n  Expected: %s\n  Actual  : %s",
    expected.StringDraft(), actual.StringDraft(),
  )
}

// Execute a procedure saved in stack.
//    m: Virtual machine.
//    run: Function which running a machine.
func prRun (m *machine.T, run func (m *machine.T)) {
  tk := m.Pop1(token.Procedure)
  proc, _ := tk.P()
  run (machine.New(m.SourceDir, m.Pmachines, proc))
}

// Execute a procedure in an isolated machine. The stack of the isolated
// machine is put en the current stack as a list.
//    m: Virtual machine
//    run: Function which running a machine.
func prData (m *machine.T, run func (m *machine.T)) {
  tk := m.Pop1(token.Procedure)
  proc, _ := tk.P()
  m2 := machine.NewIsolate(m.SourceDir, m.Pmachines, proc)
  run(m2)
  m.Push(token.NewL(*m2.Stack, m.MkPos()))
}

// Process an 'import'.
//    m: Virtual machine.
//    run: Function which running a machine.
func prImport (m *machine.T, run func (m *machine.T)) {
  tk := m.Pop()
  sourceKv, err := imports.ReadSymbol(tk)
  if err != nil {
    m.Fail(err.Error())
  }
  source := sourceKv.Value
  sourcef := path.Clean(path.Join(m.SourceDir, source.String()))

  onWay := imports.IsOnWay(source)
  _, imported := imports.Get(source)
  if !onWay && !imported {
    f := sourcef + ".dms"
    if _, err := os.Stat(f); err != nil {
      m.Failf("File '%v' not found.", f)
      return
    }

    imports.PutOnWay(source)
    bs, err := ioutil.ReadFile(f)
    if err != nil {
      panic(err)
    }
    rd := reader.New(source.String(), string(bs))

    prg, ok := rd.Process().P()
    if !ok {
      m.Fail("Reader process does not return a Program")
    }

    m2 := machine.NewIsolate(path.Dir(f), []*machine.T{}, prg)
    run(m2)

    imports.Add(source, m2.Heap)
    imports.QuitOnWay(source)
  } else if (onWay) {
    m.Failf("Cyclic import of '%v'", source)
  }
  m.ImportsAdd(source)
}

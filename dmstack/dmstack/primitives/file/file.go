// Copyright 11-Jan-2021 ºDeme
// GNU General Public License - V3 <http://www.gnu.org/licenses/>

// File procedures.
package file

import (
	"bufio"
	"fmt"
	"github.com/dedeme/dmstack/machine"
	"github.com/dedeme/dmstack/operator"
	"github.com/dedeme/dmstack/primitives/it"
	"github.com/dedeme/dmstack/symbol"
	"github.com/dedeme/dmstack/token"
	"io"
	"io/ioutil"
	"os"
)

// Buffer to read binary files.
const BufferSize = 8192

// Auxiliar function
func excf(m *machine.T, template string, values ...interface{}) {
	panic(&machine.Error{
		m, machine.EFile(), fmt.Sprintf(template, values...),
	})
}

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

// Auxiliar function
func dir(m *machine.T, d string) []os.FileInfo {
	fis, err := ioutil.ReadDir(d)
	if err != nil {
		excf(m, "Fail reading the directory '%v'.", d)
	}
	return fis
}

// Returns the working directory.
//    m: Virtual machine.
func prCwd(m *machine.T) {
	p, err := os.Getwd()
	if err != nil {
		excf(m, "Fail changing work directory to '%v'.", p)
	}
	pushStr(m, p)
}

// Changes the working directory to 'd'.
//    m: Virtual machine.
func prCd(m *machine.T) {
	d := popStr(m)
	err := os.Chdir(d)
	if err != nil {
		excf(m, "Fail changing work directory to '%v'.", d)
	}
}

// Creates the directory 'd'. If 'd' already exists, it does nothing.
//    m: Virtual machine.
func prMkdir(m *machine.T) {
	d := popStr(m)
	err := os.MkdirAll(d, os.FileMode(0755))
	if err != nil {
		excf(m, "Fail creating the directory '%v'.", d)
	}
}

// Returns an array with the name of files in 'd'.
//    m: Virtual machine.
func prDir(m *machine.T) {
	pos := m.MkPos()
	d := popStr(m)
	var r []*token.T
	for _, e := range dir(m, d) {
		r = append(r, token.NewS(e.Name(), pos))
	}
	m.Push(token.NewA(r, pos))
}

// Returns 'true' if 'd' exists and is a directory.
//    m: Virtual machine.
func prIsDirectory(m *machine.T) {
	d := popStr(m)
	r := false
	if info, err := os.Stat(d); err == nil && info.IsDir() {
		r = true
	}
	m.Push(token.NewB(r, m.MkPos()))
}

// Returns 'true' if 'p' exists.
//    m: Virtual machine.
func prExists(m *machine.T) {
	p := popStr(m)
	r := false
	if _, err := os.Stat(p); err == nil {
		r = true
	}
	m.Push(token.NewB(r, m.MkPos()))
}

// Remove 'p' from file system. If 'p' does not exists, it does nothing.
//    m: Virtual machine.
func prDel(m *machine.T) {
	p := popStr(m)
	err := os.RemoveAll(p)
	if err != nil {
		excf(m, "Fail deleting file '%v'.", p)
	}
}

// Rename old as new.
//    m: Virtual machine.
func prRename(m *machine.T) {
	new := popStr(m)
	old := popStr(m)
	err := os.Rename(old, new)
	if err != nil {
		excf(m, "Fail renaming file '%v'.", old)
	}
}

// link old as new.
//    m: Virtual machine.
func prLink(m *machine.T) {
	new := popStr(m)
	old := popStr(m)
	err := os.Symlink(old, new)
	if err != nil {
		excf(m, "Fail linking file '%v'.", err)
	}
}

// copy old in new.
//    m: Virtual machine.
func prCopy(m *machine.T) {
	new := popStr(m)
	old := popStr(m)

	source, err := os.Open(old)
	if err != nil {
		excf(m, "Fail openning source file '%v' for copy.", old)
	}
	target, err := os.Create(new)
	if err != nil {
		excf(m, "Fail openning target file '%v' for copy.", new)
	}

	defer source.Close()
	defer target.Close()

	buf := make([]byte, BufferSize)
	for {
		n, err := source.Read(buf)
		if err != nil && err != io.EOF {
			excf(m, "Fail reading source file '%v' in copy.", old)
		}
		if n == 0 {
			break
		}
		if _, err := target.Write(buf[:n]); err != nil {
			excf(m, "Fail writing target file '%v' in copy.", new)
		}
	}
}

// Create a temporary file named fname in 'dp' and returns its path.
// If dp == "" the file is created in the system temporary directory.
//    m: Virtual machine.
func prTmp(m *machine.T) {
	fname := popStr(m)
	dp := popStr(m)
	f, err := ioutil.TempFile(dp, fname)
	if err != nil {
		excf(m, "Fail creating temporary file '%v/%v'.", dp, fname)
	}
	pushStr(m, f.Name())
}

// Returns 'true' if 'p' is a regular file (nither directory, link, nor pipe).
//    m: Virtual machine.
func prIsRegular(m *machine.T) {
	p := popStr(m)
	f, err := os.Lstat(p)
	if err != nil {
		excf(m, "Fail reading file '%v' properties.", p)
	}
	m.Push(token.NewB(f.Mode().IsRegular(), m.MkPos()))
}

// Returns 'true' if 'p' is a link.
//    m: Virtual machine.
func prIsLink(m *machine.T) {
	p := popStr(m)
	f, err := os.Lstat(p)
	if err != nil {
		excf(m, "Fail reading file '%v' properties.", p)
	}
	m.Push(token.NewB(f.Mode()&os.ModeSymlink != 0, m.MkPos()))
}

// Returns the date of the last modification of 'p'.
//    m: Virtual machine.
func prModified(m *machine.T) {
	p := popStr(m)
	f, err := os.Lstat(p)
	if err != nil {
		excf(m, "Fail reading file '%v' properties.", p)
	}
	m.Push(token.NewN(operator.Date_, f.ModTime(), m.MkPos()))
}

// Returns the size of 'p'.
//    m: Virtual machine.
func prSize(m *machine.T) {
	p := popStr(m)
	f, err := os.Lstat(p)
	if err != nil {
		excf(m, "Fail reading file '%v' properties.", p)
	}
	m.Push(token.NewI(f.Size(), m.MkPos()))
}

// Writes 's' in 'p'. If 'p' exists, it is overwritten.
//    m: Virtual machine.
func prWrite(m *machine.T) {
	s := popStr(m)
	p := popStr(m)
	err := ioutil.WriteFile(p, []byte(s), 0755)
	if err != nil {
		excf(m, "Fail writing file '%v'.", p)
	}
}

// Appends 's' to 'p'.
//    m: Virtual machine.
func prAppend(m *machine.T) {
	s := popStr(m)
	p := popStr(m)
	f, err := os.OpenFile(p, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0755)
	if err != nil {
		excf(m, "Fail openning file '%v' to append.", p)
	}
	defer f.Close()
	if _, err := f.WriteString(s); err != nil {
		excf(m, "Fail appending file '%v'.", p)
	}
}

// Read 'p'.
//    m: Virtual machine.
func prRead(m *machine.T) {
	p := popStr(m)
	bs, err := ioutil.ReadFile(p)
	if err != nil {
		excf(m, "Fail reading file '%v'.", p)
	}
	pushStr(m, string(bs))
}

// Open 'p' to append.
//    m: Virtual machine.
func prAopen(m *machine.T) {
	p := popStr(m)
	f, err := os.OpenFile(p, os.O_APPEND|os.O_WRONLY, 0755)
	if err != nil {
		excf(m, "Fail openning file '%v' for appending.", p)
	}
	m.Push(token.NewN(operator.File_, f, m.MkPos()))
}

// Open 'p' to read.
//    m: Virtual machine.
func prRopen(m *machine.T) {
	p := popStr(m)
	f, err := os.Open(p)
	if err != nil {
		excf(m, "Fail openning file '%v' for reading.", p)
	}
	m.Push(token.NewN(operator.File_, f, m.MkPos()))
}

// Open 'p' to append.
//    m: Virtual machine.
func prWopen(m *machine.T) {
	p := popStr(m)
	f, err := os.Create(p)
	if err != nil {
		excf(m, "Fail openning file '%v' for writing.", p)
	}
	m.Push(token.NewN(operator.File_, f, m.MkPos()))
}

// Close 'p'.
//    m: Virtual machine.
func prClose(m *machine.T) {
	tk := m.PopT(token.Native)
	o, f, _ := tk.N()
	if o != operator.File_ {
		m.Failt("\n Expected: File object.\n  Actual  : '%v'.", o)
	}

	err := f.(*os.File).Close()
	if err != nil {
		excf(m, "Fail closing file '%v'.", tk.StringDraft())
	}
}

// Read next 'BufferSize' bytes of 'p'. If next data is less than 'BufferSize',
// it returns a blob with the rest. When all the file is read, returns a
// blob of size 0.
//    m: Virtual machine.
func prReadBin(m *machine.T) {
	tk1 := m.PopT(token.Native)
	o, f, _ := tk1.N()
	if o != operator.File_ {
		m.Failt("\n Expected: File object.\n  Actual  : '%v'.", o)
	}

	bs := make([]byte, BufferSize)
	sz, err := f.(*os.File).Read(bs)
	if err != nil && err != io.EOF {
		excf(m, "Fail reading '%v'.", tk1.StringDraft())
	}

	m.Push(token.NewN(operator.Blob_, bs[:sz], m.MkPos()))
}

// Returns a File and a Iterator over lines of 'p'.
// File is intendet to close it after using the Iterator.
// Lines are read without end of line.
//    m: Virtual machine.
func prLines(m *machine.T) {
	pos := m.MkPos()
	p := popStr(m)
	f, err := os.Open(p)
	if err != nil {
		excf(m, "Fail openning file '%v' for reading.", p)
	}
	scanner := bufio.NewScanner(f)
	i := it.New(func() (tk *token.T, more bool) {
		more = scanner.Scan()
		if more {
			tk = token.NewS(scanner.Text(), pos)
			return
		}
		if err := scanner.Err(); err != nil {
			excf(m, "Fail reading '%v' as text file.", p)
		}
		return
	})

	m.Push(token.NewN(operator.File_, f, m.MkPos()))
	m.Push(token.NewN(operator.Iterator_, i, m.MkPos()))
}

// Writes a Blob in 'f'
//    m: Virtual machine.
func prWriteBin(m *machine.T) {
	tk2 := m.PopT(token.Native)
	o, b, _ := tk2.N()
	if o != operator.Blob_ {
		m.Failt("\n Expected: Blob object.\n  Actual  : '%v'.", o)
	}
	tk1 := m.PopT(token.Native)
	o, f, _ := tk1.N()
	if o != operator.File_ {
		m.Failt("\n Expected: File object.\n  Actual  : '%v'.", o)
	}

	_, err := f.(*os.File).Write(b.([]byte))
	if err != nil {
		excf(m, "Fail writing blob in '%v'.", tk1.StringDraft())
	}
}

// Writes a String in 'f'
//    m: Virtual machine.
func prWriteText(m *machine.T) {
	s := popStr(m)
	tk1 := m.PopT(token.Native)
	o, f, _ := tk1.N()
	if o != operator.File_ {
		m.Failt("\n Expected: File object.\n  Actual  : '%v'.", o)
	}

	_, err := f.(*os.File).Write([]byte(s))
	if err != nil {
		excf(m, "Fail writing text in '%v'.", tk1.StringDraft())
	}
}

// Processes file procedures.
//    m   : Virtual machine.
//    proc: Procedure
func Proc(m *machine.T, proc symbol.T) {
	switch proc {
	case symbol.New("cwd"):
		prCwd(m)
	case symbol.New("cd"):
		prCd(m)
	case symbol.New("mkdir"):
		prMkdir(m)
	case symbol.New("dir"):
		prDir(m)
	case symbol.New("isDirectory"):
		prIsDirectory(m)
	// ----
	case symbol.New("exists"):
		prExists(m)
	case symbol.New("del"):
		prDel(m)
	case symbol.New("rename"):
		prRename(m)
	case symbol.New("link"):
		prLink(m)
	case symbol.New("copy"):
		prCopy(m)
	case symbol.New("tmp"):
		prTmp(m)
	// ----
	case symbol.New("isRegular"):
		prIsRegular(m)
	case symbol.New("isLink"):
		prIsLink(m)
	case symbol.New("modified"):
		prModified(m)
	case symbol.New("size"):
		prSize(m)
	// ----
	case symbol.New("write"):
		prWrite(m)
	case symbol.New("append"):
		prAppend(m)
	case symbol.New("read"):
		prRead(m)
	// ----
	case symbol.New("aopen"):
		prAopen(m)
	case symbol.New("ropen"):
		prRopen(m)
	case symbol.New("wopen"):
		prWopen(m)
	case symbol.New("close"):
		prClose(m)
	// ----
	case symbol.New("readBin"):
		prReadBin(m)
	case symbol.New("lines"):
		prLines(m)
	case symbol.New("writeBin"):
		prWriteBin(m)
	case symbol.New("writeText"):
		prWriteText(m)
	default:
		m.Failt("'file' does not contains the procedure '%v'.", proc.String())
	}
}

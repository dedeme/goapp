// Copyright 03-Mar-2022 ºDeme
// GNU General Public License - V3 <http://www.gnu.org/licenses/>

package runner

import (
	"bufio"
	"crypto/rand"
	"encoding/base64"
	"github.com/dedeme/kut/builtin/bfail"
	"github.com/dedeme/kut/builtin/bfunction"
	"github.com/dedeme/kut/expression"
	"github.com/dedeme/kut/runner/file"
	"io"
	"io/ioutil"
	"os"
	"os/user"
	"strings"
)

type fileT struct {
	f    *os.File
	scan *bufio.Scanner
}

// \s -> <file>
func fileAopen(args []*expression.T) (ex *expression.T, err error) {
	switch path := (args[0].Value).(type) {
	case string:
		f, er := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, 0755)
		if er != nil {
			err = bfail.Mk(er.Error())
			return
		}

		ex = expression.MkFinal(&fileT{f, nil})
	default:
		err = bfail.Type(args[0], "string")
	}
	return
}

// \s -> ()
func fileCd(args []*expression.T) (ex *expression.T, err error) {
	switch fpath := (args[0].Value).(type) {
	case string:
		err = os.Chdir(fpath)
	default:
		err = bfail.Type(args[0], "string")
	}
	return
}

// \<file> -> i
func fileClose(args []*expression.T) (ex *expression.T, err error) {
	switch f := (args[0].Value).(type) {
	case *fileT:
		err = f.f.Close()
	default:
		err = bfail.Type(args[0], "<file>")
	}
	return
}

// \<file> -> b
func fileCheckType(args []*expression.T) (ex *expression.T, err error) {
	switch (args[0].Value).(type) {
	case *fileT:
		ex = expression.MkFinal(true)
	default:
		ex = expression.MkFinal(false)
	}
	return
}

// \s, s -> ()
func fileCopy(args []*expression.T) (ex *expression.T, err error) {
	switch source := (args[0].Value).(type) {
	case string:
		switch target := (args[1].Value).(type) {
		case string:
			err = file.Copy(source, target)
		default:
			err = bfail.Type(args[1], "string")
		}
	default:
		err = bfail.Type(args[0], "string")
	}
	return
}

// \s -> ()
func fileDel(args []*expression.T) (ex *expression.T, err error) {
	switch fpath := (args[0].Value).(type) {
	case string:
		os.RemoveAll(fpath)
	default:
		err = bfail.Type(args[0], "string")
	}
	return
}

// \s -> [s...]
func fileDir(args []*expression.T) (ex *expression.T, err error) {
	switch dir := (args[0].Value).(type) {
	case string:
		fis, er := ioutil.ReadDir(dir)
		if er != nil {
			err = bfail.Mk(er.Error())
			return
		}
		var exs []*expression.T
		for _, fi := range fis {
			exs = append(exs, expression.MkFinal(fi.Name()))
		}
		ex = expression.MkFinal(exs)
	default:
		err = bfail.Type(args[0], "string")
	}
	return
}

// \s -> b
func fileExists(args []*expression.T) (ex *expression.T, err error) {
	switch path := (args[0].Value).(type) {
	case string:
		if _, er := os.Stat(path); er == nil {
			ex = expression.MkFinal(true)
		} else {
			ex = expression.MkFinal(false)
		}
	default:
		err = bfail.Type(args[0], "string")
	}
	return
}

// \-> b
func fileHome(args []*expression.T) (ex *expression.T, err error) {
	var u *user.User
	u, err = user.Current()
	if err == nil {
		ex = expression.MkFinal(u.HomeDir)
	}
	return
}

// \s -> b
func fileIsDirectory(args []*expression.T) (ex *expression.T, err error) {
	switch path := (args[0].Value).(type) {
	case string:
		if info, err := os.Stat(path); err == nil && info.IsDir() {
			ex = expression.MkFinal(true)
		} else {
			ex = expression.MkFinal(false)
		}
	default:
		err = bfail.Type(args[0], "string")
	}
	return
}

// \s -> ()
func fileMkdir(args []*expression.T) (ex *expression.T, err error) {
	switch path := (args[0].Value).(type) {
	case string:
		os.MkdirAll(path, os.FileMode(0755))
	default:
		err = bfail.Type(args[0], "string")
	}
	return
}

// \s -> s
func fileRead(args []*expression.T) (ex *expression.T, err error) {
	switch path := (args[0].Value).(type) {
	case string:
		var bs []byte
		bs, err = ioutil.ReadFile(path)
		if err == nil {
			ex = expression.MkFinal(string(bs))
		}
	default:
		err = bfail.Type(args[0], "string")
	}
	return
}

// \<file>, i -> [s]|[]
func fileReadBin(args []*expression.T) (ex *expression.T, err error) {
	switch f := (args[0].Value).(type) {
	case *fileT:
		switch lg := (args[1].Value).(type) {
		case int64:
			bs := make([]byte, lg)
			var n int
			n, err = f.f.Read(bs)
			if n == 0 && err == io.EOF {
				err = nil
			}
			if err == nil {
				bs2 := make([]byte, n)
				for i := 0; i < n; i++ {
					bs2[i] = bs[i]
				}
				ex = expression.MkFinal(bs2)
			}
		default:
			bfail.Type(args[1], "int")
		}
	default:
		err = bfail.Type(args[0], "<file>")
	}
	return
}

// \<file> -> [s]|[]
func fileReadLine(args []*expression.T) (ex *expression.T, err error) {
	switch f := (args[0].Value).(type) {
	case *fileT:
		if f.scan.Scan() {
			ex = expression.MkFinal(
				[]*expression.T{expression.MkFinal(f.scan.Text())})
		} else {
			ex = expression.MkFinal([]*expression.T{})
		}
	default:
		err = bfail.Type(args[0], "<file>")
	}
	return
}

// \s, s -> ()
func fileRename(args []*expression.T) (ex *expression.T, err error) {
	switch oldPath := (args[0].Value).(type) {
	case string:
		switch newPath := (args[1].Value).(type) {
		case string:
			err = os.Rename(oldPath, newPath)
		default:
			err = bfail.Type(args[1], "string")
		}
	default:
		err = bfail.Type(args[0], "string")
	}
	return
}

// \s -> <file>
func fileRopen(args []*expression.T) (ex *expression.T, err error) {
	switch path := (args[0].Value).(type) {
	case string:
		f, er := os.Open(path)
		if er != nil {
			err = bfail.Mk(er.Error())
			return
		}

		ex = expression.MkFinal(&fileT{f, bufio.NewScanner(f)})
	default:
		err = bfail.Type(args[0], "string")
	}
	return
}

// \s -> i
func fileSize(args []*expression.T) (ex *expression.T, err error) {
	switch path := (args[0].Value).(type) {
	case string:
		if info, er := os.Stat(path); err == nil {
			ex = expression.MkFinal(info.Size())
		} else {
			err = bfail.Mk(er.Error())
		}
	default:
		err = bfail.Type(args[0], "string")
	}
	return
}

// \s -> i
func fileTm(args []*expression.T) (ex *expression.T, err error) {
	switch path := (args[0].Value).(type) {
	case string:
		if info, er := os.Stat(path); err == nil {
			ex = expression.MkFinal(info.ModTime().UnixMilli())
		} else {
			err = bfail.Mk(er.Error())
		}
	default:
		err = bfail.Type(args[0], "string")
	}
	return
}

// \s, s -> s
func fileTmp(args []*expression.T) (ex *expression.T, err error) {
	switch dpath := (args[0].Value).(type) {
	case string:
		switch fpath := (args[1].Value).(type) {
		case string:
			for {
				lg := 8
				arr := make([]byte, lg)
				_, err = rand.Read(arr)

				if dpath == "" {
					dpath = "."
				}
				fpath = dpath + "/" + fpath +
					strings.ReplaceAll(base64.StdEncoding.EncodeToString(arr)[:lg], "/", "+")
				if _, er := os.Stat(fpath); er != nil {
					break
				}
			}
			ex = expression.MkFinal(fpath)
		default:
			err = bfail.Type(args[1], "string")
		}
	default:
		err = bfail.Type(args[0], "string")
	}
	return
}

// \-> b
func fileWd(args []*expression.T) (ex *expression.T, err error) {
	var d string
	d, err = os.Getwd()
	if err == nil {
		ex = expression.MkFinal(d)
	}
	return
}

// \s -> <file>
func fileWopen(args []*expression.T) (ex *expression.T, err error) {
	switch path := (args[0].Value).(type) {
	case string:
		f, er := os.Create(path)
		if er != nil {
			err = bfail.Mk(er.Error())
			return
		}

		ex = expression.MkFinal(&fileT{f, nil})
	default:
		err = bfail.Type(args[0], "string")
	}
	return
}

// \s, s -> ()
func fileWrite(args []*expression.T) (ex *expression.T, err error) {
	switch path := (args[0].Value).(type) {
	case string:
		switch text := (args[1].Value).(type) {
		case string:
			err = ioutil.WriteFile(path, []byte(text), 0755)
		default:
			err = bfail.Type(args[1], "string")
		}
	default:
		err = bfail.Type(args[0], "string")
	}
	return
}

// \<file>, s -> ()
func fileWriteText(args []*expression.T) (ex *expression.T, err error) {
	switch f := (args[0].Value).(type) {
	case *fileT:
		switch text := (args[1].Value).(type) {
		case string:
			_, err = f.f.WriteString(text)
		default:
			err = bfail.Type(args[1], "string")
		}
	default:
		err = bfail.Type(args[0], "<file>")
	}
	return
}

// \<file>, <bytes> -> ()
func fileWriteBin(args []*expression.T) (ex *expression.T, err error) {
	switch f := (args[0].Value).(type) {
	case *fileT:
		switch bs := (args[1].Value).(type) {
		case []byte:
			_, err = f.f.Write(bs)
		default:
			err = bfail.Type(args[1], "string")
		}
	default:
		err = bfail.Type(args[0], "<file>")
	}
	return
}

func fileGet(fname string) (fn *bfunction.T, ok bool) {
	ok = true
	switch fname {
	case "aopen":
		fn = bfunction.New(1, fileAopen)
	case "cd":
		fn = bfunction.New(1, fileCd)
	case "checkType":
		fn = bfunction.New(1, fileCheckType)
	case "close":
		fn = bfunction.New(1, fileClose)
	case "copy":
		fn = bfunction.New(2, fileCopy)
	case "del":
		fn = bfunction.New(1, fileDel)
	case "dir":
		fn = bfunction.New(1, fileDir)
	case "exists":
		fn = bfunction.New(1, fileExists)
	case "home":
		fn = bfunction.New(0, fileHome)
	case "isDirectory":
		fn = bfunction.New(1, fileIsDirectory)
	case "mkdir":
		fn = bfunction.New(1, fileMkdir)
	case "read":
		fn = bfunction.New(1, fileRead)
	case "readBin":
		fn = bfunction.New(2, fileReadBin)
	case "readLine":
		fn = bfunction.New(1, fileReadLine)
	case "rename":
		fn = bfunction.New(2, fileRename)
	case "ropen":
		fn = bfunction.New(1, fileRopen)
	case "size":
		fn = bfunction.New(1, fileSize)
	case "tm":
		fn = bfunction.New(1, fileTm)
	case "tmp":
		fn = bfunction.New(2, fileTmp)
	case "wd":
		fn = bfunction.New(0, fileWd)
	case "wopen":
		fn = bfunction.New(1, fileWopen)
	case "write":
		fn = bfunction.New(2, fileWrite)
	case "writeBin":
		fn = bfunction.New(2, fileWriteBin)
	case "writeText":
		fn = bfunction.New(2, fileWriteText)
	default:
		ok = false
	}

	return
}

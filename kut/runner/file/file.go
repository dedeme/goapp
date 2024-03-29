// Copyright 31-Aug-2017 ºDeme
// GNU General Public License - V3 <http://www.gnu.org/licenses/>

// Utilities to easy file management
package file

import (
	"archive/zip"
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/user"
	"path"
	"path/filepath"
	"strings"
)

// UserDir returns the name of user dir
func HomeDir() string {
	u, _ := user.Current()
	return u.HomeDir
}

// Exists returns true if path actually exists in file system
func Exists(path string) bool {
	if _, err := os.Stat(path); err == nil {
		return true
	}
	return false
}

// Is directory return true if path exists and is a directory
func IsDirectory(path string) bool {
	if info, err := os.Stat(path); err == nil && info.IsDir() {
		return true
	}
	return false
}

func LastModification(path string) int64 {
	info, err := os.Stat(path)
	if err != nil {
		panic(err)
	}
	return info.ModTime().Unix()
}

func Size(path string) int64 {
	info, err := os.Stat(path)
	if err != nil {
		panic(err)
	}
	return info.Size()
}

// Mkdir makes a directory
func Mkdir(f string) {
	os.Mkdir(f, os.FileMode(0755))
}

// Mkdirs makes a directory and its parents
func Mkdirs(f string) {
	os.MkdirAll(f, os.FileMode(0755))
}

// TempDir makes a directorio in the temporal directory.
func TempDir(prefix string) string {
	name, err := ioutil.TempDir("", prefix)
	if err != nil {
		panic(err)
	}
	return name
}

// List return the list of files of a directory
func List(dir string) []os.FileInfo {
	fis, err := ioutil.ReadDir(dir)
	if err != nil {
		panic(err)
	}
	return fis
}

// TempFile creates a file in 'dir'. If 'dir' is "" file is created in the
// temporal directory.
func TempFile(dir string, prefix string) *os.File {
	f, err := ioutil.TempFile(dir, prefix)
	if err != nil {
		panic(err)
	}
	return f
}

// Rename changes the name of a file or directory
func Rename(oldname, newname string) {
	err := os.Rename(oldname, newname)
	if err != nil {
		panic(err)
	}
}

// Remove removes path and all its subdirectories.
func Remove(path string) error {
	return os.RemoveAll(path)
}

// OpenRead opens path for read.
func OpenRead(path string) *os.File {
	f, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	return f
}

// Read reads a file completely. (File is open and closed)
func ReadAllBin(path string) []byte {
	bs, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	return bs
}

// ReadAll reads a data file completely. (File is open and closed)
func ReadAll(path string) string {
	return string(ReadAllBin(path))
}

// Lines are read without end of line.
// If 'f' returns 'true', reading is stoped.
func Lines(path string, f func(s string) bool) {
	file := OpenRead(path)
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if f(scanner.Text()) {
			break
		}
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}
}

// OpenRead opens path for write.
func OpenWrite(path string) *os.File {
	f, err := os.Create(path)
	if err != nil {
		panic(err)
	}
	return f
}

// OpenRead opens path for append.
func OpenAppend(path string) *os.File {
	f, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, 0755)
	if err != nil {
		panic(err)
	}
	return f
}

// WriteAll writes data overwriting 'file'. (File is open and closed)
func WriteAllBin(path string, data []byte) {
	err := ioutil.WriteFile(path, data, 0755)
	if err != nil {
		panic(err)
	}
}

// WriteAll writes a text overwriting 'file'. (File is open and closed)
func WriteAll(path, text string) {
	WriteAllBin(path, []byte(text))
}

// Write  writes a text in 'file'
func Write(file *os.File, text string) {
	_, err := file.WriteString(text)
	if err != nil {
		panic(err)
	}
}

// Write  writes binary data in 'file'
func WriteBin(file *os.File, data []byte) {
	_, err := file.Write(data)
	if err != nil {
		panic(err)
	}
}

//  Copy a file.
//    source: Can be a regular file or a directory.
//    target: - If source is a directory, target must be the target parent
//              directory.
//            - If source is a regular file, target can be the target parent
//              directory or a regular file.
//  NOTE: Target will be overwritte if already exists.
func Copy(source, target string) (err error) {
	defer func() {
		if e := recover(); e != nil {
			switch e.(type) {
			case error:
				err = e.(error)
			default:
				err = fmt.Errorf("Fail copying '%v' to '%v'", source, target)
			}
		}
	}()

	if IsDirectory(source) {
		p := path.Join(target, path.Base(source))
		if Exists(p) {
			if !IsDirectory(p) {
				err = fmt.Errorf(
					"Copying '%v' to '%v', when the later exists and is not a directory",
					source, p,
				)
				return
			}
			Remove(p)
		}
		Mkdirs(p)

		for _, e := range List(source) {
			if err = Copy(path.Join(source, e.Name()), p); err != nil {
				return
			}
		}
		return
	}

	if IsDirectory(target) {
		target = path.Join(target, path.Base(source))
	}

	sourcef, err := os.Open(source)
	if err != nil {
		return
	}

	targetf, err := os.Create(target)
	if err != nil {
		return
	}

	defer sourcef.Close()
	defer targetf.Close()

	buf := make([]byte, 8192)
	for {
		var n int
		n, err = sourcef.Read(buf)
		if err != nil && err != io.EOF {
			return
		}
		if n == 0 {
			err = nil
			break
		}

		if _, err = targetf.Write(buf[:n]); err != nil {
			return
		}
	}

	return
}

// Zip a file or directory
func Zip(source, target string) error {
	zipfile, err := os.Create(target)
	if err != nil {
		return err
	}
	defer zipfile.Close()

	archive := zip.NewWriter(zipfile)
	defer archive.Close()

	info, err := os.Stat(source)
	if err != nil {
		return err
	}

	var baseDir string
	if info.IsDir() {
		baseDir = filepath.Base(source)
	}

	filepath.Walk(source, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}

		if baseDir != "" {
			header.Name = filepath.Join(baseDir, strings.TrimPrefix(path, source))
		}

		if info.IsDir() {
			header.Name += "/"
		} else {
			header.Method = zip.Deflate
		}

		writer, err := archive.CreateHeader(header)
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()
		_, err = io.Copy(writer, file)
		return err
	})

	return err
}

// Unzip a file or directory
func Unzip(archive, target string) error {
	reader, err := zip.OpenReader(archive)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(target, 0755); err != nil {
		return err
	}

	for _, file := range reader.File {
		path := filepath.Join(target, file.Name)
		if file.FileInfo().IsDir() {
			os.MkdirAll(path, file.Mode())
			continue
		}

		fileReader, err := file.Open()
		if err != nil {
			return err
		}
		defer fileReader.Close()

		targetFile, err := os.OpenFile(
			path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
		if err != nil {
			return err
		}
		defer targetFile.Close()

		if _, err := io.Copy(targetFile, fileReader); err != nil {
			return err
		}
	}

	return nil
}

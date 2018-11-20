package fileutils

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"

	"io/ioutil"
)

func TestReadLines(t *testing.T) {
	testFile := "./testfile"
	defer func() {
		os.Remove(testFile)
	}()

	var testLines []string
	for i := 0; i < 10; i++ {
		testLines = append(testLines, fmt.Sprintf("Hello %d", i))
	}
	data := strings.Join(testLines, "\n")
	ioutil.WriteFile(testFile, []byte(data), 0666)

	lines, err := ReadLines(testFile)
	if err != nil || !reflect.DeepEqual(lines, testLines) {
		t.Fatalf("ReadLines failed err:%v lines:%v", err, lines)
	}

	// File which does not exist should return nil lines and valid error
	if lines, err = ReadLines("abcd"); err == nil || lines != nil {
		t.Fatalf("ReadLines failed to bad file err:%v lines:%v", err, lines)
	}
}

func TestFileExists(t *testing.T) {
	tests := []struct {
		name   string
		exists bool
		err    error
	}{
		{"file_test.go", true, nil},
		{"audio", false, nil},
		//TODO simulate failure case
	}

	for _, test := range tests {
		if ok, err := Exists(test.name); ok != test.exists || err != test.err {
			t.Fatalf("Test:%v failed ok:%t err:%v", test, ok, err)
		}
	}
}

func TestReadDirByName(t *testing.T) {
	testDir := "/tmp/_filetest"
	os.RemoveAll(testDir)
	os.MkdirAll(testDir, 0777)

	files := []string{"9", "3", "6"}
	for _, f := range files {
		if err := ioutil.WriteFile(filepath.Join(testDir, f), []byte("0"), 666); err != nil {
			t.Fatalf("Failed to create test file:%s err:%v", f, err)
		}
	}

	var (
		fis []os.FileInfo
		err error
	)
	if fis, err = ReadDir(testDir, SortByName); err != nil {
		t.Fatalf("ReadDirByName failed err:%v", err)
	}
	if !verifyFiles([]string{"3", "6", "9"}, fis) {
		t.Fatalf("ReadDirByName incorrect output fis:%+v", fis)
	}
}

func TestReadDirByTime(t *testing.T) {
	testDir := "/tmp/_filetest"
	os.RemoveAll(testDir)
	os.MkdirAll(testDir, 0777)

	files := []string{"22", "1", "2"}
	for _, f := range files {
		if err := ioutil.WriteFile(filepath.Join(testDir, f), []byte("0"), 666); err != nil {
			t.Fatalf("Failed to create test file:%s err:%v", f, err)
		}
		time.Sleep(500 * time.Millisecond)
	}

	var (
		fis []os.FileInfo
		err error
	)
	if fis, err = ReadDir(testDir, SortByMTime); err != nil {
		t.Fatalf("ReadDirByTime failed err:%v", err)
	}
	if !verifyFiles([]string{"22", "1", "2"}, fis) {
		t.Fatalf("ReadDirByTime incorrect output fis:%+v", fis)
	}
}

func TestGlob(t *testing.T) {
	testDir := "/tmp/_filetest"
	os.RemoveAll(testDir)
	os.MkdirAll(testDir, 0777)

	files := []string{"c.logg", "b.log", "a.log"}
	for _, f := range files {
		if err := ioutil.WriteFile(filepath.Join(testDir, f), []byte("0"), 666); err != nil {
			t.Fatalf("Failed to create test file:%s err:%v", f, err)
		}
	}

	fileInfos, err := Glob(testDir + "/*log")
	if err != nil {
		t.Fatalf("Glob failed err:%v", err)
	}
	Sort(fileInfos, SortByName)
	verifyFiles([]string{"a.log", "b.log"}, fileInfos)
}

func verifyFiles(files []string, fis []os.FileInfo) bool {
	for i, file := range files {
		if fis[i].Name() != file {
			return false
		}
	}

	return true
}

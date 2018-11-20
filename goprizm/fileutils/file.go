package fileutils

import (
	"bufio"
	"encoding/json"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
)

// ReadLines reads all lines from given file in a string slice.
func ReadLines(name string) ([]string, error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var lines []string
	scan := bufio.NewScanner(f)
	for scan.Scan() {
		lines = append(lines, scan.Text())
	}

	if err := scan.Err(); err != nil {
		return nil, err
	}

	return lines, nil
}

// Read the file contents as bytes from the given file
func ReadBytes(name string) ([]byte, error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, err
	}

	defer f.Close()

	return ioutil.ReadAll(f)
}

// Exists return true if file is present.
func Exists(name string) (bool, error) {
	_, err := os.Stat(name)
	if err == nil {
		return true, nil
	}

	if os.IsNotExist(err) {
		return false, nil
	}

	return false, err
}

// ReadJSON reads given file content and json unmarshal it as given object.
func ReadJSON(file string, obj interface{}) error {
	bytes, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(bytes, obj); err != nil {
		return fmt.Errorf("%s(%v)", file, err)
	}

	return nil
}

// SaveJSON save the given object to file
func SaveJSON(file string, obj interface{}) error {
	data, err := json.MarshalIndent(obj, "", "    ")
	if err != nil {
		return nil
	}

	return ioutil.WriteFile(file, data, 0644)
}

// ReadYAML reads given file content and yaml unmarhal it as given object
func ReadYAML(file string, obj interface{}) error {
	bytes, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(bytes, obj)
	if err != nil {
		return fmt.Errorf("%s(%v)", file, err)
	}

	return nil
}

// SaveYAML save the given object to file
func SaveYAML(file string, obj interface{}) error {
	data, err := yaml.Marshal(obj)
	if err != nil {
		return nil
	}

	return ioutil.WriteFile(file, data, 0644)
}

// FileCmpFunc compare 2 FileInfos
type FileCmpFunc func(os.FileInfo, os.FileInfo) bool

// ReadDir lists files in a directory ordered by given compare func.
func ReadDir(dir string, cmp FileCmpFunc) (files []os.FileInfo, err error) {
	if files, err = ioutil.ReadDir(dir); err != nil {
		return
	}

	Sort(files, cmp)
	return files, nil
}

// Sort sorts FileInfos using comparison func provided
func Sort(files []os.FileInfo, cmp FileCmpFunc) {
	sort.Sort(fileInfos{files, cmp})
}

//SortByName is the FileCmpFunc used to sort FileInfos by file name
func SortByName(f0, f1 os.FileInfo) bool {
	return f0.Name() < f1.Name()
}

//SortByMTime is the FileCmpFunc used to sort FileInfos by file modified timestamp
func SortByMTime(f0, f1 os.FileInfo) bool {
	return f0.ModTime().Before(f1.ModTime())
}

// fileInfos implements sort interface using given cmpFunc
type fileInfos struct {
	files   []os.FileInfo // list of files
	cmpFunc FileCmpFunc   // custom comparison func
}

func (fis fileInfos) Len() int {
	return len(fis.files)
}

func (fis fileInfos) Swap(i, j int) {
	fis.files[i], fis.files[j] = fis.files[j], fis.files[i]
}

func (fis fileInfos) Less(i, j int) bool {
	return fis.cmpFunc(fis.files[i], fis.files[j])
}

func Glob(pattern string) ([]os.FileInfo, error) {
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return nil, err
	}

	fileInfos := make([]os.FileInfo, len(matches))
	for i, file := range matches {
		fi, err := os.Stat(file)
		if err != nil {
			return nil, err
		}

		fileInfos[i] = fi
	}

	return fileInfos, nil
}

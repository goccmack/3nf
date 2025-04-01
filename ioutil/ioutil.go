/*
Package ioutil contains functions for writing directories and files.
*/
package ioutil

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// FilePermission given to all non-exectable files
const filePermission = 0644

// File permission given to all created directories and executable files
const exePermission = 0755

// Exist returns true if path exists, otherwise false.
func Exist(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func CaseInsensitiveGetFileName(path string) (string, error) {
	base := filepath.Base(path)
	dir := filepath.Dir(path)
	realBase := ""
	files, err := os.ReadDir(dir)
	if err != nil {
		return "", err
	}
	for _, file := range files {
		if strings.EqualFold(base, file.Name()) {
			realBase = file.Name()
			break
		}
	}
	if realBase == "" {
		return "", os.ErrNotExist
	}
	return filepath.Join(dir, realBase), nil
}

// MkdirAll makes all the directories in path.
func MkdirAll(path string) error {
	if path == "" {
		return nil
	}
	return os.MkdirAll(path, exePermission)
}

// WriteFile creates all the non-existent directories in path before writing
// data to a non-executable file, path.
func WriteFile(path string, data []byte) error {
	dir, _ := filepath.Split(path)
	if err := MkdirAll(dir); err != nil {
		return fmt.Errorf("error creating directory %s: %s", dir, err)
	}
	if err := os.WriteFile(path, data, filePermission); err != nil {
		return fmt.Errorf("error writing file %s: %s", path, err)
	}
	return nil
}

// WriteExeFile creates all the non-existent directories in path before writing
// data to an executable file, path.
func WriteExeFile(path string, data []byte) error {
	dir, _ := filepath.Split(path)
	if err := MkdirAll(dir); err != nil {
		return fmt.Errorf("error creating directory %s: %s", dir, err)
	}
	if err := os.WriteFile(path, data, exePermission); err != nil {
		return fmt.Errorf("error writing file %s: %s", path, err)
	}
	return nil
}

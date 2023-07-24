//go:build windows

package filehandler

import (
	"path/filepath"
	"syscall"
)

func isHiddenFile(filename string, directory string) (bool, error) {
	filePath := filepath.Join(directory, filename)
	pointer, err := syscall.UTF16PtrFromString(filePath)
	if err != nil {
		return false, err
	}
	attributes, err := syscall.GetFileAttributes(pointer)
	if err != nil {
		return false, err
	}
	return attributes&syscall.FILE_ATTRIBUTE_HIDDEN != 0, nil
}

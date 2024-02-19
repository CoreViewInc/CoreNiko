package io

import (
	"os"
	"path/filepath"
	"io"
)


func CopyDir(src string, dest string) error {
	// Get properties of source directory
	srcInfo, err := os.Stat(src)
	if err != nil {
		return err
	}

	// Create the destination directory
	err = os.MkdirAll(dest, srcInfo.Mode())
	if err != nil {
		return err
	}

	directory, _ := os.Open(src)
	objects, err := directory.Readdir(-1)

	for _, obj := range objects {
		srcFilePointer := filepath.Join(src, obj.Name())
		destFilePointer := filepath.Join(dest, obj.Name())

		if obj.IsDir() {
			// Create sub-directories, recursively.
			err = CopyDir(srcFilePointer, destFilePointer)
			if err != nil {
				return err
			}
		} else {
			// Perform the file copy.
			err = CopyFile(srcFilePointer, destFilePointer)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func CopyFile(src, dst string) error {
	var err error
	var srcfd *os.File
	var dstfd *os.File
	var srcinfo os.FileInfo

	if srcfd, err = os.Open(src); err != nil {
		return err
	}
	defer srcfd.Close()

	if dstfd, err = os.Create(dst); err != nil {
		return err
	}
	defer dstfd.Close()

	if _, err = io.Copy(dstfd, srcfd); err != nil {
		return err
	}
	if srcinfo, err = os.Stat(src); err != nil {
		return err
	}
	return os.Chmod(dst, srcinfo.Mode())
}
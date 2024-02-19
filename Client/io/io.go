package io

import (
	"io"
	"os"
	"path/filepath"
)

func CopyDir(src string, dest string) error {
	// Get properties of source directory
	srcInfo, err := os.Stat(src)
	if err != nil {
		return err
	}

	// Create the destination directory with the same permissions as the source
	err = os.MkdirAll(dest, srcInfo.Mode())
	if err != nil {
		return err
	}

	entries, err := os.ReadDir(src)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		destPath := filepath.Join(dest, entry.Name())

		if entry.IsDir() {
			// Recursively copy sub-directories
			err = CopyDir(srcPath, destPath)
			if err != nil {
				return err
			}
		} else {
			// Copy the file
			err = CopyFile(srcPath, destPath)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func CopyFile(src, dest string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	// Get the file information
	srcInfo, err := srcFile.Stat()
	if err != nil {
		return err
	}

	// Ensure the destination directory exists
	destDir := filepath.Dir(dest)
	err = os.MkdirAll(destDir, srcInfo.Mode().Perm())
	if err != nil {
		return err
	}

	// Create the destination file
	destFile, err := os.OpenFile(dest, os.O_RDWR|os.O_CREATE|os.O_TRUNC, srcInfo.Mode())
	if err != nil {
		return err
	}
	defer destFile.Close()

	// Copy the file content
	_, err = io.Copy(destFile, srcFile)
	return err
}
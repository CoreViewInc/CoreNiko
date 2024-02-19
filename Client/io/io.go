package io

import (
	"io"
	"os"
	"path/filepath"
)

// CopyDir recursively copies a directory tree, attempting to preserve permissions.
// The dest parameter should be the base destination directory where src structure will start.
func CopyDir(src string, dest string) error {
	// Define the source base path for path concatenation
	srcBase := filepath.Dir(src)

	// Get properties of source directory
	srcInfo, err := os.Stat(src)
	if err != nil {
		return err
	}

	// Calculate the destination directory path
	relPath, err := filepath.Rel(srcBase, src)
	if err != nil {
		return err
	}
	destPath := filepath.Join(dest, relPath)

	// Create the destination directory with the same permissions as the source
	err = os.MkdirAll(destPath, srcInfo.Mode())
	if err != nil {
		return err
	}

	entries, err := os.ReadDir(src)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())

		if entry.IsDir() {
			// Recursively copy sub-directories
			err = CopyDir(srcPath, dest)
			if err != nil {
				return err
			}
		} else {
			// Copy the file
			err = CopyFile(srcPath, dest)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// CopyFile copies a single file from src to dest, including the full src directory structure.
func CopyFile(src string, dest string) error {
	// Open the source file
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

	// Ensure the src is not a directory
	if srcInfo.IsDir() {
		return &os.PathError{Op: "copy", Path: src, Err: os.ErrInvalid}
	}

	// Find the root directory of the src to determine the relative path
	srcRootPath := filepath.VolumeName(src)
	if srcRootPath == "" { // On UNIX-like systems, the root is '/'
		srcRootPath = "/"
	} else { // On Windows, it may be something like 'C:\'
		srcRootPath = srcRootPath + "\\"
	}
	
	// Calculate the relative path
	relPath, err := filepath.Rel(srcRootPath, src)
	if err != nil {
		return err
	}

	// Create the full destination path by appending the relative path to dest
	destPath := filepath.Join(dest, relPath)

	// Ensure the destination directory exists
	destDir := filepath.Dir(destPath)
	err = os.MkdirAll(destDir, srcInfo.Mode().Perm())
	if err != nil {
		return err
	}

	// Create the destination file
	destFile, err := os.Create(destPath)
	if err != nil {
		return err
	}
	defer destFile.Close()

	// Copy the file content
	_, err = io.Copy(destFile, srcFile)
	return err
}
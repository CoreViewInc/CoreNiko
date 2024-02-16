package io

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

)

type IO struct {
}

// New creates a new FileCopier.
func New() *IO {
	return &IO{}
}

// CopyFileToZip copies a single file to a zip writer.
func (fc *IO) CopyFileToZip(src string, zw *zip.Writer, basePath string) error {
	fileToZip, err := os.Open(src)
	if err != nil {
		return err
	}
	defer fileToZip.Close()

	// Get the basename of the source file for the zip
	relativePath, err := filepath.Rel(basePath, src)
	if err != nil {
		return err
	}
	zipFileWriter, err := zw.Create(relativePath)
	if err != nil {
		return err
	}

	_, err = io.Copy(zipFileWriter, fileToZip)
	return err
}

// CopyDirToZip copies a directory and its sub-directories into a zip file.
func (fc *IO) CopyDirToZip(src string, dstZipPath string) error {
	// Create a new zip archive.
	zipFile, err := os.Create(dstZipPath)
	if err != nil {
		return err
	}
	defer zipFile.Close()

	zw := zip.NewWriter(zipFile)
	defer zw.Close()

	// Walk the directory tree
	err = filepath.Walk(src, func(path string, info os.FileInfo, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		// Check if it’s a directory, if so, skip it.
		if info.IsDir() {
			return nil
		}
		// Detect if file is regular. Non-regular files (like symbolic links) are skipped.
		if !info.Mode().IsRegular() {
			fmt.Printf("Non-regular file skipped: %s\n", path)
			return nil
		}

		return fc.CopyFileToZip(path, zw, src)
	})

	return err
}

// Unzip extracts a zip file's contents into a specified destination directory.
func (fc *IO) Unzip(srcZipPath, dstDir string, ignoredPaths ...string) error {
	r, err := zip.OpenReader(srcZipPath)
	if err != nil {
		return err
	}
	defer r.Close()

	ignoredSet := make(map[string]struct{})
	for _, ignoredPath := range ignoredPaths {
		ignoredSet[filepath.Clean(ignoredPath)] = struct{}{}
	}

	for _, file := range r.File {
		fPath := filepath.Join(dstDir, file.Name)
		relativePath, err := filepath.Rel(dstDir, fPath)
		if err != nil {
			return err
		}

		// Ignore specified paths
		if _, found := ignoredSet[relativePath]; found {
			continue
		}

		// Check for ZipSlip (Directory traversal)
		if !strings.HasPrefix(fPath, filepath.Clean(dstDir)+string(os.PathSeparator)) {
			return fmt.Errorf("invalid file path: %s", fPath)
		}

		if file.FileInfo().IsDir() {
			if err = os.MkdirAll(fPath, file.Mode()); err != nil {
				return err
			}
			continue
		}

		if err = os.MkdirAll(filepath.Dir(fPath), os.ModePerm); err != nil {
			return err
		}

		fmt.Println("extracting :"+fPath)
		outFile, err := os.OpenFile(fPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
		if err != nil {
			return err
		}

		rc, err := file.Open()
		if err != nil {
			outFile.Close()
			return err
		}

		_, err = io.Copy(outFile, rc)

		closeErr := outFile.Close()
		rc.Close()

		if err != nil {
			return err
		}

		if closeErr != nil {
			return closeErr
		}
	}

	return nil
}
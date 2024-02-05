package compress

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"github.com/spf13/afero"
	"io"
	"os"
	"path/filepath"
)

type Compressor struct {
	Fs        afero.Fs
	SourceDir string
}

func New(fs afero.Fs, sourceDir string) *Compressor {
	return &Compressor{Fs: fs, SourceDir: sourceDir}
}

func (c *Compressor) Compress(outputPath string) error {
	outFile, err := c.Fs.Create(outputPath)
	if err != nil {
		return fmt.Errorf("error creating archive file: %v", err)
	}
	defer outFile.Close()

	gzipWriter := gzip.NewWriter(outFile)
	defer gzipWriter.Close()

	tarWriter := tar.NewWriter(gzipWriter)
	defer tarWriter.Close()

	err = afero.Walk(c.Fs, c.SourceDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if path == c.SourceDir {
			return nil
		}

		header, err := tar.FileInfoHeader(info, "")
		if err != nil {
			return err
		}

		header.Name, err = filepath.Rel(c.SourceDir, path)
		if err != nil {
			return err
		}

		if err := tarWriter.WriteHeader(header); err != nil {
			return err
		}

		if !info.IsDir() {
			file, err := c.Fs.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()

			if _, err := io.Copy(tarWriter, file); err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return err
	}

	if err := tarWriter.Close(); err != nil {
		return err
	}

	if err := gzipWriter.Close(); err != nil {
		return err
	}

	return nil
}
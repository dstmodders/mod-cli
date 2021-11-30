package workshop

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/mattn/go-zglob"
)

// Controller is the interface that wraps the Workshop methods.
type Controller interface {
	IsPathIgnored(string) bool
	GetFiles() ([]string, int64, error)
	DestDirExists() bool
	MakeDestDir() error
	MakeDestFile(string) error
	CopyFiles() error
	CountDestItems() (int, error)
	Files() []string
	FilesSize() int64
	RelSrcPath() string
	AbsSrcPath() string
	RelDestPath() string
	AbsDestPath() string
}

// Workshop represents the workshop-related data.
type Workshop struct {
	// Ignore holds a list of files to ignore.
	Ignore []string

	files       []string
	filesSize   int64
	relSrcPath  string
	absSrcPath  string
	relDestPath string
	absDestPath string
	destDirName string
}

// New creates a new Workshop instance.
func New(src, dest string) (*Workshop, error) {
	relSrcPath, err := filepath.Rel(src, src)
	if err != nil {
		return nil, err
	}

	relDestPath, err := filepath.Rel(src, dest)
	if err != nil {
		return nil, err
	}

	absSrcPath, err := filepath.Abs(src)
	if err != nil {
		return nil, err
	}

	absDestPath, err := filepath.Abs(dest)
	if err != nil {
		return nil, err
	}

	destDirName := filepath.Base(relDestPath)

	return &Workshop{
		Ignore: []string{
			".*",
			"Makefile",
			"codecov.yml",
			"config.ld",
			"lcov.info",
			"luacov.*",
			"spec/",
		},
		absDestPath: absDestPath,
		absSrcPath:  absSrcPath,
		destDirName: destDirName,
		relDestPath: relDestPath,
		relSrcPath:  relSrcPath,
	}, nil
}

// IsPathIgnored checks if the provided path is ignored based on Ignore.
func (w *Workshop) IsPathIgnored(path string) bool {
	matched, _ := zglob.Match(w.destDirName+"**/*", path)
	if matched {
		return true
	}

	for _, ignore := range w.Ignore {
		ignore = strings.TrimSuffix(ignore, "/")
		matched, _ = zglob.Match("**/"+ignore+"**/*", path)
		if matched {
			return true
		}
	}
	return false
}

// GetFiles gets a list of files and their size in total from the source path
// based on Ignore.
func (w *Workshop) GetFiles() (files []string, size int64, err error) {
	if err := filepath.Walk(w.relSrcPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() || w.IsPathIgnored(path) {
			return nil
		}

		files = append(files, path)
		size += info.Size()
		return nil
	}); err != nil {
		return files, size, err
	}

	w.files = files
	w.filesSize = size

	return files, size, err
}

// DestDirExists checks if destination directory exists.
func (w *Workshop) DestDirExists() bool {
	stat, _ := os.Stat(w.absDestPath)
	return stat != nil && stat.IsDir()
}

// MakeDestDir makes an empty destination directory.
func (w *Workshop) MakeDestDir() error {
	return os.MkdirAll(w.absDestPath, os.ModePerm)
}

// MakeDestFile makes an empty destination file.
func (w *Workshop) MakeDestFile(name string) error {
	return os.MkdirAll(filepath.Dir(filepath.Join(w.absDestPath, name)), os.ModePerm)
}

// CopyFiles copies all files retrieved earlier using GetFiles to the
// destination path.
func (w *Workshop) CopyFiles() error {
	if len(w.files) == 0 {
		return errors.New("no files to copy")
	}

	for _, file := range w.files {
		stat, err := os.Stat(file)
		if err != nil {
			return err
		}

		if !stat.Mode().IsRegular() {
			return fmt.Errorf("%s is not a regular file", stat.Name())
		}

		src, err := os.Open(file)
		if err != nil {
			_ = src.Close()
			return err
		}

		if err := w.MakeDestFile(src.Name()); err != nil {
			return err
		}

		dest, err := os.Create(filepath.Join(w.relDestPath, src.Name()))
		if err != nil {
			_ = src.Close()
			return err
		}

		_, err = io.Copy(dest, src)
		if err != nil {
			_ = src.Close()
			_ = dest.Close()
			return err
		}

		_ = src.Close()
		_ = dest.Close()
	}

	return nil
}

// CountDestItems counts the total number of items within a destination
// directory.
func (w *Workshop) CountDestItems() (int, error) {
	if !w.DestDirExists() {
		return 0, nil
	}

	files, err := os.ReadDir(w.absDestPath)
	if err != nil {
		return 0, err
	}

	return len(files), nil
}

// Files gets a list of files retrieved earlier using GetFiles.
func (w *Workshop) Files() []string {
	return w.files
}

// FilesSize gets the total size of all files retrieved earlier using GetFiles.
func (w *Workshop) FilesSize() int64 {
	return w.filesSize
}

// RelSrcPath gets a relative source path.
func (w *Workshop) RelSrcPath() string {
	return w.relSrcPath
}

// AbsSrcPath gets an absolute source path.
func (w *Workshop) AbsSrcPath() string {
	return w.absSrcPath
}

// RelDestPath gets a relative destination path.
func (w *Workshop) RelDestPath() string {
	return w.relDestPath
}

// AbsDestPath gets an absolute destination path.
func (w *Workshop) AbsDestPath() string {
	return w.absDestPath
}

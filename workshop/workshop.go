// Package workshop has been designed to prepare a mod directory or archive for
// Steam Workshop. It allows including only the essential files based on ignore
// list. In the future, it may also include the features to automatically
// publish your mod.
package workshop

import (
	"archive/zip"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/dstmodders/mod-cli/dir"
)

// Controller is the interface that wraps the Workshop methods.
type Controller interface {
	SetIgnore([]string)
	IsPathIgnored(string) bool
	GetFiles() ([]string, int64, error)
	DestDirExists() bool
	MakeDestDir() error
	MakeDestFile(string) error
	CopyFiles() error
	ZipFiles() error
	CountDestItems() (int, error)
	Files() []string
	FilesSize() int64
	RelSrcPath() string
	AbsSrcPath() string
	RelDestPath() string
	AbsDestPath() string
	PrintFiles()
}

// Workshop represents a workshop-related data.
type Workshop struct {
	files       []string
	filesSize   int64
	srcDir      dir.Dir
	relDestPath string
	absDestPath string
	destDirName string
}

// New creates a new Workshop instance.
func New(src, dest string) (*Workshop, error) {
	srcDir, err := dir.New(src)
	if err != nil {
		return nil, err
	}

	relDestPath, err := filepath.Rel(src, dest)
	if err != nil {
		return nil, err
	}

	absDestPath, err := filepath.Abs(dest)
	if err != nil {
		return nil, err
	}

	destDirName := filepath.Base(relDestPath)

	return &Workshop{
		srcDir:      *srcDir,
		absDestPath: absDestPath,
		destDirName: destDirName,
		relDestPath: relDestPath,
	}, nil
}

// SetIgnore sets ignore list.
func (w *Workshop) SetIgnore(ignore []string) {
	w.srcDir.SetIgnore(ignore)
}

// IsPathIgnored checks if the provided path is ignored.
func (w *Workshop) IsPathIgnored(path string) bool {
	return w.srcDir.IsPathIgnored(path)
}

// GetFiles gets a list of files and their size in total from the source path
// based on ignore list.
func (w *Workshop) GetFiles() ([]string, int64, error) {
	files, size, err := w.srcDir.ListFiles()
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

// ZipFiles create an archive of all files retrieved earlier using GetFiles.
func (w *Workshop) ZipFiles() error {
	if len(w.files) == 0 {
		return errors.New("no files to zip")
	}

	archive, err := os.Create(w.relDestPath + ".zip")
	if err != nil {
		return err
	}

	zipWriter := zip.NewWriter(archive)
	defer func(archive *os.File, zipWriter *zip.Writer) {
		_ = archive.Close()
		_ = zipWriter.Close()
	}(archive, zipWriter)

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

		w, err := zipWriter.Create(file)
		if err != nil {
			_ = src.Close()
			return err
		}

		if _, err := io.Copy(w, src); err != nil {
			_ = src.Close()
			return err
		}

		_ = src.Close()
	}

	_ = zipWriter.Close()
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
	return w.srcDir.RelPath()
}

// AbsSrcPath gets an absolute source path.
func (w *Workshop) AbsSrcPath() string {
	return w.srcDir.AbsPath()
}

// RelDestPath gets a relative destination path.
func (w *Workshop) RelDestPath() string {
	return w.relDestPath
}

// AbsDestPath gets an absolute destination path.
func (w *Workshop) AbsDestPath() string {
	return w.absDestPath
}

// PrintFiles prints a list of files retrieved earlier using GetFiles.
func (w *Workshop) PrintFiles() {
	for _, file := range w.files {
		fmt.Println(file)
	}
}

// Package dir has been designed to list files within a certain directory
// excluding files from ignore list.
package dir

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/mattn/go-zglob"
)

// Controller is the interface that wraps the Dir methods.
type Controller interface {
	SetIgnore([]string)
	Ignore() []string
	AbsPath() string
	RelPath() string
	Base() string
	Dir() string
	IsPathIgnored(string) bool
	ListFiles(...string) ([]string, int64, error)
}

// Dir represents a working directory.
type Dir struct {
	absPath string
	ignore  []string
	relPath string
}

// New creates a new Dir instance.
func New(path string) (*Dir, error) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}

	relPath, err := filepath.Rel(absPath, absPath)
	if err != nil {
		return nil, err
	}

	return &Dir{
		absPath: absPath,
		relPath: relPath,
	}, nil
}

// SetIgnore sets ignore list.
func (d *Dir) SetIgnore(ignore []string) {
	d.ignore = ignore
}

// Ignore gets ignore list.
func (d *Dir) Ignore() []string {
	return d.ignore
}

// AbsPath returns an absolute path.
func (d *Dir) AbsPath() string {
	return d.absPath
}

// RelPath returns a relative path.
func (d *Dir) RelPath() string {
	return d.relPath
}

// Base returns a base path.
func (d *Dir) Base() string {
	return filepath.Base(d.absPath)
}

// Dir returns a directory.
func (d *Dir) Dir() string {
	return filepath.Dir(d.absPath)
}

// IsPathIgnored checks if the provided path is ignored based on ignore.
func (d *Dir) IsPathIgnored(path string) bool {
	for _, ignore := range d.ignore {
		hasPrefix := strings.HasPrefix(ignore, "/")
		ignore = strings.TrimPrefix(ignore, "/")
		ignore = strings.TrimSuffix(ignore, "/")

		prefix := "**/"
		if hasPrefix {
			prefix = ""
		}

		matched, _ := zglob.Match(prefix+ignore+"**/*", path)
		if matched {
			return true
		}
	}
	return false
}

// ListFiles lists all files and their size in total from the based on ignore.
func (d *Dir) ListFiles(ext ...string) (path []string, size int64, err error) {
	if err := filepath.Walk(d.relPath, func(p string, i os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if i.IsDir() || d.IsPathIgnored(p) {
			return nil
		}

		if len(ext) > 0 {
			isFound := false
			for _, e := range ext {
				if filepath.Ext(p) == e {
					isFound = true
				}
			}

			if !isFound {
				return nil
			}
		}

		path = append(path, p)
		size += i.Size()
		return nil
	}); err != nil {
		return path, size, err
	}
	return path, size, err
}

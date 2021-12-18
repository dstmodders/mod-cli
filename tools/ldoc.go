package tools

import (
	"errors"
	"strings"
)

// LDoc represents an LDoc tool.
type LDoc struct {
	Tool
}

// NewLDoc creates a new LDoc instance.
func NewLDoc() (*LDoc, error) {
	tool, err := NewTool("LDoc", "ldoc")
	if err != nil {
		return nil, err
	}
	return &LDoc{
		Tool: *tool,
	}, nil
}

func (l *LDoc) parseVersion(str string) (string, error) {
	match := versionRegex.FindString(str)
	if len(match) == 0 {
		return "", errors.New("not found")
	}
	return strings.TrimSpace(match), nil
}

// LoadVersion loads an LDoc version.
func (l *LDoc) LoadVersion() (string, error) {
	cmd := l.ExecCommand("--version")
	out, _ := cmd.CombinedOutput()

	ver, err := l.parseVersion(string(out))
	if err != nil {
		return ver, err
	}
	l.version = ver

	return ver, nil
}

package tools

import (
	"errors"
	"strings"
)

// Ktools represents a Ktools tool.
type Ktools struct {
	Tool
}

// NewKtools creates a new Ktools instance.
func NewKtools(name, cmd string) (*Ktools, error) {
	tool, err := NewTool(name, cmd)
	if err != nil {
		return nil, err
	}
	return &Ktools{
		Tool: *tool,
	}, nil
}

func (k *Ktools) parseVersion(str string) (string, error) {
	match := versionRegex.FindString(str)
	if len(match) == 0 {
		return "", errors.New("not found")
	}
	return strings.TrimSpace(match), nil
}

// LoadVersion loads a Ktools version.
func (k *Ktools) LoadVersion() (string, error) {
	cmd := k.ExecCommand("--version")
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}

	ver, err := k.parseVersion(string(out))
	if err != nil {
		return ver, err
	}
	k.version = ver

	return ver, nil
}

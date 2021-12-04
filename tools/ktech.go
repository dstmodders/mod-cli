package tools

import (
	"errors"
	"strings"
)

// Ktech represents a Ktech tool.
type Ktech struct {
	Tool
}

// NewKtech creates a new Ktech instance.
func NewKtech() *Ktech {
	return &Ktech{
		Tool: *NewTool("ktech", "ktech"),
	}
}

func (k *Ktech) parseVersion(str string) (string, error) {
	match := versionRegex.FindString(str)
	if len(match) == 0 {
		return "", errors.New("not found")
	}
	return strings.TrimSpace(match), nil
}

// LoadVersion loads a ktech version.
func (k *Ktech) LoadVersion() (string, error) {
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

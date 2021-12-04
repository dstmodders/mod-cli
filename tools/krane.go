package tools

import (
	"errors"
	"strings"
)

// Krane represents a Krane tool.
type Krane struct {
	Tool
}

// NewKrane creates a new Krane instance.
func NewKrane() *Krane {
	return &Krane{
		Tool: *NewTool("krane", "krane"),
	}
}

func (k *Krane) parseVersion(str string) (string, error) {
	match := versionRegex.FindString(str)
	if len(match) == 0 {
		return "", errors.New("not found")
	}
	return strings.TrimSpace(match), nil
}

// LoadVersion loads a krane version.
func (k *Krane) LoadVersion() (string, error) {
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

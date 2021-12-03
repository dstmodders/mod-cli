package tools

import (
	"errors"
	"strings"
)

// Busted represents a Busted tool.
type Busted struct {
	Tool
}

// NewBusted creates a new Busted instance.
func NewBusted() *Busted {
	return &Busted{
		Tool: *NewTool("Busted", "busted"),
	}
}

func (b *Busted) parseVersion(str string) (string, error) {
	return strings.TrimSpace(str), nil
}

// LoadVersion loads a Busted version.
func (b *Busted) LoadVersion() (string, error) {
	cmd := b.ExecCommand("--version")
	out, err := cmd.Output()
	if err != nil {
		return "", errors.New("not found")
	}

	ver, err := b.parseVersion(string(out))
	if err != nil {
		return ver, err
	}
	b.version = ver

	return ver, nil
}

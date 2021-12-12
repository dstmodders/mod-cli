package tools

import (
	"strings"
)

// Prettier represents a Prettier tool.
type Prettier struct {
	Tool
}

// NewPrettier creates a new Prettier instance.
func NewPrettier() *Prettier {
	return &Prettier{
		Tool: *NewTool("Prettier", "prettier"),
	}
}

func (d *Prettier) parseVersion(str string) (string, error) {
	return strings.TrimSpace(str), nil
}

// LoadVersion loads a Prettier version.
func (d *Prettier) LoadVersion() (string, error) {
	cmd := d.ExecCommand("--version")
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}

	ver, err := d.parseVersion(string(out))
	if err != nil {
		return ver, err
	}
	d.version = ver

	return ver, nil
}

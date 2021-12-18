package tools

import (
	"strings"
)

// Prettier represents a Prettier tool.
type Prettier struct {
	Tool
}

// NewPrettier creates a new Prettier instance.
func NewPrettier() (*Prettier, error) {
	tool, err := NewTool("Prettier", "prettier")
	if err != nil {
		return nil, err
	}
	return &Prettier{
		Tool: *tool,
	}, nil
}

func (p *Prettier) parseVersion(str string) (string, error) {
	return strings.TrimSpace(str), nil
}

// LoadVersion loads a Prettier version.
func (p *Prettier) LoadVersion() (string, error) {
	cmd := p.ExecCommand("--version")
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}

	ver, err := p.parseVersion(string(out))
	if err != nil {
		return ver, err
	}
	p.version = ver

	return ver, nil
}

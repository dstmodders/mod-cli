package tools

import (
	"bytes"
	"errors"
	"io"
	"strings"
)

// StyLua represents a StyLua tool.
type StyLua struct {
	Tool
}

// NewStyLua creates a new StyLua instance.
func NewStyLua() *StyLua {
	return &StyLua{
		Tool: *NewTool("StyLua", "stylua"),
	}
}

func (s *StyLua) parseVersion(str string) (string, error) {
	words := strings.Split(str, " ")
	if len(words) == 0 {
		return "", errors.New("not found")
	}
	return strings.TrimSpace(words[1]), nil
}

// LoadVersion loads a StyLua version.
func (s *StyLua) LoadVersion() (string, error) {
	cmd := s.ExecCommand("--version")

	stdout, _ := cmd.StdoutPipe()

	if err := cmd.Start(); err != nil {
		return "", err
	}

	buf := bytes.NewBuffer(make([]byte, 10))
	_, err := io.Copy(buf, stdout)
	if err != nil {
		return "", err
	}

	ver, err := s.parseVersion(buf.String())
	if err != nil {
		return ver, err
	}
	s.version = ver

	return ver, nil
}

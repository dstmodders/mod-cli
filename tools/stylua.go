package tools

import (
	"bufio"
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
func NewStyLua() (*StyLua, error) {
	tool, err := NewTool("StyLua", "stylua")
	if err != nil {
		return nil, err
	}
	return &StyLua{
		Tool: *tool,
	}, nil
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

// Check checks formatting in the provided files.
func (s *StyLua) Check(arg ...string) (result Format, err error) {
	if len(arg) == 0 {
		files, _, _ := s.workingDir.ListFiles(".lua")
		arg = append(arg, files...)
	}

	a := []string{"-c"}
	a = append(a, arg...)

	cmd := s.ExecCommand(a...)
	stdout, _ := cmd.StdoutPipe()

	if err := cmd.Start(); err != nil {
		return result, err
	}

	scanner := bufio.NewScanner(stdout)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		line := scanner.Text()
		line = ansiRegex.ReplaceAllString(line, "")
		if strings.HasPrefix(line, "Diff in ") {
			str := strings.TrimPrefix(line, "Diff in ")
			str = strings.TrimSuffix(str, ":")
			result.Files = append(result.Files, FormatFile{
				Path: strings.TrimSpace(str),
			})
		}
	}

	_ = cmd.Wait()

	return result, nil
}

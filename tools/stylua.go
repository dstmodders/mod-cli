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

	// DefaultExt holds default extensions: ".lua".
	DefaultExt []string
}

// NewStyLua creates a new StyLua instance.
func NewStyLua() (*StyLua, error) {
	tool, err := NewTool("StyLua", "stylua")
	if err != nil {
		return nil, err
	}
	return &StyLua{
		Tool:       *tool,
		DefaultExt: []string{".lua"},
	}, nil
}

func (s *StyLua) parseVersion(str string) (string, error) {
	words := strings.Split(str, " ")
	if len(words) == 0 {
		return "", errors.New("not found")
	}
	return strings.TrimSpace(words[1]), nil
}

func (s *StyLua) prepareArg(write bool, files ...string) []string {
	if len(files) == 0 {
		f, _, _ := s.workingDir.ListFiles(s.DefaultExt...)
		files = append(files, f...)
	}
	var a []string
	if !write {
		a = append(a, "-c")
	}
	return append(a, files...)
}

// LoadVersion loads a StyLua version.
func (s *StyLua) LoadVersion() (string, error) {
	cmd := s.ExecCommand("--version")

	stdout, _ := cmd.StdoutPipe()

	if err := cmd.Start(); err != nil {
		return "", err
	}

	buf := bytes.NewBuffer([]byte{})
	_, err := io.Copy(buf, stdout)
	if err != nil {
		return "", err
	}

	str := buf.String()
	if len(str) == 0 {
		return "", errors.New("no output")
	}

	ver, err := s.parseVersion(str)
	if err != nil {
		return ver, err
	}
	s.version = ver

	return ver, nil
}

// Check checks formatting in the provided files.
func (s *StyLua) Check(arg ...string) (result Format, err error) {
	cmd := s.ExecCommand(s.prepareArg(false, arg...)...)
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
				Path:  strings.TrimSpace(str),
				State: FileStateWarning,
			})
		}
	}

	_ = cmd.Wait()

	return result, nil
}

// Fix fixes formatting in the provided files.
func (s *StyLua) Fix(arg ...string) (result Format, err error) {
	result, err = s.Check(arg...)
	if err != nil {
		return result, err
	}

	cmd := s.ExecCommand(s.prepareArg(true, arg...)...)

	if err := cmd.Start(); err != nil {
		return result, err
	}

	for i := 0; i < len(result.Files); i++ {
		result.Files[i].State = FileStateSuccess
	}

	_ = cmd.Wait()

	return result, nil
}

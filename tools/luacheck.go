package tools

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"strconv"
	"strings"
)

// Luacheck represents a Luacheck tool.
type Luacheck struct {
	Tool
}

// NewLuacheck creates a new Luacheck instance.
func NewLuacheck() (*Luacheck, error) {
	tool, err := NewTool("Luacheck", "luacheck")
	if err != nil {
		return nil, err
	}
	return &Luacheck{
		Tool: *tool,
	}, nil
}

func (l *Luacheck) parseVersion(str string) (string, error) {
	s := strings.Split(str, "\n")
	if len(s) == 0 {
		return "", errors.New("not found")
	}
	result := s[0]
	result = strings.ReplaceAll(result, "Luacheck: ", "")
	result = strings.TrimSpace(result)
	return result, nil
}

// LoadVersion loads a Luacheck version.
func (l *Luacheck) LoadVersion() (string, error) {
	cmd := l.ExecCommand("--version")

	stdout, _ := cmd.StdoutPipe()

	if err := cmd.Start(); err != nil {
		return "", err
	}

	buf := bytes.NewBuffer(make([]byte, 10))
	_, err := io.Copy(buf, stdout)
	if err != nil {
		return "", err
	}

	ver, err := l.parseVersion(buf.String())
	if err != nil {
		return ver, err
	}
	l.version = ver

	return ver, nil
}

// Lint lints provided files.
func (l *Luacheck) Lint(arg ...string) (result Lint, err error) {
	if len(arg) == 0 {
		files, _, _ := l.workingDir.ListFiles(".lua")
		arg = append(arg, files...)
	}

	cmd := l.ExecCommand(arg...)
	stdout, _ := cmd.StdoutPipe()

	if err := cmd.Start(); err != nil {
		return result, err
	}

	scanner := bufio.NewScanner(stdout)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		line := scanner.Text()
		line = ansiRegex.ReplaceAllString(line, "")
		matches := luacheckRegex.FindStringSubmatch(line)
		if len(matches) > 0 {
			issues, err := strconv.Atoi(matches[2])
			if err != nil {
				return result, err
			}

			result.Files = append(result.Files, LintFile{
				Path:   strings.TrimSpace(matches[1]),
				Issues: issues,
			})
		}
	}

	_ = cmd.Wait()

	return result, nil
}

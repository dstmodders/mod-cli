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

	buf := bytes.NewBuffer([]byte{})
	_, err := io.Copy(buf, stdout)
	if err != nil {
		return "", err
	}

	str := buf.String()
	if len(str) == 0 {
		return "", errors.New("no output")
	}

	ver, err := l.parseVersion(buf.String())
	if err != nil {
		return ver, err
	}
	l.version = ver

	return ver, nil
}

//nolint:funlen
// Lint lints provided files.
func (l *Luacheck) Lint(arg ...string) (result Lint, err error) {
	var stdoutLines []string
	var file LintFile

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
		stdoutLines = append(stdoutLines, line)
		line = ansiRegex.ReplaceAllString(line, "")

		if luacheckIssueRegex.MatchString(line) {
			matches := luacheckIssueRegex.FindStringSubmatch(line)
			if len(matches) != 5 {
				continue
			}

			name := cleanString(matches[1])

			if file.Path != name {
				result.Files = append(result.Files, file)
				file = LintFile{
					Path:   name,
					Issues: []LintFileIssue{},
				}
			}
			file.Path = name

			startLine, err := strconv.Atoi(matches[2])
			if err != nil {
				return result, err
			}

			endLine, err := strconv.Atoi(matches[3])
			if err != nil {
				return result, err
			}

			file.Issues = append(file.Issues, LintFileIssue{
				Name:        name,
				StartLine:   startLine,
				EndLine:     endLine,
				Description: cleanString(matches[4]),
			})
		}
	}

	if len(result.Files) > 1 {
		result.Files = result.Files[1:]
	}

	result.Stdout = strings.Join(stdoutLines, "\n")
	result.Stdout = strings.TrimSpace(result.Stdout)

	_ = cmd.Wait()

	return result, nil
}

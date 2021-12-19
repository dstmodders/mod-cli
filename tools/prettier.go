package tools

import (
	"bufio"
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

// Check checks formatting in the provided files.
func (p *Prettier) Check(arg ...string) (result Format, err error) {
	if len(arg) == 0 {
		files, _, _ := p.workingDir.ListFiles(".md", ".xml", ".yml")
		arg = append(arg, files...)
	}

	a := []string{"--list-different"}
	a = append(a, arg...)

	cmd := p.ExecCommand(a...)
	stdout, _ := cmd.StdoutPipe()

	if err := cmd.Start(); err != nil {
		return result, err
	}

	scanner := bufio.NewScanner(stdout)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		line := scanner.Text()
		line = ansiRegex.ReplaceAllString(line, "")
		result.Files = append(result.Files, FormatFile{
			Path: strings.TrimSpace(line),
		})
	}

	_ = cmd.Wait()

	return result, nil
}

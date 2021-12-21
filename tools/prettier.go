package tools

import (
	"bufio"
	"strings"
)

// Prettier represents a Prettier tool.
type Prettier struct {
	Tool

	// DefaultExt holds default extensions: ".md", ".xml", ".yml".
	DefaultExt []string

	// ListDifferent sets whether only different files should be listed.
	ListDifferent bool
}

// NewPrettier creates a new Prettier instance.
func NewPrettier() (*Prettier, error) {
	tool, err := NewTool("Prettier", "prettier")
	if err != nil {
		return nil, err
	}
	return &Prettier{
		Tool:          *tool,
		DefaultExt:    []string{".md", ".xml", ".yml"},
		ListDifferent: true,
	}, nil
}

func (p *Prettier) parseVersion(str string) (string, error) {
	return strings.TrimSpace(str), nil
}

func (p *Prettier) prepareArg(write bool, files ...string) []string {
	if len(files) == 0 {
		f, _, _ := p.workingDir.ListFiles(p.DefaultExt...)
		files = append(files, f...)
	}

	var a []string

	if p.ListDifferent {
		a = append(a, "--list-different")
	}

	if write {
		a = append(a, "-w")
	}

	return append(a, files...)
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
	cmd := p.ExecCommand(p.prepareArg(false, arg...)...)
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
			Path:  strings.TrimSpace(line),
			State: FileStateWarning,
		})
	}

	_ = cmd.Wait()

	return result, nil
}

// Fix fixes formatting in the provided files.
func (p *Prettier) Fix(arg ...string) (result Format, err error) {
	cmd := p.ExecCommand(p.prepareArg(true, arg...)...)
	stdout, _ := cmd.StdoutPipe()

	if err := cmd.Start(); err != nil {
		return result, err
	}

	scanner := bufio.NewScanner(stdout)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		line := scanner.Text()
		line = ansiRegex.ReplaceAllString(line, "")
		splitted := strings.Split(line, " ")
		result.Files = append(result.Files, FormatFile{
			Path:  strings.TrimSpace(splitted[0]),
			State: FileStateSuccess,
		})
	}

	_ = cmd.Wait()

	return result, nil
}

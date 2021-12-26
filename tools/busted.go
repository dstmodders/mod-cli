package tools

import (
	"bufio"
	"errors"
	"fmt"
	"strings"
)

// Busted represents a Busted tool.
type Busted struct {
	Tool
}

// NewBusted creates a new Busted instance.
func NewBusted() (*Busted, error) {
	tool, err := NewTool("Busted", "busted")
	if err != nil {
		return nil, err
	}
	return &Busted{
		Tool: *tool,
	}, nil
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

// Test runs tests.
func (b *Busted) Test() (result Lint, err error) {
	cmd := b.ExecCommand(".")
	stdout, _ := cmd.StdoutPipe()

	if err := cmd.Start(); err != nil {
		return result, err
	}

	scanner := bufio.NewScanner(stdout)
	scanner.Split(bufio.ScanRunes)
	for scanner.Scan() {
		r := scanner.Text()
		r = ansiRegex.ReplaceAllString(r, "")
		fmt.Print(r)
	}

	_ = cmd.Wait()

	return result, nil
}

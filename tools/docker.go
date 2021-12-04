package tools

import (
	"errors"
	"fmt"
	"strings"
)

// Docker represents a Docker tool.
type Docker struct {
	Tool
}

// NewDocker creates a new Docker instance.
func NewDocker() *Docker {
	return &Docker{
		Tool: *NewTool("Docker", "docker"),
	}
}

func (d *Docker) parseVersion(str string) (string, error) {
	match := dockerVersionRegex.FindStringSubmatch(str)
	if len(match) != 3 {
		return "", errors.New("not found")
	}
	version := strings.TrimSpace(match[1])
	build := strings.TrimSpace(match[2])
	return fmt.Sprintf("%s-%s", version, build), nil
}

// LoadVersion loads a Docker version.
func (d *Docker) LoadVersion() (string, error) {
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

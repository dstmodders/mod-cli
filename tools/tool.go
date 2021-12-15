package tools

import (
	"os/exec"
	"strings"
)

// Tooler is the interface that wraps the Tool methods.
type Tooler interface {
	exec(string, ...string) *exec.Cmd
	execDocker(...string) *exec.Cmd
	parseVersion(string) (string, error)
	Name() string
	Path() string
	Version() string
	SetDockerized(*Dockerized) error
	SetRunInDocker(bool)
	ExecCommand(...string) *exec.Cmd
	LookPath() (string, error)
	LoadVersion() (string, error)
}

// Tool represents a single tool.
type Tool struct {
	Cmd           string
	CmdArgs       []string
	CmdDockerArgs []string
	dockerized    *Dockerized
	name          string
	path          string
	runInDocker   bool
	version       string
}

// NewTool creates a new Tool instance.
func NewTool(name, cmd string) *Tool {
	return &Tool{
		Cmd:         cmd,
		CmdArgs:     []string{},
		dockerized:  NewDockerized(),
		name:        name,
		runInDocker: false,
	}
}

func (t *Tool) exec(name string, arg ...string) *exec.Cmd {
	a := t.CmdArgs
	a = append(a, arg...)
	return exec.Command(name, a...)
}

func (t *Tool) execDocker(arg ...string) *exec.Cmd {
	a := t.dockerized.Args()
	a = append(a, arg...)
	return exec.Command("docker", a...)
}

// Name returns a name of the tool.
func (t *Tool) Name() string {
	return t.name
}

// Path returns a path of the tool found earlier using LookPath.
func (t *Tool) Path() string {
	return t.path
}

// Version returns a version of the tool for both direct and Dockerized usage.
func (t *Tool) Version() string {
	return t.version
}

// SetDockerized sets Dockerized.
func (t *Tool) SetDockerized(dockerized *Dockerized) error {
	t.dockerized = dockerized
	_, err := t.dockerized.PrepareArgs()
	return err
}

// SetRunInDocker sets whether a tool should be run in Docker.
func (t *Tool) SetRunInDocker(runInDocker bool) {
	t.runInDocker = runInDocker
}

// ExecCommand executes the command with the passed arguments either directly or
// through a Docker container.
func (t *Tool) ExecCommand(arg ...string) *exec.Cmd {
	if t.runInDocker {
		a := []string{t.Cmd}
		a = append(a, arg...)
		return t.execDocker(a...)
	}
	return t.exec(t.Cmd, arg...)
}

// LookPath looks for a direct path of the tool which can be retrieved later
// using Path.
func (t *Tool) LookPath() (string, error) {
	if t.runInDocker {
		cmd := t.execDocker("which", t.Cmd)
		out, err := cmd.Output()
		t.path = strings.TrimSpace(string(out))
		return t.path, err
	}

	path, err := exec.LookPath(t.Cmd)
	if err != nil {
		return "", err
	}

	t.path = path
	return path, nil
}

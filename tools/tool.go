package tools

import "os/exec"

// Tooler is the interface that wraps the Tool methods.
type Tooler interface {
	parseVersion(string) (string, error)
	Name() string
	Path() string
	Version() string
	ExecCommand(...string) *exec.Cmd
	LookPath() (string, error)
	LoadVersion() (string, error)
}

// Tool represents a single tool.
type Tool struct {
	Cmd           string
	CmdArgs       []string
	CmdDockerArgs []string
	RunInDocker   bool
	name          string
	path          string
	version       string
}

// NewTool creates a new Tool instance.
func NewTool(name, cmd string) *Tool {
	return &Tool{
		Cmd:     cmd,
		CmdArgs: []string{},
		CmdDockerArgs: []string{
			"run",
			"--rm",
			"-u",
			"dst-mod",
			"dstmodders/dst-mod:latest",
		},
		RunInDocker: false,
		name:        name,
	}
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

// ExecCommand executes the command with the passed arguments either directly or
// through a Docker container.
func (t *Tool) ExecCommand(arg ...string) *exec.Cmd {
	name := t.Cmd
	a := t.CmdArgs
	if t.RunInDocker {
		name = "docker"
		a = t.CmdDockerArgs
		a = append(a, t.Cmd)
	}
	a = append(a, arg...)
	return exec.Command(name, a...)
}

// LookPath looks for a direct path of the tool which can be retrieved later
// using Path.
func (t *Tool) LookPath() (string, error) {
	path, err := exec.LookPath(t.Cmd)
	if err != nil {
		return "", err
	}
	t.path = path
	return path, nil
}

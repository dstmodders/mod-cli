package tools

import (
	"os"
	"os/exec"
	"strings"

	"github.com/dstmodders/mod-cli/dir"
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
	DockerImage() string
	SetDockerImage(image string)
	IsDockerImageAvailable() bool
	PullDockerImage() bool
	SetIgnore([]string)
	SetRunInDocker(bool)
	ExecCommand(...string) *exec.Cmd
	LookPath() (string, error)
	ExistsOnSystem() bool
	ExistsInDocker() bool
	LoadVersion() (string, error)
}

// Tool represents a single tool.
type Tool struct {
	// Cmd holds a command.
	Cmd string

	// CmdArgs holds a command arguments.
	CmdArgs []string

	dockerized  *Dockerized
	name        string
	path        string
	runInDocker bool
	version     string
	workingDir  dir.Dir
}

// Format represents a formatting result.
type Format struct {
	Files []FormatFile
}

// FormatFile represents a single formatting file.
type FormatFile struct {
	// Path holds a file path.
	Path string

	// State holds a file state.
	State FileState
}

// Lint represents a linting result.
type Lint struct {
	Files  []LintFile
	Stdout string
}

// LintFile represents a single linting file.
type LintFile struct {
	// Path holds a file path.
	Path string

	// State holds a file state.
	State FileState

	// Issues holds the found issues.
	Issues []LintFileIssue
}

// LintFileIssue represents a single issue in a file.
type LintFileIssue struct {
	// Name holds a file name.
	Name string

	// StartLine holds the start line number which corresponds to this issue.
	StartLine int

	// EndLine holds the end line number which corresponds to this issue.
	EndLine int

	// Issue holds the issue description.
	Description string
}

// FileState represents a single linting or formatting file result state.
type FileState int

const (
	// FileStateWarning represents a state of the file having some issues.
	FileStateWarning FileState = iota

	// FileStateSuccess represents a state of the file after successful fix.
	FileStateSuccess
)

// NewTool creates a new Tool instance.
func NewTool(name, cmd string) (*Tool, error) {
	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	workingDir, err := dir.New(wd)
	if err != nil {
		return nil, err
	}

	return &Tool{
		Cmd:         cmd,
		CmdArgs:     []string{},
		dockerized:  NewDockerized(),
		name:        name,
		runInDocker: false,
		workingDir:  *workingDir,
	}, nil
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

// SetDockerized sets dockerized.
func (t *Tool) SetDockerized(dockerized *Dockerized) error {
	t.dockerized = dockerized
	_, err := t.dockerized.PrepareArgs()
	return err
}

// DockerImage gets the Docker image.
func (t *Tool) DockerImage() string {
	return t.dockerized.Image
}

// SetDockerImage sets the Docker image.
func (t *Tool) SetDockerImage(image string) {
	t.dockerized.Image = image
}

// IsDockerImageAvailable checks if the Docker image is available.
func (t *Tool) IsDockerImageAvailable() bool {
	return t.dockerized.IsImageAvailable()
}

// PullDockerImage pull the Docker image.
func (t *Tool) PullDockerImage() bool {
	return t.dockerized.PullImage()
}

// SetIgnore sets ignore list.
func (t *Tool) SetIgnore(ignore []string) {
	t.workingDir.SetIgnore(ignore)
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

// ExistsOnSystem checks if the tool exists on the system.
func (t *Tool) ExistsOnSystem() bool {
	_, err := exec.LookPath(t.Cmd)
	return err == nil
}

// ExistsInDocker checks if the tool exists inside a Docker container.
func (t *Tool) ExistsInDocker() bool {
	cmd := t.execDocker("which", t.Cmd)
	out, err := cmd.Output()
	if err != nil {
		return false
	}
	t.path = strings.TrimSpace(string(out))
	return len(t.path) > 0
}

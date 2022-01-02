package tools

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// Dockerized represents a Docker arguments to run a tool.
type Dockerized struct {
	// Image hold the image name and tag.
	//
	// Default: dstmodders/dst-mod:latest
	Image string

	// Remove sets whether a container should be removed right after running.
	//
	// Default: true
	Remove bool

	// User holds a username or UID to run a container as.
	//
	// Default: dst-mod
	User string

	// Volume holds a volume to mount. By default, points to a working directory.
	Volume string

	args []string
}

// NewDockerized creates a new Dockerized instance.
func NewDockerized() *Dockerized {
	d := &Dockerized{
		Image:  "dstmodders/dst-mod:latest",
		Remove: true,
		User:   "dst-mod",
	}
	_, _ = d.PrepareArgs()
	return d
}

// Args returns arguments prepared using PrepareArgs.
func (d *Dockerized) Args() []string {
	return d.args
}

// IsImageAvailable checks whether an image is available locally.
func (d *Dockerized) IsImageAvailable() bool {
	cmd := exec.Command("docker", "image", "ls")

	stdout, _ := cmd.StdoutPipe()

	if err := cmd.Start(); err != nil {
		return false
	}

	image := d.Image
	tag := "latest"

	split := strings.Split(d.Image, ":")
	if len(split) == 2 {
		image = split[0]
		tag = split[1]
	}

	scanner := bufio.NewScanner(stdout)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, image) && strings.Contains(line, tag) {
			return true
		}
	}

	_ = cmd.Wait()

	return false
}

// PrepareArgs prepare arguments to return later using Args.
func (d *Dockerized) PrepareArgs() (result []string, err error) {
	volume := d.Volume

	if len(volume) == 0 {
		volume, err = os.Getwd()
		if err != nil {
			return result, err
		}
	}

	result = []string{"run"}

	if d.Remove {
		result = append(result, "--rm")
	}

	if len(d.User) > 0 {
		result = append(result, "-u")
		result = append(result, d.User)
	}

	if len(volume) > 0 {
		base := filepath.Base(volume)
		result = append(result, "-v")
		result = append(result, fmt.Sprintf("%s:/opt/%s:z", volume, base))
		result = append(result, "-w")
		result = append(result, fmt.Sprintf("/opt/%s", base))
	}

	result = append(result, d.Image)
	d.args = result

	return result, err
}

// PullImage pulls an image.
func (d *Dockerized) PullImage() bool {
	//nolint:gosec
	cmd := exec.Command("docker", "pull", d.Image)

	if err := cmd.Start(); err != nil {
		return false
	}

	if err := cmd.Wait(); err != nil {
		return false
	}

	return true
}

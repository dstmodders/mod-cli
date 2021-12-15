package tools

import (
	"fmt"
	"os"
	"path"
)

// Dockerized represents a Docker arguments to run a tool.
type Dockerized struct {
	Image  string
	Remove bool
	User   string
	args   []string
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

// PrepareArgs prepare arguments to return later using Args.
func (d *Dockerized) PrepareArgs() (result []string, err error) {
	dir, err := os.Getwd()
	if err != nil {
		return result, err
	}

	base := path.Base(dir)
	result = []string{"run"}

	if d.Remove {
		result = append(result, "--rm")
	}

	if len(d.User) > 0 {
		result = append(result, "-u")
		result = append(result, d.User)
	}

	if len(d.User) > 0 {
		result = append(result, "-v")
		result = append(result, fmt.Sprintf("%s:/opt/%s", dir, base))
		result = append(result, "-w")
		result = append(result, fmt.Sprintf("/opt/%s", base))
	}

	result = append(result, d.Image)
	d.args = result

	return result, err
}

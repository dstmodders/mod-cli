package tools

import (
	"fmt"
	"os"
	"path/filepath"
)

// Dockerized represents a Docker arguments to run a tool.
type Dockerized struct {
	// Image hold the image name and tag. By default: dstmodders/dst-mod:latest.
	Image string

	// Remove sets whether a container should be removed right after running. By
	// default: true.
	Remove bool

	// User holds a username or UID to run a container as. By default: dst-mod.
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
		result = append(result, fmt.Sprintf("%s:/opt/%s", volume, base))
		result = append(result, "-w")
		result = append(result, fmt.Sprintf("/opt/%s", base))
	}

	result = append(result, d.Image)
	d.args = result

	return result, err
}

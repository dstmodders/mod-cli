// Package tools has been designed to run different tools either directly or
// though Docker.
package tools

import "regexp"

const regexSemVer string = `v?([0-9]+)(\.[0-9]+)?(\.[0-9]+)?` +
	`(-([0-9A-Za-z\-]+(\.[0-9A-Za-z\-]+)*))?` +
	`(\+([0-9A-Za-z\-]+(\.[0-9A-Za-z\-]+)*))?`

const regexDockerVer string = `Docker version (.*), build (.*)`

var versionRegex *regexp.Regexp
var dockerVersionRegex *regexp.Regexp

func init() {
	versionRegex = regexp.MustCompile(regexSemVer)
	dockerVersionRegex = regexp.MustCompile(regexDockerVer)
}

// Controller is the interface that wraps the Tools methods.
type Controller interface {
	SetToolsRunInDocker(bool)
	LookPaths()
	LoadVersions()
}

// Tools represents all the supported tools.
type Tools struct {
	Busted   *Busted
	Docker   *Docker
	Krane    *Krane
	Ktech    *Ktech
	LDoc     *LDoc
	Luacheck *Luacheck
	Prettier *Prettier
	all      []Tooler
}

// New creates a new Tools instance.
func New() *Tools {
	busted := NewBusted()
	docker := NewDocker()
	krane := NewKrane()
	ktech := NewKtech()
	ldoc := NewLDoc()
	luacheck := NewLuacheck()
	prettier := NewPrettier()

	return &Tools{
		Busted:   busted,
		Docker:   docker,
		Krane:    krane,
		Ktech:    ktech,
		LDoc:     ldoc,
		Luacheck: luacheck,
		Prettier: prettier,
		all: []Tooler{
			busted,
			docker,
			krane,
			ktech,
			ldoc,
			luacheck,
			prettier,
		},
	}
}

// SetToolsRunInDocker sets all tools to be run in Docker.
func (t *Tools) SetToolsRunInDocker(runInDocker bool) {
	for _, tool := range t.all {
		tool.SetRunInDocker(runInDocker)
	}
}

// LookPaths looks for paths of all tools.
func (t *Tools) LookPaths() {
	for _, tool := range t.all {
		_, _ = tool.LookPath()
	}
}

// LoadVersions loads versions of all tools.
func (t *Tools) LoadVersions() {
	for _, tool := range t.all {
		_, _ = tool.LoadVersion()
	}
}

// Package tools has been designed to run different tools either directly or
// though Docker.
package tools

import (
	"regexp"
	"strings"
)

// Taken from: https://github.com/acarl005/stripansi
const ansiSemVer string = "[\u001B\u009B][[\\]()#;?]*(?:(?:(?:[a-zA-Z\\d]*(?:;[a-zA-Z\\d]*)*)?\u0007)|(?:(?:\\d{1,4}(?:;\\d{0,4})*)?[\\dA-PRZcf-ntqry=><~]))"

const luacheckIssue string = `(.*):(\d*):(\d*):\s*(.*)`

const regexDockerVer string = `Docker version (.*), build (.*)`

const regexSemVer string = `v?([0-9]+)(\.[0-9]+)?(\.[0-9]+)?` +
	`(-([0-9A-Za-z\-]+(\.[0-9A-Za-z\-]+)*))?` +
	`(\+([0-9A-Za-z\-]+(\.[0-9A-Za-z\-]+)*))?`

var ansiRegex *regexp.Regexp
var luacheckIssueRegex *regexp.Regexp
var versionRegex *regexp.Regexp
var dockerVersionRegex *regexp.Regexp

func init() {
	ansiRegex = regexp.MustCompile(ansiSemVer)
	luacheckIssueRegex = regexp.MustCompile(luacheckIssue)
	versionRegex = regexp.MustCompile(regexSemVer)
	dockerVersionRegex = regexp.MustCompile(regexDockerVer)
}

func cleanString(str string) string {
	return strings.TrimSpace(ansiRegex.ReplaceAllString(str, ""))
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
	StyLua   *StyLua
	all      []Tooler
}

// New creates a new Tools instance.
func New() (*Tools, error) {
	busted, err := NewBusted()
	if err != nil {
		return nil, err
	}

	docker, err := NewDocker()
	if err != nil {
		return nil, err
	}

	krane, err := NewKrane()
	if err != nil {
		return nil, err
	}

	ktech, err := NewKtech()
	if err != nil {
		return nil, err
	}

	ldoc, err := NewLDoc()
	if err != nil {
		return nil, err
	}

	luacheck, err := NewLuacheck()
	if err != nil {
		return nil, err
	}

	prettier, err := NewPrettier()
	if err != nil {
		return nil, err
	}

	stylua, err := NewStyLua()
	if err != nil {
		return nil, err
	}

	return &Tools{
		Busted:   busted,
		Docker:   docker,
		Krane:    krane,
		Ktech:    ktech,
		LDoc:     ldoc,
		Luacheck: luacheck,
		Prettier: prettier,
		StyLua:   stylua,
		all: []Tooler{
			busted,
			docker,
			krane,
			ktech,
			ldoc,
			luacheck,
			prettier,
			stylua,
		},
	}, nil
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

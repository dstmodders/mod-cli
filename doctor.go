package main

import (
	"fmt"

	"github.com/dstmodders/mod-cli/tools"
	"github.com/fatih/color"
	"github.com/matishsiao/goInfo"
)

type Doctor struct {
	cfg  *Config
	exec *tools.Tools
}

func NewDoctor(cfg *Config) *Doctor {
	return &Doctor{
		cfg: cfg,
	}
}

func (d *Doctor) printOS() {
	printTitle("OS")

	gi, err := goInfo.GetInfo()
	if err != nil {
		return
	}

	printNameValue("GoOS", gi.GoOS)
	printNameValue("Kernel", gi.Kernel)
	printNameValue("Core", gi.Core)
	printNameValue("Platform", gi.Platform)
	printNameValue("OS", gi.OS)
	printNameValue("CPUs", gi.CPUs)
}

func (d *Doctor) printTool(tool tools.Tooler) {
	name := tool.Name()
	path := tool.Path()

	if len(path) == 0 {
		printNameValue(name, color.RedString("not found"))
		return
	}

	ver := tool.Version()
	if len(ver) == 0 {
		ver = "-"
	}

	printNameValue(name, fmt.Sprintf("%s | %s", path, ver))
}

func (d *Doctor) printTools() {
	e := d.exec

	printTitle("Tools | System")

	e.LookPaths()
	e.LoadVersions()

	d.printTool(e.Busted)
	d.printTool(e.Docker)
	d.printTool(e.LDoc)
	d.printTool(e.Luacheck)
	d.printTool(e.Prettier)
	d.printTool(e.StyLua)
	d.printTool(e.Krane)
	d.printTool(e.Ktech)
	fmt.Println()

	printTitle("Tools | Dockerized")

	e.SetToolsRunInDocker(true)
	e.LookPaths()
	e.LoadVersions()

	d.printTool(e.Busted)
	d.printTool(e.LDoc)
	d.printTool(e.Luacheck)
	d.printTool(e.Prettier)
	d.printTool(e.StyLua)
	d.printTool(e.Krane)
	d.printTool(e.Ktech)
}

func (d *Doctor) print() error {
	d.exec = tools.New()

	d.printOS()
	fmt.Println()

	d.printTools()

	return nil
}

func (d *Doctor) run() error {
	return d.print()
}

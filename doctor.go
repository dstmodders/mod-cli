package main

import (
	"fmt"

	"github.com/dstmodders/mod-cli/tools"
	"github.com/fatih/color"
	"github.com/matishsiao/goInfo"
)

type Doctor struct {
	cfg   *Config
	tools *tools.Tools
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
	t := d.tools

	printTitle("Tools | System")

	t.LookPaths()
	t.LoadVersions()

	d.printTool(t.Busted)
	d.printTool(t.Docker)
	d.printTool(t.LDoc)
	d.printTool(t.Luacheck)
	d.printTool(t.Prettier)
	d.printTool(t.StyLua)
	d.printTool(t.Krane)
	d.printTool(t.Ktech)
	fmt.Println()

	printTitle("Tools | Dockerized")

	t.SetToolsRunInDocker(true)
	t.LookPaths()
	t.LoadVersions()

	d.printTool(t.Busted)
	d.printTool(t.LDoc)
	d.printTool(t.Luacheck)
	d.printTool(t.Prettier)
	d.printTool(t.StyLua)
	d.printTool(t.Krane)
	d.printTool(t.Ktech)
}

func (d *Doctor) print() error {
	t, err := tools.New()
	if err != nil {
		return err
	}
	d.tools = t

	d.printOS()
	fmt.Println()

	d.printTools()

	return nil
}

func (d *Doctor) run() error {
	return d.print()
}

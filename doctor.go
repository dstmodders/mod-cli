package main

import (
	"github.com/matishsiao/goInfo"
)

type Doctor struct {
	cfg *Config
}

func NewDoctor(cfg *Config) *Doctor {
	return &Doctor{
		cfg: cfg,
	}
}

func printOS() {
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

func (d *Doctor) print() error {
	printOS()
	return nil
}

func (d *Doctor) run() error {
	return d.print()
}

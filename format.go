package main

import (
	"fmt"

	"github.com/dstmodders/mod-cli/tools"
	"github.com/fatih/color"
)

type Format struct {
	cfg   *Config
	tools *tools.Tools
}

func NewFormat(cfg *Config) (*Format, error) {
	t, err := tools.New()
	if err != nil {
		return nil, err
	}

	t.Prettier.SetIgnore(cfg.Format.Prettier.Ignore)
	if cfg.Format.Prettier.Docker {
		t.Prettier.SetRunInDocker(true)
	}

	t.StyLua.SetIgnore(cfg.Format.StyLua.Ignore)
	if cfg.Format.StyLua.Docker {
		t.StyLua.SetRunInDocker(true)
	}

	return &Format{
		cfg:   cfg,
		tools: t,
	}, nil
}

func (l *Format) printFormat(format tools.Format) {
	if len(format.Files) == 0 {
		fmt.Println("No issues found")
		return
	}

	for _, file := range format.Files {
		fmt.Printf("%s %s\n", color.YellowString("warning"), file.Path)
	}
}

func (l *Format) runPrettier() error {
	format, err := l.tools.Prettier.Check()
	if err != nil {
		return err
	}
	l.printFormat(format)
	return nil
}

func (l *Format) runStyLua() error {
	format, err := l.tools.StyLua.Check()
	if err != nil {
		return err
	}
	l.printFormat(format)
	return nil
}

func (l *Format) run() {
	if !l.cfg.Format.Prettier.Enabled && !l.cfg.Format.StyLua.Enabled {
		fmt.Println("")
	}

	if l.cfg.Format.Prettier.Enabled {
		printTitle("Prettier")
		_ = l.runPrettier()
	}

	if l.cfg.Format.StyLua.Enabled {
		if l.cfg.Format.Prettier.Enabled {
			fmt.Println()
		}
		printTitle("StyLua")
		_ = l.runStyLua()
	}
}

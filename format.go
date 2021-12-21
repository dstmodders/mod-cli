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

	var state string

	for _, file := range format.Files {
		switch file.State {
		case tools.FileStateSuccess:
			state = color.GreenString("success")
		default:
			state = color.YellowString("warning")
		}
		fmt.Printf("%s %s\n", state, file.Path)
	}
}

func (l *Format) runPrettier() (err error) {
	checkIfToolExists(l.tools.Docker, l.tools.Prettier)

	var format tools.Format

	if l.cfg.Format.Prettier.Fix {
		format, err = l.tools.Prettier.Fix()
	} else {
		format, err = l.tools.Prettier.Check()
	}

	if err != nil {
		return err
	}

	l.printFormat(format)
	return nil
}

func (l *Format) runStyLua() (err error) {
	checkIfToolExists(l.tools.Docker, l.tools.StyLua)

	var format tools.Format

	if l.cfg.Format.StyLua.Fix {
		format, err = l.tools.StyLua.Fix()
	} else {
		format, err = l.tools.StyLua.Check()
	}

	if err != nil {
		return err
	}

	l.printFormat(format)
	return nil
}

func (l *Format) run() {
	if !l.cfg.Format.Prettier.Enabled && !l.cfg.Format.StyLua.Enabled {
		fatalError("both Prettier and StyLua are disabled. Enable at least one of them first")
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

package main

import (
	"errors"
	"fmt"

	"github.com/dstmodders/mod-cli/tools"
	"github.com/fatih/color"
)

type Format struct {
	canRunPrettier bool
	canRunStyLua   bool
	cfg            *Config
	tools          *tools.Tools
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

func (f *Format) checkTools() {
	var errPrettier, errStyLua error

	err := errors.New("neither Prettier nor StyLua are available")

	if !f.cfg.Format.Prettier.Enabled && !f.cfg.Format.StyLua.Enabled {
		fatalError("both Prettier and StyLua are disabled. Enable at least one of them first")
	}

	if f.cfg.Format.Prettier.Enabled {
		errPrettier = checkIfToolExists(f.tools.Docker, f.tools.Prettier)
		if errPrettier == nil {
			f.canRunPrettier = true
			err = nil
		}
	}

	if f.cfg.Format.StyLua.Enabled {
		errStyLua = checkIfToolExists(f.tools.Docker, f.tools.StyLua)
		if errStyLua == nil {
			f.canRunStyLua = true
			err = nil
		}
	}

	if f.canRunStyLua || f.canRunPrettier {
		if errPrettier != nil {
			printWarning(errPrettier)
		}

		if errStyLua != nil {
			printWarning(errStyLua)
		}
	}

	if err != nil {
		fatalError(err)
	}
}

func (f *Format) printFormat(format tools.Format) {
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

func (f *Format) runPrettier() (err error) {
	var format tools.Format

	if f.cfg.Format.Prettier.Fix {
		format, err = f.tools.Prettier.Fix()
	} else {
		format, err = f.tools.Prettier.Check()
	}

	if err != nil {
		return err
	}

	printTitle("Prettier")
	f.printFormat(format)
	return nil
}

func (f *Format) runStyLua() (err error) {
	var format tools.Format

	if f.cfg.Format.StyLua.Fix {
		format, err = f.tools.StyLua.Fix()
	} else {
		format, err = f.tools.StyLua.Check()
	}

	if err != nil {
		return err
	}

	printTitle("StyLua")
	f.printFormat(format)
	return nil
}

func (f *Format) run() {
	f.checkTools()

	if f.canRunPrettier && f.cfg.Format.Prettier.Enabled {
		_ = f.runPrettier()
	}

	if f.canRunStyLua && f.cfg.Format.StyLua.Enabled {
		if f.canRunPrettier && f.cfg.Format.Prettier.Enabled {
			fmt.Println()
		}
		_ = f.runStyLua()
	}
}

package main

import (
	"errors"
	"fmt"

	"github.com/dstmodders/mod-cli/tools"
	"github.com/fatih/color"
)

type Lint struct {
	Full           bool
	canRunLuacheck bool
	cfg            *Config
	tools          *tools.Tools
}

func NewLint(cfg *Config) (*Lint, error) {
	t, err := tools.New()
	if err != nil {
		return nil, err
	}

	t.Luacheck.SetIgnore(cfg.Lint.Luacheck.Ignore)
	if cfg.Lint.Luacheck.Docker {
		t.Luacheck.SetRunInDocker(true)
	}

	return &Lint{
		cfg:   cfg,
		tools: t,
	}, nil
}

func (l *Lint) checkTools() {
	var errLuacheck error

	//nolint:stylecheck
	//goland:noinspection ALL
	err := errors.New("Luacheck is not available")

	if !l.cfg.Lint.Luacheck.Enabled {
		fatalError("Luacheck is disabled. Enable it first")
	}

	if l.cfg.Lint.Luacheck.Enabled {
		errLuacheck = checkIfToolExists(l.tools.Docker, l.tools.Luacheck)
		if errLuacheck == nil {
			l.canRunLuacheck = true
			err = nil
		}
	}

	if l.canRunLuacheck {
		if errLuacheck != nil {
			printWarning(errLuacheck)
		}
	}

	if err != nil {
		fatalError(err)
	}
}

func (l *Lint) printLint(lint tools.Lint) {
	if l.Full {
		fmt.Println(lint.Stdout)
		return
	}

	if len(lint.Files) == 0 {
		fmt.Println("No issues found")
		return
	}

	total := len(lint.Files)
	for i := 0; i < total; i++ {
		file := lint.Files[i]

		issues := "issues"
		if len(file.Issues) == 1 {
			issues = "issue"
		}
		issues = fmt.Sprintf("%d %s", len(file.Issues), issues)

		fmt.Printf(
			"%s %s %s\n",
			color.YellowString("warning"),
			file.Path,
			color.YellowString(issues),
		)

		fmt.Println()
		for _, issue := range file.Issues {
			fmt.Printf(
				"    %s:%d:%d: %s\n",
				issue.Name,
				issue.StartLine,
				issue.EndLine,
				issue.Description,
			)
		}

		if i < (total - 1) {
			fmt.Println()
		}
	}
}

func (l *Lint) runLuacheck() error {
	lint, err := l.tools.Luacheck.Lint()
	if err != nil {
		return err
	}

	printTitle("Luacheck")
	l.printLint(lint)
	return nil
}

func (l *Lint) run() {
	l.checkTools()

	if l.canRunLuacheck && l.cfg.Lint.Luacheck.Enabled {
		_ = l.runLuacheck()
	}
}

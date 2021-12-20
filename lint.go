package main

import (
	"fmt"

	"github.com/dstmodders/mod-cli/tools"
	"github.com/fatih/color"
)

type Lint struct {
	cfg   *Config
	tools *tools.Tools
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

func (l *Lint) printLint(lint tools.Lint) {
	if len(lint.Files) == 0 {
		fmt.Println("No issues found")
		return
	}

	for _, file := range lint.Files {
		issues := "issues"
		if file.Issues == 1 {
			issues = "issue"
		}
		issues = fmt.Sprintf("%d %s", file.Issues, issues)
		fmt.Printf(
			"%s %s %s\n",
			color.YellowString("warning"),
			file.Path,
			color.YellowString(issues),
		)
	}
}

func (l *Lint) runLuacheck() error {
	lint, err := l.tools.Luacheck.Lint()
	if err != nil {
		return err
	}
	l.printLint(lint)
	return nil
}

func (l *Lint) run() {
	if l.cfg.Lint.Luacheck.Enabled {
		printTitle("Luacheck")
		_ = l.runLuacheck()
	}
}

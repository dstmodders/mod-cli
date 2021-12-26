package main

import (
	"errors"

	"github.com/dstmodders/mod-cli/tools"
)

type Test struct {
	canRunBusted bool
	cfg          *Config
	tools        *tools.Tools
}

func NewTest(cfg *Config) (*Test, error) {
	t, err := tools.New()
	if err != nil {
		return nil, err
	}

	return &Test{
		cfg:   cfg,
		tools: t,
	}, nil
}

func (t *Test) checkTools() {
	var errBusted error

	//goland:noinspection ALL
	err := errors.New("Busted is not available")

	errBusted = checkIfToolExists(t.tools.Docker, t.tools.Busted)
	if errBusted == nil {
		t.canRunBusted = true
		err = nil
	}

	if t.canRunBusted {
		if errBusted != nil {
			printWarning(errBusted)
		}
	}

	if err != nil {
		fatalError(err)
	}
}

func (t *Test) runBusted() error {
	_, err := t.tools.Busted.Test()
	if err != nil {
		return err
	}
	return nil
}

func (t *Test) run() {
	t.checkTools()
	_ = t.runBusted()
}

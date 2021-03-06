package main

import (
	"fmt"
	"path"

	"github.com/dstmodders/mod-cli/workshop"
	"github.com/manifoldco/promptui"
)

type Workshop struct {
	cfg      *Config
	destName string
	list     bool
	mod      *Mod
	path     string
	workshop *workshop.Workshop
	zip      bool
}

func NewWorkshop(cfg *Config) *Workshop {
	return &Workshop{
		cfg: cfg,
	}
}

func (w *Workshop) printInfo() {
	printTitle("Info")
	printNameValue("Name", w.mod.name)
	printNameValue("Version", w.mod.version)
}

func (w *Workshop) printPaths() {
	printTitle("Paths")
	printNameValue("Source", w.workshop.AbsSrcPath())
	printNameValue("Destination", w.workshop.AbsDestPath())
}

func (w *Workshop) printFiles() {
	printTitle(fmt.Sprintf("Files | Total: %d", len(w.workshop.Files())))
	w.workshop.PrintFiles()
}

func (w *Workshop) printDefault() error {
	w.printInfo()
	fmt.Println()

	w.printPaths()
	fmt.Println()

	w.printFiles()
	return nil
}

func (w *Workshop) copy() error {
	total, err := w.workshop.CountDestItems()
	if err != nil {
		return err
	}

	if total > 0 {
		prompt := promptui.Prompt{
			Label:     "Destination directory already exists. Override",
			Default:   "y",
			IsConfirm: true,
		}

		_, err := prompt.Run()
		if err != nil {
			return nil
		}
	}

	if w.zip {
		if err := w.workshop.ZipFiles(); err != nil {
			return err
		}
	} else {
		if err := w.workshop.CopyFiles(); err != nil {
			return err
		}
	}

	fmt.Println("Done")
	return nil
}

func (w *Workshop) run() error {
	w.mod = NewMod()
	if err := w.mod.Load(w.path); err != nil {
		return err
	}

	ws, err := workshop.New(w.mod.pathAbs, path.Join(w.mod.pathAbs, w.destName))
	if err != nil {
		return err
	}

	ignore := w.cfg.Workshop.Ignore
	ignore = append(ignore, w.destName)
	ws.SetIgnore(ignore)

	w.workshop = ws

	if _, _, err := ws.GetFiles(); err != nil {
		return err
	}

	if w.list {
		w.workshop.PrintFiles()
		return nil
	}

	if err := w.printDefault(); err != nil {
		return err
	}

	fmt.Println("---")
	if err := w.copy(); err != nil {
		return err
	}

	return nil
}

package main

import (
	"errors"
	"os"
	"path"
	"path/filepath"

	"github.com/dstmodders/mod-cli/modinfo"
)

type Mod struct {
	modinfo *modinfo.ModInfo
	name    string
	path    string
	pathAbs string
	version string
}

func NewMod() *Mod {
	return &Mod{}
}

func (m *Mod) Load(modPath string) error {
	stat, err := os.Stat(modPath)
	if err != nil {
		return err
	}

	if !stat.IsDir() {
		return errors.New("not a directory")
	}

	absPath, err := filepath.Abs(stat.Name())
	if err != nil {
		return err
	}

	info := modinfo.New()
	if err := info.Load(path.Join(absPath, "modinfo.lua")); err != nil {
		return err
	}

	name := info.FieldByName("name")
	if name == nil {
		return errors.New("mod info field name doesn't have any value")
	}

	v := info.FieldByName("version")
	if v == nil {
		return errors.New("mod info field version doesn't have any value")
	}

	m.modinfo = info
	m.name = name.String()
	m.path = modPath
	m.pathAbs = absPath
	m.version = v.String()

	return nil
}

package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/dstmodders/mod-cli/cmd/changelog"
	"github.com/dstmodders/mod-cli/cmd/info"
	"github.com/dstmodders/mod-cli/config"
	"github.com/dstmodders/mod-cli/modinfo"
	"github.com/yuin/gopher-lua"
	"gopkg.in/alecthomas/kingpin.v2"
)

const (
	Version = "0.1"
)

var (
	// kingpin
	app = kingpin.New(
		filepath.Base(os.Args[0]),
		"Different modding tools by Depressed DST Modders for Klei's game Don't Starve Together.",
	)

	cfg = config.New()

	appConfig = app.Flag("config", "Path to configurations file.").Short('c').Default(".modcli").String()

	changelogCmd             = app.Command("changelog", "Changelog tools.")
	changelogCmdPath         = changelogCmd.Arg("path", "Path.").Default("CHANGELOG.md").String()
	changelogCmdCount        = changelogCmd.Flag("count", "Show total number of releases.").Bool()
	changelogCmdFirst        = changelogCmd.Flag("first", "Show first release.").Short('f').Bool()
	changelogCmdLatest       = changelogCmd.Flag("latest", "Show latest release.").Short('l').Bool()
	changelogCmdList         = changelogCmd.Flag("list", "Show list of releases without changes.").Bool()
	changelogCmdListVersions = changelogCmd.Flag("list-versions", "Show list of versions.").Bool()

	infoCmd                      = app.Command("info", "Mod info tools.")
	infoCmdPath                  = infoCmd.Arg("path", "Path to modinfo.lua.").Default("modinfo.lua").String()
	infoCmdCompatability         = infoCmd.Flag("compatibility", "Show compatability fields.").Bool()
	infoCmdConfiguration         = infoCmd.Flag("configuration", "Show configuration options with their default values.").Bool()
	infoCmdConfigurationMarkdown = infoCmd.Flag("configuration-markdown", "Show configuration options with their default values as a Markdown table.").Short('m').Bool()
	infoCmdDescription           = infoCmd.Flag("description", "Show description.").Short('d').Bool()
	infoCmdFirstLine             = infoCmd.Flag("first-line", "Show first lines for values.").Short('f').Bool()
	infoCmdGeneral               = infoCmd.Flag("general", "Show general fields.").Short('g').Bool()
	infoCmdNames                 = infoCmd.Flag("names", "Show variable names or options data instead of their descriptions.").Short('n').Bool()
	infoCmdOther                 = infoCmd.Flag("other", "Show other fields.").Short('o').Bool()
)

func fatalError(msg string, args ...interface{}) {
	if len(args) > 0 {
		fmt.Printf("[error] %s (%s)\n", msg, args[0].(error).Error())
	} else {
		fmt.Printf("[error] %s\n", msg)
	}
	os.Exit(1)
}

func loadConfig() {
	if len(*appConfig) > 0 && *appConfig != ".modcli" {
		errMsg := "failed to load config"
		stat, err := os.Stat(*appConfig)

		if errors.Is(err, os.ErrNotExist) {
			fatalError(
				errMsg,
				errors.New(fmt.Sprintf("open %s: no such file", *appConfig)),
			)
		}

		if stat.IsDir() {
			fatalError(
				errMsg,
				errors.New(fmt.Sprintf("open %s: expected file but got directory", *appConfig)),
			)
		}
	}

	if _, err := os.Stat(*appConfig); err == nil {
		errMsg := "failed to load config"

		f, err := os.OpenFile(*appConfig, os.O_RDONLY, 0644)
		if err != nil {
			fatalError(errMsg, err)
		}

		if err := cfg.Load(f); err != nil {
			fatalError(errMsg, err)
		}
	}
}

func runChangelog() error {
	path := *changelogCmdPath

	c := changelog.New()
	c.Count = *changelogCmdCount
	c.First = *changelogCmdFirst
	c.Latest = *changelogCmdLatest
	c.List = *changelogCmdList
	c.ListVersions = *changelogCmdListVersions

	if err := c.Load(path); err != nil {
		return err
	}

	if err := c.Print(); err != nil {
		return err
	}

	return nil
}

func runInfo() error {
	path := *infoCmdPath

	L := lua.NewState()
	defer L.Close()
	if err := L.DoFile(path); err != nil {
		return err
	}

	m := modinfo.New()
	if err := m.Load(path); err != nil {
		return err
	}

	i := info.New(m)
	i.Compatability = *infoCmdCompatability
	i.Configuration = *infoCmdConfiguration
	i.ConfigurationMarkdown = *infoCmdConfigurationMarkdown
	i.Description = *infoCmdDescription
	i.FirstLine = *infoCmdFirstLine
	i.General = *infoCmdGeneral
	i.Names = *infoCmdNames
	i.Other = *infoCmdOther

	if err := i.Print(); err != nil {
		return err
	}

	return nil
}

func init() {
	// kingpin
	app.UsageTemplate(kingpin.DefaultUsageTemplate).Version(Version)
	app.HelpFlag.Short('h')
	app.VersionFlag.Short('v')
}

func main() {
	command, err := app.Parse(os.Args[1:])
	if err != nil {
		fatalError("failed to parse arguments", err)
	}

	loadConfig()

	// commands
	switch kingpin.MustParse(command, err) {
	case changelogCmd.FullCommand():
		if err := runChangelog(); err != nil {
			fatalError("failed to run changelog command", err)
		}
	case infoCmd.FullCommand():
		if err := runInfo(); err != nil {
			fatalError("failed to run info command", err)
		}
	}
}

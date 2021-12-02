package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/alecthomas/kingpin.v2"
)

var version string

var (
	// kingpin
	app = kingpin.New(
		filepath.Base(os.Args[0]),
		"Different modding tools by Depressed DST Modders for Klei's game Don't Starve Together.",
	)

	cfg = NewConfig()

	appConfig = app.Flag("config", "Path to configuration file.").Short('c').Default(".modcli").String()

	changelogCmd             = app.Command("changelog", "Changelog tools.")
	changelogCmdPath         = changelogCmd.Arg("path", "Path.").Default("CHANGELOG.md").String()
	changelogCmdCount        = changelogCmd.Flag("count", "Show total number of releases.").Bool()
	changelogCmdFirst        = changelogCmd.Flag("first", "Show first release.").Short('f').Bool()
	changelogCmdLatest       = changelogCmd.Flag("latest", "Show latest release.").Short('l').Bool()
	changelogCmdList         = changelogCmd.Flag("list", "Show list of releases without changes.").Bool()
	changelogCmdListVersions = changelogCmd.Flag("list-versions", "Show list of versions.").Bool()

	doctorCmd = app.Command("doctor", "Check health of this CLI app.")

	infoCmd                      = app.Command("info", "Mod info tools.")
	infoCmdPath                  = infoCmd.Arg("path", "Path to modinfo.lua.").Default("modinfo.lua").String()
	infoCmdCompatibility         = infoCmd.Flag("compatibility", "Show compatibility fields.").Bool()
	infoCmdConfiguration         = infoCmd.Flag("configuration", "Show configuration options with their default values.").Bool()
	infoCmdConfigurationMarkdown = infoCmd.Flag("configuration-markdown", "Show configuration options with their default values as a Markdown table.").Short('m').Bool()
	infoCmdDescription           = infoCmd.Flag("description", "Show description.").Short('d').Bool()
	infoCmdFirstLine             = infoCmd.Flag("first-line", "Show first lines for values.").Short('f').Bool()
	infoCmdGeneral               = infoCmd.Flag("general", "Show general fields.").Short('g').Bool()
	infoCmdNames                 = infoCmd.Flag("names", "Show variable names or options data instead of their descriptions.").Short('n').Bool()
	infoCmdOther                 = infoCmd.Flag("other", "Show other fields.").Short('o').Bool()

	workshopCmd     = app.Command("workshop", "Steam Workshop tools.")
	workshopCmdPath = workshopCmd.Arg("path", "Path to mod directory.").Default(".").ExistingDir()
	workshopCmdList = workshopCmd.Flag("list", "Show only files that are going to be included.").Short('l').Bool()
	workshopCmdName = workshopCmd.Flag("name", "Name of destination directory/archive.").Default("workshop").Short('n').String()
	workshopCmdZip  = workshopCmd.Flag("zip", "Create a ZIP archive instead.").Short('z').Bool()
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
				fmt.Errorf("open %s: no such file", *appConfig),
			)
		}

		if stat.IsDir() {
			fatalError(
				errMsg,
				fmt.Errorf("open %s: expected file but got directory", *appConfig),
			)
		}
	}

	if _, err := os.Stat(*appConfig); err == nil {
		errMsg := "failed to load config"

		f, err := os.OpenFile(*appConfig, os.O_RDONLY, 0644)
		if err != nil {
			fatalError(errMsg, err)
		}

		if err := cfg.load(f); err != nil {
			fatalError(errMsg, err)
		}
	}
}

func runChangelog() {
	c := NewChangelog()
	c.Count = *changelogCmdCount
	c.First = *changelogCmdFirst
	c.Latest = *changelogCmdLatest
	c.List = *changelogCmdList
	c.ListVersions = *changelogCmdListVersions

	if err := c.run(*changelogCmdPath); err != nil {
		fatalError("failed to run changelog command", err)
	}
}

func runDoctor() {
	d := NewDoctor(cfg)

	if err := d.run(); err != nil {
		fatalError("failed to run doctor command", err)
	}
}

func runInfo() {
	i := NewInfo()
	i.Compatibility = *infoCmdCompatibility
	i.Configuration = *infoCmdConfiguration
	i.ConfigurationMarkdown = *infoCmdConfigurationMarkdown
	i.Description = *infoCmdDescription
	i.FirstLine = *infoCmdFirstLine
	i.General = *infoCmdGeneral
	i.Names = *infoCmdNames
	i.Other = *infoCmdOther

	if err := i.run(*infoCmdPath); err != nil {
		fatalError("failed to run info command", err)
	}
}

func runWorkshop() {
	w := NewWorkshop(cfg)
	w.destName = *workshopCmdName
	w.list = *workshopCmdList
	w.path = *workshopCmdPath
	w.zip = *workshopCmdZip

	if err := w.run(); err != nil {
		fatalError("failed to run workshop command", err)
	}
}

func init() {
	// kingpin
	app.UsageTemplate(kingpin.DefaultUsageTemplate).Version(version)
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
		runChangelog()
	case doctorCmd.FullCommand():
		runDoctor()
	case infoCmd.FullCommand():
		runInfo()
	case workshopCmd.FullCommand():
		runWorkshop()
	}
}

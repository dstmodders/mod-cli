package main

import (
	"fmt"
	"strings"

	md "github.com/fbiville/markdown-table-formatter/pkg/markdown"
	lua "github.com/yuin/gopher-lua"

	"github.com/dstmodders/mod-cli/modinfo"
)

type Info struct {
	Compatibility         bool
	Configuration         bool
	ConfigurationMarkdown bool
	Description           bool
	FirstLine             bool
	General               bool
	Names                 bool
	Other                 bool
	modinfo               *modinfo.ModInfo
}

func NewInfo() *Info {
	return &Info{
		Compatibility: true,
		Configuration: true,
		Description:   true,
		General:       true,
		Other:         true,
	}
}

func (i *Info) printGlobal(name string) {
	g := i.modinfo.FieldByName(name)
	n := g.Description
	if i.Names {
		n = g.Name
	}
	fmt.Printf("%s: %s\n", n, g.String())
}

func (i *Info) printGeneral() {
	i.printGlobal("name")
	i.printGlobal("author")
	i.printGlobal("version")
	i.printGlobal("api_version")
}

func (i *Info) printDescription() {
	fmt.Println(i.modinfo.FieldByName("description"))
}

func (i *Info) printCompatibility() {
	i.printGlobal("dont_starve_compatible")
	i.printGlobal("dst_compatible")
	i.printGlobal("reign_of_giants_compatible")
}

func (i *Info) printConfigurationOption(option modinfo.Option) {
	if option.Hover != "" {
		n := option.Label
		if i.Names {
			if (!i.General && !i.Description && !i.Compatibility && !i.Configuration && !i.Other) ||
				(i.General || i.Description || i.Compatibility || i.Other) {
				n = "configuration_options." + option.Name
			} else {
				n = option.Name
			}
			fmt.Printf("%s: %s\n", n, option.Default.Data)
			return
		}
		fmt.Printf("%s: %s\n", n, option.Default.Description)
	}
}

func (i *Info) printConfiguration() {
	for _, option := range i.modinfo.ConfigurationOptions.Values {
		i.printConfigurationOption(option)
	}
}

func (i *Info) printConfigurationMarkdown() error {
	var data [][]string

	for _, option := range i.modinfo.ConfigurationOptions.Values {
		if option.Hover != "" {
			label := option.Label
			defaultStr := option.Default.String()
			if i.Names {
				label = option.Name
				defaultStr = option.Default.DataString()
			}

			hoverSlice := strings.Split(option.Hover, "\n")
			hover := strings.TrimSpace(hoverSlice[0])
			if !i.FirstLine {
				hover = strings.TrimSpace(strings.Join(hoverSlice, "<br />"))
			}

			data = append(data, []string{
				fmt.Sprintf("**%s**", label),
				fmt.Sprintf("_%s_", defaultStr),
				hover,
			})
		}
	}

	basicTable, err := md.NewTableFormatterBuilder().
		WithPrettyPrint().
		Build("Configuration", "Default", "Description").
		Format(data)

	if err != nil {
		return err
	}

	fmt.Println(strings.TrimSpace(basicTable))
	return nil
}

func (i *Info) printOther() {
	i.printGlobal("icon")
	i.printGlobal("icon_atlas")
	i.printGlobal("forum_thread")
	i.printGlobal("priority")
	i.printGlobal("folder_name")
}

func (i *Info) print() error { //nolint:funlen,gocyclo
	if i.ConfigurationMarkdown {
		if err := i.printConfigurationMarkdown(); err != nil {
			return err
		}
		return nil
	}

	if i.General {
		if i.Description || i.Compatibility || i.Configuration || i.Other {
			printTitle("General")
		}
		i.printGeneral()
	}

	if i.Description {
		if i.General || i.Compatibility || i.Configuration || i.Other {
			if i.General {
				fmt.Println()
			}
			printTitle("Description")
		}
		i.printDescription()
	}

	if i.Compatibility {
		if i.General || i.Description || i.Configuration || i.Other {
			if i.General || i.Description {
				fmt.Println()
			}
			printTitle("Compatibility")
		}
		i.printCompatibility()
	}

	if i.Configuration {
		if i.General || i.Description || i.Compatibility || i.Other {
			if i.General || i.Description || i.Compatibility {
				fmt.Println()
			}
			printTitle("Configuration")
		}
		i.printConfiguration()
	}

	if i.Other {
		if i.General || i.Description || i.Compatibility || i.Configuration {
			if i.General || i.Description || i.Compatibility || i.Configuration {
				fmt.Println()
			}
			printTitle("Other")
		}
		i.printOther()
	}

	if i.General || i.Description || i.Compatibility || i.Configuration || i.Other {
		return nil
	}

	printTitle("General")
	i.printGeneral()
	fmt.Println()

	printTitle("Description")
	i.printDescription()
	fmt.Println()

	printTitle("Compatibility")
	i.printCompatibility()
	fmt.Println()

	printTitle("Configuration")
	i.printConfiguration()
	fmt.Println()

	printTitle("Other")
	i.printOther()

	return nil
}

func (i *Info) run(path string) error {
	L := lua.NewState()
	defer L.Close()
	if err := L.DoFile(path); err != nil {
		return err
	}

	m := modinfo.New()
	if err := m.Load(path); err != nil {
		return err
	}

	i.modinfo = m

	if err := i.print(); err != nil {
		return err
	}

	return nil
}

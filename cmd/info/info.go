package info

import (
	"fmt"
	"strings"

	md "github.com/fbiville/markdown-table-formatter/pkg/markdown"

	"github.com/dstmodders/mod-cli/modinfo"
)

type Info struct {
	Compatability         bool
	Configuration         bool
	ConfigurationMarkdown bool
	Description           bool
	FirstLine             bool
	General               bool
	Names                 bool
	Other                 bool
	modinfo               *modinfo.ModInfo
}

func NewInfo(modinfo *modinfo.ModInfo) *Info {
	return &Info{
		Compatability: true,
		Configuration: true,
		Description:   true,
		General:       true,
		Other:         true,
		modinfo:       modinfo,
	}
}

func (i *Info) PrintGlobal(name string) {
	g := i.modinfo.FieldByName(name)
	n := g.Description
	if i.Names {
		n = g.Name
	}
	fmt.Printf("%s: %s\n", n, g.String())
}

func (i *Info) PrintConfigurationOption(option modinfo.Option) {
	if option.Hover != "" {
		n := option.Label
		if i.Names {
			if (!i.General && !i.Description && !i.Compatability && !i.Configuration && !i.Other) ||
				(i.General || i.Description || i.Compatability || i.Other) {
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

func (i *Info) PrintGeneral() {
	i.PrintGlobal("name")
	i.PrintGlobal("author")
	i.PrintGlobal("version")
	i.PrintGlobal("api_version")
}

func (i *Info) PrintDescription() {
	fmt.Println(i.modinfo.FieldByName("description"))
}

func (i *Info) PrintCompatability() {
	i.PrintGlobal("dont_starve_compatible")
	i.PrintGlobal("dst_compatible")
	i.PrintGlobal("reign_of_giants_compatible")
}

func (i *Info) PrintConfiguration() {
	for _, option := range i.modinfo.ConfigurationOptions.Values {
		i.PrintConfigurationOption(option)
	}
}

func (i *Info) PrintConfigurationMarkdown() error {
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

func (i *Info) PrintOther() {
	i.PrintGlobal("icon")
	i.PrintGlobal("icon_atlas")
	i.PrintGlobal("forum_thread")
	i.PrintGlobal("priority")
	i.PrintGlobal("folder_name")
}

func (i *Info) Print() error {
	if i.ConfigurationMarkdown {
		if err := i.PrintConfigurationMarkdown(); err != nil {
			return err
		}
		return nil
	}

	if i.General {
		if i.Description || i.Compatability || i.Configuration || i.Other {
			fmt.Printf("[GENERAL]\n\n")
		}
		i.PrintGeneral()
	}

	if i.Description {
		if i.General || i.Compatability || i.Configuration || i.Other {
			if i.General {
				fmt.Println()
			}
			fmt.Printf("[DESCRIPTION]\n\n")
		}
		i.PrintDescription()
	}

	if i.Compatability {
		if i.General || i.Description || i.Configuration || i.Other {
			if i.General || i.Description {
				fmt.Println()
			}
			fmt.Printf("[COMPATABILITY]\n\n")
		}
		i.PrintCompatability()
	}

	if i.Configuration {
		if i.General || i.Description || i.Compatability || i.Other {
			if i.General || i.Description || i.Compatability {
				fmt.Println()
			}
			fmt.Printf("[CONFIGURATION]\n\n")
		}
		i.PrintConfiguration()
	}

	if i.Other {
		if i.General || i.Description || i.Compatability || i.Configuration {
			if i.General || i.Description || i.Compatability || i.Configuration {
				fmt.Println()
			}
			fmt.Printf("[OTHER]\n\n")
		}
		i.PrintOther()
	}

	if i.General || i.Description || i.Compatability || i.Configuration || i.Other {
		return nil
	}

	fmt.Printf("[GENERAL]\n\n")
	i.PrintGeneral()
	fmt.Println()

	fmt.Printf("[DESCRIPTION]\n\n")
	i.PrintDescription()
	fmt.Println()

	fmt.Printf("[COMPATABILITY]\n\n")
	i.PrintCompatability()
	fmt.Println()

	fmt.Printf("[CONFIGURATION]\n\n")
	i.PrintConfiguration()
	fmt.Println()

	fmt.Printf("[OTHER]\n\n")
	i.PrintOther()

	return nil
}

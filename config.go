package main

import (
	"errors"
	"fmt"
	"os"
	"reflect"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Format   ConfigFormat
	Lint     ConfigLint
	Workshop ConfigWorkshop
	file     *os.File
	yaml     ConfigYAML
}

type ConfigFormat struct {
	Prettier ConfigTool
	StyLua   ConfigTool
}

type ConfigLint struct {
	Luacheck ConfigTool
}

type ConfigTool struct {
	Docker  bool
	Enabled bool
	Fix     bool
	Ignore  []string
}

type ConfigWorkshop struct {
	Ignore []string
}

func NewConfig() *Config {
	tool := ConfigTool{
		Ignore: []string{
			".idea/",
		},
		Enabled: true,
	}

	return &Config{
		Format: ConfigFormat{
			Prettier: tool,
			StyLua:   tool,
		},
		Lint: ConfigLint{
			Luacheck: ConfigTool{
				Enabled: true,
			},
		},
		Workshop: ConfigWorkshop{
			Ignore: []string{
				".*",
				"CHANGELOG.md",
				"Makefile",
				"README.md",
				"codecov.yml",
				"config.ld",
				"docs/",
				"lcov.info",
				"luacov.*",
				"modicon.png",
				"preview.gif",
				"readme/",
				"spec/",
			},
		},
	}
}

func (c *Config) errorValue(name string, args ...interface{}) error {
	errStr := fmt.Sprintf("invalid YAML %s value", name)
	if len(args) > 0 {
		return fmt.Errorf("%s: %s", errStr, args[0].(string))
	}
	return errors.New(errStr)
}

func (c *Config) errorExpected(name, expected string, value interface{}) error {
	t := reflect.TypeOf(value).String()
	switch value.(type) {
	case []interface{}:
		t = "sequence"
	case map[interface{}]interface{}:
		t = "mapping"
	}
	expected = fmt.Sprintf("expected %s but got %s", expected, t)
	return c.errorValue(name, expected)
}

func (c *Config) toBool(name string, value interface{}, dest *bool) error {
	expected := "bool"
	switch val := value.(type) {
	case bool:
		*dest = val
		return nil
	default:
		return c.errorExpected(name, expected, value)
	}
}

func (c *Config) toSequence(name string, value interface{}, dest *[]string) error {
	switch val := value.(type) {
	case []interface{}:
		if len(val) > 0 {
			*dest = []string{}
			for _, str := range val {
				*dest = append(*dest, str.(string))
			}
		}
		return nil
	case nil:
		return nil
	default:
		return c.errorExpected(name, "null or sequence", value)
	}
}

func (c *Config) parseYAMLTool(name string, value interface{}, dest *ConfigTool) error {
	switch val := value.(type) {
	case map[interface{}]interface{}:
		dest.Enabled = true

		if val["docker"] != nil {
			if err := c.toBool(name, val["docker"], &dest.Docker); err != nil {
				return err
			}
		}

		if val["fix"] != nil {
			if err := c.toBool(name, val["fix"], &dest.Fix); err != nil {
				return err
			}
		}

		if val["ignore"] != nil {
			if err := c.toSequence(name, val["ignore"], &dest.Ignore); err != nil {
				return err
			}
		}

		return nil
	case bool:
		if err := c.toBool(name, val, &dest.Enabled); err != nil {
			return err
		}
		return nil
	case nil:
		return nil
	default:
		return c.errorExpected(name, "bool or mapping", value)
	}
}

func (c *Config) parseYAMLFormat() error {
	switch val := c.yaml.Format.(type) {
	case map[interface{}]interface{}:
		if err := c.parseYAMLTool("format.prettier", val["prettier"], &c.Format.Prettier); err != nil {
			return err
		}

		if err := c.parseYAMLTool("format.stylua", val["stylua"], &c.Format.StyLua); err != nil {
			return err
		}

		return nil
	case nil:
		return nil
	default:
		return c.errorExpected("format", "mapping", c.yaml.Format)
	}
}

func (c *Config) parseYAMLLint() error {
	switch val := c.yaml.Lint.(type) {
	case map[interface{}]interface{}:
		if err := c.parseYAMLTool("lint.luacheck", val["luacheck"], &c.Lint.Luacheck); err != nil {
			return err
		}
		return nil
	case nil:
		return nil
	default:
		return c.errorExpected("lint", "mapping", c.yaml.Lint)
	}
}

func (c *Config) parseYAMLWorkshop() error {
	switch val := c.yaml.Workshop.(type) {
	case map[interface{}]interface{}:
		if err := c.toSequence("workshop.ignore", val["ignore"], &c.Workshop.Ignore); err != nil {
			return err
		}
		return nil
	case nil:
		return nil
	default:
		return c.errorExpected("workshop", "mapping", c.yaml.Workshop)
	}
}

func (c *Config) load(file *os.File) error {
	yml := NewYAML()
	if err := yml.unmarshalFile(file); err != nil {
		return errors.New("not YAML format")
	}

	c.file = file
	c.yaml = *yml

	if err := c.parseYAMLFormat(); err != nil {
		return err
	}

	if err := c.parseYAMLLint(); err != nil {
		return err
	}

	if err := c.parseYAMLWorkshop(); err != nil {
		return err
	}

	return nil
}

type ConfigYAML struct {
	Format   interface{} `yaml:"format"`
	Lint     interface{} `yaml:"lint"`
	Workshop interface{} `yaml:"workshop"`
}

func NewYAML() *ConfigYAML {
	return &ConfigYAML{}
}

func (c *ConfigYAML) unmarshalFile(file *os.File) error {
	stat, err := file.Stat()
	if err != nil {
		return err
	}

	in := make([]byte, stat.Size())
	_, err = file.Read(in)
	if err != nil {
		return err
	}

	return yaml.Unmarshal(in, &c)
}

package main

import (
	"errors"
	"fmt"
	"os"
	"reflect"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Workshop ConfigWorkshop
	file     *os.File
	yaml     ConfigYAML
}

type ConfigWorkshop struct {
	Ignore []string
}

func NewConfig() *Config {
	return &Config{
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

	if err := c.parseYAMLWorkshop(); err != nil {
		return err
	}

	return nil
}

type ConfigYAML struct {
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

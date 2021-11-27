package config

import (
	"errors"
	"fmt"
	"os"
	"reflect"
)

type Controller interface {
	Load(path string) error
}

type Config struct {
	Workshop Workshop
	file     *os.File
	yaml     YAML
}

func New() *Config {
	return &Config{
		Workshop: Workshop{},
	}
}

func invalidValueError(name string, args ...interface{}) error {
	errStr := fmt.Sprintf("invalid YAML %s value", name)
	if len(args) > 0 {
		return errors.New(fmt.Sprintf("%s: %s", errStr, args[0].(string)))
	}
	return errors.New(errStr)
}

func (c *Config) parseYAMLIgnore(name string, value map[interface{}]interface{}) error {
	expectedStr := "expected null or sequence"
	for k, v := range value {
		if k == "ignore" {
			switch v.(type) {
			case []interface{}:
				w := v.([]interface{})
				for _, str := range w {
					c.Workshop.Ignore = append(c.Workshop.Ignore, str.(string))
				}
				return nil
			case map[interface{}]interface{}:
				expectedStr = fmt.Sprintf("%s but got mapping", expectedStr)
				return invalidValueError(name, expectedStr)
			case nil:
				return nil
			default:
				expectedStr = fmt.Sprintf("%s but got %s", expectedStr, reflect.TypeOf(v))
				return invalidValueError(name, expectedStr)
			}
		}
	}
	return nil
}

func (c *Config) parseYAMLWorkshop() error {
	switch c.yaml.Workshop.(type) {
	case map[interface{}]interface{}:
		return c.parseYAMLIgnore(
			"workshop.ignore",
			c.yaml.Workshop.(map[interface{}]interface{}),
		)
	case nil:
		return nil
	default:
		return invalidValueError("workshop", fmt.Sprintf(
			"expected mapping but got %s",
			reflect.TypeOf(c.yaml.Workshop),
		))
	}
}

func (c *Config) Load(file *os.File) error {
	yml := NewYAML()
	if err := yml.UnmarshalFile(file); err != nil {
		return errors.New("not YAML format")
	}

	c.file = file
	c.yaml = *yml

	if err := c.parseYAMLWorkshop(); err != nil {
		return err
	}

	return nil
}

type Workshop struct {
	Ignore []string
}

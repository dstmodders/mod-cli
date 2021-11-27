package config

import (
	"os"

	"gopkg.in/yaml.v2"
)

type YAML struct {
	Workshop interface{} `yaml:"workshop"`
}

func NewYAML() *YAML {
	return &YAML{}
}

func (y *YAML) UnmarshalFile(file *os.File) error {
	stat, err := file.Stat()
	if err != nil {
		return err
	}

	in := make([]byte, stat.Size())
	_, err = file.Read(in)
	if err != nil {
		return err
	}

	return yaml.Unmarshal(in, &y)
}

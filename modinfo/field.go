package modinfo

type Field struct {
	Name        string
	Description string
	Value       interface{}
	IsRequired  bool
}

func NewField(name, description string, isRequired bool) *Field {
	return &Field{
		Name:        name,
		Description: description,
		IsRequired:  isRequired,
	}
}

func (g Field) String() string {
	return InterfaceToString(g.Value)
}

type ConfigurationOptions struct {
	Values []Option
}

func NewConfigurationOptions() *ConfigurationOptions {
	return &ConfigurationOptions{}
}

type Option struct {
	Default *OptionDefault
	Hover   string
	Label   string
	Name    string
	Options string
}

func NewOption() *Option {
	return &Option{
		Default: &OptionDefault{},
	}
}

type OptionDefault struct {
	Data        interface{}
	Description string
}

func (c *OptionDefault) DataString() string {
	return InterfaceToString(c.Data)
}

func (c *OptionDefault) String() string {
	return c.Description
}

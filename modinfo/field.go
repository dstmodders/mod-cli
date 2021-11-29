package modinfo

// Field represents a single global in modinfo.lua.
type Field struct {
	// Name is the original global name like "name", "description", "api_version",
	// etc.
	Name string

	// Description is the global description in a human-friendly format. It could
	// be used to  describe what the global does in a short form.
	Description string

	// Value holds the original global value.
	Value interface{}

	// IsRequired marks whether the global is required to be in modinfo.lua.
	IsRequired bool
}

// NewField creates a new Field instance.
func NewField(name, description string, isRequired bool) *Field {
	return &Field{
		Name:        name,
		Description: description,
		IsRequired:  isRequired,
	}
}

// String returns a string representation of a Field value which is its value.
func (g *Field) String() string {
	return InterfaceToString(g.Value)
}

// ConfigurationOptions represents "configuration_options" in modinfo.lua.
type ConfigurationOptions struct {
	// Values holds a list of all existing options.
	Values []Option
}

// NewConfigurationOptions creates a new ConfigurationOptions instance.
func NewConfigurationOptions() *ConfigurationOptions {
	return &ConfigurationOptions{}
}

// Option represents a single option for "configuration_options" in modinfo.lua.
type Option struct {
	// Label is the original label value.
	//
	// configuration_options.<option>.label
	Label string

	// Name is the original name value.
	//
	// configuration_options.<option>.name
	Name string

	// Default is the original default value.
	//
	// configuration_options.<option>.default
	Default *OptionDefault

	// Default is the original hover value.
	//
	// configuration_options.<option>.hover
	Hover string
}

// NewOption creates a new Option instance.
func NewOption() *Option {
	return &Option{
		Default: &OptionDefault{},
	}
}

// OptionDefault represents an option default value.
type OptionDefault struct {
	// Description is the original description value.
	//
	// configuration_options.<option>.default.description
	Description string

	// Description is the original data value.
	//
	// configuration_options.<option>.default.data
	Data interface{}
}

// DataString returns a string representation of an OptionDefault data.
func (c *OptionDefault) DataString() string {
	return InterfaceToString(c.Data)
}

// String returns a string representation of a OptionDefault value which is its
// description.
func (c *OptionDefault) String() string {
	return c.Description
}

// Package modinfo has been designed to interpretate modinfo.lua and give access
// to the supported values.
package modinfo

import (
	"fmt"

	lua "github.com/yuin/gopher-lua"
)

// Controller is the interface that wraps the ModInfo methods.
type Controller interface {
	FieldByName(string) *Field
	Load(string) error
}

// ModInfo represents modinfo.lua data.
type ModInfo struct {
	// General holds a map of all general fields. It includes all the required
	// ones except compatibility ones.
	General map[string]*Field

	// Compatibility holds a map of all required compatibility fields.
	Compatibility map[string]*Field

	// ConfigurationOptions holds "configuration_options".
	ConfigurationOptions *ConfigurationOptions

	// Other holds all other fields that are not required.
	Other map[string]*Field

	lState *lua.LState
}

// New creates a new ModInfo instance.
func New() *ModInfo {
	return &ModInfo{
		General: map[string]*Field{
			"api_version": NewField("api_version", "API Version", true),
			"author":      NewField("author", "Author", true),
			"description": NewField("description", "Description", true),
			"name":        NewField("name", "Title", true),
			"version":     NewField("version", "Version", true),
		},
		Compatibility: map[string]*Field{
			"dont_starve_compatible":     NewField("dont_starve_compatible", "Don't Starve Compatible", true),
			"dst_compatible":             NewField("dst_compatible", "Don't Starve Together Compatible", true),
			"reign_of_giants_compatible": NewField("reign_of_giants_compatible", "Reign Of Giants Compatible", true),
			"shipwrecked_compatible":     NewField("shipwrecked_compatible", "Shipwrecked Compatible", true),
		},
		Other: map[string]*Field{
			"all_clients_require_mod": NewField("all_clients_require_mod", "All Clients Require Mod", false),
			"client_only_mod":         NewField("client_only_mod", "Client Only Mod", false),
			"folder_name":             NewField("folder_name", "Folder Title", false),
			"forum_thread":            NewField("forum_thread", "Forum Thread", false),
			"icon":                    NewField("icon", "Icon", false),
			"icon_atlas":              NewField("icon_atlas", "Icon Atlas", false),
			"priority":                NewField("priority", "Priority", false),
		},
	}
}

func (m *ModInfo) lvBool(lv lua.LValue) (bool, error) {
	if v, ok := lv.(lua.LBool); ok {
		return bool(v), nil
	}

	if lua.LVIsFalse(lv) {
		return false, nil
	}

	return false, fmt.Errorf("not a bool")
}

func (m *ModInfo) lvNumber(lv lua.LValue) (float64, error) {
	if v, ok := lv.(lua.LNumber); ok {
		return float64(v), nil
	}

	if lua.LVIsFalse(lv) {
		return 0, nil
	}

	return 0, fmt.Errorf("not a number")
}

func (m *ModInfo) lvString(lv lua.LValue) (string, error) {
	if v, ok := lv.(lua.LString); ok {
		return string(v), nil
	}

	if lua.LVIsFalse(lv) {
		return "", nil
	}

	return "", fmt.Errorf("not a string")
}

func (m *ModInfo) lvStringField(lv lua.LValue, name string) (string, error) {
	if v, ok := m.lState.GetField(lv, name).(lua.LString); ok {
		return string(v), nil
	}

	if lua.LVIsFalse(lv) {
		return "", nil
	}

	return "", fmt.Errorf("%s is not a string", name)
}

func (m *ModInfo) globalBool(name string) (bool, error) {
	lv := m.lState.GetGlobal(name)
	return m.lvBool(lv)
}

func (m *ModInfo) globalFloat64(name string) (float64, error) {
	lv := m.lState.GetGlobal(name)
	return m.lvNumber(lv)
}

func (m *ModInfo) globalInt(name string) (int, error) {
	lv := m.lState.GetGlobal(name)
	f, err := m.lvNumber(lv)
	return int(f), err
}

func (m *ModInfo) globalString(name string) (string, error) {
	lv := m.lState.GetGlobal(name)
	return m.lvString(lv)
}

func (m *ModInfo) getConfigurationOptions() (*ConfigurationOptions, error) {
	result := NewConfigurationOptions()

	lvConfigurationOptions := m.lState.GetGlobal("configuration_options")
	//nolint:nestif
	if lvConfigurationOptionsTbl, ok := lvConfigurationOptions.(*lua.LTable); ok {
		lvConfigurationOptionsTblLen := m.lState.ObjLen(lvConfigurationOptionsTbl)
		for i := 1; i <= lvConfigurationOptionsTblLen; i++ {
			co := NewOption()
			lvConfigurationOption := m.lState.RawGet(lvConfigurationOptionsTbl, lua.LNumber(i))

			globalsStr := map[string]*string{
				"label": &co.Label,
				"name":  &co.Name,
				"hover": &co.Hover,
			}

			for global, field := range globalsStr {
				val, err := m.lvStringField(lvConfigurationOption, global)
				if err != nil {
					return nil, err
				}

				if field != nil {
					*field = val
				}
			}

			lvOptions := m.lState.GetField(lvConfigurationOption, "options")
			lvDefault := m.lState.GetField(lvConfigurationOption, "default")

			if lvOptionsTbl, ok := lvOptions.(*lua.LTable); ok {
				lvOptionsTblLen := m.lState.ObjLen(lvOptionsTbl)
				d := &OptionDefault{}
				for j := 1; j <= lvOptionsTblLen; j++ {
					lvOptionsTblOption := m.lState.RawGet(lvOptionsTbl, lua.LNumber(j))
					lvOptionsTblOptionData := m.lState.GetField(lvOptionsTblOption, "data")
					if lvOptionsTblOptionData == lvDefault {
						d.Data = lvOptionsTblOptionData.String()
						d.Description = m.lState.GetField(lvOptionsTblOption, "description").String()
					}
				}
				co.Default = d
			}

			result.Values = append(result.Values, *co)
		}
	}

	return result, nil
}

func (m *ModInfo) setBoolValues() error {
	globals, err := m.FieldsByName([]string{
		"all_clients_require_mod",
		"client_only_mod",
		"dont_starve_compatible",
		"dst_compatible",
		"reign_of_giants_compatible",
		"shipwrecked_compatible",
	})

	if err != nil {
		return err
	}

	for _, global := range globals {
		if global != nil {
			val, err := m.globalBool(global.Name)
			if err != nil {
				return err
			}
			global.Value = val
		}
	}

	return nil
}

func (m *ModInfo) setFloat64Values() error {
	globals, err := m.FieldsByName([]string{"priority"})
	if err != nil {
		return err
	}

	for _, global := range globals {
		if global != nil {
			val, err := m.globalFloat64(global.Name)
			if err != nil {
				return err
			}
			global.Value = val
		}
	}

	return nil
}

func (m *ModInfo) setIntValues() error {
	globals, err := m.FieldsByName([]string{"api_version"})
	if err != nil {
		return err
	}

	for _, global := range globals {
		if global != nil {
			val, err := m.globalInt(global.Name)
			if err != nil {
				return err
			}
			global.Value = val
		}
	}

	return nil
}

func (m *ModInfo) setStringValues() error {
	globals, err := m.FieldsByName([]string{
		"description",
		"folder_name",
		"forum_thread",
		"icon",
		"icon_atlas",
		"name",
		"version",
		"author",
	})

	if err != nil {
		return err
	}

	for _, global := range globals {
		if global != nil {
			val, err := m.globalString(global.Name)
			if err != nil {
				return err
			}
			global.Value = val
		}
	}

	return nil
}

func (m *ModInfo) setValues() error {
	if err := m.setBoolValues(); err != nil {
		return err
	}

	if err := m.setFloat64Values(); err != nil {
		return err
	}

	if err := m.setIntValues(); err != nil {
		return err
	}

	if err := m.setStringValues(); err != nil {
		return err
	}

	co, err := m.getConfigurationOptions()
	if err != nil {
		return err
	}

	m.ConfigurationOptions = co

	return nil
}

// Load loads modinfo.lua files from the provided path and sets all the
// supported values.
func (m *ModInfo) Load(path string) error {
	m.lState = lua.NewState()
	defer m.lState.Close()

	if err := m.lState.DoFile(path); err != nil {
		return err
	}

	if err := m.setValues(); err != nil {
		return err
	}

	return nil
}

// FieldByName returns a single Field based on its original global name.
func (m *ModInfo) FieldByName(name string) (*Field, error) {
	if val, ok := m.General[name]; ok {
		return val, nil
	}

	if val, ok := m.Compatibility[name]; ok {
		return val, nil
	}

	if val, ok := m.Other[name]; ok {
		return val, nil
	}

	return nil, fmt.Errorf("field %s is not supported", name)
}

// FieldsByName returns multiple Field structs based on the provided global
// names.
func (m *ModInfo) FieldsByName(names []string) (result []*Field, err error) {
	for _, name := range names {
		val, err := m.FieldByName(name)
		if err != nil {
			return result, err
		}
		result = append(result, val)
	}
	return result, nil
}

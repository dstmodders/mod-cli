package modinfo

import "strconv"

func InterfaceToString(value interface{}) string {
	switch value.(type) {
	case bool:
		val := value.(bool)
		if val {
			return "yes"
		}
		return "no"
	case float64:
		val := value.(float64)
		return strconv.FormatFloat(val, 'f', -1, 64)
	case int:
		val := value.(int)
		return strconv.Itoa(val)
	case string:
		val := value.(string)
		if len(val) == 0 {
			return "-"
		}
		return val
	}
	return "-"
}

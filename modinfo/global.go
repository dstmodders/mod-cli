package modinfo

import "strconv"

func InterfaceToString(value interface{}) string {
	switch val := value.(type) {
	case bool:
		if val {
			return "yes"
		}
		return "no"
	case float64:
		return strconv.FormatFloat(val, 'f', -1, 64)
	case int:
		return strconv.Itoa(val)
	case string:
		if len(val) == 0 {
			return "-"
		}
		return val
	}
	return "-"
}

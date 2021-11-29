package main

import (
	"fmt"
	"strconv"
	"strings"
)

func printTitle(str string) {
	fmt.Printf("[%s]\n\n", strings.ToUpper(str))
}

func printNameValue(name string, value interface{}) {
	v := "-"

	switch val := value.(type) {
	case bool:
		if val {
			v = "Yes"
			break
		}
		v = "No"
	case int:
		v = strconv.Itoa(val)
	case string:
		v = val
	}

	fmt.Printf("%s: %s\n", name, v)
}

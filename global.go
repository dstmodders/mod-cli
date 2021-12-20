package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/fatih/color"
)

func fatalError(msg string, args ...interface{}) {
	if len(args) > 0 {
		color.Red("Error: %s (%s)", msg, args[0].(error).Error())
	} else {
		color.Red("Error: %s", msg)
	}
	os.Exit(1)
}

func printTitle(str string) {
	fmt.Printf("[%s]\n\n", strings.ToUpper(str))
}

func printNameValue(name string, value interface{}) {
	v := "-"

	switch val := value.(type) {
	case bool:
		if val {
			v = "true"
			break
		}
		v = "false"
	case int:
		v = strconv.Itoa(val)
	case string:
		v = val
	}

	fmt.Printf("%s: %s\n", name, v)
}

package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/dstmodders/mod-cli/tools"
	"github.com/fatih/color"
)

func fatalError(msg string, args ...interface{}) {
	if len(args) > 0 {
		color.Red("Error: %s (%s)\n", msg, args[0].(error).Error())
	} else {
		color.Red("Error: %s\n", msg)
	}
	os.Exit(1)
}

func checkIfToolExists(docker tools.Tooler, tool tools.Tooler) {
	if !tool.ExistsOnSystem() {
		if !docker.ExistsOnSystem() {
			fatalError(fmt.Sprintf(
				"neither %s nor %s are available on the system. Install at least %s",
				tool.Name(),
				docker.Name(),
				docker.Name(),
			))
		}

		if !tool.ExistsInDocker() {
			fatalError(fmt.Sprintf(
				"%s is not available neither on the system nor in Docker",
				tool.Name(),
			))
		}

		tool.SetRunInDocker(true)
	}
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

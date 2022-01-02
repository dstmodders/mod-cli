package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/dstmodders/mod-cli/tools"
	"github.com/fatih/color"
)

func printError(err interface{}, args ...interface{}) {
	msg := ""
	switch e := err.(type) {
	case error:
		msg = e.Error()
	case string:
		msg = e
	}

	if len(args) > 0 {
		color.Red("Error: %s (%s)\n", msg, args[0].(error).Error())
	} else {
		color.Red("Error: %s\n", msg)
	}
}

//nolint:unparam
func printWarning(err interface{}, args ...interface{}) {
	msg := ""
	switch e := err.(type) {
	case error:
		msg = e.Error()
	case string:
		msg = e
	}

	if len(args) > 0 {
		color.Yellow("Warning: %s (%s)\n", msg, args[0].(error).Error())
	} else {
		color.Yellow("Warning: %s\n", msg)
	}
}

func fatalError(err interface{}, args ...interface{}) {
	printError(err, args...)
	os.Exit(1)
}

func checkIfToolExists(docker tools.Tooler, tool tools.Tooler) error {
	if !tool.ExistsOnSystem() {
		if !docker.ExistsOnSystem() {
			return fmt.Errorf(
				"neither %s nor %s are available on the system. Install at least %s",
				tool.Name(),
				docker.Name(),
				docker.Name(),
			)
		} else if !tool.IsDockerImageAvailable() {
			fmt.Printf("Pulling %s Docker image. It may take a few minutes...\n", tool.DockerImage())
			tool.PullDockerImage()
		}

		if !tool.ExistsInDocker() {
			return fmt.Errorf(
				"%s is not available neither on the system nor through Docker",
				tool.Name(),
			)
		}

		tool.SetRunInDocker(true)
	}

	return nil
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

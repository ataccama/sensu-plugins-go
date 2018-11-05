package sensuutil

import (
	"fmt"
	"os"
	"strings"

  "github.com/spf13/cobra"
)

// Exit method for all sensu checks that will print the output and desired exit code
// To use, pass it in the state you want and an optional text you would want outputted with the check.
//Ex. sensuutil.Exit("warn", "this kinda sucks")
//    sensuutil.Exit("critical", variable)
// A list of error codes currently supported can be found in common.go
func Exit(args ...interface{}) {
	var exitCode int
	output := ""

	if len(args) == 0 {
		panic("Not enough parameters.")
	}

	for i, p := range args {
		switch i {
		case 0: // name
			param, ok := p.(string)
			if !ok {
				panic("1st parameter not type string.")
			}

			for k := range MonitoringErrorCodes {
				if k == strings.ToUpper(param) {
					exitCode = MonitoringErrorCodes[k]
				}
			}

		case 1: // optional text
			param, ok := p.(string)
			if !ok {
				panic("2nd parameter not type string.")
			}
			output = param

		default:
			panic("Incorrect parameters")
		}
	}

	fmt.Printf("%v\n", output)
	os.Exit(exitCode)
}

// Cmd wrapper for initializing common sensu check plugin
//
func Cmd(name string) *cobra.Command {
  return &cobra.Command{
    Use: name,
  }
}

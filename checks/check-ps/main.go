package main

import (
	"fmt"
	"os"
	"regexp"
	"sensuatc/sensuutil"
	"strconv"
	"strings"

	"github.com/shirou/gopsutil/process"
	"github.com/spf13/cobra"
)

// Argument variables and Cobra rootCmd
//
var (
	cmd    *cobra.Command
	filter string
)

// Initialize CLI arguments
func init() {
	cmd = sensuutil.Cmd("check-ps", check)
	cmd.Flags().StringVarP(&filter, "filter", "f", "", "Search for a process")
}

// Execute the check
func main() {
	cmd.Execute()
}

// Check function
//
func check(c *cobra.Command, args []string) {
	if filter == "" {
		sensuutil.Exit("CONFIGERROR", "Filter argument ins mandatory (-f)")
	}

	var matching []string

	ps, err := process.Processes()
	if err != nil {
		sensuutil.Exit("RUNTIMEERROR", "Can't get process list")
	}

	for _, p := range ps {
		n, err := p.Cmdline()
		if err != nil {
			continue
		}

		if ok, _ := regexp.MatchString(filter, n); ok {
			if int(p.Pid) == os.Getpid() {
				continue
			}
			matching = append(matching, strconv.Itoa(int(p.Pid)))
		}
	}

	if len(matching) == 0 {
		sensuutil.Exit("CRITICAL", "No matching process")
	}

	sensuutil.Exit("OK", fmt.Sprintf("Found PIDs: %s", strings.Join(matching, ", ")))
}

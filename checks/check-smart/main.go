package main

import (
	"fmt"
	"sensuatc/sensuutil"

	"github.com/jaypipes/ghw"
	"github.com/spf13/cobra"
)

// Argument variables and Cobra rootCmd
//
var (
	cmd *cobra.Command
)

// Initialize CLI arguments
func init() {
	cmd = sensuutil.Cmd("check-smart", check)
}

// Execute the check
func main() {
	cmd.Execute()
}

// Check function
//
func check(c *cobra.Command, args []string) {
	state := "CRITICAL"

	d, _ := ghw.Block()
	for _, disk := range d.Disks {
		smartOk("/dev/"+disk.Name, &state)
	}

	sensuutil.Exit("OK", "")
}

func smartOk(disk string, state *string) (bool, error) {
	fmt.Println(disk)
	// TODO(@blufor): This needs more work
	// https://github.com/dswarbrick/smart/blob/master/cmd/smartctl/smartctl.go
	return true, nil
}

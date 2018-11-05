package main

import (
	"fmt"
	"sensuatc/sensuutil"

	sd "github.com/coreos/go-systemd/dbus"
	"github.com/spf13/cobra"
)

// Argument variables and Cobra rootCmd
//
var (
	cmd      *cobra.Command
	unitName string
)

// Initialize CLI arguments
//
func init() {
	cmd = sensuutil.Cmd("check-systemd-unit", check)
	cmd.Flags().StringVarP(&unitName, "unit-name", "u", "", "SystemD unit name (full) to check")
}

// Execute the check
//
func main() {
	cmd.Execute()
}

// Check function
//
func check(c *cobra.Command, args []string) {
	if unitName == "" {
		sensuutil.Exit("CONFIGERROR", "Unit name is mandatory (-u)")
	}
	sensuutil.Exit(checkSystemdService(unitName))
}

// Wrapper for SystemD services
//
func checkSystemdService(name string) (string, string) {
	systemd, err := sd.NewSystemConnection()
	if err != nil {
		return "RUNTIMEERROR", fmt.Sprint(err)
	}
	defer systemd.Close()

	units, err := systemd.ListUnitsByNames([]string{name})
	if err != nil {
		return "RUNTIMEERROR", fmt.Sprint(err)
	}

	if len(units) > 1 {
		return "CONFIGERROR", "multiple units match - refine your filter"
	} else if len(units) == 0 {
		return "RUNTIMEERROR", "no such unit"
	}

	message := fmt.Sprintf("%s is %s", units[0].Name, units[0].SubState)

	switch units[0].ActiveState {
	case "active":
		return "OK", message
	case "activating":
		return "WARNING", message
	case "failed":
		return "CRITICAL", message
	default:
		return "UNKNOWN", message
	}
}

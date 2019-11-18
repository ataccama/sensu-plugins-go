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

	_unitName, unitState, unitSubState, err := getSystemdService(unitName)
	if err != nil {
		sensuutil.Exit("UNKNOWN", err.Error())
	}

	if _unitName == "" {
		sensuutil.Exit("UNKNOWN", fmt.Sprintf("Unit %s not found", unitName))
	}

	var outputState string
	switch unitState {
	case "active":
		outputState = "OK"
	case "activating":
		outputState = "WARNING"
	case "failed":
		outputState = "CRITICAL"
	default:
		outputState = "UNKNOWN"
	}

	sensuutil.Exit(outputState, fmt.Sprintf("%s is %s (%s)", _unitName, unitState, unitSubState))
}

// Wrapper for SystemD services
//
func getSystemdService(name string) (string, string, string, error) {
	systemd, err := sd.NewSystemConnection()
	if err != nil {
		return "", "", "", err
	}
	defer systemd.Close()

	units, err := systemd.ListUnits()
	if err != nil {
		return "", "", "", err
	}

	for _, u := range units {
		if name == u.Name {
			return u.Name, u.ActiveState, u.SubState, nil
		}
	}

	return "", "", "", nil
}

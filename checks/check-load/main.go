package main

import (
	"fmt"
	"sensuatc/sensuutil"

	"github.com/mackerelio/go-osstat/cpu"
	"github.com/mackerelio/go-osstat/loadavg"
	"github.com/spf13/cobra"
)

// Argument variables and Cobra rootCmd
//
var (
	cmd             *cobra.Command
	normalizeToCPUs bool
	threshCrit      uint64
	threshWarn      uint64
	timeWindow      uint8
)

// Initialize CLI arguments
func init() {
	cmd = sensuutil.Cmd("check-load", check)
	cmd.Flags().BoolVarP(&normalizeToCPUs, "normalize-to-cpus", "n", true, "Divide the value by number of CPUs")
	cmd.Flags().Uint64VarP(&threshCrit, "critical", "c", 2, "Critical level threshold")
	cmd.Flags().Uint64VarP(&threshWarn, "warning", "w", 1, "Warning level threshold")
	cmd.Flags().Uint8VarP(&timeWindow, "time-window", "t", 5, "Load averages time window (1/5/15)")
}

// Execute the check
func main() {
	cmd.Execute()
}

// Check function
//
func check(c *cobra.Command, args []string) {
	load, err := loadavg.Get()
	if err != nil {
		sensuutil.Exit("runtimeerror", err)
	}

	processor, err := cpu.Get()
	if err != nil {
		sensuutil.Exit("runtimeerror", err)
	}

	var l float64
	switch timeWindow {
	case 1:
		l = load.Loadavg1
	case 5:
		l = load.Loadavg5
	case 15:
		l = load.Loadavg15
	default:
		sensuutil.Exit("CONFIGERROR", fmt.Sprintf("%d minutes is not a valid loadavg time window. Only 1/5/15 is accepted", timeWindow))
	}

	var ln float64
	if normalizeToCPUs {
		ln = l / float64(processor.CPUCount)
	} else {
		ln = l
	}

	message := fmt.Sprintf("Load averages: %.2f/%.2f/%.2f (normalized %dm load: %.2f)", load.Loadavg1, load.Loadavg5, load.Loadavg15, timeWindow, ln)

	if ln >= float64(threshCrit) {
		sensuutil.Exit("CRITICAL", message)
	} else if ln >= float64(threshWarn) {
		sensuutil.Exit("WARNING", message)
	} else {
		sensuutil.Exit("OK", message)
	}
}

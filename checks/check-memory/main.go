package main

import (
  "fmt"
  "math"
  "sensuatc/sensuutil"

  "github.com/spf13/cobra"
  "github.com/mackerelio/go-osstat/memory"
)

// Argument variables and Cobra rootCmd
//
var (
  cmd *cobra.Command
  threshWarn uint64
  threshCrit uint64
)

// Initialize CLI arguments
func init() {
  cmd = sensuutil.Cmd("check-memory", check)
  cmd.Flags().Uint64VarP(&threshWarn, "warning", "w", 80, "Warning level threshold")
  cmd.Flags().Uint64VarP(&threshCrit, "critical", "c", 90, "Critical level threshold")
}

// Execute the check
func main() {
	cmd.Execute()
}

// Check function
//
func check(c *cobra.Command, args []string) {
  m, err := memory.Get()
  if err != nil {
    sensuutil.Exit("runtimeerror", err)
  }

  used := float64(m.Total - m.Free - m.Cached)
  cached := float64(m.Cached)
  pctUsed := float64(used*100)/float64(m.Total)
  pctCached := cached*float64(100)/float64(m.Total)
  message := fmt.Sprintf("%.2f%% (%.2fGB) memory used\n%.2f%% (%.2fGB) used by cache", pctUsed, used/math.Pow(1024,3), pctCached, cached/math.Pow(1024,3))

  if pctUsed >= float64(threshCrit) {
    sensuutil.Exit("CRITICAL", message)
  } else if pctUsed >= float64(threshWarn) {
    sensuutil.Exit("WARNING", message)
  } else {
    sensuutil.Exit("OK", message)
  }
}

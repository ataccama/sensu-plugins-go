package main

import (
  "fmt"
  "math"
  "sensuatc/sensuutil"

  "github.com/mackerelio/go-osstat/memory"
)

var (
  ThreshWarn uint64
  ThreshCrit uint64
)

func init() {
  cmd := sensuutil.Cmd("check-memory")
  cmd.Flags().Uint64VarP(&ThreshWarn, "warning", "w", 80, "Warning level threshold")
  cmd.Flags().Uint64VarP(&ThreshCrit, "critical", "s", 90, "Critical level threshold")
}

func main() {
  m, err := memory.Get()
  if err != nil {
    sensuutil.Exit("runtimeerror", err)
  }

  used := float64(m.Total - m.Free - m.Cached)
  cached := float64(m.Cached)
  pctUsed := float64(used*100)/float64(m.Total)
  pctCached := cached*float64(100)/float64(m.Total)
  message := fmt.Sprintf("%.2f%% (%.2fGB) memory used\n%.2f%% (%.2fGB) used by cache", pctUsed, used/math.Pow(1024,3), pctCached, cached/math.Pow(1024,3))

  if pctUsed >= float64(ThreshCrit) {
    sensuutil.Exit("CRITICAL", message)
  } else if pctUsed >= float64(ThreshWarn) {
    sensuutil.Exit("WARNING", message)
  } else {
    sensuutil.Exit("OK", message)
  }
}
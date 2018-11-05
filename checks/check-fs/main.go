package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"sensuatc/sensuutil"
	"strings"
	"syscall"

	"github.com/spf13/cobra"
)

// Argument variables and Cobra rootCmd
//
var (
	cmd        *cobra.Command
	threshCrit uint64
	threshWarn uint64
)

// Initialize CLI arguments
func init() {
	cmd = sensuutil.Cmd("check-fs", check)
	cmd.Flags().Uint64VarP(&threshCrit, "critical", "c", 90, "Critical level threshold")
	cmd.Flags().Uint64VarP(&threshWarn, "warning", "w", 80, "Warning level threshold")
}

// Execute the check
func main() {
	cmd.Execute()
}

// Check function
//
func check(c *cobra.Command, args []string) {
	mounts, err := readMounts()
	if err != nil {
		sensuutil.Exit("RUNTIMEERROR", err)
	}

	var state string
	var message []string

	for _, mount := range mounts {
		if mount.usedPct() >= float64(threshCrit) {
			state = "CRITICAL"
			message = append(message, fmt.Sprintf("%s %.2f%% used (CRITICAL)", mount.mountPoint, mount.usedPct()))
		} else if mount.usedPct() >= float64(threshWarn) {
			message = append(message, fmt.Sprintf("%s %.2f%% used (WARNING)", mount.mountPoint, mount.usedPct()))
			if state == "" {
				state = "WARNING"
			}
		} else {
			message = append(message, fmt.Sprintf("%s %.2f%% used (OK)", mount.mountPoint, mount.usedPct()))
		}
	}

	if state == "" {
		state = "OK"
	}

	sensuutil.Exit(state, strings.Join(message, ", "))
}

// Mount type
//
type mount struct {
	dev        string
	mountPoint string
	fsType     string
	stats      syscall.Statfs_t
}

// Create slice of *mount instances
//
func readMounts() ([]*mount, error) {
	var mounts []*mount
	file, err := os.Open("/proc/self/mounts")
	if err != nil {
		return mounts, err
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		if matches, _ := regexp.Match("(^ext|^btrfs|^xfs)", []byte(fields[2])); matches {
			fs := syscall.Statfs_t{}
			err := syscall.Statfs(fields[1], &fs)
			if err != nil {
				continue
			}
			mounts = append(mounts, &mount{
				dev:        fields[0],
				mountPoint: fields[1],
				fsType:     fields[2],
				stats:      fs,
			})
		} else {
			continue
		}
	}

	return mounts, nil
}

// Return total bytes for a mountpoint
//
func (m *mount) totalSpace() uint64 {
	return m.stats.Blocks * uint64(m.stats.Bsize)
}

// Return free bytes for a mountpoint
//
func (m *mount) freeSpace() uint64 {
	return m.stats.Bfree * uint64(m.stats.Bsize)
}

// Return used bytes for a mountpoint
//
func (m *mount) usedSpace() uint64 {
	return (m.stats.Blocks - m.stats.Bfree) * uint64(m.stats.Bsize)
}

// Return free space in percentage
func (m *mount) freePct() float64 {
	return float64(m.freeSpace()) / float64(m.totalSpace()) * float64(100)
}

// Return used space in percentage
func (m *mount) usedPct() float64 {
	return float64(m.usedSpace()) / float64(m.totalSpace()) * float64(100)
}

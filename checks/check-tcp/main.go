package main

import (
	"fmt"
	"net"
	"sensuatc/sensuutil"
	"time"

	"github.com/spf13/cobra"
)

// Argument variables and Cobra rootCmd
//
var (
	cmd     *cobra.Command
	host    string
	port    uint
	timeout uint
)

// Initialize CLI arguments
func init() {
	cmd = sensuutil.Cmd("check-tcp", check)
	cmd.Flags().StringVarP(&host, "host", "H", "127.0.0.1", "Host/IP to connect to")
	cmd.Flags().UintVarP(&port, "port", "p", 0, "Port to connect to")
	cmd.Flags().UintVarP(&timeout, "timeout", "t", 5, "Timeout for connection in seconds")
}

// Execute the check
func main() {
	cmd.Execute()
}

// Check function
//
func check(c *cobra.Command, args []string) {
	if port == 0 {
		sensuutil.Exit("CONFIGERROR", "You must specify at least port")
	}
	tOutDuration, err := time.ParseDuration(fmt.Sprintf("%ds", timeout))
	if err != nil {
		sensuutil.Exit("RUNTIMEERROR", "Couldn't parse timeout duration")
	}
	ok, took, message := measure(host, port, tOutDuration)
	if ok {
		if took <= tOutDuration {
			sensuutil.Exit("OK", fmt.Sprintf("%s in %s", message, took))
		} else {
			sensuutil.Exit("WARNING", fmt.Sprintf("%s in %s (>%s)", message, took, tOutDuration))
		}
	}
	sensuutil.Exit("CRITICAL", message)
}

func measure(host string, port uint, tOut time.Duration) (bool, time.Duration, string) {
	t1 := time.Now()
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", host, port), tOut)
	if err != nil {
		return false, time.Since(t1), fmt.Sprint(err)
	}
	conn.Close()
	return true, time.Since(t1), "Successfully connected"
}

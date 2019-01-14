package main

import (
	"fmt"
	"os"
	"sort"
	//	"regexp"
	"sensuatc/sensuutil"
	"strings"

	"github.com/hashicorp/consul/api"
	"github.com/spf13/cobra"
)

// Argument variables and Cobra rootCmd
//
var (
	cmd      *cobra.Command
	service  []string
	hostname string
	address  string
)

// Initialize CLI arguments
func init() {
	hostnameAuto, err := os.Hostname()
	if err != nil {
		sensuutil.Exit("CRITICAL", err)
	}

	cmd = sensuutil.Cmd("check-consul-service", check)
	cmd.Flags().StringSliceVarP(&service, "service", "s", []string{}, "Name of the service (can be repeated)")
	cmd.Flags().StringVarP(&hostname, "hostname", "H", hostnameAuto, "Hostname of the cluster member where to run the check")
	cmd.Flags().StringVarP(&address, "address", "a", "127.0.0.1:8500", "Address of the API endpoint")
}

// Execute the check
func main() {
	cmd.Execute()
}

// Check function
//
func check(c *cobra.Command, args []string) {
	// Consul connection setup
	cc := &api.Config{
		Address: address,
		Scheme:  "http",
	}
	client, err := api.NewClient(cc)
	if err != nil {
		sensuutil.Exit("CRITICAL", "Can't connect to the Consul API endpoint")
	}

	// Query health of services
	states := make(map[string][]string)

	checks, _, err := client.Health().Node(hostname, nil)
	if err != nil {
		sensuutil.Exit("CRITICAL", "Can't get health of services")
	}

	for _, c := range checks {
		name := c.ServiceName
		if name == "" {
			name = c.Name
		}

		if len(service) > 0 {
			for _, s := range service {
				if s == name {
					states[c.Status] = append(states[c.Status], name)
				}
			}
		} else {
			states[c.Status] = append(states[c.Status], name)
		}
	}

	// Decide the state
	var exitState string
	var exitMsg []string
	if len(states["critical"]) > 0 {
		exitState = "CRITICAL"
	} else if len(states["warning"]) > 0 {
		exitState = "WARNING"
	} else {
		exitState = "OK"
	}

	// Format message
	for st, sv := range states {
		sort.Strings(sv)
		exitMsg = append(exitMsg, fmt.Sprintf("%s: %s", strings.ToUpper(st), strings.Join(sv, ", ")))
	}

	// Return the result
	sensuutil.Exit(exitState, strings.Join(exitMsg, "\n"))
}

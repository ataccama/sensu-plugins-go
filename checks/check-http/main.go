package main

import (
	"fmt"
	"net/http"
	"sensuatc/sensuutil"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

// Argument variables and Cobra rootCmd
//
var (
	cmd        *cobra.Command
	address    string
	port       uint
	timeout    uint
	path       string
	host       string
	method     string
	data       string
	threshCrit string
	threshWarn string
)

// Initialize CLI arguments
func init() {
	cmd = sensuutil.Cmd("check-ps", check)
	cmd.Flags().StringVarP(&address, "address", "a", "127.0.0.1", "Address/hostname to connect to")
	cmd.Flags().UintVarP(&port, "port", "P", 80, "Port to connect to")
	cmd.Flags().UintVarP(&timeout, "timeout", "t", 5, "Connection timeout in seconds")
	cmd.Flags().StringVarP(&path, "path", "p", "/", "URI path")
	cmd.Flags().StringVarP(&host, "host", "H", "", "HTTP Host header value")
	cmd.Flags().StringVarP(&method, "method", "m", "GET", "HTTP method to call")
	cmd.Flags().StringVarP(&data, "data", "d", "", "Send data with request")
	cmd.Flags().StringVarP(&threshCrit, "critical", "c", "400-600", "HTTP codes regarded as critical state")
	cmd.Flags().StringVarP(&threshWarn, "warning", "w", "300-399", "HTTP codes regarded as warning state")
}

// Execute the check
func main() {
	cmd.Execute()
}

// Check function
//
func check(c *cobra.Command, args []string) {
	tDuration, err := time.ParseDuration(fmt.Sprintf("%ds", timeout))
	if err != nil {
		sensuutil.Exit("CONFIGERROR", "Invalid duration")
	}

	client := &http.Client{
		Timeout: tDuration,
	}

	url := fmt.Sprintf("http://%s:%d%s", address, port, path)

	req, err := http.NewRequest(method, url, strings.NewReader(data))
	if err != nil {
		sensuutil.Exit("RUNTIMEERROR", "Can't assemble HTTP request")
	}

	resp, err := client.Do(req)
	if err != nil {
		sensuutil.Exit("CRITICAL", fmt.Sprint(err))
	}

	critCodes := generateRespCodes(threshCrit)
	warnCodes := generateRespCodes(threshWarn)

	if compareResult(resp.StatusCode, critCodes) {
		sensuutil.Exit("CRITICAL", fmt.Sprint(resp.Status))
	} else if compareResult(resp.StatusCode, warnCodes) {
		sensuutil.Exit("WARNING", fmt.Sprint(resp.Status))
	}

	sensuutil.Exit("OK", fmt.Sprint(resp.Status))

}

func generateRespCodes(r string) []uint {
	splitStr := strings.Split(r, "-")
	numBeg, _ := strconv.Atoi(splitStr[0])
	numEnd, _ := strconv.Atoi(splitStr[1])
	var retVal []uint
	for i := numBeg; i <= numEnd; i++ {
		retVal = append(retVal, uint(i))
	}
	return retVal
}

func compareResult(respCode int, codes []uint) bool {
	for i := range codes {
		if i == respCode {
			return true
		}
	}
	return false
}

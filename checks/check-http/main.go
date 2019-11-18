package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"sensuatc/sensuutil"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/tidwall/gjson"
)

// Argument variables and Cobra rootCmd
//
var (
	cmd        *cobra.Command
	url        string
	address    string
	port       uint
	scheme     bool
	timeout    uint
	timedOk    uint
	path       string
	method     string
	sendData   string
	userAuth   string
	header     []string
	threshCrit string
	threshWarn string
	jsonQuery  string
	jsonResult string
	regex      string
)

// Initialize CLI arguments
func init() {
	cmd = sensuutil.Cmd("check-http", check)
	cmd.Flags().StringVarP(&url, "url", "u", "", "Full URL")
	cmd.Flags().StringVarP(&address, "address", "a", "127.0.0.1", "Address/hostname to connect to")
	cmd.Flags().UintVarP(&port, "port", "P", 80, "Port to connect to")
	cmd.Flags().BoolVarP(&scheme, "secure", "s", false, "Enable for HTTPS")
	cmd.Flags().UintVarP(&timeout, "timeout", "t", 5000, "Connection timeout in ms")
	cmd.Flags().UintVarP(&timedOk, "timed-ok", "T", 2500, "How many ms are considered OK to take for the response")
	cmd.Flags().StringVarP(&path, "path", "p", "/", "URI path")
	cmd.Flags().StringVarP(&method, "method", "m", "GET", "HTTP method to use")
	cmd.Flags().StringVarP(&sendData, "data", "d", "", "Send data with the request")
	cmd.Flags().StringVarP(&userAuth, "user-auth", "U", "", "HTTP Basic authentication in form of [username]:[password]")
	cmd.Flags().StringArrayVarP(&header, "header", "H", []string{}, "Set header <can be repeated>")
	cmd.Flags().StringVarP(&threshCrit, "critical", "c", "400-599", "HTTP codes regarded as critical state")
	cmd.Flags().StringVarP(&threshWarn, "warning", "w", "300-399", "HTTP codes regarded as warning state")
	cmd.Flags().StringVarP(&regex, "regex", "r", "", "Regex to match the output against")
	cmd.Flags().StringVarP(&jsonQuery, "json-query", "j", "", "Query on the resulting JSON")
	cmd.Flags().StringVarP(&jsonResult, "json-result", "J", "", "Result the JSON query you expect to return")
}

// Execute the check
func main() {
	cmd.Long = `HTTP Check

Two ways of specifying target:
- full URL (-u flag)
- compose it from parts (-s, -a, -P, -p)

You can also change your request method (-m), send data with your request (-d) and add multiple headers (-H).
HTTP Basic auth is provided as well (-u).

There are 3 exclusive modes of response validation:
- HTTP reponse code (-w & -c)
- RegEx matching (-r)
- JSON data query (-j & -J; see https://github.com/tidwall/gjson/blob/master/SYNTAX.md for query language doc)

Additionally there's a "tripwire" timer (-T) that bumps any OK state to WARNING if the request takes longer.
`
	cmd.Execute()
}

// Check function
//
func check(c *cobra.Command, args []string) {
	// initial status
	status := "OK"
	var statusNote string

	// HTTP request
	code, data, duration, err := doHTTPRequest()
	if err != nil {
		sensuutil.Exit("UNKNOWN", fmt.Sprintf("HTTP Request failed (duration %v):\n\n%v", duration, err.Error()))
	}

	// Warning tripwire timer
	if duration >= time.Duration(int64(timedOk))*time.Millisecond {
		status = "WARNING"
		statusNote = " is too long"
	}

	// JSON query check
	if jsonQuery != "" {
		jsonData := gjson.GetBytes(data, jsonQuery).Raw
		if jsonResult != jsonData {
			sensuutil.Exit("CRITICAL", fmt.Sprintf("Unexpected data found\n\n%s\n\n(duration %v%s)", jsonData, duration, statusNote))
		}
		sensuutil.Exit(status, fmt.Sprintf("OK, Expected data found (duration %v%s)", duration, statusNote))
	}

	// RegEx matching check
	if regex != "" {
		re, err := regexp.Compile(regex)
		if err != nil {
			sensuutil.Exit("UNKNOWN", fmt.Sprintf("Failed to compile regular expression:\n\n%v", err.Error()))
		}
		if !re.Match(data) {
			sensuutil.Exit("CRITICAL", fmt.Sprintf("Regular expression %s doesn't match the data\n\n%s\n\n(duration %v%s)", regex, data, duration, statusNote))
		}
		sensuutil.Exit(status, fmt.Sprintf("OK, regular expression matches the data (duration %v%s)", duration, statusNote))
	}

	// HTTP response code check
	critCodes := generateRespCodes(threshCrit)
	warnCodes := generateRespCodes(threshWarn)
	if compareResultCode(code, critCodes) {
		sensuutil.Exit("CRITICAL", fmt.Sprintf("HTTP %d - critical status code (duration %v%s)", code, duration, statusNote))
	} else if compareResultCode(code, warnCodes) {
		sensuutil.Exit("WARNING", fmt.Sprintf("HTTP %d - warning status code (duration %v%s)", code, duration, statusNote))
	}
	sensuutil.Exit(status, fmt.Sprintf("HTTP %d (duration %v%s)", code, duration, statusNote))
}

// HTTP request wrapper that returns response code,
// response data, time it took and eventually error
//
func doHTTPRequest() (uint, []byte, time.Duration, error) {
	// set timeout duration
	tDuration := time.Duration(int64(timeout)) * time.Millisecond

	// Construct URL
	if url == "" {
		var _s string
		if scheme {
			_s = "https"
		} else {
			_s = "http"
		}
		url = fmt.Sprintf("%s://%s:%d%s", _s, address, port, path)
	}

	// Crepare HTTP client
	client := &http.Client{
		Timeout: tDuration,
	}

	// Construct HTTP request
	req, err := http.NewRequest(method, url, strings.NewReader(sendData))
	if err != nil {
		return 0, nil, 0, err
	}

	// Set some agent ID
	req.Header.Set("User-Agent", "Sensu HTTP check")

	// Prepare HTTP headers
	for _, h := range header {
		_h := strings.SplitN(h, ":", 2)
		req.Header.Set(_h[0], strings.TrimSpace(_h[1]))
	}

	// HTTP Basic auth if needed
	if userAuth != "" {
		_ua := strings.SplitN(userAuth, ":", 2)
		req.SetBasicAuth(_ua[0], _ua[1])
	}

	// Do the HTTP request
	t := time.Now()
	resp, err := client.Do(req)
	if err != nil {
		return 0, nil, time.Since(t), err
	}

	// Read reponse body
	defer resp.Body.Close()
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, nil, time.Since(t), err
	}
	return uint(resp.StatusCode), bodyBytes, time.Since(t), nil
}

// Handy func to generate array of ints from a range
//
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

// Even handier func to compare the returned code
// with the prepared tables
//
func compareResultCode(respCode uint, codes []uint) bool {
	for i := range codes {
		if uint(i) == respCode {
			return true
		}
	}
	return false
}

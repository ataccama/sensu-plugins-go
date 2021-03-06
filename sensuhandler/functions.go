package sensuhandler

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	"sensuatc/sensuutil"
)

// AcquireUchiwa returns an uchiwa url for the node alerting
// It requires uchiwa to be running in a consul environment
func (e EnvDetails) AcquireUchiwa(h string, c string) string {

	tags := e.Sensu.Consul.Tags
	dc := e.Sensu.Consul.Datacenter

	url := "<https://" + tags + ".uchiwa.service" + "." + dc + ".consul/#/client/" + dc + "/" + h + "?check=" + c + "|" + c + ">"
	return url
}

// CleanOutput will shorten the output to a more manageable length
func CleanOutput(output string) string {
	return strings.Split(output, ":")[0]
}

// EventName generates a simple string for use by elasticsearch and internal logging of all monitoring alerts.
func EventName(client string, check string) string {
	return client + "_" + check
}

// AcquireMonitoredInstance sets the correct device that is being monitored. In the case of snmp trap collection, containers,
// or appliance monitoring the device running the sensu-client may not be the device actually being monitored.
func (e SensuEvent) AcquireMonitoredInstance() string {
	if e.Check.Source != "" {
		return e.Check.Source
	}
	return e.Client.Name
}

// AcquireThreshold will get the current threshold for the alert state.
func (e SensuEvent) AcquireThreshold() string {
	var w, c string

	if e.Check.Thresholds.Warning != -1 {
		w = strconv.Itoa(e.Check.Thresholds.Warning)
	}
	if e.Check.Thresholds.Critical != -1 {
		c = strconv.Itoa(e.Check.Thresholds.Critical)
	}

	switch e.Check.Status {
	case 0:
		if w != "" { // this is stupid and ugly, fix it
			if c != "" {
				return "Warning Threshold: " + w + " Critical Threshold: " + c
			}
		}
		return "No thresholds set"
	case 1:
		if w != "" {
			return "Warning Threshold: " + w
		}
		return "No " + strings.ToLower(DefineStatus(1)) + " threshold set"
	case 2:
		if c != "" {
			return "Critical Threshold: " + c
		}
		return "No " + strings.ToLower(DefineStatus(2)) + " threshold set"
	case 3:
		return "No " + strings.ToLower(DefineStatus(3)) + " threshold set"
	default:
		return "No threshold information"
	}
}

// SetColor is used to set the correct notification color for a given status. By setting it in a single place for all alerts
// we ensure a measure of cohesiveness across various notification channels.
func SetColor(status int) string {
	switch status {
	case 0:
		return NotificationColor["green"]
	case 1:
		return NotificationColor["yellow"]
	case 2:
		return NotificationColor["red"]
	case 3:
		return NotificationColor["orange"]
	default:
		return NotificationColor["orange"]
	}
}

// DefineSensuEnv sets the environment that the machine is running in based upon values
// dropped via Oahi during the Chef run.
func DefineSensuEnv(env string) string {
	switch env {
	case "prd":
		return "Prod "
	case "dev":
		return "Dev "
	case "stg":
		return "Stg "
	case "vagrant":
		return "Vagrant "
	default:
		return "Test "
	}
}

// DefineStatus converts the check result status from an integer to a string.
func DefineStatus(status int) string {
	eCode := "UNDEFINED_STATUS"

	for k, v := range sensuutil.MonitoringErrorCodes {
		if status == v {
			eCode = k
		}
	}
	return eCode
}

// CreateCheckName creates a monitor name that is easily searchable in ES using different
// levels of granularity.
func CreateCheckName(check string) string {
	return strings.Replace(check, "-", ".", -1)
}

// DefineCheckStateDuration calculates how long a monitor has been in a given state.
func DefineCheckStateDuration() int {
	return 0
}

// AcquirePlaybook will return the check playbook
func (e SensuEvent) AcquirePlaybook() string {
	if e.Check.Playbook != "" {
		return "<" + e.Check.Playbook + "|" + e.Check.Name + ">"
	}
	return "No playbook is available"
}

// AcquireSensuEvent reads in the check result generated by Sensu via stdin and
// drops it into a statically defined struct.
func (e SensuEvent) AcquireSensuEvent() *SensuEvent {
	results, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(results, &e)
	if err != nil {
		panic(err)
	}
	return &e
}

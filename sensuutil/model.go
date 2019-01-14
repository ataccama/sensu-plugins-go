package sensuutil

// ObjectWrapper ..
//
type ObjectWrapper struct {
	Type string `json:"type"`
	APIVersion string `json:"api_version"`
	Spec Event `json:"spec"`
	ObjectMeta `json:"metadata"`
}


// Event ...
//
type Event struct {
	Timestamp  int64       `json:"timestamp,omitempty"`
	Entity     *Entity     `json:"entity,omitempty"`
	Check      *Check      `json:"check,omitempty"`
	Metrics    *Metrics    `json:"metrics,omitempty"`
	ObjectMeta *ObjectMeta `json:"metadata"`
}

// Entity ...
//
type Entity struct {
	EntityClass    string         `json:"entity_class"`
	System         System         `json:"system"`
	Subscriptions  []string       `json:"subscriptions"`
	LastSeen       int64          `json:"last_seen"`
	Deregister     bool           `json:"deregister"`
	Deregistration Deregistration `json:"deregistration"`
	User           string         `json:"user,omitempty"`
	Redact         []string       `json:"redact,omitempty"`
	ObjectMeta     `json:"metadata,omitempty"`
}

// Check ...
//
type Check struct {
	Command              string          `json:"command,omitempty"`
	Handlers             []string        `json:"handlers"`
	HighFlapThreshold    uint32          `json:"high_flap_threshold"`
	Interval             uint32          `json:"interval"`
	LowFlapThreshold     uint32          `json:"low_flap_threshold"`
	Publish              bool            `json:"publish"`
	RuntimeAssets        []string        `json:"runtime_assets"`
	Subscriptions        []string        `json:"subscriptions"`
	ProxyEntityName      string          `json:"proxy_entity_name"`
	CheckHooks           []HookList      `json:"check_hooks"`
	Stdin                bool            `json:"stdin"`
	Subdue               *TimeWindowWhen `json:"subdue"`
	Cron                 string          `json:"cron,omitempty"`
	Ttl                  int64           `json:"ttl"`
	Timeout              uint32          `json:"timeout"`
	ProxyRequests        *ProxyRequests  `json:"proxy_requests,omitempty"`
	RoundRobin           bool            `json:"round_robin"`
	Duration             float64         `json:"duration,omitempty"`
	Executed             int64           `json:"executed"`
	History              []CheckHistory  `json:"history"`
	Issued               int64           `json:"issued"`
	Output               string          `json:"output"`
	State                string          `json:"state,omitempty"`
	Status               uint32          `json:"status"`
	TotalStateChange     uint32          `json:"total_state_change"`
	LastOK               int64           `json:"last_ok"`
	Occurrences          int64           `json:"occurrences"`
	OccurrencesWatermark int64           `json:"occurrences_watermark"`
	Silenced             []string        `json:"silenced,omitempty"`
	Hooks                []*Hook         `json:"hooks,omitempty"`
	OutputMetricFormat   string          `json:"output_metric_format"`
	OutputMetricHandlers []string        `json:"output_metric_handlers"`
	EnvVars              []string        `json:"env_vars"`
	ObjectMeta           `json:"metadata,omitempty"`
}

// CheckHistory ...
//
type CheckHistory struct {
	Status   uint32 `json:"status"`
	Executed int64  `json:"executed"`
}

// ProxyRequests ...
//
type ProxyRequests struct {
	EntityAttributes []string `json:"entity_attributes"`
	Splay            bool     `json:"splay"`
	SplayCoverage    uint32   `json:"splay_coverage"`
}

// Metrics ...
//
type Metrics struct{}

// Deregistration ...
//
type Deregistration struct {
  Handler string `json:"handler,omitempty"`
}

// ObjectMeta ...
//
type ObjectMeta struct {
	Name        string            `json:"name,omitempty"`
	Namespace   string            `json:"namespace,omitempty"`
	Labels      map[string]string `json:"labels,omitempty"`
	Annotations map[string]string `json:"annotations,omitempty"`
}

// System ...
//
type System struct {
	Hostname        string  `json:"hostname,omitempty"`
	OS              string  `json:"os,omitempty"`
	Platform        string  `json:"platform,omitempty"`
	PlatformFamily  string  `json:"platform_family,omitempty"`
	PlatformVersion string  `json:"platform_version,omitempty"`
	Network         Network `json:"network"`
	Arch            string  `json:"arch,omitempty"`
}

// Network ...
//
type Network struct {
	Interfaces []NetworkInterface `json:"interfaces"`
}

// NetworkInterface ...
//
type NetworkInterface struct {
	Name      string   `json:"name,omitempty"`
	MAC       string   `json:"mac,omitempty"`
	Addresses []string `json:"addresses"`
}

// TimeWindowWhen ...
//
type TimeWindowWhen struct {
	Days TimeWindowDays `json:"days"`
}

// TimeWindowDays ...
//
type TimeWindowDays struct {
	All       []*TimeWindowTimeRange `json:"all,omitempty"`
	Sunday    []*TimeWindowTimeRange `json:"sunday,omitempty"`
	Monday    []*TimeWindowTimeRange `json:"monday,omitempty"`
	Tuesday   []*TimeWindowTimeRange `json:"tuesday,omitempty"`
	Wednesday []*TimeWindowTimeRange `json:"wednesday,omitempty"`
	Thursday  []*TimeWindowTimeRange `json:"thursday,omitempty"`
	Friday    []*TimeWindowTimeRange `json:"friday,omitempty"`
	Saturday  []*TimeWindowTimeRange `json:"saturday,omitempty"`
}

// TimeWindowTimeRange ...
//
type TimeWindowTimeRange struct {
	Begin string `json:"begin"`
	End   string `json:"end"`
}

// Hook ...
//
type Hook struct{}

// HookList ...
//
type HookList struct{}

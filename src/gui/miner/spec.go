package miner

// SupportedMiners contains a list of the currently supported miners
var SupportedMiners = []string{"xlarig"}

// Config holds miner specific configuration information
type Config struct {
	// Type of miner
	Type string `json:"type"`
	// Path to the selected miner's executable
	Path string `json:"path"`
	// Endpoint of the miner's JSON API
	Endpoint string `json:"endpoint"`
}

// ProcessingConfig holds the config for the miner's processing setup
// TODO: Right now this is only for CPU threads and will be extended into
// full CPU/GPU config
type ProcessingConfig struct {
	// Type of miner
	Type string `json:"type"`
	// Threads is the amount of CPU threads
	Threads uint16 `json:"threads"`
	// MaxThreads is the maximum threads as read by runtime.NumCPU
	MaxThreads uint16 `json:"max_threads"`
	// MaxUsage is the maximum CPU usage in percentage the miner should
	// attempt to use.
	// Currently only supported by xmrig CPU backend
	MaxUsage uint8 `json:"max_usage"`
}

// Stats contains the miner statistics required by the front-end
type Stats struct {
	// Hashrate is the current miner hashrate
	Hashrate float64 `json:"hashrate"`
	// HashrateHuman is the H/s, KH/s or MH/s representation of hashrate
	HashrateHuman string `json:"hashrate_human"`
	// CurrentDifficulty as set by the pool
	CurrentDifficulty int `json:"current_difficulty"`
	// SharesGood is the good shares counter
	SharesGood int `json:"shares_good"`
	// SharesGood is the bad shares counter
	SharesBad int `json:"shares_bad"`
	// Uptime for the miner in seconds
	Uptime int `json:"uptime"`
	// UptimeHuman is the human readable version of uptime, ex. 10 minutes
	UptimeHuman string `json:"uptime_human"`
	// Errors is a list of errors that have occurred
	Errors []string `json:"errors"`
	// UpdateGraph is set to true if the stats graph should be updated
	UpdateGraph bool `json:"update_graph"`
	// Address contains the Scala address we are mining to
	// TODO: This should be somewhere else, it's not stats!
	Address string `json:"address"`
}

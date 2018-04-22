package miner

// SupportedMiners contains a list of the currently supported miners
var SupportedMiners = []string{"xmr-stak", "xmrig"}

// Config holds miner specific configuration information
type Config struct {
	// Type of miner
	Type string `json:"type"`
	// Path to the selected miner's executable
	Path string `json:"path"`
	// Endpoint of the miner's JSON API
	Endpoint string `json:"endpoint"`
}

// Stats contains the miner statistics required by the front-end
type Stats struct {
	// Hashrate is the current miner hashrate
	Hashrate float64
	// HashrateHuman is the H/s, KH/s or MH/s representation of hashrate
	HashrateHuman string
	// CurrentDifficulty as set by the pool
	CurrentDifficulty int
	// SharesGood is the good shares counter
	SharesGood int
	// SharesGood is the bad shares counter
	SharesBad int
	// Uptime for the miner in seconds
	Uptime int
	// UptimeHuman is the human readable version of uptime, ex. 10 minutes
	UptimeHuman string
	// Errors is a list of errors that have occurred
	Errors []string
}

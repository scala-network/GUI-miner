package miner

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"runtime"
)

// Xmrig implements the miner interface for the xmrig miner, including
// xmrig-amd and xmrig-nvidia
// https://github.com/xmrig/xmrig
// https://github.com/xmrig/xmrig-amd
// https://github.com/xmrig/xmrig-nvidia
type Xmrig struct {
	Base
	name         string
	endpoint     string
	lastHashrate float64
}

// XmrigConfig is the config.json structure for Xmrig
// Generated with https://mholt.github.io/json-to-go/
type XmrigConfig struct {
	Algo        string            `json:"algo"`
	Av          int               `json:"av"`
	Background  bool              `json:"background"`
	Colors      bool              `json:"colors"`
	CPUAffinity interface{}       `json:"cpu-affinity"`
	CPUPriority interface{}       `json:"cpu-priority"`
	DonateLevel int               `json:"donate-level"`
	LogFile     interface{}       `json:"log-file"`
	MaxCPUUsage int               `json:"max-cpu-usage"`
	PrintTime   int               `json:"print-time"`
	Retries     int               `json:"retries"`
	RetryPause  int               `json:"retry-pause"`
	Safe        bool              `json:"safe"`
	Syslog      bool              `json:"syslog"`
	Threads     int               `json:"threads"`
	Pools       []XmrigPoolConfig `json:"pools"`
	API         XmrigAPIConfig    `json:"api"`
}

// XmrigPoolConfig contains the configuration for a pool in Xmrig
type XmrigPoolConfig struct {
	URL       string `json:"url"`
	User      string `json:"user"`
	Pass      string `json:"pass"`
	Keepalive bool   `json:"keepalive"`
	Nicehash  bool   `json:"nicehash"`
	Variant   int    `json:"variant"`
}

// XmrigAPIConfig contains the Xmrig API config
type XmrigAPIConfig struct {
	Port        int         `json:"port"`
	AccessToken interface{} `json:"access-token"`
	WorkerID    interface{} `json:"worker-id"`
}

// XmrigResponse contains the data from xmrig API
// Generated with https://mholt.github.io/json-to-go/
type XmrigResponse struct {
	ID       string `json:"id"`
	WorkerID string `json:"worker_id"`
	Version  string `json:"version"`
	Kind     string `json:"kind"`
	Ua       string `json:"ua"`
	CPU      struct {
		Brand   string `json:"brand"`
		Aes     bool   `json:"aes"`
		X64     bool   `json:"x64"`
		Sockets int    `json:"sockets"`
	} `json:"cpu"`
	Algo        string `json:"algo"`
	Hugepages   bool   `json:"hugepages"`
	DonateLevel int    `json:"donate_level"`
	Hashrate    struct {
		Total   []float64   `json:"total"`
		Highest float64     `json:"highest"`
		Threads [][]float64 `json:"threads"`
	} `json:"hashrate"`
	Results struct {
		DiffCurrent int      `json:"diff_current"`
		SharesGood  int      `json:"shares_good"`
		SharesTotal int      `json:"shares_total"`
		AvgTime     int      `json:"avg_time"`
		HashesTotal int      `json:"hashes_total"`
		Best        []int    `json:"best"`
		ErrorLog    []string `json:"error_log"`
	} `json:"results"`
	Connection struct {
		Pool     string   `json:"pool"`
		Uptime   int      `json:"uptime"`
		Ping     int      `json:"ping"`
		Failures int      `json:"failures"`
		ErrorLog []string `json:"error_log"`
	} `json:"connection"`
}

// NewXmrig creates a new xmrig miner instance
func NewXmrig(config Config) (*Xmrig, error) {

	endpoint := config.Endpoint
	if endpoint == "" {
		endpoint = "http://127.0.0.1:16000"
	}

	miner := Xmrig{
		name:     "xmrig",
		endpoint: endpoint,
	}
	miner.Base.executableName = filepath.Base(config.Path)
	miner.Base.executablePath = filepath.Dir(config.Path)

	return &miner, nil
}

// WriteConfig writes the miner's configuration in the xmrig format
func (miner *Xmrig) WriteConfig(
	poolEndpoint string,
	walletAddress string) error {

	defaultConfig := miner.defaultConfig(poolEndpoint, walletAddress)
	configBytes, err := json.Marshal(defaultConfig)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(
		filepath.Join(miner.Base.executablePath, "config.json"),
		configBytes,
		0644)
	if err != nil {
		return err
	}

	return nil
}

// GetName returns the name of the miner
func (miner *Xmrig) GetName() string {
	return miner.name
}

// GetLastHashrate returns the last reported hashrate
func (miner *Xmrig) GetLastHashrate() float64 {
	return miner.lastHashrate
}

// GetStats returns the current miner stats
func (miner *Xmrig) GetStats() (Stats, error) {
	var stats Stats
	var xmrigStats XmrigResponse
	resp, err := http.Get(miner.endpoint)
	if err != nil {
		return stats, err
	}
	err = json.NewDecoder(resp.Body).Decode(&xmrigStats)
	if err != nil {
		return stats, err
	}

	var hashrate float64
	if len(xmrigStats.Hashrate.Total) > 0 {
		hashrate = xmrigStats.Hashrate.Total[0]
	}
	miner.lastHashrate = hashrate

	var errors []string
	/*
		TODO: I noticed errors are not reported in the xmrig API. To replicate,
		use an invalid Stellite address with a pool that checks the address. In
		the command line you'll notice errors printed, but not added in the API.
		ApiState.cpp::getConnection and getResults functions might give some clues
		to getting it fixed.
	*/
	/*
		if len(xmrigStats.Connection.ErrorLog) > 0 {
			for _, err := range xmrigStats.Connection.ErrorLog {
				errors = append(errors, fmt.Sprintf("%s",
					err.Text,
				))
			}
		}
		if len(xmrigStats.Results.ErrorLog) > 0 {
			for _, err := range xmrigStats.Results.ErrorLog {
				errors = append(errors, fmt.Sprintf("(%d) %s",
					err.Count,
					err.Text,
				))
			}
		}
	*/
	stats = Stats{
		Hashrate:          hashrate,
		HashrateHuman:     HumanizeHashrate(hashrate),
		CurrentDifficulty: xmrigStats.Results.DiffCurrent,
		Uptime:            xmrigStats.Connection.Uptime,
		UptimeHuman:       HumanizeTime(xmrigStats.Connection.Uptime),
		SharesGood:        xmrigStats.Results.SharesGood,
		SharesBad:         xmrigStats.Results.SharesTotal - xmrigStats.Results.SharesGood,
		Errors:            errors,
	}
	return stats, nil
}

// defaultConfig returns a default setup for Xmrig
func (miner *Xmrig) defaultConfig(
	poolEndpoint string,
	walletAddress string) XmrigConfig {

	runInBackground := true
	// On Mac OSX xmrig doesn't run is we fork the process to the background and
	// xmrig forks to the background again
	if runtime.GOOS == "darwin" {
		runInBackground = false
	}

	return XmrigConfig{
		Algo:        "cryptonight",
		Av:          0,
		Background:  runInBackground,
		Colors:      true,
		CPUAffinity: nil,
		CPUPriority: nil,
		DonateLevel: 1,
		LogFile:     nil,
		MaxCPUUsage: 80,
		PrintTime:   3600,
		Retries:     5,
		RetryPause:  5,
		Safe:        false,
		Syslog:      false,
		Threads:     0,
		Pools: []XmrigPoolConfig{
			{
				URL:       poolEndpoint,
				User:      walletAddress,
				Pass:      "Stellite GUI Miner",
				Keepalive: true,
				Nicehash:  false,
				Variant:   1,
			},
		},
		API: XmrigAPIConfig{
			Port:        16000,
			AccessToken: nil,
			WorkerID:    nil,
		},
	}
}

package miner

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"runtime"
	"strings"
)

// Xmrig implements the miner interface for the xmrig miner, including
// xmrig-amd and xmrig-nvidia
// https://github.com/xmrig/xmrig
// https://github.com/xmrig/xmrig-amd
// https://github.com/xmrig/xmrig-nvidia
type Xmrig struct {
	Base
	name             string
	endpoint         string
	lastHashrate     float64
	resultStatsCache XmrigResponse
	isGPU            bool
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
	MaxCPUUsage uint8             `json:"max-cpu-usage"`
	PrintTime   int               `json:"print-time"`
	Retries     int               `json:"retries"`
	RetryPause  int               `json:"retry-pause"`
	Safe        bool              `json:"safe"`
	Syslog      bool              `json:"syslog"`
	Threads     uint16            `json:"threads"`
	Pools       []XmrigPoolConfig `json:"pools"`
	API         XmrigAPIConfig    `json:"api"`
}

// XmrigGPUConfig is the config.json structure for Xmrig's GPU
// Generated with https://mholt.github.io/json-to-go/
type XmrigGPUConfig struct {
	Algo        string            `json:"algo"`
	Av          int               `json:"av"`
	Background  bool              `json:"background"`
	Colors      bool              `json:"colors"`
	CPUAffinity interface{}       `json:"cpu-affinity"`
	CPUPriority interface{}       `json:"cpu-priority"`
	DonateLevel int               `json:"donate-level"`
	LogFile     interface{}       `json:"log-file"`
	MaxCPUUsage uint8             `json:"max-cpu-usage"`
	PrintTime   int               `json:"print-time"`
	Retries     int               `json:"retries"`
	RetryPause  int               `json:"retry-pause"`
	Safe        bool              `json:"safe"`
	Syslog      bool              `json:"syslog"`
	Threads     []struct{}        `json:"threads"`
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
	// xmrig appends either nvidia or amd to the miner if it's GPU only
	// just make sure that it's not the platform name containing amd64
	if (strings.Contains(config.Path, "nvidia") ||
		strings.Contains(config.Path, "amd")) &&
		strings.Contains(config.Path, "amd64") == false {
		fmt.Println("SETTING TO GPU")
		fmt.Println(config.Path)
		miner.isGPU = true
		miner.name += "-gpu"
	}
	miner.Base.executableName = filepath.Base(config.Path)
	miner.Base.executablePath = filepath.Dir(config.Path)

	return &miner, nil
}

// WriteConfig writes the miner's configuration in the xmrig format
func (miner *Xmrig) WriteConfig(
	poolEndpoint string,
	walletAddress string,
	processingConfig ProcessingConfig) error {

	var err error
	var configBytes []byte
	if miner.isGPU {
		defaultConfig := miner.createGPUConfig(
			poolEndpoint,
			walletAddress)
		configBytes, err = json.Marshal(defaultConfig)
		if err != nil {
			return err
		}
	} else {
		defaultConfig := miner.createConfig(
			poolEndpoint,
			walletAddress,
			processingConfig)
		configBytes, err = json.Marshal(defaultConfig)
		if err != nil {
			return err
		}
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

// GetProcessingConfig returns the current miner processing config
// TODO: Currently only CPU threads, extend this to full CPU/GPU config
func (miner *Xmrig) GetProcessingConfig() ProcessingConfig {

	// Get max CPU usage from the config file
	configBytes, err := ioutil.ReadFile(
		filepath.Join(miner.Base.executablePath, "config.json"))
	if err != nil {
		return ProcessingConfig{}
	}

	// xmrig's threads field is not an int when it's GPU only so we need to use
	// a defferent config structure
	if miner.isGPU {
		var config XmrigGPUConfig
		err = json.Unmarshal(configBytes, &config)
		if err != nil {
			return ProcessingConfig{}
		}
		return ProcessingConfig{
			MaxUsage:   config.MaxCPUUsage,
			Threads:    uint16(len(miner.resultStatsCache.Hashrate.Threads)),
			MaxThreads: uint16(runtime.NumCPU()),
			Type:       miner.name,
		}
	}

	var config XmrigConfig
	err = json.Unmarshal(configBytes, &config)
	if err != nil {
		return ProcessingConfig{}
	}
	return ProcessingConfig{
		MaxUsage:   config.MaxCPUUsage,
		Threads:    uint16(len(miner.resultStatsCache.Hashrate.Threads)),
		MaxThreads: uint16(runtime.NumCPU()),
		Type:       miner.name,
	}
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
	miner.resultStatsCache = xmrigStats
	return stats, nil
}

// createConfig returns creates the config for Xmrig
func (miner *Xmrig) createConfig(
	poolEndpoint string,
	walletAddress string,
	processingConfig ProcessingConfig) XmrigConfig {

	runInBackground := true
	// On Mac OSX xmrig doesn't run is we fork the process to the background and
	// xmrig forks to the background again
	// Seems like xmrig doesn't like running GPU in the background
	if runtime.GOOS == "darwin" || miner.isGPU {
		runInBackground = false
	}

	config := XmrigConfig{
		Algo:        "cryptonight",
		Av:          0,
		Background:  runInBackground,
		Colors:      true,
		CPUAffinity: nil,
		CPUPriority: nil,
		DonateLevel: 1,
		LogFile:     nil,
		MaxCPUUsage: processingConfig.MaxUsage,
		PrintTime:   3600,
		Retries:     5,
		RetryPause:  5,
		Safe:        false,
		Syslog:      false,
		Threads:     processingConfig.Threads,
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

	return config
}

// createGPUConfig returns creates the config for Xmrig GPU setups
func (miner *Xmrig) createGPUConfig(
	poolEndpoint string,
	walletAddress string) XmrigGPUConfig {

	config := XmrigGPUConfig{
		Algo:        "cryptonight",
		Av:          0,
		Background:  false,
		Colors:      true,
		CPUAffinity: nil,
		CPUPriority: nil,
		DonateLevel: 1,
		LogFile:     nil,
		PrintTime:   3600,
		Retries:     5,
		RetryPause:  5,
		Safe:        false,
		Syslog:      false,
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

	return config
}

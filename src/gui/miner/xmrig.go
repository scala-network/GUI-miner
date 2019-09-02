package miner

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"runtime"
	"strings"
	"os"
	// "fmt"
	// "bytes"

	"github.com/sirupsen/logrus"
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
	logger           *logrus.Entry
}

// XmrigConfig is the config.json structure for Xmrig
// Generated with https://mholt.github.io/json-to-go/
type XmrigConfig struct {
	API             XmrigAPIConfig      `json:"api"`
	HTTP            XmrigHttpConfig     `json:"http"`
	Autosave        bool                `json:"autosave"`
	Version         int                 `json:"version"`
	Background      bool                `json:"background"`
	Colors          bool                `json:"colors"`
	RandomX         XmrigRandomXConfig  `json:"randomx"`
	Cpu             XmrigCpuConfig      `json:"cpu"`
	DonateLevel     int                 `json:"donate-level"`
	DonateOverProxy int                 `json:"donate-over-proxy"`
	LogFile         interface{}         `json:"log-file"`
	Pools           []XmrigPoolConfig   `json:"pools"`
	PrintTime       int                 `json:"print-time"`
	Retries         int                 `json:"retries"`
	RetryPause      int                 `json:"retry-pause"`
	Syslog          bool                `json:"syslog"`
	UserAgent       interface{}         `json:"user-agent"`
	Watch           bool                `json:"watch"`
}

// XmrigGPUConfig is the config.json structure for Xmrig's GPU
// Generated with https://mholt.github.io/json-to-go/
type XmrigGPUConfig struct {
	// Av          int         `json:"av"`
	// Background  bool        `json:"background"`
	// Colors      bool        `json:"colors"`
	// CPUAffinity interface{} `json:"cpu-affinity"`
	// CPUPriority interface{} `json:"cpu-priority"`
	// DonateLevel int         `json:"donate-level"`
	// LogFile     interface{} `json:"log-file"`
	MaxCPUUsage uint8       `json:"max-cpu-usage"`
	// PrintTime   int         `json:"print-time"`
	// Retries     int         `json:"retries"`
	// RetryPause  int         `json:"retry-pause"`
	// Safe        bool        `json:"safe"`
	// Syslog      bool        `json:"syslog"`
	// TODO: This is the only difference between GPU and CPU, the threads
	// structure in the config. I need to merge the the config structures into
	// one with an omitempty tag on each and only fill in the one needed
	// Threads []struct{}        `json:"threads"`
	// Pools   []XmrigPoolConfig `json:"pools"`
	// API     XmrigAPIConfig    `json:"api"`
}

// XmrigAPIConfig contains the Xmrig API config
type XmrigAPIConfig struct {
	Id          interface{} `json:"id"`
	WorkerID    interface{} `json:"worker-id"`
}

// XmrigHttpConfig contains the Xmrig HTTP config
type XmrigHttpConfig struct {
	Enabled     bool        `json:"enabled"`
	Host        string      `json:"host"`
	Port        int         `json:"port"`
	AccessToken interface{} `json:"access-token"`
	Restricted  bool        `json:"restricted"`
}

// XmrigCpuConfig contains the Xmrig CPU config
type XmrigCpuConfig struct {
	Enabled    bool        `json:"enabled"`
	HugePages  bool        `json:"huge-pages"`
	HwAes      interface{} `json:"hw-aes"`
	Priority   interface{} `json:"priority"`
	Asm        bool        `json:"asm"`
	Argon2Impl interface{} `json:"argon2-impl"`
	Argon2     []int       `json:"argon2"`
	Cn         [][]int     `json:"cn"`
	CnHeavy    [][]int     `json:"cn-heavy"`
	CnLite     [][]int     `json:"cn-lite"`
	CnPico     [][]int     `json:"cn-pico"`
	CnGpu      []int       `json:"cn/gpu"`
	Rx         []int       `json:"rx"`
	RxWow      []int       `json:"rx/wow"`
	Cn0        bool        `json:"cn/0"`
	CnLite0    bool        `json:"cn-lite/0"`
}

// XmrigRandomXConfig contains the Xmrig RandomX config
type XmrigRandomXConfig struct {
	Init int  `json:"init"`
	Numa bool `json:"numa"`
}

// XmrigPoolConfig contains the configuration for a pool in Xmrig
type XmrigPoolConfig struct {
	// Algo           interface{} `json:"algo"`
	Algo           string      `json:"algo"`
	URL            string      `json:"url"`
	User           string      `json:"user"`
	Pass           string      `json:"pass"`
	RigId          string      `json:"rig-id"`
	Nicehash       bool        `json:"nicehash"`
	Keepalive      bool        `json:"keepalive"`
	Enabled        bool        `json:"enabled"`
	Tls            bool        `json:"tls"`
	TlsFingerprint interface{} `json:"tls-fingerprint"`
	Daemon         bool        `json:"daemon"`
}

// XmrigResponse contains the data from xmrig API
// Generated with https://mholt.github.io/json-to-go/
type XmrigResponse struct {
	ID       string   `json:"id"`
	WorkerID string   `json:"worker_id"`
	Uptime   int      `json:"uptime"`
	Features []string `json:"features"`
	Results struct {
		DiffCurrent int64    `json:"diff_current"`
		SharesGood  int      `json:"shares_good"`
		SharesTotal int      `json:"shares_total"`
		AvgTime     int      `json:"avg_time"`
		HashesTotal int      `json:"hashes_total"`
		Best        []int    `json:"best"`
		ErrorLog    []string `json:"error_log"`
	} `json:"results"`
	Algo        string `json:"algo"`
	Connection struct {
		Pool           string      `json:"pool"`
		Ip             string      `json:"ip"`
		Uptime         int         `json:"uptime"`
		Ping           int         `json:"ping"`
		Failures       int         `json:"failures"`
		Tls            interface{} `json:"tls"`
		TlsFingerprint interface{} `json:"tls-fingerprint"`
		ErrorLog       []string    `json:"error_log"`
	} `json:"connection"`
	Version  string   `json:"version"`
	Kind     string   `json:"kind"`
	Ua       string   `json:"ua"`
	CPU      struct {
		Brand    string `json:"brand"`
		Aes      bool   `json:"aes"`
		Avx2     bool   `json:"avx2"`
		X64      bool   `json:"x64"`
		L2       int    `json:"l2"`
		L3       int    `json:"l3"`
		Cores    int    `json:"cores"`
		Threads  int    `json:"threads"`
		Packages int    `json:"packages"`
		Nodes    int    `json:"nodes"`
		Backend  string `json:"backend"`
		Assembly string `json:"assembly"`
	} `json:"cpu"`
	Hugepages   bool   `json:"hugepages"`
	DonateLevel int    `json:"donate_level"`
	Hashrate    struct {
		Total   []float64   `json:"total"`
		Highest float64     `json:"highest"`
		Threads [][]float64 `json:"threads"`
	} `json:"hashrate"`
}

// NewXmrig creates a new xmrig miner instance
func NewXmrig(config Config) (*Xmrig, error) {

	endpoint := config.Endpoint
	if endpoint == "" {
		endpoint = "http://127.0.0.1:16000/1/summary"
	}

	miner := Xmrig{
		// We've switched back to the original miner XMRIG but we will 
		// keep an eye on it to make sure the compatibility works for future update
		name:     "xmrig",
		endpoint: endpoint,
	}
	// xmrig appends either nvidia or amd to the miner if it's GPU only
	// just make sure that it's not the platform name containing amd64
	if (strings.Contains(config.Path, "nvidia") ||
		strings.Contains(config.Path, "amd")) &&
		strings.Contains(config.Path, "amd64") == false {
		miner.isGPU = true
		miner.name += "-gpu"
	}
	miner.Base.executableName = filepath.Base(config.Path)
	miner.Base.executablePath = filepath.Dir(config.Path)

	// Setup the logging, by default we log to stdout
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "Jan 02 15:04:05",
	})
	logrus.SetLevel(logrus.DebugLevel)
	logrus.SetOutput(os.Stdout)

	// Setting the WithFields now will ensure all log entries from this point
	// includes the fields
	miner.logger = logrus.WithFields(logrus.Fields{
		"service": "XmrStak",
	})

	return &miner, nil
}

// WriteConfig writes the miner's configuration in the xmrig format
func (miner *Xmrig) WriteConfig(
	poolEndpoint string,
	walletAddress string,
	coinAlgorithm string,
	XmrigAlgo string,
	XmrigVariant string,
	processingConfig ProcessingConfig) error {

	var err error
	var configBytes []byte
	if miner.isGPU {
		defaultConfig := miner.createGPUConfig(
			poolEndpoint,
			walletAddress,
			XmrigAlgo,
			XmrigVariant)
		configBytes, err = json.MarshalIndent(defaultConfig, "", "  ")
		if err != nil {
			return err
		}
	} else {
		defaultConfig := miner.createConfig(
			poolEndpoint,
			walletAddress,
			XmrigAlgo,
			XmrigVariant,
			processingConfig)
		configBytes, err = json.MarshalIndent(defaultConfig, "", "  ")
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
	// Reset hashrate
	miner.lastHashrate = 0.00
	return nil
}

// GetProcessingConfig returns the current miner processing config
// TODO: Currently only CPU threads, extend this to full CPU/GPU config
func (miner *Xmrig) GetProcessingConfig() ProcessingConfig {

	// Get max CPU usage from the config file
	configBytes, err := ioutil.ReadFile(
		filepath.Join(miner.Base.executablePath, "config.json"))
	if err != nil {
		return ProcessingConfig{
			MaxUsage:     0,
			Threads:      uint16(len(miner.resultStatsCache.Hashrate.Threads)),
			MaxThreads:   uint16(runtime.NumCPU()),
			Type:         miner.name,
			HardwareType: 1,
		}
	}

	// xmrig's threads field is not an int when it's GPU only so we need to use
	// a different config structure
	if miner.isGPU {
		var config XmrigGPUConfig
		err = json.Unmarshal(configBytes, &config)
		if err != nil {
			return ProcessingConfig{
				MaxUsage:     0,
				Threads:      uint16(len(miner.resultStatsCache.Hashrate.Threads)),
				MaxThreads:   uint16(runtime.NumCPU()),
				Type:         miner.name,
				HardwareType: 1,
			}
		}

		return ProcessingConfig{
			MaxUsage:     config.MaxCPUUsage,
			Threads:      uint16(len(miner.resultStatsCache.Hashrate.Threads)),
			MaxThreads:   uint16(runtime.NumCPU()),
			Type:         miner.name,
			HardwareType: 1,
		}
	}

	var config XmrigConfig
	err = json.Unmarshal(configBytes, &config)
	if err != nil {
		return ProcessingConfig{
			MaxUsage:     0,
			Threads:      uint16(len(miner.resultStatsCache.Hashrate.Threads)),
			MaxThreads:   uint16(runtime.NumCPU()),
			Type:         miner.name,
			HardwareType: 1,
		}
	}

	return ProcessingConfig{
		// MaxUsage:     config.MaxCPUUsage,
		MaxUsage:     100, // a small hack FTM
		Threads:      uint16(len(miner.resultStatsCache.Hashrate.Threads)),
		MaxThreads:   uint16(runtime.NumCPU()),
		Type:         miner.name,
		HardwareType: 1,
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
		use an invalid Bloc address with a pool that checks the address. In
		the command line you'll notice errors printed, but not added in the API.
		ApiState.cpp::getConnection and getResults functions might give some clues
		to getting it fixed.
		Issue reported: https://github.com/xmrig/xmrig/issues/589
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
	XmrigAlgo string,
	XmrigVariant string,
	processingConfig ProcessingConfig) XmrigConfig {

	// runInBackground := true
	// On Mac OSX xmrig doesn't run if we fork the process to the background and
	// xmrig forks to the background again
	// Seems like xmrig doesn't like running GPU in the background
	// if runtime.GOOS == "darwin" || miner.isGPU {
		// runInBackground = false
	// }

	config := XmrigConfig{
		API: XmrigAPIConfig{
			Id:       nil,
			WorkerID: nil,
		},
		HTTP: XmrigHttpConfig{
			Enabled:     true,
			Host:        "127.0.0.1",
			Port:        16000,
			AccessToken: nil,
			Restricted:  true,
		},
		Autosave:   true,
		Version:    1,
		// Background: runInBackground,
		Background: false,
		Colors:     true,
		RandomX: XmrigRandomXConfig{
			Init: -1,
			Numa: true,
		},
		Cpu: XmrigCpuConfig{
			Enabled:    true,
			HugePages:  true,
			HwAes:      nil,
			Priority:   nil,
			Asm:        true,
			Argon2Impl: nil,
			Argon2:     []int{0,1},
			Cn:         [][]int{
				[]int{1,0}, []int{1,1},
			},
			CnHeavy:    [][]int{
				[]int{1,0}, []int{1,1},
			},
			CnLite:     [][]int{
				[]int{1,0}, []int{1,1},
			},
			CnPico:     [][]int{
				[]int{2,0}, []int{2,1},
			},
			CnGpu:     []int{0,1},
			Rx:        []int{0,1},
			RxWow:     []int{0,1},
			Cn0:       false,
			CnLite0:   false,
		},
		DonateLevel:     2,
		DonateOverProxy: 1,
		LogFile:         nil,
		Pools: []XmrigPoolConfig{
			{
				Algo:           XmrigAlgo,
				URL:            poolEndpoint,
				User:           walletAddress,
				Pass:           "x",
				RigId:          "BLOC GUI Miner",
				Nicehash:       false,
				Keepalive:      true,
				Enabled:        true,
				Tls:            false,
				TlsFingerprint: nil,
				Daemon:         false,
			},
		},
		PrintTime:   60,
		Retries:     5,
		RetryPause:  5,
		Syslog:      false,
		UserAgent:   nil,
		Watch:       true,
	}

	return config
}

// createGPUConfig returns creates the config for Xmrig GPU setups
func (miner *Xmrig) createGPUConfig(
	poolEndpoint string,
	walletAddress string,
	XmrigAlgo string,
	XmrigVariant string) XmrigGPUConfig {

	config := XmrigGPUConfig{
		// Algo:        XmrigAlgo,
		// Av:          0,
		// Background:  false,
		// Colors:      true,
		// CPUAffinity: nil,
		// CPUPriority: nil,
		// DonateLevel: 2,
		// LogFile:     nil,
		// PrintTime:   3600,
		// Retries:     5,
		// RetryPause:  5,
		// Safe:        false,
		// Syslog:      false,
		// Pools: []XmrigPoolConfig{
			// {
				// URL:       poolEndpoint,
				// User:      walletAddress,
				// Pass:      "BLOC GUI Miner",
				// Keepalive: true,
				// Nicehash:  false,
				// Variant:   XmrigVariant,
			// },
		// },
		// API: XmrigAPIConfig{
			// Port:        16000,
			// AccessToken: nil,
			// WorkerID:    nil,
		// },
	}

	return config
}

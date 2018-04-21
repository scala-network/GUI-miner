package miner

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// XmrStak implements the miner interface for the xmr-stak miner
// https://github.com/fireice-uk/xmr-stak
type XmrStak struct {
	Base
	name         string
	endpoint     string
	lastHashrate float64
}

// XmrStakResponse contains the data from xmr-stak API
// Generated with https://mholt.github.io/json-to-go/
type XmrStakResponse struct {
	Version  string `json:"version"`
	Hashrate struct {
		Threads [][]interface{} `json:"threads"`
		Total   []float64       `json:"total"`
		Highest float64         `json:"highest"`
	} `json:"hashrate"`
	Results struct {
		DiffCurrent int     `json:"diff_current"`
		SharesGood  int     `json:"shares_good"`
		SharesTotal int     `json:"shares_total"`
		AvgTime     float64 `json:"avg_time"`
		HashesTotal int     `json:"hashes_total"`
		Best        []int   `json:"best"`
		ErrorLog    []struct {
			Count    int    `json:"count"`
			LastSeen int    `json:"last_seen"`
			Text     string `json:"text"`
		} `json:"error_log"`
	} `json:"results"`
	Connection struct {
		Pool     string `json:"pool"`
		Uptime   int    `json:"uptime"`
		Ping     int    `json:"ping"`
		ErrorLog []struct {
			LastSeen int    `json:"last_seen"`
			Text     string `json:"text"`
		} `json:"error_log"`
	} `json:"connection"`
}

// NewXmrStak creates a new xmr-stak miner instance
func NewXmrStak(config Config) (*XmrStak, error) {

	endpoint := config.Endpoint
	if endpoint == "" {
		endpoint = "http://127.0.0.1:16000/api.json"
	}

	miner := XmrStak{
		name:     "xmr-stak",
		endpoint: endpoint,
	}
	miner.Base.executableName = miner.name
	miner.Base.executablePath = config.Path

	return &miner, nil
}

// GetName returns the name of the miner
func (miner *XmrStak) GetName() string {
	return miner.name
}

// GetLastHashrate returns the last reported hashrate
func (miner *XmrStak) GetLastHashrate() float64 {
	return miner.lastHashrate
}

// GetStats returns the current miner stats
func (miner *XmrStak) GetStats() (Stats, error) {
	var stats Stats
	var xmrStats XmrStakResponse
	resp, err := http.Get(miner.endpoint)
	if err != nil {
		return stats, err
	}
	err = json.NewDecoder(resp.Body).Decode(&xmrStats)
	if err != nil {
		return stats, err
	}

	var hashrate float64
	if len(xmrStats.Hashrate.Total) > 0 {
		hashrate = xmrStats.Hashrate.Total[0]
	}
	miner.lastHashrate = hashrate

	var errors []string
	if len(xmrStats.Connection.ErrorLog) > 0 {
		for _, err := range xmrStats.Connection.ErrorLog {
			errors = append(errors, fmt.Sprintf("%s: %s",
				HumanizeTime(err.LastSeen),
				err.Text,
			))
		}
	}
	if len(xmrStats.Results.ErrorLog) > 0 {
		for _, err := range xmrStats.Results.ErrorLog {
			errors = append(errors, fmt.Sprintf("%s: (%d) %s",
				HumanizeTime(err.LastSeen),
				err.Count,
				err.Text,
			))
		}
	}

	stats = Stats{
		Hashrate:          hashrate,
		HashrateHuman:     HumanizeHashrate(hashrate),
		CurrentDifficulty: xmrStats.Results.DiffCurrent,
		Uptime:            xmrStats.Connection.Uptime,
		UptimeHuman:       HumanizeTime(xmrStats.Connection.Uptime),
		SharesGood:        xmrStats.Results.SharesGood,
		SharesBad:         xmrStats.Results.SharesTotal - xmrStats.Results.SharesGood,
		Errors:            errors,
	}
	return stats, nil
}

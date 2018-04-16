package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
)

// PoolData contains the parsed JSON data from the pool list
type PoolData struct {
	ID             int    `json:"id"`
	Rank           int    `json:"rank"`
	APIType        string `json:"api_type"`
	Name           string `json:"name"`
	URL            string `json:"url"`
	Endpoint       string `json:"endpoint"`
	Hashrate       string `json:"hashrate"`
	Miners         int    `json:"miners"`
	LastBlock      string `json:"last_block"`
	Config         string `json:"config"`
	IsEnabled      int    `json:"is_enabled"`
	DisplayInMiner int    `json:"display_in_miner"`
	LastUpdate     string `json:"last_update"`
}

type Stats struct {
	Pool struct {
		ID             int    `json:"id"`
		Rank           int    `json:"rank"`
		APIType        string `json:"api_type"`
		Name           string `json:"name"`
		URL            string `json:"url"`
		Endpoint       string `json:"endpoint"`
		Hashrate       int    `json:"hashrate"`
		Miners         int    `json:"miners"`
		LastBlock      string `json:"last_block"`
		Config         string `json:"config"`
		IsEnabled      int    `json:"is_enabled"`
		DisplayInMiner int    `json:"display_in_miner"`
		LastUpdate     string `json:"last_update"`
	} `json:"pool"`
	Circulation string `json:"circulation"`
	LastBlock   struct {
		ID         int    `json:"id"`
		Height     string `json:"height"`
		Difficulty string `json:"difficulty"`
		TxCount    int    `json:"tx_count"`
		Reward     string `json:"reward"`
		Timestamp  string `json:"timestamp"`
	} `json:"last_block"`
	VolumeCrex      string `json:"volume_crex"`
	VolumeTradeogre string `json:"volume_tradeogre"`
	Volume          string `json:"volume"`
	Price           string `json:"price"`
	MarketCap       string `json:"market_cap"`
	Records         struct {
		Price  string `json:"price"`
		Volume string `json:"volume"`
	} `json:"records"`
	XtlPerDay string  `json:"xtl_per_day"`
	Hashrate  float64 `json:"hashrate"`
}

// GUIInitConfig holds the config entered in the GUI
type GUIInitConfig struct {
	Address string `json:"address"`
	PoolID  int    `json:"pool"`
	Mid     string `json:"mid"`
}

// XmrStakResponse contains the data from xmr-stak API
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
	Address string `json:"address"`
}

// Helper contains some helper functions
type Helper struct {
	MinerAPI string
}

// GetPoolList returns the list of pools available to the GUI miner
func (h *Helper) GetPoolList() ([]PoolData, error) {
	var pools []PoolData
	resp, err := http.Get(fmt.Sprintf("%s/pool-list", h.MinerAPI))
	if err != nil {
		return pools, err
	}
	err = json.NewDecoder(resp.Body).Decode(&pools)
	if err != nil {
		return pools, err
	}

	return pools, nil
}

// GetPool returns a single pool's information
func (h *Helper) GetPool(id int) (PoolData, error) {
	var pool PoolData
	resp, err := http.Get(fmt.Sprintf("%s/pool/%d", h.MinerAPI, id))
	if err != nil {
		return pool, err
	}
	err = json.NewDecoder(resp.Body).Decode(&pool)
	if err != nil {
		return pool, err
	}

	return pool, nil
}

// HumanizeHashrate turns 1000 into 1 KH/s
func (h *Helper) HumanizeHashrate(hashrate string) string {
	hashval, err := strconv.ParseFloat(hashrate, 64)
	if err != nil {
		return "0 H/s"
	}
	if hashval > 1000000.00 {
		return fmt.Sprintf("%.2f MH/s", (hashval / 1000000))
	}
	if hashval > 1000.00 {
		return fmt.Sprintf("%.2f KH/s", (hashval / 1000))
	}
	return fmt.Sprintf("%.2f H/s", hashval)
}

// SaveConfig saves the configuration to disk
func (h *Helper) SaveConfig(config GUIInitConfig) error {

	configBytes, err := json.Marshal(&config)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile("config.json", configBytes, 0644)
	if err != nil {
		return err
	}
	return nil
}

// GetXmrStats returns the local xmr-stak hashrate
func (h *Helper) GetXmrStats() (XmrStakResponse, error) {
	fmt.Println("Getting xmr stats")
	var xmrResponse XmrStakResponse
	resp, err := http.Get("http://127.0.0.1:16000/api.json")
	if err != nil {
		return xmrResponse, err
	}
	err = json.NewDecoder(resp.Body).Decode(&xmrResponse)
	if err != nil {
		return xmrResponse, err
	}

	return xmrResponse, nil
}

// GetStats returns stats for the interface
func (h *Helper) GetStats(poolID int, hashrate float64, mid string) (string, error) {
	fmt.Println("Getting stats")
	// http://stellite.live.local/miner/stats?pool=1&hr=300&mid=minerXXX
	resp, err := http.Get(fmt.Sprintf("%s/stats?pool=%d&hr=%.2f&mid=%s", h.MinerAPI, poolID, hashrate, mid))
	if err != nil {
		return "", err
	}
	statBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(statBytes), nil
}

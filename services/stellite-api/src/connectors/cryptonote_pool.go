package connectors

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

type CryptonotePoolResult struct {
	Config struct {
		Ports []struct {
			Port       int    `json:"port"`
			Difficulty int    `json:"difficulty"`
			Desc       string `json:"desc"`
		} `json:"ports"`
		HashrateWindow       int     `json:"hashrateWindow"`
		Fee                  float64 `json:"fee"`
		Coin                 string  `json:"coin"`
		CoinUnits            int     `json:"coinUnits"`
		CoinDifficultyTarget int     `json:"coinDifficultyTarget"`
		Symbol               string  `json:"symbol"`
		Depth                int     `json:"depth"`
		Donation             struct {
		} `json:"donation"`
		Version             string `json:"version"`
		MinPaymentThreshold int    `json:"minPaymentThreshold"`
		DenominationUnit    int    `json:"denominationUnit"`
	} `json:"config"`
	System struct {
		Load        []float64 `json:"load"`
		NumberCores int       `json:"number_cores"`
	} `json:"system"`
	Pool struct {
		Stats struct {
			LastBlockFound string `json:"lastBlockFound"`
		} `json:"stats"`
		Blocks          []string `json:"blocks"`
		TotalBlocks     int      `json:"totalBlocks"`
		Payments        []string `json:"payments"`
		TotalPayments   int      `json:"totalPayments"`
		TotalMinersPaid int      `json:"totalMinersPaid"`
		Miners          uint32   `json:"miners"`
		Hashrate        uint32   `json:"hashrate"`
		RoundHashes     int64    `json:"roundHashes"`
		LastBlockFound  string   `json:"lastBlockFound"`
	} `json:"pool"`
	Charts struct {
		Hashrate   [][]int `json:"hashrate"`
		Workers    [][]int `json:"workers"`
		Difficulty [][]int `json:"difficulty"`
	} `json:"charts"`
	Network struct {
		Difficulty uint64 `json:"difficulty"`
		Height     uint32 `json:"height"`
		Timestamp  int    `json:"timestamp"`
		Reward     int    `json:"reward"`
		Hash       string `json:"hash"`
	} `json:"network"`
}

// CryptonotePool implements fetching pool stats from a pool running the
// cryptonote pool software
// http://cryptonotemining.org
type CryptonotePool struct {
	Endpoint string
}

// GetStats returns the stats for the given pool
func (pool *CryptonotePool) GetStats() (PoolStats, error) {
	stats := PoolStats{}

	resp, err := http.Get(pool.Endpoint)
	if err != nil {
		return stats, err
	}
	var result CryptonotePoolResult
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return stats, err
	}

	stats.Hashrate = result.Pool.Hashrate
	stats.Miners = result.Pool.Miners
	fmt.Println("HERE")
	result.Pool.LastBlockFound = result.Pool.LastBlockFound[:len(result.Pool.LastBlockFound)-3]
	timestamp, err := strconv.ParseInt(result.Pool.LastBlockFound, 10, 64)
	if err != nil {
		return stats, err
	}
	stats.Height = result.Network.Height
	stats.LastBlockTime = time.Unix(timestamp, 0)
	return stats, nil
}

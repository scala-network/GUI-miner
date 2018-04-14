package connectors

import (
	"encoding/json"
	"net/http"
	"time"
)

type NodeJSPoolResult struct {
	PoolList       []string `json:"pool_list"`
	PoolStatistics struct {
		HashRate           uint32 `json:"hashRate"`
		Miners             uint32 `json:"miners"`
		TotalHashes        int64  `json:"totalHashes"`
		LastBlockFoundTime int64  `json:"lastBlockFoundTime"`
		LastBlockFound     int    `json:"lastBlockFound"`
		TotalBlocksFound   int    `json:"totalBlocksFound"`
		TotalMinersPaid    int    `json:"totalMinersPaid"`
		TotalPayments      int    `json:"totalPayments"`
		RoundHashes        int    `json:"roundHashes"`
	} `json:"pool_statistics"`
	LastPayment int `json:"last_payment"`
}

// NodeJSPool implements fetching pool stats from a pool running the
// nodejs-pool software
// https://github.com/Snipa22/nodejs-pool
type NodeJSPool struct {
	Endpoint string
}

// GetStats returns the stats for the given pool
func (pool *NodeJSPool) GetStats() (PoolStats, error) {
	stats := PoolStats{}

	resp, err := http.Get(pool.Endpoint)
	if err != nil {
		return stats, err
	}

	var result NodeJSPoolResult
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return stats, err
	}

	stats.Hashrate = result.PoolStatistics.HashRate
	stats.Miners = result.PoolStatistics.Miners
	stats.LastBlockTime = time.Unix(result.PoolStatistics.LastBlockFoundTime, 0)
	return stats, nil
}

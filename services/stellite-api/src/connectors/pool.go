package connectors

import "time"

// PoolStats contains all the information from a pool's public exchange
type PoolStats struct {
	Hashrate      uint32
	Height        uint32
	Miners        uint32
	LastBlockTime time.Time
	LastPayment   time.Time
}

// Pool is the interface for getting pool information from
// a mining pool
type Pool interface {
	GetStats() (PoolStats, error)
}

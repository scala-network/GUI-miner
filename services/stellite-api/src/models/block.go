package models

import "time"

// Block is the blockchain block model
type Block struct {
	ID         int64
	Height     int64
	Difficulty int64
	TxCount    int64
	Reward     float64
	Timestamp  time.Time
}

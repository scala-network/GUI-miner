package models

import "time"

// Pool model
type Pool struct {
	ID         uint32
	Rank       uint32
	APIType    string
	Name       string
	URL        string
	Endpoint   string
	Hashrate   uint32
	Miners     uint32
	LastBlock  time.Time
	IsEnabled  uint
	LastUpdate time.Time
}

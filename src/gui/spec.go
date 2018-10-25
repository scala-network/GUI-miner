package gui

import (
	"time"

	"github.com/furiousteam/gui-miner/src/gui/miner"
)

// Config contains the basic configuration for a miner
type Config struct {
	// APIEndpoint is the web endpoint where stats and pools are retrieved from
	APIEndpoint string `json:"api_endpoint"`
	// Address to mine to
	Address string `json:"address"`
	// PoolID selected on startup
	PoolID int `json:"pool"`
	// Mid is the miner identifier
	Mid string `json:"mid"`
	// Miner is the config for the miner
	Miner miner.Config `json:"miner"`
}

// ElectronMessage is marshalled and sent to the UI
type ElectronMessage struct {
	Data string `json:"data"`
}

// frontendConfig is received from the miner's config page
type frontendConfig struct {
	Address string `json:"address"`
	Pool    int    `json:"pool"`
	Threads uint16 `json:"threads"`
	MaxCPU  uint8  `json:"max_cpu"`
}

// PoolData contains the parsed JSON data from the pool list
type PoolData struct {
	ID             int    `json:"id"`
	Rank           int    `json:"rank"`
	APIType        string `json:"api_type"`
	Name           string `json:"name"`
	URL            string `json:"url"`
	Endpoint       string `json:"endpoint"`
	Hashrate       string `json:"hashrate"`
	Miners         string `json:"miners"`
	LastBlock      string `json:"last_block"`
	Config         string `json:"config"`
	IsEnabled      int    `json:"is_enabled"`
	DisplayInMiner int    `json:"display_in_miner"`
	LastUpdate     string `json:"last_update"`
}

// GlobalStats contains the current stats for the network,
// trading and selected mining pool
type GlobalStats struct {
	Pool        PoolData `json:"pool"`
	Circulation string   `json:"circulation"`
	LastBlock   struct {
		ID         int    `json:"id"`
		Height     int    `json:"height"`
		Difficulty int    `json:"difficulty"`
		TxCount    int    `json:"tx_count"`
		Reward     string `json:"reward"`
		Timestamp  string `json:"timestamp"`
	} `json:"last_block"`
	Difficulty      string `json:"difficulty"`
	Height          string `json:"height"`
	VolumeCrex      string `json:"volume_stex"`
	VolumeTradeogre string `json:"volume_tradeogre"`
	Volume          string `json:"volume"`
	Price           string `json:"price"`
	MarketCap       string `json:"market_cap"`
	Records         struct {
		Price  string `json:"price"`
		Volume string `json:"volume"`
	} `json:"records"`
	BlocPerDay string `json:"bloc_per_day"`
	Hashrate  string `json:"hashrate"`
	// PoolHTML is injected before sending the update to the front-end. Avoids
	// having to send extra packets
	PoolHTML string `json:"pool_html"`
}

// Announcement is the structure returned is an announcement is made
// available
type Announcement struct {
	ID         int       `json:"id"`
	Text       string    `json:"text"`
	Link       string    `json:"link"`
	DateString string    `json:"date"`
	Date       time.Time `json:"-"`
	Ann        bool      `json:"ann"`
}

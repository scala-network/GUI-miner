package gui

import (
	"time"

	"github.com/furiousteam/BLOC-GUI-Miner/src/gui/miner"
)

// Config contains the basic configuration for a miner
type Config struct {
	// APIEndpoint is the web endpoint where stats and pools are retrieved from
	APIEndpoint string `json:"api_endpoint"`
	// CoinType is the type of the coin the miner is currently minning
	CoinType string `json:"coin_type"`
	// CoinAlgo is the algo of the coin the miner is currently minning
	CoinAlgo string `json:"coin_algo"`
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
	Address  string `json:"address"`
	Pool     int    `json:"pool"`
	Threads  uint16 `json:"threads"`
	MaxCPU   uint8  `json:"max_cpu"`
	CoinType string `json:"coin_type"`
	CoinAlgo string `json:"coin_algo"`
}

// coinsContentJson is received from github "BLOC-GUI-Miner/coins/content.json"
type coinsContentJson struct {
	Coins []struct {
		CoinType string `json:"coin_type"`
		CoinAlgo string `json:"coin_algo"`
	} `json:"coins"`
	Names             map[string]interface{} `json:"names"`
	Abbr2             map[string]interface{} `json:"abbr"`
	AddressPrefix     map[string]interface{} `json:"address_prefix"`
	AddressValidation map[string]interface{} `json:"address_validation"`
	MainBg            map[string]interface{} `json:"mainBg"`
	Logo              map[string]interface{} `json:"logo"`
	DownloadPage      map[string]interface{} `json:"downloadPage"`
	SocialLinks       map[string]interface{} `json:"socialLinks"`
	NetworkLinks      map[string]interface{} `json:"networkLinks"`
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
	Fee            string `json:"fee"`
	Miners         string `json:"miners"`
	Payout         string `json:"payout"`
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
	Ticker      string   `json:"abbreviation"`
	Supply      string   `json:"maximum_supply"`
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
	Volumes         []struct {
		Name  string `json:"name"`
		Value string `json:"value"`
		Unit  string `json:"unit"`
	} `json:"volumes"`
	Prices         []struct {
		Name  string `json:"name"`
		Value string `json:"value"`
		Unit  string `json:"unit"`
	} `json:"prices"`
	Volume          string `json:"volume"`
	Price           string `json:"price"`
	MarketCap       string `json:"market_cap"`
	Records         struct {
		Price  string `json:"price"`
		Volume string `json:"volume"`
	} `json:"records"`
	CoinsPerDay string `json:"coins_per_day"`
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

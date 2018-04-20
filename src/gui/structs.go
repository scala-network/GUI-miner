package gui

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
	Miners         int    `json:"miners"`
	LastBlock      string `json:"last_block"`
	Config         string `json:"config"`
	IsEnabled      int    `json:"is_enabled"`
	DisplayInMiner int    `json:"display_in_miner"`
	LastUpdate     string `json:"last_update"`
}

// Stats contains the current stats for the network, trading and selected pool
type Stats struct {
	Pool        PoolData `json:"pool"`
	Circulation string   `json:"circulation"`
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
	XtlPerDay string `json:"xtl_per_day"`
	Hashrate  string `json:"hashrate"`
	// PoolHTML is injected before sending the update to the front-end. Avoids
	// having to send to packets
	PoolHTML string `json:"pool_html"`
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
	Address     string `json:"address"`
	UpdateGraph bool   `json:"update_graph"`
}

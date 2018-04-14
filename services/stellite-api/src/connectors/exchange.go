package connectors

// Ticker holds the information for exchange trades
type Ticker struct {
	Last      float64
	High      float64
	Low       float64
	VolumeBTC float64
}

// Exchange is the interface for getting trade information from
// an exchange
type Exchange interface {
	GetName() string
	GetTicker() (Ticker, error)
}

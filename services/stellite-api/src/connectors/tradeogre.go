package connectors

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

// TradeOgreTicker implements the TradeOgre ticker response format
type TradeOgreTicker struct {
	InitialPrice string `json:"initialprice"`
	Price        string `json:"price"`
	High         string `json:"high"`
	Low          string `json:"low"`
	Volume       string `json:"volume"`
}

// TradeOgre retrieves trade information from https://tradeogre.com/
type TradeOgre struct {
	Endpoint string
}

// GetName returns the name of the exchange
func (exchange *TradeOgre) GetName() string {
	return "TradeOgre"
}

// GetTicker fetches the latest trade information for the exchange and
// trading pair
func (exchange *TradeOgre) GetTicker() (Ticker, error) {
	var ticker Ticker

	response, err := http.Get(fmt.Sprintf(exchange.Endpoint, "BTC", "XTL"))
	if err != nil {
		return ticker, err
	}

	var tradeOgreTicker TradeOgreTicker
	err = json.NewDecoder(response.Body).Decode(&tradeOgreTicker)
	if err != nil {
		return ticker, err
	}

	ticker.Last, err = strconv.ParseFloat(tradeOgreTicker.Price, 64)
	if err != nil {
		return ticker, nil
	}

	ticker.High, err = strconv.ParseFloat(tradeOgreTicker.High, 64)
	if err != nil {
		return ticker, nil
	}

	ticker.Low, err = strconv.ParseFloat(tradeOgreTicker.Low, 64)
	if err != nil {
		return ticker, nil
	}

	ticker.VolumeBTC, err = strconv.ParseFloat(tradeOgreTicker.Volume, 64)
	if err != nil {
		return ticker, nil
	}

	return ticker, nil
}

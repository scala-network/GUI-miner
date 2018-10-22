package miner_test

import (
	"testing"

	"github.com/furiousteam/gui-miner/src/gui/miner"
)

// TestHumanizeTime tests if conversion from seconds to mintes and hours
// is correct
func TestHumanizeTime(t *testing.T) {
	tests := map[int]string{
		0:    "0 seconds",
		1:    "1 second",
		2:    "2 seconds",
		60:   "1 minute",
		80:   "1 minute",
		120:  "2 minutes",
		160:  "2 minutes",
		3600: "1 hour",
		4000: "1 hour",
		7200: "2 hours",
		8000: "2 hours",
	}

	for seconds, expected := range tests {
		actual := miner.HumanizeTime(seconds)
		if actual != expected {
			t.Errorf("Incorrect result for %d second(s). Got '%s', expected '%s'",
				seconds,
				actual,
				expected)
		}
	}
}

// TestHumanizeHashrate tests if the hashrate to H/s, KH/s and MH/s is correct
func TestHumanizeHashrate(t *testing.T) {
	tests := map[float64]string{
		0:       "0.00 H/s",
		1:       "1.00 H/s",
		500.2:   "500.20 H/s",
		1000:    "1.00 KH/s",
		1500:    "1.50 KH/s",
		18300:   "18.30 KH/s",
		250400:  "250.40 KH/s",
		1300000: "1.30 MH/s",
	}

	for hashrate, expected := range tests {
		actual := miner.HumanizeHashrate(hashrate)
		if actual != expected {
			t.Errorf("Incorrect result for %.2f hashes/s. Got '%s', expected '%s'",
				hashrate,
				actual,
				expected)
		}
	}
}

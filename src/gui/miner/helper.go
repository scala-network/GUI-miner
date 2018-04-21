package miner

import (
	"fmt"
	"strings"
)

// CreateMiner creates a supported miner from the given configuration
func CreateMiner(config Config) (Miner, error) {
	switch strings.ToLower(config.Type) {
	case "xmr-stak":
		return NewXmrStak(config)
	}

	return nil,
		fmt.Errorf(
			"'%s' is not a supported miner. Supported miners are %s",
			config.Type,
			strings.Join(SupportedMiners, ","))
}

// HumanizeHashrate returns the H/s, KH/s or MH/s representation of hashrate
func HumanizeHashrate(hashrate float64) string {
	if hashrate >= 1000000 {
		return fmt.Sprintf("%.2f MH/s", hashrate/1000000)
	}
	if hashrate >= 1000 {
		return fmt.Sprintf("%.2f KH/s", hashrate/1000)
	}
	return fmt.Sprintf("%.0f H/s", hashrate)
}

// HumanizeTime turns seconds into minutes, hours, etc
func HumanizeTime(seconds int) string {
	var humanTime string
	minutes := seconds / 60
	hours := minutes / 60

	if hours > 0 {
		humanTime = fmt.Sprintf("%d hour", hours)
		if hours != 1 {
			humanTime += "s"
		}
		return humanTime
	}
	if minutes > 0 {
		humanTime = fmt.Sprintf("%d minute", minutes)
		if minutes != 1 {
			humanTime += "s"
		}
		return humanTime
	}

	humanTime = fmt.Sprintf("%d second", seconds)
	if seconds != 1 {
		humanTime += "s"
	}
	return humanTime
}

package miner

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"runtime"
	"strings"
)

// CreateMiner creates a supported miner from the given configuration
func CreateMiner(config Config) (Miner, error) {
	switch strings.ToLower(config.Type) {
	case "xmr-stak":
		return NewXmrStak(config)
	case "xmrig":
		return NewXmrig(config)
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
	return fmt.Sprintf("%.2f H/s", hashrate)
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

// DetermineMinerType checks the given path for supported miners
// and returns the type of miner and path to the executable
func DetermineMinerType(dir string) (string, string, error) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return "", "", err
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}
		if runtime.GOOS == "windows" {
			if strings.Contains(strings.ToLower(file.Name()), "exe") == false {
				continue
			}
		} else {
			if strings.Contains(file.Mode().Perm().String(), "x") == false {
				continue
			}
		}

		fileName := strings.ToLower(file.Name())
		for _, supportedMiner := range SupportedMiners {
			if strings.Contains(fileName, supportedMiner) {
				return supportedMiner, filepath.Join(dir, fileName), nil
			}
		}
	}
	return "", "", fmt.Errorf("No supported miner was found in '%s'", dir)
}

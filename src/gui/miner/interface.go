package miner

// Miner defines the required behaviour to be implemented by a miner to
// work with the GUI
type Miner interface {
	// Start the miner
	Start() error
	// Stop the miner
	Stop() error
	// WriteConfig writes the miner's configuration to the file format as
	// specified by the miner
	WriteConfig(
		poolEndpoint string,
		walletAddress string,
		coinAlgorithm string,
		processingConfig ProcessingConfig) error
	// GetProcessingConfig returns the current miner processing config
	// TODO: Currently only CPU threads, extend this to full CPU/GPU config
	GetProcessingConfig() ProcessingConfig
	// GetName returns the name of the miner
	GetName() string
	// GetLastHashrate returns the last reported hashrate
	GetLastHashrate() float64
	// GetStats returns the current miner stats
	GetStats() (Stats, error)
}

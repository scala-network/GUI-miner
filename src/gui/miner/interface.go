package miner

// Miner defines the required behaviour to be implemented by a miner to
// work with the GUI
type Miner interface {
	// Start the miner
	Start() error
	// Stop the miner
	Stop() error
}

package miner

import (
	"os/exec"
)

// Base implements core functionality common to all miners
type Base struct {
	executableName string
	executablePath string
	command        *exec.Cmd
}

// Start the miner
func (b *Base) Start() error {
	params := []string{}
	commandName := b.executableName
	commandDir := b.executablePath
	b.command = exec.Command(commandName, params...)
	b.command.Dir = commandDir
	return b.command.Start()
}

// Stop the miner
func (b *Base) Stop() error {
	if b.command != nil {
		return b.command.Process.Kill()
	}
	return nil
}

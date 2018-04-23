package miner

import (
	"fmt"
	"os"
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
	commandName := fmt.Sprintf(".%c%s", os.PathSeparator, b.executableName)
	fmt.Println("CommandName:", commandName)
	commandDir := b.executablePath
	fmt.Println("CommandPath:", commandDir)
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

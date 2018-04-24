package miner

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	ps "github.com/mitchellh/go-ps"
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
	commandDir := b.executablePath
	b.command = exec.Command(commandName, params...)
	b.command.Dir = commandDir
	return b.command.Start()
}

// Stop the miner
func (b *Base) Stop() error {
	if b.command != nil {
		// Some of the miners fork in a way that we loose track of the actual
		// miner's pid. To make sure the miner is stopped, we find all processes
		// that match the original executable name
		processes, err := ps.Processes()
		if err != nil {
			// If for some reason we can't get the process list, we use the
			// standard kill available
			return b.command.Process.Kill()
		}
		for _, process := range processes {
			if strings.Contains(strings.ToLower(process.Executable()), b.executableName) {
				p, err := os.FindProcess(process.Pid())
				if err != nil {
					// If the process is in the list, but we can't find it by Pid, then
					// it probably died or something weird is going on
					return err
				}
				// Kill the process we found, then continue searching - just in case
				// there is still others lingering around. Not worried about any errors
				// here since there is nothing we can do about it at this point
				_ = p.Kill()
			}
		}

	}
	return nil
}

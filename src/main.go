// Package main implements the communication and control parts of the
// GUI miner. It constructs and launches the Electron front-end
package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"runtime"

	"github.com/donovansolms/stellite-gui-miner/src/gui"
)

// AppName is injected by the Astilectron packager
var AppName string

// BuiltAt is injected by the Astilectron packager
var BuiltAt string

// main implements the main runnable of the application
func main() {
	// Grab the command-line flags
	debug := flag.Bool("d", false, "Enable debug mode")
	flag.Parse()

	// We need to get the acutal working directory to ensure proper operation
	// on all platforms

	// TODO: Add Back HACK os.Executable()
	workingDir, err := os.Executable()
	//workingDir, err := os.Getwd()
	if err != nil {
		log.Fatalf("Can't read current directory: %s", err)
	}
	workingDir = filepath.Dir(workingDir)
	if err != nil {
		log.Fatalf("Can't format current directory: %s", err)
	}
	if runtime.GOOS == "darwin" {
		// Mac executes from within the .app/Content/MacOS folder, this moves
		// the folder back to the actual app
		workingDir, err = filepath.Abs(workingDir + "/../../..")
		if err != nil {
			log.Fatalf("Can't update current directory: %s", err)
		}
	}

	var config *gui.Config
	var apiEndpoint string
	fileBytes, err := ioutil.ReadFile(filepath.Join(workingDir, "config.json"))
	if err == nil {
		err = json.Unmarshal(fileBytes, &config)
		if err != nil {
			panic(err)
		}
		apiEndpoint = config.APIEndpoint
	} else {
		config = nil
		// Not set yet, set to default
		apiEndpoint = "http://stellite.live.local/miner"
		//apiEndpoint = "https://www.stellite.live/miner"
	}

	// Create the miner
	// AppName, Asset and RestoreAssets are injected by the bundler
	gui, err := gui.New(
		AppName,
		config,
		Asset,
		RestoreAssets,
		apiEndpoint,
		workingDir,
		*debug,
	)
	if err != nil {
		// Setting the output to stdout so the user can see the error
		log.SetOutput(os.Stdout)
		log.Fatalf("Unable to set up miner: %s", err)
	}

	err = gui.Run()
	if err != nil {
		// Setting the output to stdout so the user can see the error
		log.SetOutput(os.Stdout)
		log.Fatalf("Unable to run miner: %s", err)
	}
}

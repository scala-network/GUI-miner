// Package main implements the communication and control parts of the
// GUI miner. It constructs and launches the Electron front-end
package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"os"

	"github.com/donovansolms/stellite-gui-miner/src/gui"
)

// AppName is injected by the Astilectron packager
var AppName string

// BuiltAt is injected by the Astilectron packager
var BuiltAt string

// main implements the main runnable of the application
func main() {
	// Grab the command-line flags
	debug := flag.Bool("d", false, "Enabled debug mode")
	flag.Parse()

	var config *gui.Config
	fileBytes, err := ioutil.ReadFile("./config.json")
	if err == nil {
		err = json.Unmarshal(fileBytes, &config)
		if err != nil {
			panic(err)
		}
	} else {
		config = nil
	}
	//apiEndpoint := "https://www.stellite.live/miner"
	apiEndpoint := "http://stellite.live.local/miner"

	// Create the miner
	// AppName, Asset and RestoreAssets are injected by the bundler
	gui, err := gui.New(
		AppName,
		config,
		Asset,
		RestoreAssets,
		apiEndpoint,
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

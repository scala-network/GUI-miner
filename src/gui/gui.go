package gui

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"time"

	astilectron "github.com/asticode/go-astilectron"
	bootstrap "github.com/asticode/go-astilectron-bootstrap"
	"github.com/donovansolms/stellite-gui-miner/src/gui/miner"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

// GUI implements the core control for the GUI miner
type GUI struct {
	// window is the main Astilectron window
	window *astilectron.Window
	// // minerCmd is a reference to the xmr-stak miner process
	// minerCmd *exec.Cmd
	// astilectronOptions holds the Astilectron options
	astilectronOptions bootstrap.Options
	// config for the miner
	config *Config
	// apiEndpoint is the web endpoint where stats and pools are retrieved from
	apiEndpoint string
	// miner is the selected miner backend as chosen by the user
	miner miner.Miner
	// logger logs to stdout
	logger *logrus.Entry
	// currentHashrate of the user if mining
	lastHashrate float64
	// captureMiningStats determines if the mining stats loop should be running
	// or not
	// captureMiningStats uint32
}

// New creates a new instance of the miner application
func New(
	appName string,
	config *Config,
	asset bootstrap.Asset,
	restoreAssets bootstrap.RestoreAssets,
	apiEndpoint string,
	isDebug bool) (*GUI, error) {

	if apiEndpoint == "" {
		return nil, errors.New("The API Endpoint must be specified")
	}

	gui := GUI{
		apiEndpoint: apiEndpoint,
		config:      config,
	}

	// If no config is specified then this is the first run
	startPage := "firstrun.html"
	if gui.config != nil {
		startPage = "index.html"
		// Already configured, set up the miner
		var err error
		gui.miner, err = miner.CreateMiner(gui.config.Miner)
		if err != nil {
			return nil,
				fmt.Errorf("Unable to use '%s' as miner: %s", gui.config.Miner.Type, err)
		}
	}

	gui.astilectronOptions = bootstrap.Options{
		Debug:         isDebug,
		Asset:         asset,
		RestoreAssets: restoreAssets,
		Homepage:      startPage,
		AstilectronOptions: astilectron.Options{
			AppName:            appName,
			AppIconDarwinPath:  "resources/icon.icns",
			AppIconDefaultPath: "resources/icon.png",
		},
		WindowOptions: &astilectron.WindowOptions{
			BackgroundColor: astilectron.PtrStr("#0B0C22"),
			Center:          astilectron.PtrBool(true),
			Height:          astilectron.PtrInt(700),
			Width:           astilectron.PtrInt(1175),
		},
		MenuOptions: []*astilectron.MenuItemOptions{{
			Label: astilectron.PtrStr("File"),
			SubMenu: []*astilectron.MenuItemOptions{
				{
					Role: astilectron.MenuItemRoleClose,
				},
			},
		}},
		// OnWait is triggered as soon as the electron window is ready and running
		OnWait: func(
			_ *astilectron.Astilectron,
			window *astilectron.Window,
			_ *astilectron.Menu,
			_ *astilectron.Tray,
			_ *astilectron.Menu) error {
			gui.window = window
			go func() {
				gui.logger.Info("Start capturing stats")
				// Run forever!
				for {
					gui.updateNetworkStats()
					time.Sleep(time.Second * 30)
				}
			}()
			return nil
		},
		MessageHandler: gui.handleElectronCommands,
	}

	// Setup the logging, by default we log to stdout
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "Jan 02 15:04:05",
	})
	logrus.SetLevel(logrus.InfoLevel)
	logrus.SetOutput(os.Stdout)
	if isDebug {
		logrus.SetLevel(logrus.DebugLevel)
	}
	// Setting the WithFields now will ensure all log entries from this point
	// includes the fields
	gui.logger = logrus.WithFields(logrus.Fields{
		"service": "stellite-gui-miner",
	})

	gui.logger.Info("Setup complete")
	return &gui, nil
}

// Run the miner!
func (gui *GUI) Run() error {
	gui.logger.Info("Starting miner")
	err := bootstrap.Run(gui.astilectronOptions)
	if err != nil {
		return err
	}
	// If xmr-stak is running, kill it
	/*if gui.minerCmd != nil {
		err = gui.minerCmd.Process.Kill()
		if err != nil {
			gui.logger.Fatalf("Unable to stop xmr-stak: %s", err)
		}
	}*/
	return nil
}

// handleElectronCommands handles the messages sent by the Electron front-end
func (gui *GUI) handleElectronCommands(
	_ *astilectron.Window,
	command bootstrap.MessageIn) (interface{}, error) {

	gui.logger.WithField(
		"command", command.Name,
	).Debug("Received command from Electron")

	// Every Electron command has a name together with a payload containing the
	// actual message
	switch command.Name {

	// Firstrun is received on the first run of the miner. We return the current
	// logged in username
	case "firstrun":
		var username string
		currentUser, err := user.Current()
		if err == nil {
			if currentUser.Name != "" {
				username = currentUser.Name
			} else if currentUser.Username != "" {
				username = currentUser.Username
			}
		}
		return username, nil

	// pool-list requests the recommended pool list from the miner API
	// and returns the rendered HTML
	case "pool-list":
		// Grab the pool list and send that to the GUI as well
		poolJSONs, err := gui.GetPoolList()
		if err != nil {
			_ = gui.sendElectronCommand("fatal_error", ElectronMessage{
				Data: fmt.Sprintf("Unable to fetch pool list from API."+
					"Please check that you are connected to the internet and try again."+
					"<br/>The error was '%s'\"}", err),
			})
			// Give the UI some time to display the message
			time.Sleep(time.Second * 15)
			gui.logger.Fatalf("Unable to fetch pool list: '%s'", err)
		}
		poolTemplate, err := gui.GetPoolTemplate(false)
		if err != nil {
			log.Fatalf("Unable to load pool template: '%s'", err)
		}
		var poolsList string
		for _, poolData := range poolJSONs {
			var templateHTML bytes.Buffer
			// Get the string time in the correct format
			//t, _ := time.Parse("2006-01-02 15:04",
			//	poolData.LastBlock[:len(poolData.LastBlock)-3])
			//since := time.Since(t)
			//poolData.LastBlock = fmt.Sprintf("%d minutes ago", int(since.Minutes()))
			err = poolTemplate.Execute(&templateHTML, poolData)
			if err != nil {
				log.Fatalf("Unable to load pool template: '%s'", err)
			}
			poolsList += templateHTML.String()
		}
		return poolsList, nil

	// configure is sent after the firstrun setup has been completed
	case "configure":
		// HACK: Adding a slight delay before switching to the mining dashboard
		// after initial setup to have the user at least see the 'configure' message
		time.Sleep(time.Second * 5)
		gui.configureMiner(command, false)
		return "Ok", nil

	// reconfigure is sent after settings are changes by the user
	case "reconfigure":
		gui.logger.Info("Reconfiguring miner")
		err := gui.stopMiner()
		if err != nil {
			_ = gui.sendElectronCommand("fatal_error", ElectronMessage{
				Data: fmt.Sprintf("Unable to stop miner for reconfigure."+
					"Please close the miner and open it again."+
					"<br/>The error was '%s'\"}", err),
			})
			// Give the UI some time to display the message
			time.Sleep(time.Second * 15)
			gui.logger.Fatalf("Unable to reconfigure miner: '%s'", err)
		}
		// TODO HACK Wait for the mining stats loop to exit
		time.Sleep(time.Second * 10)
		gui.configureMiner(command, true)
		gui.startMiner()
		gui.logger.Info("Miner reconfigured")

	// miner_start is sent after configuration or when the user
	// clicks 'start mining'
	case "miner_start":
		gui.startMiner()

	// miner_stop is sent whenever the user clicks 'stop mining'
	case "miner_stop":
		_ = gui.stopMiner()
	}
	return nil, fmt.Errorf("'%s' is an unknown command", command.Name)
}

// sendElectronCommand sends the given data to Electron under the command name
func (gui *GUI) sendElectronCommand(
	name string,
	data ElectronMessage) error {
	dataBytes, err := json.Marshal(&data)
	if err != nil {
		return err
	}
	return bootstrap.SendMessage(gui.window, name, string(dataBytes))
}

// configureMiner creates the xmr-stak configuration to use
func (gui *GUI) configureMiner(command bootstrap.MessageIn, isReconfigure bool) {
	err := json.Unmarshal(command.Payload, &gui.config)
	if err != nil {
		_ = gui.sendElectronCommand("fatal_error", ElectronMessage{
			Data: fmt.Sprintf("Unable to configure miner."+
				"Please check your configuration is value."+
				"<br/>The error was '%s'\"}", err),
		})
		// Give the UI some time to display the message
		time.Sleep(time.Second * 15)
		gui.logger.Fatalf("Unable to configure miner: '%s'", err)
	}
	if isReconfigure == false {
		gui.config.Mid = uuid.New().String()
	}
	err = gui.SaveConfig(*gui.config)
	if err != nil {
		_ = gui.sendElectronCommand("fatal_error", ElectronMessage{
			Data: fmt.Sprintf("Unable to configure miner."+
				"Please check that you can write to the miner's installation path."+
				"<br/>The error was '%s'\"}", err),
		})
		// Give the UI some time to display the message
		time.Sleep(time.Second * 15)
		gui.logger.Fatalf("Unable to configure miner: '%s'", err)
	}

	err = ioutil.WriteFile("./xmr-stak/config.txt",
		[]byte(gui.GetXmrStakConfig()),
		0644)
	if err != nil {
		_ = gui.sendElectronCommand("fatal_error", ElectronMessage{
			Data: fmt.Sprintf("Unable to configure miner."+
				"Please check that you can write to the miner's installation path."+
				"<br/>The error was '%s'\"}", err),
		})
		// Give the UI some time to display the message
		time.Sleep(time.Second * 15)
		gui.logger.Fatalf("Unable to configure miner: '%s'", err)
	}

	var poolInfo PoolData
	poolInfo, err = gui.GetPool(gui.config.PoolID)
	if err != nil {
		_ = gui.sendElectronCommand("fatal_error", ElectronMessage{
			Data: fmt.Sprintf("Unable to configure miner."+
				"Please check that you are connected to the internet."+
				"<br/>The error was '%s'\"}", err),
		})
		// Give the UI some time to display the message
		time.Sleep(time.Second * 15)
		gui.logger.Fatalf("Unable to configure miner: '%s'", err)
	}

	err = ioutil.WriteFile("./xmr-stak/pools.txt",
		[]byte(gui.GetXmrStakPoolConfig(poolInfo.Config, gui.config.Address)),
		0644)
	if err != nil {
		_ = gui.sendElectronCommand("fatal_error", ElectronMessage{
			Data: fmt.Sprintf("Unable to configure miner."+
				"Please check that you can write to the miner's installation path."+
				"<br/>The error was '%s'\"}", err),
		})
		// Give the UI some time to display the message
		time.Sleep(time.Second * 15)
		gui.logger.Fatalf("Unable to configure miner: '%s'", err)
	}
}

// startMiner starts the xmr-stak miner
func (gui *GUI) startMiner() {
	err := gui.miner.Start()
	if err != nil {
		_ = gui.sendElectronCommand("fatal_error", ElectronMessage{
			Data: fmt.Sprintf("Unable to start '%s' miner, please check that you "+
				"can run the miner from your installation directory."+
				"<br/>The error was '%s'\"}", gui.miner.GetName(), err),
		})
		// Give the UI some time to display the message
		time.Sleep(time.Second * 15)
		gui.logger.Fatalf("Error starting '%s': %s", gui.miner.GetName(), err)
	}
	gui.logger.Info("Started '%s' miner", gui.miner.GetName())
	// TODO go gui.updateMiningStatsLoop()
}

// stopMiner stops the xmr-stak miner
func (gui *GUI) stopMiner() error {
	err := gui.miner.Stop()
	if err != nil {
		_ = gui.sendElectronCommand("fatal_error", ElectronMessage{
			Data: fmt.Sprintf("Unable to stop miner.."+
				"Please close the GUI miner and open it again."+
				"<br/>The error was '%s'\"}", err),
		})
		gui.logger.Errorf("Unable to stop miner '%s': %s", gui.miner.GetName(), err)
		return err
	}
	gui.logger.Info("Stopped '%s' miner", gui.miner.GetName())
	return nil
}

// updateNetworkStats is a single stat update for network and payment info
func (gui *GUI) updateNetworkStats() {
	gui.logger.WithField(
		"hashrate", gui.lastHashrate,
	).Debug("Fetching network stats")
	// On firstrun we fon't have a config yet
	if gui.config == nil {
		gui.logger.Warning("No config set yet")
		return
	}
	stats, err := gui.GetStats(gui.config.PoolID, gui.lastHashrate, gui.config.Mid)
	if err != nil {
		gui.logger.Warningf("Unable to get network stats: %s", err)
	} else {
		err := bootstrap.SendMessage(gui.window, "network_stats", stats)
		if err != nil {
			gui.logger.Errorf("Unable to send stats to front-end: %s", err)
		}
	}
}

// updateMiningStats retrieves the miner's stats from xmr-stak and updates
// the front-end
func (gui *GUI) updateMiningStatsLoop() {
	/*atomic.StoreUint32(&gui.captureMiningStats, 1)
	lastGraphUpdate := time.Now()
	for atomic.LoadUint32(&gui.captureMiningStats) == 1 {
		gui.logger.Debug("Fetching mining stats")
		xmrStats, err := gui.GetXmrStats()
		if err != nil {
			gui.logger.Warningf("Unable to get mining stats, xmr-stak not available yet?: %s", err)
		} else {
			if len(xmrStats.Hashrate.Total) > 0 {
				gui.lastHashrate = xmrStats.Hashrate.Total[0]
				// The first time we get a hashrate, update the XTL amount so that the
				// user doesn't think it doesn't work
				gui.updateNetworkStats()
			}
			xmrStats.Address = gui.config.Address

			if time.Since(lastGraphUpdate).Minutes() >= 1 {
				lastGraphUpdate = time.Now()
				xmrStats.UpdateGraph = true
			}
			statBytes, _ := json.Marshal(&xmrStats)
			err = bootstrap.SendMessage(gui.window, "miner_stats", string(statBytes))
			if err != nil {
				gui.logger.Errorf("Unable to send miner stats to front-end: %s", err)
			}
		}
		time.Sleep(time.Second * 10)
	}*/
	gui.logger.Debug("Stopped fetching mining stats")
}

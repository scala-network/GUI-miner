package gui

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"runtime"
	"time"

	astilectron "github.com/asticode/go-astilectron"
	bootstrap "github.com/asticode/go-astilectron-bootstrap"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/furiousteam/gui-miner/src/gui/miner"
)

// GUI implements the core control for the GUI miner
type GUI struct {
	// window is the main Astilectron window
	window *astilectron.Window
	// astilectronOptions holds the Astilectron options
	astilectronOptions bootstrap.Options
	// config for the miner
	config *Config
	// miner is the selected miner backend as chosen by the user
	miner miner.Miner
	// logger logs to stdout
	logger *logrus.Entry
	// workingDir holds the current working directory
	workingDir string
	// currentHashrate of the user if mining
	lastHashrate float64
	// miningStatsTicker controls the interval for fetching mining stats from
	// the selected miner
	miningStatsTicker *time.Ticker
	// networkStatsTicker controls the interval for fetching network, trading
	// and other stats
	networkStatsTicker *time.Ticker
	// annTicker controls the interval for checking for announements
	annTicker *time.Ticker
	// annProcessed keeps track of announcements already processed and displayed
	annProcessed map[int]bool
}

// New creates a new instance of the miner application
func New(
	appName string,
	config *Config,
	asset bootstrap.Asset,
	restoreAssets bootstrap.RestoreAssets,
	apiEndpoint string,
	workingDir string,
	isDebug bool) (*GUI, error) {

	if apiEndpoint == "" {
		return nil, errors.New("The API Endpoint must be specified")
	}

	gui := GUI{
		config:       config,
		workingDir:   workingDir,
		annProcessed: make(map[int]bool),
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
	} else {
		// Nothing has been configured yet, set some defaults
		gui.config = &Config{
			APIEndpoint: apiEndpoint,
			Mid:         uuid.New().String(),
		}
	}
	var menu []*astilectron.MenuItemOptions

	// Setup the logging, by default we log to stdout
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "Jan 02 15:04:05",
	})
	logrus.SetLevel(logrus.InfoLevel)

	logrus.SetOutput(os.Stdout)

	// Create the window options
	windowOptions := astilectron.WindowOptions{
		// If frame is false, the window frame is removed. If isDebug is true,
		// we show the frame to have debugging options available
		Frame:           astilectron.PtrBool(isDebug),
		BackgroundColor: astilectron.PtrStr("#0B0C22"),
		Center:          astilectron.PtrBool(true),
		Height:          astilectron.PtrInt(700),
		Width:           astilectron.PtrInt(1175),
	}

	if isDebug {
		logrus.SetLevel(logrus.DebugLevel)
		debugLog, err := os.OpenFile(
			filepath.Join(gui.workingDir, "debug.log"),
			os.O_CREATE|os.O_TRUNC|os.O_WRONLY,
			0644)
		if err != nil {
			panic(err)
		}
		// TODO: logrus.SetOutput(debugLog)
		_ = debugLog

		// We only show the menu bar in debug mode
		menu = append(menu, &astilectron.MenuItemOptions{
			Label: astilectron.PtrStr("File"),
			SubMenu: []*astilectron.MenuItemOptions{
				{
					Role: astilectron.MenuItemRoleClose,
				},
			},
		})
	}
	// To make copy and paste work on Mac, the copy and paste entries need to
	// be defined, the alternative is to implement the clipboard API
	// https://github.com/electron/electron/blob/master/docs/api/clipboard.md
	if runtime.GOOS == "darwin" {
		menu = append(menu, &astilectron.MenuItemOptions{
			Label: astilectron.PtrStr("Edit"),
			SubMenu: []*astilectron.MenuItemOptions{
				{
					Role: astilectron.MenuItemRoleCut,
				},
				{
					Role: astilectron.MenuItemRoleCopy,
				},
				{
					Role: astilectron.MenuItemRolePaste,
				},
				{
					Role: astilectron.MenuItemRoleSelectAll,
				},
			},
		})

		windowOptions.Frame = astilectron.PtrBool(isDebug)
		windowOptions.TitleBarStyle = astilectron.PtrStr("hidden")
	}

	// Setting the WithFields now will ensure all log entries from this point
	// includes the fields
	gui.logger = logrus.WithFields(logrus.Fields{
		"service": "bloc-gui-miner",
	})

	gui.astilectronOptions = bootstrap.Options{
		Debug:         isDebug,
		Asset:         asset,
		RestoreAssets: restoreAssets,
		Windows: []*bootstrap.Window{{
			Homepage:       startPage,
			MessageHandler: gui.handleElectronCommands,
			Options:        &windowOptions,
		}},
		AstilectronOptions: astilectron.Options{
			AppName:            appName,
			AppIconDarwinPath:  "resources/icon.icns",
			AppIconDefaultPath: "resources/icon.png",
		},
		// TODO: Fix this tray to display nicely
		/*TrayOptions: &astilectron.TrayOptions{
			Image:   astilectron.PtrStr("/static/i/miner-logo.png"),
			Tooltip: astilectron.PtrStr(appName),
		},*/
		MenuOptions: menu,
		// OnWait is triggered as soon as the electron window is ready and running
		OnWait: func(
			_ *astilectron.Astilectron,
			windows []*astilectron.Window,
			_ *astilectron.Menu,
			_ *astilectron.Tray,
			_ *astilectron.Menu) error {
			gui.window = windows[0]
			gui.miningStatsTicker = time.NewTicker(time.Second * 5)
			gui.logger.Info("Start capturing mining stats")
			go gui.updateMiningStatsLoop()
			gui.networkStatsTicker = time.NewTicker(time.Minute * 2)
			go func() {
				for _ = range gui.networkStatsTicker.C {
					gui.updateNetworkStats()
				}
			}()
			gui.annTicker = time.NewTicker(time.Hour)
			go func() {
				for _ = range gui.annTicker.C {
					gui.checkAnnouncement()
				}
			}()
			// Trigger a network stats update as soon as we start
			gui.updateNetworkStats()
			// Check for any initial announcement
			gui.checkAnnouncement()
			return nil
		},
	}

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
	err = gui.stopMiner()
	if err != nil {
		return err
	}
	gui.miningStatsTicker.Stop()
	gui.networkStatsTicker.Stop()
	gui.annTicker.Stop()
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

	// get-miner-path is requested so the UI can show the path to exclude
	// in antivirus software
	case "get-miner-path":
		return filepath.Join(gui.workingDir, "miner"), nil

	// pool-list requests the recommended pool list from the miner API
	// and returns the rendered HTML
	case "pool-list":
		// Grab the pool list and send that to the GUI as well
		poolJSONs, err := gui.GetPoolList()
		if err != nil {
			_ = gui.sendElectronCommand("fatal_error", ElectronMessage{
				Data: fmt.Sprintf("Unable to fetch pool list from API."+
					"Please check that you are connected to the internet and try again."+
					"<br/>The error was '%s'", err),
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
		for i, poolData := range poolJSONs {
			var templateHTML bytes.Buffer
			err = poolTemplate.Execute(&templateHTML, poolData)
			if err != nil {
				log.Fatalf("Unable to load pool template: '%s'", err)
			}
			// TODO: This is a dirty way to only show the top 3 and reveal the rest
			// when needed. An API that implements paging is needed to fix this
			if i == 3 {
				poolsList += "<a href=\"#\" id=\"show_pool_list\">Show all</a>"
				poolsList += "<div id=\"pool_list_bottom\" class=\"dn\">"
			}
			poolsList += templateHTML.String()
		}
		// TODO: Part of the hack above
		poolsList += "</div>"
		return poolsList, nil

	// get-processing-config returns the current miner's processing config
	case "get-processing-config":
		if gui.miner == nil {
			_ = gui.sendElectronCommand("fatal_error", ElectronMessage{
				Data: fmt.Sprintf("Unable to fetch miner config." +
					"Please check that your miner is working and running."),
			})
			return "", nil
		}
		// Call the stats method to get processing info first, this causes the
		// stats to be cached by the miner
		_, _ = gui.miner.GetStats()
		processingConfig := gui.miner.GetProcessingConfig()
		configBytes, err := json.Marshal(processingConfig)
		if err != nil {
			_ = gui.sendElectronCommand("fatal_error", ElectronMessage{
				Data: fmt.Sprintf("Unable to fetch miner config."+
					"Please check that your miner is working and running."+
					"<br/>The error was '%s'", err),
			})
		}
		return string(configBytes), nil

	// configure is sent after the firstrun setup has been completed
	case "configure":
		// HACK: Adding a slight delay before switching to the mining dashboard
		// after initial setup to have the user at least see the 'configure' message
		time.Sleep(time.Second * 3)
		gui.configureMiner(command)
		return "Ok", nil

	// reconfigure is sent after settings are changes by the user
	case "reconfigure":
		gui.logger.Info("Reconfiguring miner")
		err := gui.stopMiner()
		if err != nil {
			_ = gui.sendElectronCommand("fatal_error", ElectronMessage{
				Data: fmt.Sprintf("Unable to stop miner for reconfigure."+
					"Please close the miner and open it again."+
					"<br/>The error was '%s'", err),
			})
			// Give the UI some time to display the message
			time.Sleep(time.Second * 15)
			gui.logger.Fatalf("Unable to reconfigure miner: '%s'", err)
		}
		gui.logger.WithField(
			"name", command.Name,
		).Debug("Received command from Electrom")
		gui.configureMiner(command)
		// Fake some time to have GUI at least display the message
		time.Sleep(time.Second * 3)
		gui.startMiner()
		gui.logger.Info("Miner reconfigured")

		gui.lastHashrate = 0.00
		// Trigger pool update
		go gui.updateNetworkStats()

		return "Ok", nil
	// miner_start is sent after configuration or when the user
	// clicks 'start mining'
	case "miner_start":
		gui.startMiner()

	// miner_stop is sent whenever the user clicks 'stop mining'
	case "miner_stop":
		err := gui.stopMiner()
		if err != nil {
			_ = gui.sendElectronCommand("fatal_error", ElectronMessage{
				Data: fmt.Sprintf("Unable to stop miner backend."+
					"Please close the miner and open it again."+
					"<br/>The error was '%s'", err),
			})
			// Give the UI some time to display the message
			time.Sleep(time.Second * 15)
			gui.logger.Fatalf("Unable to reconfigure miner: '%s'", err)
		}
	}
	return nil, fmt.Errorf("'%s' is an unknown command", command.Name)
}

// configureMiner creates the miner configuration to use
func (gui *GUI) configureMiner(command bootstrap.MessageIn) {
	gui.logger.Info("Configuring miner")

	var newConfig frontendConfig
	err := json.Unmarshal(command.Payload, &newConfig)
	if err != nil {
		_ = gui.sendElectronCommand("fatal_error", ElectronMessage{
			Data: fmt.Sprintf("Unable to configure miner."+
				"Please check your configuration is valid."+
				"<br/>The error was '%s'", err),
		})
		// Give the UI some time to display the message
		time.Sleep(time.Second * 15)
		gui.logger.Fatalf("Unable to configure miner: '%s'", err)
	}
	gui.config.Address = newConfig.Address
	gui.config.PoolID = newConfig.Pool

	scanPath := filepath.Join(gui.workingDir, "miner")
	// TODO: Fix own miner paths option
	/*if gui.config.Miner.Path != "" {
		//scanPath = path.Base(gui.config.Miner.Path)
	}*/
	gui.logger.WithField(
		"scan_path", scanPath,
	).Debug("Determining miner type")

	// Determine the type of miner bundled
	minerType, minerPath, err := miner.DetermineMinerType(scanPath)
	if err != nil {
		_ = gui.sendElectronCommand("fatal_error", ElectronMessage{
			Data: fmt.Sprintf("Unable to configure miner."+
				"Could not determine the miner type."+
				"<br/>The error was '%s'", err),
		})
		// Give the UI some time to display the message
		time.Sleep(time.Second * 15)
		gui.logger.Fatalf("Unable to configure miner: '%s'", err)
	}

	// Write config for this miner
	gui.config.Miner = miner.Config{
		Type: minerType,
		Path: minerPath,
	}
	gui.logger.WithFields(logrus.Fields{
		"path": minerPath,
		"type": minerType,
	}).Debug("Creating miner")
	gui.miner, err = miner.CreateMiner(gui.config.Miner)
	if err != nil {
		_ = gui.sendElectronCommand("fatal_error", ElectronMessage{
			Data: fmt.Sprintf("Unable to configure miner."+
				"<br/>The error was '%s'", err),
		})
		// Give the UI some time to display the message
		time.Sleep(time.Second * 15)
		gui.logger.Fatalf("Unable to configure miner: '%s'", err)
	}

	// The pool API returns the low-end hardware host:port config for pool
	gui.logger.Debug("Getting pool information")
	poolInfo, err := gui.GetPool(gui.config.PoolID)
	if err != nil {
		_ = gui.sendElectronCommand("fatal_error", ElectronMessage{
			Data: fmt.Sprintf("Unable to configure miner."+
				"Please check that you are connected to the internet."+
				"<br/>The error was '%s'", err),
		})
		// Give the UI some time to display the message
		time.Sleep(time.Second * 15)
		gui.logger.Fatalf("Unable to configure miner: '%s'", err)
	}

	// Write the config for the specified miner
	gui.logger.Debug("Writing miner config")

	err = gui.miner.WriteConfig(
		poolInfo.Config,
		gui.config.Address,
		miner.ProcessingConfig{
			Threads:  newConfig.Threads,
			MaxUsage: newConfig.MaxCPU,
		})
	if err != nil {
		_ = gui.sendElectronCommand("fatal_error", ElectronMessage{
			Data: fmt.Sprintf("Unable to configure miner."+
				"Please check that you are connected to the internet."+
				"<br/>The error was '%s'", err),
		})
		// Give the UI some time to display the message
		time.Sleep(time.Second * 15)
		gui.logger.Fatalf("Unable to configure miner: '%s'", err)
	}

	// Save the core miner config
	gui.logger.Debug("Writing GUI config")
	err = gui.SaveConfig(*gui.config)
	if err != nil {
		_ = gui.sendElectronCommand("fatal_error", ElectronMessage{
			Data: fmt.Sprintf("Unable to configure miner."+
				"Please check that you can write to the miner's installation path."+
				"<br/>The error was '%s'", err),
		})
		// Give the UI some time to display the message
		time.Sleep(time.Second * 15)
		gui.logger.Fatalf("Unable to configure miner: '%s'", err)
	}
	gui.logger.WithFields(logrus.Fields{
		"type": minerType,
	}).Info("Miner configured")
}

// startMiner starts the miner
func (gui *GUI) startMiner() {
	err := gui.miner.Start()
	if err != nil {
		_ = gui.sendElectronCommand("fatal_error", ElectronMessage{
			Data: fmt.Sprintf("Unable to start '%s' miner, please check that you "+
				"can run the miner from your installation directory."+
				"<br/>The error was '%s'", gui.miner.GetName(), err),
		})
		// Give the UI some time to display the message
		time.Sleep(time.Second * 15)
		gui.logger.Fatalf("Error starting '%s': %s", gui.miner.GetName(), err)
	}
	gui.logger.Infof("Started '%s' miner", gui.miner.GetName())
}

// stopMiner stops the miner
func (gui *GUI) stopMiner() error {
	if gui.miner == nil {
		return nil
	}
	err := gui.miner.Stop()
	if err != nil {
		_ = gui.sendElectronCommand("fatal_error", ElectronMessage{
			Data: fmt.Sprintf("Unable to stop miner.."+
				"Please close the GUI miner and open it again."+
				"<br/>The error was '%s'", err),
		})
		gui.logger.Errorf("Unable to stop miner '%s': %s", gui.miner.GetName(), err)
		return err
	}
	gui.logger.Infof("Stopped '%s' miner", gui.miner.GetName())
	return nil
}

// updateNetworkStats is a single stat update for network and payment info
func (gui *GUI) updateNetworkStats() {
	gui.logger.WithField(
		"hashrate", gui.lastHashrate,
	).Debug("Fetching network stats")
	// On firstrun we won't have a config yet
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

// updateMiningStats retrieves the miner's stats and updates
// the front-end
func (gui *GUI) updateMiningStatsLoop() {
	lastGraphUpdate := time.Now()
	for _ = range gui.miningStatsTicker.C {
		if gui.miner == nil {
			// No miner set up yet.. wait more
			gui.logger.Debug("Miner not set up yet, try again later")
			continue
		}
		gui.logger.Debug("Fetching mining stats")
		stats, err := gui.miner.GetStats()
		if err != nil {
			gui.logger.Debugf("Unable to get mining stats, miner not available yet?: %s", err)
		} else {
			if gui.lastHashrate == 0 && stats.Hashrate > 0 {
				gui.lastHashrate = stats.Hashrate
				// The first time we get a hashrate, update the BLOC amount so that the
				// user doesn't think it doesn't work
				gui.updateNetworkStats()
			}
			gui.lastHashrate = stats.Hashrate
			stats.Address = gui.config.Address

			if time.Since(lastGraphUpdate).Minutes() >= 1 {
				lastGraphUpdate = time.Now()
				stats.UpdateGraph = true
			}
			statBytes, _ := json.Marshal(&stats)
			err = bootstrap.SendMessage(gui.window, "miner_stats", string(statBytes))
			if err != nil {
				gui.logger.Errorf("Unable to send miner stats to front-end: %s", err)
			}
		}
	}
}

// checkAnnouncement checks for a new announcement and sends it to the
// front-end if one is available
func (gui *GUI) checkAnnouncement() {
	ann, err := gui.GetAnnouncement()
	if err != nil {
		gui.logger.Warningf("Unable to fetch announcements: %s", err)
		return
	}
	if ann.Ann == false {
		gui.logger.Debug("No new announcements are available")
		return
	}

	gui.logger.WithFields(logrus.Fields{
		"id":   ann.ID,
		"date": ann.Date,
		"link": ann.Link,
	}).Debug("Announcement fetched")

	// Only show announcements we haven't shown before in this session
	if _, ok := gui.annProcessed[ann.ID]; !ok {
		ann.DateString = fmt.Sprintf(
			"%s at %s UTC",
			ann.Date.Format("Monday January 2"),
			ann.Date.Format("15:04"))
		err = gui.sendElectronCommand("ann", ann)
		if err != nil {
			gui.logger.Warningf("Unable to send ANN to Electron: %s", err)
		}
		gui.annProcessed[ann.ID] = true

		gui.logger.WithFields(logrus.Fields{
			"id":   ann.ID,
			"date": ann.Date,
			"link": ann.Link,
		}).Info("New announcement available")
	} else {
		gui.logger.WithFields(logrus.Fields{
			"id":   ann.ID,
			"date": ann.Date,
			"link": ann.Link,
		}).Debug("Announcement already displayed in this session")
	}
}

// sendElectronCommand sends the given data to Electron under the command name
func (gui *GUI) sendElectronCommand(
	name string,
	data interface{}) error {
	dataBytes, err := json.Marshal(&data)
	if err != nil {
		return err
	}
	return bootstrap.SendMessage(gui.window, name, string(dataBytes))
}

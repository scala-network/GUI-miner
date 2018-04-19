package miner

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"os/user"
	"runtime"
	"sync/atomic"
	"syscall"
	"time"

	astilectron "github.com/asticode/go-astilectron"
	bootstrap "github.com/asticode/go-astilectron-bootstrap"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

// Miner implements the core control for the GUI miner
type Miner struct {
	// window is the main Astilectron window
	window *astilectron.Window
	// minerCmd is a reference to the xmr-stak miner process
	minerCmd *exec.Cmd
	// astilectronOptions holds the Astilectron options
	astilectronOptions bootstrap.Options
	// config for the miner
	config *GUIConfig
	// apiEndpoint is the web endpoint where stats and pools are retrieved from
	apiEndpoint string
	// logger is the logger for the application
	logger *logrus.Entry
	// currentHashrate of the user if mining
	lastHashrate float64
	// captureMiningStats determines if the mining stats loop should be running
	// or not
	captureMiningStats uint32
}

// New creates a new instance of the miner application
func New(
	appName string,
	config *GUIConfig,
	asset bootstrap.Asset,
	restoreAssets bootstrap.RestoreAssets,
	apiEndpoint string,
	isDebug bool) (*Miner, error) {

	if apiEndpoint == "" {
		return nil, errors.New("The API Endpoint must be specified")
	}
	// If no config is specified then this is the first run
	startPage := "index.html"
	if config == nil {
		startPage = "firstrun.html"
	}

	miner := Miner{
		apiEndpoint: apiEndpoint,
		config:      config,
	}

	miner.astilectronOptions = bootstrap.Options{
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
			miner.window = window
			go func() {
				miner.logger.Info("Start capturing stats")
				// Run forever!
				for {
					miner.updateNetworkStats()
					time.Sleep(time.Second * 30)
				}
			}()
			return nil
		},
		MessageHandler: miner.handleElectronCommands,
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
	miner.logger = logrus.WithFields(logrus.Fields{
		"service": "stellite-gui-miner",
	})

	miner.logger.Info("Setup complete")
	return &miner, nil
}

// Run the miner!
func (m *Miner) Run() error {
	m.logger.Info("Starting miner")
	err := bootstrap.Run(m.astilectronOptions)
	if err != nil {
		return err
	}
	// If xmr-stak is running, kill it
	if m.minerCmd != nil {
		err = m.minerCmd.Process.Kill()
		if err != nil {
			m.logger.Fatalf("Unable to stop xmr-stak: %s", err)
		}
	}
	return nil
}

// handleElectronCommands handles the messages sent by the Electron front-end
func (m *Miner) handleElectronCommands(
	_ *astilectron.Window,
	command bootstrap.MessageIn) (interface{}, error) {

	m.logger.WithField(
		"command", command.Name,
	).Debug("Received command from Electron")

	// Every Electron command has a name together with a payload containing the
	// actual message
	switch command.Name {

	// Firstrun is received on the first run of the miner. We return the current
	// logged in user's name
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
		poolJSONs, err := m.GetPoolList()
		if err != nil {
			_ = bootstrap.SendMessage(
				m.window,
				"fatal_error",
				fmt.Sprintf("{\"message\": \"Unable to fetch pool list from API."+
					"Please check that you are connected to the internet."+
					"<br/>The error was '%s'\"}", err))
			// Give the UI some time to display the message
			time.Sleep(time.Second * 15)
			m.logger.Fatalf("Unable to fetch pool list: '%s'", err)
		}
		poolTemplate, err := m.GetPoolTemplate(false)
		if err != nil {
			log.Fatalf("Unable to load pool template: '%s'", err)
		}
		var poolsList string
		for _, poolData := range poolJSONs {
			var templateHTML bytes.Buffer
			// Get the string time in the correct format
			t, _ := time.Parse("2006-01-02 15:04",
				poolData.LastBlock[:len(poolData.LastBlock)-3])
			since := time.Since(t)
			poolData.LastBlock = fmt.Sprintf("%d minutes ago", int(since.Minutes()))
			poolData.Hashrate = m.HumanizeHashrate(poolData.Hashrate)
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
		m.configureMiner(command, false)
		return "Ok", nil

	// reconfigure is sent after settings are changes by the user
	case "reconfigure":
		m.logger.Info("Reconfiguring miner")
		err := m.stopMiner()
		if err != nil {
			_ = bootstrap.SendMessage(
				m.window,
				"fatal_error",
				fmt.Sprintf("{\"message\": \"Unable to stop miner for reconfigure."+
					"Please close the miner and open it again"+
					"<br/>The error was '%s'\"}", err))
			// Give the UI some time to display the message
			time.Sleep(time.Second * 15)
			m.logger.Fatalf("Unable to reconfigure miner: '%s'", err)
		}
		// Wait for the mining stats loop to exit
		time.Sleep(time.Second * 10)
		m.configureMiner(command, true)
		m.startMiner()
		m.logger.Info("Miner reconfigured")

	// miner_start is sent after configuration or when the user
	// clicks 'start mining'
	case "miner_start":
		m.startMiner()

	// miner_stop is sent whenever the user clicks 'stop mining'
	case "miner_stop":
		return nil, m.stopMiner()
	}
	return nil, fmt.Errorf("'%s' is an unknown command", command.Name)
}

// configureMiner creates the xmr-stak configuration to use
func (m *Miner) configureMiner(command bootstrap.MessageIn, isReconfigure bool) {
	err := json.Unmarshal(command.Payload, &m.config)
	if err != nil {
		_ = bootstrap.SendMessage(
			m.window,
			"fatal_error",
			fmt.Sprintf("{\"message\": \"Unable to configure miner."+
				"Please check that you are connected to the internet."+
				"<br/>The error was '%s'\"}", err))
		// Give the UI some time to display the message
		time.Sleep(time.Second * 15)
		m.logger.Fatalf("Unable to configure miner: '%s'", err)
	}
	if isReconfigure == false {
		m.config.Mid = uuid.New().String()
	}
	err = m.SaveConfig(*m.config)
	if err != nil {
		_ = bootstrap.SendMessage(
			m.window,
			"fatal_error",
			fmt.Sprintf("{\"message\": \"Unable to configure miner."+
				"Please check that you are connected to the internet."+
				"<br/>The error was '%s'\"}", err))
		// Give the UI some time to display the message
		time.Sleep(time.Second * 15)
		m.logger.Fatalf("Unable to configure miner: '%s'", err)
	}

	err = ioutil.WriteFile("./xmr-stak/config.txt",
		[]byte(m.GetXmrStakConfig()),
		0644)
	if err != nil {
		_ = bootstrap.SendMessage(
			m.window,
			"fatal_error",
			fmt.Sprintf("{\"message\": \"Unable to configure miner."+
				"Please check that you are connected to the internet."+
				"<br/>The error was '%s'\"}", err))
		// Give the UI some time to display the message
		time.Sleep(time.Second * 15)
		m.logger.Fatalf("Unable to configure miner: '%s'", err)
	}

	var poolInfo PoolData
	poolInfo, err = m.GetPool(m.config.PoolID)
	if err != nil {
		_ = bootstrap.SendMessage(
			m.window,
			"fatal_error",
			fmt.Sprintf("{\"message\": \"Unable to configure miner."+
				"Please check that you are connected to the internet."+
				"<br/>The error was '%s'\"}", err))
		// Give the UI some time to display the message
		time.Sleep(time.Second * 15)
		m.logger.Fatalf("Unable to configure miner: '%s'", err)
	}

	err = ioutil.WriteFile("./xmr-stak/pools.txt",
		[]byte(m.GetXmrStakPoolConfig(poolInfo.Config, m.config.Address)),
		0644)
	if err != nil {
		_ = bootstrap.SendMessage(
			m.window,
			"fatal_error",
			fmt.Sprintf("{\"message\": \"Unable to configure miner."+
				"Please check that you are connected to the internet."+
				"<br/>The error was '%s'\"}", err))
		// Give the UI some time to display the message
		time.Sleep(time.Second * 15)
		m.logger.Fatalf("Unable to configure miner: '%s'", err)
	}
}

// startMiner starts the xmr-stak miner
func (m *Miner) startMiner() {
	params := []string{}
	commandName := "./xmr-stak"
	commandDir := "./xmr-stak"
	if runtime.GOOS == "windows" {
		commandName = ".\\xmr-stak.exe"
		commandDir = ".\\xmr-stak"
		m.minerCmd.SysProcAttr = &syscall.SysProcAttr{
			HideWindow: true,
		}
	}
	m.minerCmd = exec.Command(commandName, params...)
	m.minerCmd.Dir = commandDir
	err := m.minerCmd.Start()
	if err != nil {
		_ = bootstrap.SendMessage(
			m.window,
			"fatal_error",
			fmt.Sprintf("{\"message\": \"Unable to start xmr-stak miner, please check that you "+
				"can run 'xmr-stak' from the installation folder."+
				"<br/>The error was '%s'\"}", err))
		// Give the UI some time to display the message
		time.Sleep(time.Second * 15)
		m.logger.Fatalf("Error running xmr-stak: %s", err)
	}
	m.logger.Info("Started xmr-stak")
	go m.updateMiningStatsLoop()
}

// stopMiner stops the xmr-stak miner
func (m *Miner) stopMiner() error {
	if m.minerCmd != nil {
		err := m.minerCmd.Process.Kill()
		if err != nil {
			m.logger.Errorf("Unable to stop xmr-stak: %s", err)
			return err
		}
		m.logger.Info("xmr-stak stopped")
		m.lastHashrate = 0.00
		// Stop the stats loop as well
		atomic.StoreUint32(&m.captureMiningStats, 0)
		return nil
	}
	m.logger.Info("xmr-stak wasn't running")
	return nil
}

// updateNetworkStats is a single stat update for network and payment info
func (m *Miner) updateNetworkStats() {
	m.logger.WithField(
		"hashrate", m.lastHashrate,
	).Debug("Fetching network stats")
	// On firstrun we fon't have a config yet
	if m.config == nil {
		m.logger.Warning("No config set yet")
		return
	}
	stats, err := m.GetStats(m.config.PoolID, m.lastHashrate, m.config.Mid)
	if err != nil {
		m.logger.Warningf("Unable to get network stats: %s", err)
	} else {
		err := bootstrap.SendMessage(m.window, "network_stats", stats)
		if err != nil {
			m.logger.Errorf("Unable to send stats to front-end: %s", err)
		}
	}
}

// updateMiningStats retrieves the miner's stats from xmr-stak and updates
// the front-end
func (m *Miner) updateMiningStatsLoop() {
	atomic.StoreUint32(&m.captureMiningStats, 1)
	lastGraphUpdate := time.Now()
	for atomic.LoadUint32(&m.captureMiningStats) == 1 {
		m.logger.Debug("Fetching mining stats")
		xmrStats, err := m.GetXmrStats()
		if err != nil {
			m.logger.Warningf("Unable to get mining stats, xmr-stak not available yet?: %s", err)
		} else {
			if len(xmrStats.Hashrate.Total) > 0 {
				m.lastHashrate = xmrStats.Hashrate.Total[0]
				// The first time we get a hashrate, update the XTL amount so that the
				// user doesn't think it doesn't work
				m.updateNetworkStats()
			}
			xmrStats.Address = m.config.Address

			if time.Since(lastGraphUpdate).Minutes() >= 1 {
				lastGraphUpdate = time.Now()
				xmrStats.UpdateGraph = true
			}
			statBytes, _ := json.Marshal(&xmrStats)
			err = bootstrap.SendMessage(m.window, "miner_stats", string(statBytes))
			if err != nil {
				m.logger.Errorf("Unable to send miner stats to front-end: %s", err)
			}
		}
		time.Sleep(time.Second * 10)
	}
	m.logger.Debug("Stopped fetching mining stats")
}

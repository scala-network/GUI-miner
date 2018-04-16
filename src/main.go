// Package main implements the communication and control parts of the
// GUI miner. It constructs and launches the Electron front-end
package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"os/user"
	"time"

	"github.com/google/uuid"

	astilectron "github.com/asticode/go-astilectron"
	bootstrap "github.com/asticode/go-astilectron-bootstrap"
	astilog "github.com/asticode/go-astilog"
	"github.com/pkg/errors"
)

var AppName string
var BuiltAt string
var debug = flag.Bool("d", false, "enables the debug mode")
var w *astilectron.Window
var helper Helper
var xmrStakCmd *exec.Cmd
var config GUIInitConfig

// main implements the main runnable of the application
func main() {

	helper = Helper{
		MinerAPI: "http://138.197.183.78/miner",
	}

	flag.Parse()
	astilog.FlagInit()

	astilog.Debugf("Running the GUI miner, built at %s", BuiltAt)

	homePage := "first-run.html"
	_, err := os.Stat("config.json")
	if err == nil {
		homePage = "index.html"

		configBytes, err := ioutil.ReadFile("./config.json")
		if err != nil {
			log.Fatalf("Opening config.json failed: %s", err)
		}
		err = json.Unmarshal(configBytes, &config)
		if err != nil {
			log.Fatalf("Parsing config.json failed: %s", err)
		}
	}

	options := bootstrap.Options{
		Debug:         *debug,
		Asset:         Asset,
		RestoreAssets: RestoreAssets,
		Homepage:      homePage,
		AstilectronOptions: astilectron.Options{
			AppName:            AppName,
			AppIconDarwinPath:  "resources/icon.icns",
			AppIconDefaultPath: "resources/icon.png",
		},
		WindowOptions: &astilectron.WindowOptions{
			BackgroundColor: astilectron.PtrStr("#0B0C22"),
			Center:          astilectron.PtrBool(true),
			Height:          astilectron.PtrInt(680),
			Width:           astilectron.PtrInt(1175),
		},
		MenuOptions: []*astilectron.MenuItemOptions{{
			Label: astilectron.PtrStr("File"),
			SubMenu: []*astilectron.MenuItemOptions{
				{
					Label: astilectron.PtrStr("About"),
					OnClick: func(e astilectron.Event) (deleteListener bool) {
						fmt.Println("ABOUT CLICKED, SEND MESSAGE")
						err := bootstrap.SendMessage(w, "event.name", "hello",
							func(m *bootstrap.MessageIn) {
								// Unmarshal payload
								var s string
								json.Unmarshal(m.Payload, &s)

								// Process message
								fmt.Println("Received ", s)
								//log.Infof("received %s", s)
							})
						if err != nil {
							fmt.Println("FOKKEN ERR", err)
						}
						fmt.Println("Message sent")
						return
					},
				},
				{
					Role: astilectron.MenuItemRoleClose,
				},
			},
		}},
		OnWait: func(
			_ *astilectron.Astilectron,
			iw *astilectron.Window,
			_ *astilectron.Menu,
			_ *astilectron.Tray,
			_ *astilectron.Menu) error {
			w = iw
			go captureStats()
			return nil
		},
		MessageHandler: handleMessages,
	}

	err = bootstrap.Run(options)
	if err != nil {
		astilog.Fatal(errors.Wrap(err, "running bootstrap failed"))
	}
	if xmrStakCmd != nil {
		xmrStakCmd.Process.Kill()
	}
}

// handleMessages handles messages
func handleMessages(
	_ *astilectron.Window,
	m bootstrap.MessageIn) (payload interface{}, err error) {
	switch m.Name {
	case "startup":
		// Check if a config.json file exists, if so, this is not the
		// first startup
		var username string
		currentUser, err := user.Current()
		if err == nil {
			if currentUser.Name != "" {
				username = currentUser.Name
			} else if currentUser.Username != "" {
				username = currentUser.Username
			}
		}
		err = bootstrap.SendMessage(w, "firstrun", username)
		if err != nil {
			log.Fatalf("Unable to send message: '%s'", err)
		}

	case "pool-list":
		// Grab the pool list and send that to the GUI as well
		poolJSONs, err := helper.GetPoolList()
		if err != nil {
			log.Fatalf("Unable to fetch pool list: '%s'", err)
		}
		poolTemplate, err := helper.GetPoolTemplate()
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
			poolData.Hashrate = helper.HumanizeHashrate(poolData.Hashrate)
			err = poolTemplate.Execute(&templateHTML, poolData)
			if err != nil {
				log.Fatalf("Unable to load pool template: '%s'", err)
			}
			poolsList += templateHTML.String()
		}
		err = bootstrap.SendMessage(w, "pool-list", poolsList)
		if err != nil {
			log.Fatalf("Unable to send pool list: '%s'", err)
		}

	case "configure":
		// HACK TEMP: Checking if GUI works, remove later
		time.Sleep(time.Second * 5)
		err = json.Unmarshal(m.Payload, &config)
		if err != nil {
			log.Fatalf("Unable to configure miner: '%s'", err)
		}
		config.Mid = uuid.New().String()
		err = helper.SaveConfig(config)
		if err != nil {
			log.Fatalf("Unable to configure miner: '%s'", err)
		}

		err = ioutil.WriteFile("./xmr-stak/config.txt",
			[]byte(helper.GetXmrStakConfig()),
			0644)
		if err != nil {
			log.Fatalf("Unable to configure miner: '%s'", err)
		}

		var poolInfo PoolData
		poolInfo, err = helper.GetPool(config.PoolID)
		if err != nil {
			log.Fatalf("Unable to configure miner: '%s'", err)
		}

		err = ioutil.WriteFile("./xmr-stak/pools.txt",
			[]byte(helper.GetXmrStakPoolConfig(poolInfo.Config, config.Address)),
			0644)
		if err != nil {
			log.Fatalf("Unable to configure miner: '%s'", err)
		}
		payload = "OK"
		return

	case "miner_start":
		// TODO: Start xmrstak
		params := []string{}
		xmrStakCmd = exec.Command("./xmr-stak", params...)
		xmrStakCmd.Dir = "./xmr-stak"
		err = xmrStakCmd.Start()
		if err != nil {
			log.Fatalf("Error running xmr-stak: %s", err)
		}
		fmt.Println("Started xmr-stak")
	case "miner_stop":
		// Stop xmrstak
		if xmrStakCmd != nil {
			err := xmrStakCmd.Process.Kill()
			if err != nil {
				fmt.Println("Kill failed:", err)
			}
		}
		fmt.Println("Stopped xmr-stak")
	case "event.name":
		// Unmarshal payload
		var s string
		if err = json.Unmarshal(m.Payload, &s); err != nil {
			payload = err.Error()
			return
		}
		payload = s + " world"
	}
	return
}

func captureStats() {
	// Grab current hashrate from xmrstak
	var hashrate float64
	for {
		xmrStats, err := helper.GetXmrStats()
		if err != nil {
			log.Printf("Unable to get xmr-stak stats, not running? '%s'", err)
		} else {
			if len(xmrStats.Hashrate.Total) > 0 {
				hashrate = xmrStats.Hashrate.Total[0]
			}
		}
		stats, err := helper.GetStats(config.PoolID, hashrate, config.Mid)
		if err != nil {
			fmt.Println("Unable to get stats! ", err)
		} else {
			err = bootstrap.SendMessage(w, "stats", stats)
			if err != nil {
				fmt.Println("Unable to send stats: ", err)
			}
		}
		xmrStats.Address = config.Address
		statBytes, _ := json.Marshal(&xmrStats)
		err = bootstrap.SendMessage(w, "miner_stats", string(statBytes))
		if err != nil {
			fmt.Println("Unable to send miner stats: ", err)
		}

		time.Sleep(time.Second * 60)
	}
}

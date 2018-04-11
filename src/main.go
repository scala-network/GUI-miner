// Package main implements the communication and control parts of the
// GUI miner. It constructs and launches the Electron front-end
package main

import (
	"encoding/json"
	"flag"
	"fmt"

	astilectron "github.com/asticode/go-astilectron"
	bootstrap "github.com/asticode/go-astilectron-bootstrap"
	astilog "github.com/asticode/go-astilog"
	"github.com/pkg/errors"
)

var AppName string
var BuiltAt string
var debug = flag.Bool("d", false, "enables the debug mode")
var w *astilectron.Window

// main implements the main runnable of the application
func main() {
	flag.Parse()
	astilog.FlagInit()

	astilog.Debugf("Running the GUI miner, built at %s", BuiltAt)

	options := bootstrap.Options{
		Debug:         *debug,
		Asset:         Asset,
		RestoreAssets: RestoreAssets,
		Homepage:      "index.html",
		AstilectronOptions: astilectron.Options{
			AppName:            AppName,
			AppIconDarwinPath:  "resources/icon.icns",
			AppIconDefaultPath: "resources/icon.png",
		},
		WindowOptions: &astilectron.WindowOptions{
			BackgroundColor: astilectron.PtrStr("#333"),
			Center:          astilectron.PtrBool(true),
			Height:          astilectron.PtrInt(700),
			Width:           astilectron.PtrInt(700),
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
			return nil
		},
		MessageHandler: handleMessages,
	}

	err := bootstrap.Run(options)
	if err != nil {
		astilog.Fatal(errors.Wrap(err, "running bootstrap failed"))
	}
}

// handleMessages handles messages
func handleMessages(
	_ *astilectron.Window,
	m bootstrap.MessageIn) (payload interface{}, err error) {
	fmt.Println("GOT MESSAGE FROM JS")
	switch m.Name {
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

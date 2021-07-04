// +build windows

//go:generate goversioninfo -icon=static/32x32.ico
//go:generate file2byteslice -input vXboxInterface.dll -output vgamepad/vxboxdll.go -package vgamepad -var xboxDLL
//go:generate cmd /c "2goarray iconBytes main < static/32x32.ico > icon.go"
//go:generate gofmt -s -w .
package main

import (
	"io/fs"
	"net/http"
	"os"
	"webvxbox/player"

	"github.com/gen2brain/beeep"
	"github.com/getlantern/systray"
)

func onReady() {
	msg := "Server running on: " + GetOutboundIP().String() + ":8080"
	beeep.Notify("WebvXbox", msg, "")
	systray.SetTooltip(msg)
	systray.SetIcon(iconBytes)
	quit := systray.AddMenuItem("Quit", "")
	go func() {
		<-quit.ClickedCh
		systray.Quit()
		os.Exit(0)
	}()
}

func onExit() {

}

func EmitConnection(player *player.Player) {
	item := systray.AddMenuItem(player.Name, "Click to disconnect")
	go func() {
		select {
		case <-item.ClickedCh:
			item.Hide()
			player.Close()
		case <-player.CloseCh:
			item.Hide()
		}
	}()
}

func main() {
	go systray.Run(onReady, onExit)
	http.HandleFunc("/xbox", xboxHandler)
	dir, _ := fs.Sub(staticFS, "static")

	http.Handle("/", http.FileServer(http.FS(dir)))
	http.ListenAndServe(*addr, nil)
}

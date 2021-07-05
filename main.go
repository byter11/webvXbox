// +build linux

package main

import (
	"flag"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"webvxbox/player"
	"webvxbox/vgamepad"
)

func EmitConnection(player *player.Player) {
	fmt.Print(player.Name, " connected.")
}

func main() {
	c := make(chan (os.Signal))
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		vgamepad.Cleanup()
		fmt.Println("Exit.")
		os.Exit(1)
	}()

	flag.Parse()
	log.SetFlags(0)
	http.HandleFunc("/xbox", xboxHandler)
	dir, _ := fs.Sub(staticFS, "static")

	http.Handle("/", http.FileServer(http.FS(dir)))
	ip := GetOutboundIP()
	log.Printf("WebSocket: ws://%s:8080/xbox", ip)
	log.Printf("WebApp: http://%s:8080", ip)
	log.Fatal(http.ListenAndServe(*addr, nil))
}

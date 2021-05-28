//go:generate file2byteslice -input vXboxInterface.dll -output vgamepad/vxboxdll.go -package vgamepad -var xboxDLL
//go:generate gofmt -s -w .
package main

import (
	"embed"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"webvxbox/vgamepad"

	"github.com/gorilla/websocket"
)

var addr = flag.String("addr", "0.0.0.0:8080", "http service address")

var upgrader = websocket.Upgrader{} // use default options

var clients = make(map[string]*vgamepad.Vgamepad)

//go:embed static/*
var staticFS embed.FS

func xboxHandler(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}

	id := r.Header["Sec-Websocket-Key"][0]

	var e error
	clients[id], e = vgamepad.New()
	if e != nil {
		delete(clients, id)
		c.Close()
		log.Print("vgamepad:", e)
		return
	}
	fmt.Println(get_agent(r.Header["User-Agent"][0]), "Connected on Port", clients[id].Port)

	defer c.Close()
	defer disconnect_controller(id)

	for {
		_, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		control := strings.Split(string(message), "|")

		if strings.HasPrefix(control[0], "Axis") {
			x, _ := strconv.Atoi(control[1])
			x = x * 32767 / 100
			y, _ := strconv.Atoi(control[2])
			y = y * 32767 / 100
			clients[id].SetAxis(control[0], x, y)
		} else {
			value, _ := strconv.Atoi(control[1])
			if strings.HasPrefix(control[0], "Trigger") {
				value = value * 255
			}
			clients[id].SetBtn(control[0], value)
		}
	}
}

func disconnect_controller(id string) {
	clients[id].UnPlug()
	delete(clients, id)
}

func get_agent(s string) string {
	i := strings.Index(s, "(")
	if i >= 0 {
		j := strings.Index(s, ")")
		if j >= 0 {
			return s[i+1 : j]
		}
	}
	return ""
}

func GetOutboundIP() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP
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

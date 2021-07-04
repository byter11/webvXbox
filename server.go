package main

import (
	"embed"
	"flag"
	"log"
	"net"
	"net/http"
	"strconv"
	"strings"
	"webvxbox/player"
)

var addr = flag.String("addr", "0.0.0.0:8080", "http service address")

//go:embed static/*
var staticFS embed.FS

func GetOutboundIP() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP
}

func xboxHandler(w http.ResponseWriter, r *http.Request) {
	player, err := player.New(w, r)
	if err != nil {
		log.Print(err)
		return
	}
	EmitConnection(player)

	defer player.Close()

	for {
		_, message, err := player.Conn.ReadMessage()
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
			player.Vg.SetAxis(control[0], x, y)
		} else {
			value, _ := strconv.Atoi(control[1])
			if strings.HasPrefix(control[0], "Trigger") {
				value = value * 255
			}
			player.Vg.SetBtn(control[0], value)
		}
	}
}

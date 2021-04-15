package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"webvxbox/vxbox"

	"github.com/gorilla/websocket"
)

var addr = flag.String("addr", "0.0.0.0:8080", "http service address")

var upgrader = websocket.Upgrader{} // use default options

var clients = make(map[string]*vxbox.Vxbox)

func xboxHandler(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}

	id := r.Header["Sec-Websocket-Key"][0]

	var e error
	clients[id], e = vxbox.New()
	if e != nil {
		delete(clients, id)
		c.Close()
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
		// log.Printf("recv:  %s", control)
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
func main() {
	flag.Parse()
	log.SetFlags(0)
	http.HandleFunc("/ter", xboxHandler)
	http.Handle("/", http.FileServer(http.Dir("static")))

	log.Fatal(http.ListenAndServe(*addr, nil))
}

package player

import (
	"net/http"
	"strings"
	"webvxbox/vgamepad"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{} // use default options

type Player struct {
	Name    string
	Conn    *websocket.Conn
	Vg      *vgamepad.Vgamepad
	CloseCh chan int
}

func New(w http.ResponseWriter, r *http.Request) (*Player, error) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return nil, err
	}
	vg, err := vgamepad.New()
	if err != nil {
		c.Close()
		return nil, err
	}
	return &Player{
		Name:    get_agent(r.Header["User-Agent"][0]),
		Conn:    c,
		Vg:      vg,
		CloseCh: make(chan int),
	}, nil
}

func (p Player) Close() error {
	p.Vg.UnPlug()
	err := p.Conn.Close()
	if err != nil {
		return err
	}
	close(p.CloseCh)
	return nil
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

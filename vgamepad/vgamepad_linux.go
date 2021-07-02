// +build linux

package vgamepad

import (
	"fmt"
	"os"
	"strings"

	"github.com/bendahl/uinput"
)

var btnMap = map[string]int{
	"BtnA":      uinput.BtnA,
	"BtnB":      uinput.BtnB,
	"BtnX":      uinput.BtnX,
	"BtnY":      uinput.BtnY,
	"BtnStart":  uinput.BtnStart,
	"BtnLB":     uinput.BtnTL,
	"BtnRB":     uinput.BtnTR,
	"TriggerL":  uinput.BtnTL2,
	"TriggerR":  uinput.BtnTR2,
	"DpadDown":  uinput.BtnDpadDown,
	"DpadUp":    uinput.BtnDpadUp,
	"DpadLeft":  uinput.BtnDpadLeft,
	"DpadRight": uinput.BtnDpadRight,
}

func Cleanup() {

}

type Vgamepad struct {
	Port int
	vg   uinput.Gamepad
}

func New() (*Vgamepad, error) {
	port := 0
	for i := 1; i <= 4; i++ {
		name := fmt.Sprintf("Gamepad %d", i)
		if !gamepadExists(name) {
			port = i
			break
		}
	}
	if port == 0 {
		return nil, fmt.Errorf("port limit exceeded")
	}

	vg, err := uinput.CreateGamepad("/dev/uinput", []byte(fmt.Sprintf("Gamepad %d", port)))
	if err != nil {
		return nil, err
	}
	return &Vgamepad{Port: port, vg: vg}, nil
}

func (v Vgamepad) SetBtn(function string, arg int) {
	if arg != 0 {
		v.vg.BtnDown(btnMap[function])
	} else {
		v.vg.BtnUp(btnMap[function])
	}
}

func (v Vgamepad) SetAxis(function string, x int, y int) {
	v.vg.SetAxis(int32(x), int32(-y))
}

func (v Vgamepad) UnPlug() {
	v.vg.Close()
}

func gamepadExists(name string) bool {
	devices, err := os.ReadFile("/proc/bus/input/devices")
	if err != nil {
		fmt.Print("Error accessing input devices file")
	}
	return strings.Contains(string(devices), name)
}

package vxbox

import (
	"errors"
	"syscall"
)

var proc *syscall.LazyProc
var control_dict = make(map[string]*syscall.LazyProc)

func init() {
	lib := syscall.NewLazyDLL("VXboxInterface.dll")

	controls := []string{
		"isControllerExists", "PlugIn", "UnPlug", "SetBtnA", "SetBtnB",
		"SetBtnX", "SetBtnY", "SetBtnStart", "SetBtnLT", "SetBtnRT",
		"SetBtnLB", "SetBtnRB", "SetTriggerL", "SetTriggerR", "SetAxisX",
		"SetAxisY", "SetAxisRx", "SetAxisRy"}

	for _, s := range controls {
		proc = lib.NewProc(s)
		control_dict[s] = proc
	}
}

type Vxbox struct {
	Port int
}

func New() (*Vxbox, error) {
	port := 0
	for i := 1; i <= 4; i++ {
		ret, _, _ := control_dict["isControllerExists"].Call(uintptr(i))
		if ret == 0 {
			port = i
			break
		}
	}
	if port == 0 {
		return nil, errors.New("Port limit exceeded")
	}
	xbox := Vxbox{Port: port}
	control_dict["PlugIn"].Call(uintptr(xbox.Port))

	return &xbox, nil
}

func (v Vxbox) SetBtn(function string, arg int) {
	control_dict["Set"+function].Call(uintptr(v.Port), uintptr(arg))
}

func (v Vxbox) SetAxis(function string, x int, y int) {
	control_dict["Set"+function+"X"].Call(uintptr(v.Port), uintptr(x))

	control_dict["Set"+function+"Y"].Call(uintptr(v.Port), uintptr(y))
}

func (v Vxbox) UnPlug() {
	control_dict["UnPlug"].Call(uintptr(v.Port))
}

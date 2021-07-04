// +build windows

package vgamepad

import (
	"errors"
	"io/ioutil"
	"os"
	"syscall"
)

var proc *syscall.LazyProc
var control_dict = make(map[string]*syscall.LazyProc)
var dllpath string

func init() {
	f, _ := ioutil.TempFile("", "vXboxInterface.*.dll")
	dllpath = f.Name()
	f.Write(xboxDLL)
	f.Close()
	lib := syscall.NewLazyDLL(dllpath)

	controls := []string{
		"isControllerExists", "PlugIn", "UnPlug", "SetBtnA", "SetBtnB",
		"SetBtnX", "SetBtnY", "SetBtnStart", "SetBtnLT", "SetBtnRT",
		"SetBtnLB", "SetBtnRB", "SetTriggerL", "SetTriggerR", "SetAxisX",
		"SetAxisY", "SetAxisRx", "SetAxisRy", "SetDpadUp", "SetDpadRight",
		"SetDpadDown", "SetDpadLeft", "SetDpad"}

	for _, s := range controls {
		proc = lib.NewProc(s)
		control_dict[s] = proc
	}

	// Changed to avoid conditionals
	control_dict["SetAxisRX"] = control_dict["SetAxisRx"]
	control_dict["SetAxisRY"] = control_dict["SetAxisRy"]
	delete(control_dict, "SetAxisRx")
	delete(control_dict, "SetAxisRy")
}

func Cleanup() {
	os.Remove(dllpath)
}

type Vgamepad struct {
	Port int
}

func New() (*Vgamepad, error) {
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
	xbox := Vgamepad{Port: port}
	control_dict["PlugIn"].Call(uintptr(xbox.Port))

	return &xbox, nil
}

func (v Vgamepad) SetBtn(function string, arg int) {
	control_dict["Set"+function].Call(uintptr(v.Port), uintptr(arg))
}

func (v Vgamepad) SetAxis(function string, x int, y int) {
	control_dict["Set"+function+"X"].Call(uintptr(v.Port), uintptr(x))

	control_dict["Set"+function+"Y"].Call(uintptr(v.Port), uintptr(y))
}

func (v Vgamepad) UnPlug() {
	control_dict["UnPlug"].Call(uintptr(v.Port))
}

func (v Vgamepad) SetDpad(direction string) {
	control_dict["SetDpad"+direction].Call(uintptr(v.Port))
}

func (v Vgamepad) SetDpadOff() {
	control_dict["SetDpadOff"].Call(uintptr(v.Port))
}

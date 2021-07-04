package vgamepad

type vxbox interface {
	SetBtn(function string, arg int)

	SetAxis(function string, x int, y int)

	UnPlug()
}

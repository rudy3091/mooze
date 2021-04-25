package v2

import "os"

type Terminal interface {
	MakeRaw()
	RestoreRaw()
	GetWindowResizeChan() (chan os.Signal, chan bool)
	HandleResize()
}

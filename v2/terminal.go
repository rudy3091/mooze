package v2

type Terminal interface {
	MakeRaw()
	RestoreRaw()
}

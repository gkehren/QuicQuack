package server

type Server interface {
	Start(address string) error
	Stop() error
	SetThroughputMode(mode bool)
}

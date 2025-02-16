package server

import "net"

type Server interface {
	Start(address string) error
	Stop() error
	HandleConnection(conn net.Conn)
}

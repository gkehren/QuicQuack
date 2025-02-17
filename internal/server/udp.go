package server

import (
	"errors"
	"log"
	"net"
)

type UDPServer struct {
	conn     *net.UDPConn
	running  bool
	stopChan chan struct{}
}

func NewUDPServer() *UDPServer {
	return &UDPServer{
		stopChan: make(chan struct{}),
	}
}

func (s *UDPServer) Start(address string) error {
	addr, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		return err
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		return err
	}
	s.conn = conn
	s.running = true

	go func() {
		for s.running {
			buffer := make([]byte, 1024)
			n, remoteAddr, err := s.conn.ReadFromUDP(buffer)
			if err != nil {
				if !s.running {
					return
				}
				if !errors.Is(err, net.ErrClosed) {
					log.Println("Error reading:", err)
				}
				continue
			}

			_, err = s.conn.WriteToUDP(buffer[:n], remoteAddr)
			if err != nil {
				if !s.running {
					return
				}
				log.Println("Error writing:", err)
				continue
			}
		}
	}()

	return nil
}

func (s *UDPServer) Stop() error {
	s.running = false
	if s.conn != nil {
		return s.conn.Close()
	}
	return nil
}

package server

import (
	"errors"
	"log"
	"net"
)

type UDPServer struct {
	conn     *net.UDPConn
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

	go func() {
		for {
			select {
			case <-s.stopChan:
				return
			default:
				buffer := make([]byte, 1024)
				n, remoteAddr, err := s.conn.ReadFromUDP(buffer)
				if err != nil {
					if errors.Is(err, net.ErrClosed) {
						return
					}
					log.Println("Error reading:", err)
					continue
				}

				_, err = s.conn.WriteToUDP(buffer[:n], remoteAddr)
				if err != nil {
					log.Println("Error writing:", err)
					continue
				}
			}
		}
	}()

	return nil
}

func (s *UDPServer) Stop() error {
	close(s.stopChan)
	if s.conn != nil {
		return s.conn.Close()
	}
	return nil
}

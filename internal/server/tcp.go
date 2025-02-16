package server

import (
	"log"
	"net"
)

type TCPServer struct {
	listener net.Listener
}

func NewTCPServer() *TCPServer {
	return &TCPServer{}
}

func (s *TCPServer) Start(address string) error {
	l, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}
	s.listener = l
	go func() {
		for {
			conn, err := s.listener.Accept()
			if err != nil {
				log.Println("Error accepting connection:", err)
				continue
			}
			go s.HandleConnection(conn)
		}
	}()
	return nil
}

func (s *TCPServer) Stop() error {
	return s.listener.Close()
}

func (s *TCPServer) HandleConnection(conn net.Conn) {
	defer conn.Close()
	buffer := make([]byte, 1024)
	for {
		n, err := conn.Read(buffer)
		if err != nil {
			return
		}
		_, err = conn.Write(buffer[:n])
		if err != nil {
			return
		}
	}
}

package server

import (
	"log"
	"net"
	"sync"
)

type TCPServer struct {
	listener net.Listener
	wg       sync.WaitGroup
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
	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		for {
			conn, err := s.listener.Accept()
			if err != nil {
				if ne, ok := err.(net.Error); ok && ne.Timeout() {
					continue
				}
				log.Println("Error accepting connection:", err)
				break
			}
			s.wg.Add(1)
			go func(c net.Conn) {
				defer s.wg.Done()
				defer c.Close()
				s.handleConnection(conn)
			}(conn)
		}
	}()
	return nil
}

func (s *TCPServer) Stop() error {
	if err := s.listener.Close(); err != nil {
		return err
	}
	s.wg.Wait()
	return nil
}

func (s *TCPServer) handleConnection(conn net.Conn) {
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

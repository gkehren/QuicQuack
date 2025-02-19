package server

import (
	"errors"
	"log"
	"net"
	"sync"
)

type TCPServer struct {
	listener       net.Listener
	wg             sync.WaitGroup
	throughputMode bool
}

func NewTCPServer() *TCPServer {
	return &TCPServer{
		throughputMode: false,
	}
}

func (s *TCPServer) SetThroughputMode(enabled bool) {
	s.throughputMode = enabled
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
				if errors.Is(err, net.ErrClosed) {
					return
				}
				log.Println("Error accepting connection:", err)
				continue
			}
			s.wg.Add(1)
			go func(c net.Conn) {
				defer s.wg.Done()
				defer c.Close()
				s.handleConnection(c)
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
		if !s.throughputMode {
			_, err = conn.Write(buffer[:n])
			if err != nil {
				return
			}
		}
	}
}

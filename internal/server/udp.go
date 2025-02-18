package server

import (
	"context"
	"errors"
	"log"
	"net"
	"sync"
	"time"
)

type UDPServer struct {
	conn   *net.UDPConn
	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup
}

func NewUDPServer() *UDPServer {
	ctx, cancel := context.WithCancel(context.Background())
	return &UDPServer{
		ctx:    ctx,
		cancel: cancel,
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

	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		for {
			s.conn.SetReadDeadline(time.Now().Add(500 * time.Millisecond))
			select {
			case <-s.ctx.Done():
				return
			default:
				buffer := make([]byte, 1024)
				n, remoteAddr, err := s.conn.ReadFromUDP(buffer)
				if err != nil {
					var netErr net.Error
					if errors.As(err, &netErr) && netErr.Timeout() {
						continue
					}
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
	s.cancel()
	s.wg.Wait()
	if s.conn != nil {
		if err := s.conn.Close(); err != nil {
			return err
		}
	}
	return nil
}

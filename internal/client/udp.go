package client

import (
	"net"
	"time"
)

type UDPClient struct {
	conn *net.UDPConn
}

func NewUDPClient() *UDPClient {
	return &UDPClient{}
}

func (c *UDPClient) Connect(address string) error {
	addr, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		return err
	}

	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		return err
	}
	c.conn = conn
	return nil
}

func (c *UDPClient) Send(message []byte) error {
	_, err := c.conn.Write(message)
	return err
}

func (c *UDPClient) Receive() ([]byte, error) {
	buffer := make([]byte, 1024)
	c.conn.SetReadDeadline(time.Now().Add(5 * time.Second))
	n, err := c.conn.Read(buffer)
	if err != nil {
		return nil, err
	}
	return buffer[:n], nil
}

func (c *UDPClient) Close() error {
	return c.conn.Close()
}

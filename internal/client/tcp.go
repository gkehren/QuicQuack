package client

import (
	"net"
	"time"
)

type TCPClient struct {
	conn net.Conn
}

func NewTCPClient() *TCPClient {
	return &TCPClient{}
}

func (c *TCPClient) Connect(address string) error {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return err
	}
	c.conn = conn
	return nil
}

func (c *TCPClient) Send(message []byte) error {
	_, err := c.conn.Write(message)
	return err
}

func (c *TCPClient) Receive() ([]byte, error) {
	buffer := make([]byte, 1024)
	c.conn.SetReadDeadline(time.Now().Add(5 * time.Second))
	n, err := c.conn.Read(buffer)
	if err != nil {
		return nil, err
	}
	return buffer[:n], nil
}

func (c *TCPClient) Close() error {
	return c.conn.Close()
}

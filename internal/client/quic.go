package client

import (
	"context"
	"crypto/tls"
	"time"

	"github.com/quic-go/quic-go"
)

type QuicClient struct {
	conn   quic.Connection
	stream quic.Stream
}

func NewQuicClient() *QuicClient {
	return &QuicClient{}
}

func (c *QuicClient) Connect(address string) error {
	tlsConf := &tls.Config{
		InsecureSkipVerify: true,
		NextProtos:         []string{"quicquack-benchmark"},
	}
	conn, err := quic.DialAddr(context.Background(), address, tlsConf, nil)
	if err != nil {
		return err
	}
	c.conn = conn

	stream, err := conn.OpenStreamSync(context.Background())
	if err != nil {
		return err
	}
	c.stream = stream
	return nil
}

func (c *QuicClient) Send(message []byte) error {
	_, err := c.stream.Write([]byte(message))
	if err != nil {
		return err
	}
	return nil
}

func (c *QuicClient) Receive() ([]byte, error) {
	buffer := make([]byte, 1024)
	c.stream.SetReadDeadline(time.Now().Add(5 * time.Second))
	n, err := c.stream.Read(buffer)
	if err != nil {
		return nil, err
	}
	return buffer[:n], nil
}

func (c *QuicClient) Close() error {
	if err := c.conn.CloseWithError(0, ""); err != nil {
		return err
	}
	if err := c.stream.Close(); err != nil {
		return err
	}
	return nil
}

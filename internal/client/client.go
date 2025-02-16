package client

type Client interface {
	Connect(address string) error
	Send(message []byte) error
	Receive() ([]byte, error)
	Close() error
}

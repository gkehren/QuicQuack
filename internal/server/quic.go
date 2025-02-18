package server

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"log"
	"math/big"
	"sync"

	"github.com/quic-go/quic-go"
)

type QuicServer struct {
	listener *quic.Listener
	wg       sync.WaitGroup
}

func NewQuicServer() *QuicServer {
	return &QuicServer{}
}

func (s *QuicServer) Start(address string) error {
	l, err := quic.ListenAddr(address, generateTLSConfig(), nil)
	if err != nil {
		return err
	}
	s.listener = l
	s.wg.Add(1)

	go func() {
		defer s.wg.Done()
		for {
			conn, err := s.listener.Accept(context.Background())
			if err != nil {
				if errors.Is(err, quic.ErrServerClosed) {
					break
				}
				log.Println("Error accepting connection:", err)
				continue
			}

			s.wg.Add(1)
			go func(conn quic.Connection) {
				defer s.wg.Done()
				stream, err := conn.AcceptStream(context.Background())
				if err != nil {
					log.Println("Error accepting QUIC strean:", err)
					return
				}
				defer stream.Close()
				defer conn.CloseWithError(0, "")
				s.handleConnection(stream)
			}(conn)
		}
	}()
	return nil
}

func (s *QuicServer) Stop() error {
	if err := s.listener.Close(); err != nil {
		return err
	}
	s.wg.Wait()
	return nil
}

func (s *QuicServer) handleConnection(stream quic.Stream) {
	buffer := make([]byte, 1024)
	for {
		n, err := stream.Read(buffer)
		if err != nil {
			return
		}
		if _, err := stream.Write(buffer[:n]); err != nil {
			return
		}
	}
}

func generateTLSConfig() *tls.Config {
	key, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		panic(err)
	}
	template := x509.Certificate{SerialNumber: big.NewInt(1)}
	certDER, err := x509.CreateCertificate(rand.Reader, &template, &template, &key.PublicKey, key)
	if err != nil {
		panic(err)
	}
	keyPEM := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(key)})
	certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certDER})

	tlsCert, err := tls.X509KeyPair(certPEM, keyPEM)
	if err != nil {
		panic(err)
	}
	return &tls.Config{
		Certificates: []tls.Certificate{tlsCert},
		NextProtos:   []string{"quicquack-benchmark"},
	}
}

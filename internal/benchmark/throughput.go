package benchmark

import (
	"QuicQuack/internal/client"
	"QuicQuack/internal/server"
	"fmt"
	"time"
)

type ThroughputBenchmark struct {
	Protocol    string
	Server      string
	Duration    time.Duration
	PayloadSize int
	BytesSent   int64
}

func NewThroughputBenchmark(protocol, server string, duration time.Duration, payloadSize int) *ThroughputBenchmark {
	return &ThroughputBenchmark{
		Protocol:    protocol,
		Server:      server,
		Duration:    duration,
		PayloadSize: payloadSize,
	}
}

func (tb *ThroughputBenchmark) Run() error {
	var s server.Server
	var c client.Client

	switch tb.Protocol {
	case "tcp":
		s = server.NewTCPServer()
		c = client.NewTCPClient()
	case "udp":
		s = server.NewUDPServer()
		c = client.NewUDPClient()
	case "quic":
		s = server.NewQuicServer()
		c = client.NewQuicClient()
	default:
		return fmt.Errorf("unsupported protocol: %s", tb.Protocol)
	}

	go func() {
		s.SetThroughputMode(true)
		if err := s.Start(tb.Server); err != nil {
			fmt.Printf("Error starting server: %v\n", err)
		}
	}()
	defer s.Stop()

	time.Sleep(1 * time.Second)

	if err := c.Connect(tb.Server); err != nil {
		return fmt.Errorf("error connecting to server: %v", err)
	}
	defer c.Close()

	payload := make([]byte, tb.PayloadSize)
	for i := 0; i < tb.PayloadSize; i++ {
		payload[i] = byte(i % 256)
	}
	start := time.Now()
	endTime := start.Add(tb.Duration)

	for time.Now().Before(endTime) {
		if err := c.Send(payload); err != nil {
			return fmt.Errorf("error sending data: %v", err)
		}
		tb.BytesSent += int64(tb.PayloadSize)
	}
	return nil
}

func (tb *ThroughputBenchmark) Results() string {
	duration := tb.Duration.Seconds()
	if duration == 0 {
		return "Benchmark duration is zero"
	}

	throughput := float64(tb.BytesSent) / duration
	throughputKB := throughput / 1024
	throughputMB := throughputKB / 1024

	return fmt.Sprintf("Throughput Results:\n"+
		"Duration: %v\n"+
		"Bytes Sent: %d\n"+
		"Throughput: %.2f bytes/second\n"+
		"Throughput: %.2f KB/second\n"+
		"Throughput: %.2f MB/second\n",
		tb.Duration, tb.BytesSent, throughput, throughputKB, throughputMB)
}

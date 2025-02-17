package benchmark

import (
	"QuicQuack/internal/client"
	"QuicQuack/internal/server"
	"fmt"
	"time"
)

type LatencyBenchmark struct {
	Protocol   string
	Server     string
	SampleSize int
	Measures   []time.Duration
}

func NewLatencyBenchmark(protocol string, server string, sampleSize int) *LatencyBenchmark {
	return &LatencyBenchmark{
		Protocol:   protocol,
		Server:     server,
		SampleSize: sampleSize,
	}
}

func (lb *LatencyBenchmark) Run() error {
	var tcpServer server.Server
	var tcpClient client.Client

	switch lb.Protocol {
	case "tcp":
		tcpServer = server.NewTCPServer()
		tcpClient = client.NewTCPClient()
	default:
		return fmt.Errorf("unsupported protocol: %s", lb.Protocol)
	}

	go func() {
		if err := tcpServer.Start(lb.Server); err != nil {
			fmt.Printf("Error starting server: %v\n", err)
		}
	}()
	defer tcpServer.Stop()

	time.Sleep(1 * time.Second) // wait for the server ready

	if err := tcpClient.Connect(lb.Server); err != nil {
		return fmt.Errorf("error connecting to server: %v", err)
	}
	defer tcpClient.Close()

	for i := 0; i < lb.SampleSize; i++ {
		start := time.Now()
		if err := tcpClient.Send([]byte("ping")); err != nil {
			return fmt.Errorf("error sending ping: %v", err)
		}
		if _, err := tcpClient.Receive(); err != nil {
			return fmt.Errorf("error receiving pong: %v", err)
		}
		lb.Measures = append(lb.Measures, time.Since(start))
	}

	return nil
}

func (lb *LatencyBenchmark) Results() string {
	if len(lb.Measures) == 0 {
		return "No results available."
	}

	total := time.Duration(0)
	for _, duration := range lb.Measures {
		total += duration
	}
	avg := total / time.Duration(len(lb.Measures))
	min := lb.Measures[0]
	max := lb.Measures[0]
	for _, duration := range lb.Measures {
		if duration < min {
			min = duration
		}
		if duration > max {
			max = duration
		}
	}

	return fmt.Sprintf("Latency Results (Sample Size: %d)\nAverage: %v\nMin: %v\nMax: %v",
		lb.SampleSize, avg, min, max)
}

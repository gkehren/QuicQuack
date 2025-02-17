package main

import (
	"QuicQuack/internal/benchmark"
	"fmt"
	"log"
)

func main() {
	lbTcp := benchmark.NewLatencyBenchmark("tcp", "localhost:8080", 10)

	if err := lbTcp.Run(); err != nil {
		log.Fatalf("Error running latency benchmark: %v", err)
	}

	fmt.Println("TCP Results:")
	fmt.Println(lbTcp.Results())

	lbUdp := benchmark.NewLatencyBenchmark("udp", "localhost:8080", 10)

	if err := lbUdp.Run(); err != nil {
		log.Fatalf("Error running latency benchmark: %v", err)
	}

	fmt.Println("UDP Results:")
	fmt.Println(lbUdp.Results())

	lbQuic := benchmark.NewLatencyBenchmark("quic", "localhost:8080", 10)

	if err := lbQuic.Run(); err != nil {
		log.Fatalf("Error running latency benchmark: %v", err)
	}

	fmt.Println("Quic Results:")
	fmt.Println(lbQuic.Results())
}

package main

import (
	"QuicQuack/internal/benchmark"
	"fmt"
	"log"
)

func main() {
	lb := benchmark.NewLatencyBenchmark("tcp", "localhost:8080", 10)

	if err := lb.Run(); err != nil {
		log.Fatalf("Error running latency benchmark: %v", err)
	}

	fmt.Println(lb.Results())
}

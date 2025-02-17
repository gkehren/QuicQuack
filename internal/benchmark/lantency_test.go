package benchmark

import (
	"testing"
	"time"
)

func TestTCPConnection(t *testing.T) {
	lb := NewLatencyBenchmark("tcp", "localhost:8080", 1)

	err := lb.Run()
	if err != nil {
		t.Fatalf("Expected no error, but got: %v", err)
	}

	if len(lb.Measures) == 0 {
		t.Fatal("Expected measures, but got none")
	}
}

func TestTCPAverageLatency(t *testing.T) {
	lb := NewLatencyBenchmark("tcp", "localhost:8080", 10)

	err := lb.Run()
	if err != nil {
		t.Fatalf("Expected no error, but got: %v", err)
	}

	total := time.Duration(0)
	for _, duration := range lb.Measures {
		total += duration
	}
	avg := total / time.Duration(len(lb.Measures))

	if avg > 3*time.Second {
		t.Fatalf("Expected avarage latency to be less than 3 second, but got: %v", avg)
	}
}

func TestTCPResultsCollection(t *testing.T) {
	lb := NewLatencyBenchmark("tcp", "localhost:8080", 5)

	err := lb.Run()
	if err != nil {
		t.Fatalf("Expected no error, but got: %v", err)
	}

	if len(lb.Measures) != lb.SampleSize {
		t.Fatalf("Expected %d results, but got %d", lb.SampleSize, len(lb.Measures))
	}
}

func TestUDPConnection(t *testing.T) {
	lb := NewLatencyBenchmark("udp", "localhost:8080", 1)

	err := lb.Run()
	if err != nil {
		t.Fatalf("Expected no error, but got: %v", err)
	}

	if len(lb.Measures) == 0 {
		t.Fatal("Expected measures, but got none")
	}
}

func TestUDPAverageLatency(t *testing.T) {
	lb := NewLatencyBenchmark("udp", "localhost:8080", 10)

	err := lb.Run()
	if err != nil {
		t.Fatalf("Expected no error, but got: %v", err)
	}

	total := time.Duration(0)
	for _, duration := range lb.Measures {
		total += duration
	}
	avg := total / time.Duration(len(lb.Measures))

	if avg > 3*time.Second {
		t.Fatalf("Expected avarage latency to be less than 3 second, but got: %v", avg)
	}
}

func TestUDPResultsCollection(t *testing.T) {
	lb := NewLatencyBenchmark("udp", "localhost:8080", 5)

	err := lb.Run()
	if err != nil {
		t.Fatalf("Expected no error, but got: %v", err)
	}

	if len(lb.Measures) != lb.SampleSize {
		t.Fatalf("Expected %d results, but got %d", lb.SampleSize, len(lb.Measures))
	}
}

func TestQuicConnection(t *testing.T) {
	lb := NewLatencyBenchmark("quic", "localhost:8080", 1)

	err := lb.Run()
	if err != nil {
		t.Fatalf("Expected no error, but got: %v", err)
	}

	if len(lb.Measures) == 0 {
		t.Fatal("Expected measures, but got none")
	}
}

func TestQuicAverageLatency(t *testing.T) {
	lb := NewLatencyBenchmark("quic", "localhost:8080", 10)

	err := lb.Run()
	if err != nil {
		t.Fatalf("Expected no error, but got: %v", err)
	}

	total := time.Duration(0)
	for _, duration := range lb.Measures {
		total += duration
	}
	avg := total / time.Duration(len(lb.Measures))

	if avg > 3*time.Second {
		t.Fatalf("Expected avarage latency to be less than 3 second, but got: %v", avg)
	}
}

func TestQuicResultsCollection(t *testing.T) {
	lb := NewLatencyBenchmark("quic", "localhost:8080", 5)

	err := lb.Run()
	if err != nil {
		t.Fatalf("Expected no error, but got: %v", err)
	}

	if len(lb.Measures) != lb.SampleSize {
		t.Fatalf("Expected %d results, but got %d", lb.SampleSize, len(lb.Measures))
	}
}

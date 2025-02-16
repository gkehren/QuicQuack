package benchmark

type Benchmarker interface {
	Run() error
	Results() string
}

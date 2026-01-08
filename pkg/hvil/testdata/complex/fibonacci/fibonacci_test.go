package fibonacci

import (
	"testing"
)

func BenchmarkFibonacci(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Fib(93)
	}
}

func Fib(n uint64) uint64 {
	counter := uint64(1)
	result := uint64(1)
	nextFib := uint64(1)

	for counter < n {
		nextFib = result + nextFib
		result = nextFib - result
		counter++
	}

	return result
}

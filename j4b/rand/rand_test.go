package rand

import (
	"math/rand"
	"testing"
	_ "unsafe"
)

func TestMathRand(t *testing.T) {
	for i := 0; i < 10000; i++ {
		print(rand.Intn(10000), " ")
	}
}

func TestFastRand(t *testing.T) {
	for i := 0; i < 10000; i++ {
		print(fastrandn(10000), " ")
	}
}
func BenchmarkMathRand(b *testing.B) {
	for i := 0; i < b.N; i++ {
		rand.Intn(b.N)
	}
}

func BenchmarkFastRand(b *testing.B) {
	for i := 0; i < b.N; i++ {
		fastrandn(uint32(b.N))
	}
}

//go:linkname fastrandn runtime.fastrandn
func fastrandn(n uint32) uint32

package append

import "testing"

func BenchmarkAppend(b *testing.B) {
	a := make([]int, 0, b.N)
	for i := 0; i < b.N; i++ {
		a = append(a, i)
	}
}
func BenchmarkCopy(b *testing.B) {
	a := make([]int, 0, b.N)
	for i := 0; i < b.N; i++ {
		a = a[:i]
	}
}

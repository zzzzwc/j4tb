package countdigits

import (
	"fmt"
	"math"
	"testing"
)

func n1(a int) (ret int) {
	for ; a >= 1e6; a /= 1e6 {
		ret += 18
	}
	for ; a >= 1e3; a /= 1e3 {
		ret += 6
	}
	for ; a >= 1; a /= 10 {
		ret++
	}
	return
}
func n2(a int) (ret int) {
	for ; a >= 1; a /= 10 {
		ret++
	}
	return
}
func BenchmarkCountDigitsUseFmt(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = len(fmt.Sprint(b.N))
	}
}

func BenchmarkCountDigitsFor10(b *testing.B) {
	for i := 0; i < b.N; i++ {
		n1(b.N)
	}
}
func BenchmarkCountDigits2For10(b *testing.B) {
	for i := 0; i < b.N; i++ {
		n2(b.N)
	}
}
func BenchmarkCountDigitsLog(b *testing.B) {
	for i := 0; i < b.N; i++ {
		math.Log10(float64(b.N))
	}
}

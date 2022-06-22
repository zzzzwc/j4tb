package slice_test

import (
	"fmt"
	"testing"
	"unsafe"
)

func TestU8Slice(t *testing.T) {
	var a uint8
	var slice []byte
	var p *uint8
	var l int
	a = 67
	p = &a
	l = 1
	slice = *(*[]byte)(unsafe.Pointer(&struct {
		p   unsafe.Pointer
		len int
		cap int
	}{
		unsafe.Pointer(p), l, l,
	}))
	println(slice[0])
}

func TestSlice(t *testing.T) {
	a := []int{1, 2, 3}
	b := [8]int{1, 2, 3}
	func(a []int) {
		a[0] = 100
	}(a)
	func(b [8]int) {
		a[0] = 100
	}(b)

	fmt.Println(a)
	fmt.Println(b)
}

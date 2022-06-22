package itoa_test

import (
	"strconv"
	"testing"
	"unsafe"
)

var (
	TowDigitsReverseMapping = make([]int, 64*1024)
)

func init() {
	for i := range TowDigitsReverseMapping {
		TowDigitsReverseMapping[i] = -1
	}
	for i := 0; i < 100; i++ {
		val := strconv.Itoa(i)
		if len(val) == 1 {
			TowDigitsReverseMapping[int(*(*int16)(unsafe.Pointer(&[2]byte{'0', val[0]})))] = i
		} else {
			TowDigitsReverseMapping[int(*(*int16)(unsafe.Pointer(&[2]byte{val[0], val[1]})))] = i
		}
	}
}

func BenchmarkReverseMapping(b *testing.B) {
	intv := 0
	for i := 0; i < b.N; i++ {
		str := strconv.Itoa(i)
		for i := 0; i < len(str); i += 2 {
			if i+1 == len(str) {
				idx := int(*(*int16)(unsafe.Pointer(&[2]byte{'0', str[i]})))
				intv *= 10
				intv += TowDigitsReverseMapping[idx]
			} else {
				idx := int(*(*int16)(unsafe.Pointer(&[2]byte{str[i], str[i+1]})))
				intv *= 100
				v := TowDigitsReverseMapping[idx]
				intv += v
			}
		}
		if intv != i {
			panic(intv)
		}
		intv = 0
	}
}

func BenchmarkStdConv(b *testing.B) {
	for i := 0; i < b.N; i++ {
		str := strconv.Itoa(i)
		v, err := strconv.Atoi(str)
		if err != nil {
			panic(err)
		}
		if v != i {
			panic(v)
		}
	}
}

func TestIface(t *testing.T) {

}

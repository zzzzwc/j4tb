package main

import (
	"bytes"
	"fmt"
	"strconv"
	"unsafe"
)

var (
	HeaderEnding            = []byte("\r\n\r\n")
	ContentLength           = []byte("Content-Length: ")
	TowDigitsReverseMapping = make([]int, 64*1024)
)

func init() {
	for i := range TowDigitsReverseMapping {
		TowDigitsReverseMapping[i] = -1
	}
	for i := 0; i < 99; i++ {
		val := strconv.Itoa(i)
		if len(val) == 1 {
			TowDigitsReverseMapping[int(*(*int16)(unsafe.Pointer(&[2]byte{val[0], 0})))] = i
		} else {
			TowDigitsReverseMapping[int(*(*int16)(unsafe.Pointer(&[2]byte{val[0], val[1]})))] = i
		}
	}
}

type (
	// Not thread safe
	HTTPDelimiter struct {
		buff  bytes.Buffer
		state HTTPDelimiterState
	}
	HTTPDelimiterState int
)

const (
	WaitHeaderEnding HTTPDelimiterState = iota
	WaitBodyEnding
)

func (h HTTPDelimiter) Write(in []byte) (outs [][]byte, err error) {
	// 1. buff is empty, try to separate a request from in directly
	var cl int
	if h.buff.Len() == 0 {
		for len(in) > 0 {
			if ending := bytes.Index(in, HeaderEnding); ending >= 0 {
				if clIndex := bytes.Index(in[:ending], ContentLength); clIndex >= 0 {
					clBuff := in[clIndex+len(ContentLength)+2:]
					cl, err = strconv.Atoi(*(*string)(unsafe.Pointer(&clBuff)))
					if err != nil {
						return nil, fmt.Errorf("parse contentlength %w", err)
					}
				}
			}
		}
	}
	_ = cl
	// TODO
	return
}

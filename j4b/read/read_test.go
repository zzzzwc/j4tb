package main_test

import (
	stdio "io"
	"runtime"
	"sync"
	"sync/atomic"
	"testing"

	"j4t/j4t/mock/io"
)

var (
	buffSize = 128 * 1024
	gos      = runtime.GOMAXPROCS(0)
)

func BenchmarkSingleRead(b *testing.B) {
	r, w, err := io.NewTcpConn()
	if err != nil {
		b.Fatal(err)
	}
	go func() {
		wb := make([]byte, buffSize)
		for i := 0; i < b.N; i++ {
			_, err := w.Write(wb)
			if err != nil {
				b.Error(err)
				return
			}
		}
	}()
	rb := make([]byte, buffSize)
	b.SetBytes(int64(buffSize))
	for rl := 0; rl < buffSize*b.N; {
		l, err := r.Read(rb)
		if err != nil {
			b.Fatal(err)
		}
		rl += l
	}
}

func BenchmarkMultiRead(b *testing.B) {
	r, w, err := io.NewTcpConn()
	if err != nil {
		b.Fatal(err)
	}
	b.N = gos * (b.N/gos + 1)
	go func() {
		wb := make([]byte, buffSize)
		for i := 0; i < b.N; i++ {
			_, err := w.Write(wb)
			if err != nil {
				b.Error(err)
				return
			}
		}
		err = w.Close()
		if err != nil {
			b.Error(err)
			return
		}
	}()
	var wg sync.WaitGroup
	var rl int64
	b.SetBytes(int64(buffSize))
	for i := 0; i < gos; i++ {
		wg.Add(1)
		go func(id int) {
			rb := make([]byte, buffSize)
			for atomic.LoadInt64(&rl) < int64(buffSize*b.N) {
				l, err := r.Read(rb)
				if err != nil {
					if err == stdio.EOF {
						break
					}
					b.Error(err)
					return
				}
				atomic.AddInt64(&rl, int64(l))
			}
			wg.Done()
		}(i)
	}
	wg.Wait()
}

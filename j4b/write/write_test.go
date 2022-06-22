package write_test

import (
	"runtime"
	"sync/atomic"
	"testing"

	"j4t/j4t/mock/io"
)

var (
	buffSize = 128 * 1024
	gos      = runtime.GOMAXPROCS(0)
)

func BenchmarkSyncWrite(b *testing.B) {
	var done int32
	s, c, err := io.NewTcpConn()
	if err != nil {
		b.Fatal(err)
	}

	b.SetBytes(int64(buffSize * 2))
	go func() {
		// TODO 13
		buf := make([]byte, buffSize)

		var err error
		for i := 0; i < b.N; i++ {
			_, err = c.Write(buf)
			if err != nil {
				break
			}
		}
		if err != nil && atomic.LoadInt32(&done) != 1 {
			b.Error(err)
			return
		}
	}()
	go func() {
		var rl int
		var err error
		var l int
		buf := make([]byte, buffSize)
		for rl < b.N*buffSize {
			l, err = s.Read(buf)
			if err != nil {
				break
			}
			rl += l
			_, err = s.Write(buf[:l])
			if err != nil {
				break
			}
		}
		if err != nil && atomic.LoadInt32(&done) != 1 {
			b.Error(err)
			return
		}
	}()

	rl := 0
	buf := make([]byte, buffSize)
	for rl < b.N*buffSize {
		l, err := c.Read(buf)
		if err != nil {
			b.Error(err)
			return
		}
		rl += l
	}
	atomic.StoreInt32(&done, 1)
	err = s.Close()
	if err != nil {
		b.Fatal(err)
	}
	err = c.Close()
	if err != nil {
		b.Fatal(err)
	}

}

func BenchmarkAsyncWrite(b *testing.B) {
	var done int32
	s, c, err := io.NewTcpConn()
	if err != nil {
		b.Fatal(err)
	}
	b.SetBytes(int64(buffSize * 2))
	go func() {
		buf := make([]byte, buffSize)

		var err error
		for i := 0; i < b.N; i++ {
			_, err = c.Write(buf)
			if err != nil {
				break
			}
		}
		if err != nil && atomic.LoadInt32(&done) != 1 {
			b.Error(err)
			return
		}
	}()
	go func() {
		var rl int
		var err error
		var l int
		buf := make([]byte, buffSize)
		for rl < b.N*buffSize {
			l, err = s.Read(buf)
			if err != nil {
				break
			}
			rl += l
			go func(l int) {
				_, err = s.Write(buf[:l])
				if err != nil && atomic.LoadInt32(&done) != 1 {
					b.Error(err)
					return
				}
			}(l)
		}
		if err != nil && atomic.LoadInt32(&done) != 1 {
			b.Error(err)
			return
		}
	}()

	rl := 0
	buf := make([]byte, buffSize)
	for rl < b.N*buffSize {
		l, err := c.Read(buf)
		if err != nil {
			b.Error(err)
			return
		}
		rl += l
	}
	atomic.StoreInt32(&done, 1)
	err = s.Close()
	if err != nil {
		b.Fatal(err)
	}
	err = c.Close()
	if err != nil {
		b.Fatal(err)
	}

}

func BenchmarkSyscallWrite(b *testing.B) {
	b.Fatalf("TODO")
}

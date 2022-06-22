package writev_write

import (
	"net"
	"testing"

	"j4t/j4t/mock/io"
)

const (
	packet = 4 * 1024
	repeat = 8
)

func BenchmarkWriteV(b *testing.B) {
	r, w, err := io.NewTcpConn()
	if err != nil {
		panic(err)
	}
	b.SetBytes(repeat * packet)
	go func() {
		buffs := make([][]byte, repeat)
		buffs1 := make([][]byte, repeat)
		for i := range buffs {
			buffs[i] = make([]byte, packet)
			buffs1[i] = buffs[i]
		}
		for i := 0; i < b.N; i++ {
			wl, err := (*net.Buffers)(&buffs).WriteTo(w)
			if wl != repeat*packet {
				panic(wl)
			}
			buffs = make([][]byte, 8)
			for i := range buffs {
				buffs[i] = buffs1[i]
				if len(buffs[i]) != packet {
					panic(buffs[i])
				}
			}
			if err != nil {
				panic(err)
			}
		}
	}()
	rBuff := make([]byte, packet*repeat)
	l := 0
	for l < b.N*packet*repeat {
		rl, err := r.Read(rBuff)
		if err != nil {
			panic(err)
		}
		l += rl
	}
	err = w.Close()
	if err != nil {
		panic(err)
	}
	err = r.Close()
	if err != nil {
		panic(err)
	}
}

func BenchmarkWrite(b *testing.B) {
	r, w, err := io.NewTcpConn()
	if err != nil {
		panic(err)
	}
	b.SetBytes(repeat * packet)
	go func() {
		buffs := make([][]byte, repeat)
		for i := range buffs {
			buffs[i] = make([]byte, packet)
		}
		for i := 0; i < b.N; i++ {
			for _, buff := range buffs {
				wl, err := w.Write(buff)
				if wl != packet {
					panic(wl)
				}
				if err != nil {
					panic(err)
				}
			}
		}
	}()
	rBuff := make([]byte, packet)
	l := 0
	for l < b.N*repeat*packet {
		rl, err := r.Read(rBuff)
		if err != nil {
			panic(err)
		}
		l += rl
	}
	err = w.Close()
	if err != nil {
		panic(err)
	}
	err = r.Close()
	if err != nil {
		panic(err)
	}
}

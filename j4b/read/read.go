package main

import (
	"io"
	"log"
	"net"
	"net/http"
	_ "net/http/pprof"
	"runtime"
)

var (
	gos = runtime.GOMAXPROCS(0)
)

type Reader struct {
	connes []net.Conn
	buff   []byte
}

func main() {
	go func() {
		http.ListenAndServe("localhost:9999", nil)
	}()
	err := Serve("localhost:9988")
	log.Fatalln(err)
	// ServeStd("localhost:9988")
}

func ServeStd(address string) {
	l, err := net.Listen("tcp", address)
	if err != nil {
		panic(err)
	}
	for {
		c, err := l.Accept()
		if err != nil {
			panic(err)
		}
		go func(c net.Conn) {
			var buf = make([]byte, 64*1024)
			for {
				_, err = c.Read(buf)

				if err == io.EOF {
					err = c.Close()
					if err != nil {
						panic(err)
					}
					return
				}
				if err != nil {
					panic(err)
				}
			}
		}(c)
	}
}

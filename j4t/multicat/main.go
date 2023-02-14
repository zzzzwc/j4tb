package main

import (
	"io"
	"log"
	"os"
	"strconv"
	"sync"
	"time"
)

func main() {
	if len(os.Args) < 3 {
		log.Fatalln("need file name or procs argument, like multicat ./test.data 8")
	}
	filename := os.Args[1]
	f, err := os.Open(filename)
	if err != nil {
		log.Fatalln("open", filename, "failed", err)
	}

	procs, err := strconv.ParseInt(os.Args[2], 10, 64)
	if err != nil {
		log.Fatalln("parse procs", os.Args[2], "failed", err)
	}
	info, err := f.Stat()
	if err != nil {
		log.Fatalln("stat", filename, "failed", err)
	}
	size := info.Size() / procs
	wg := &sync.WaitGroup{}
	var offset int64
	var read = Item{Name: "read_mem", CalSpeed: true}
	go (&Monitor{Items: []*Item{&read}}).Monit(3 * time.Second)
	for i := 0; i < int(procs); i++ {
		wg.Add(1)
		var loffset = offset
		go func() {
			defer wg.Done()
			f, err := os.Open(filename)
			if err != nil {
				log.Fatalln("open", filename, "failed", err)
			}
			_, err = f.Seek(loffset, io.SeekStart)
			buf := make([]byte, 4*1024*1024)
			for l := 0; l < int(size); {

				n, err := io.ReadFull(f, buf)
				if err == io.ErrUnexpectedEOF {
					break
				}
				if err != nil {
					log.Fatalln("read", filename, "failed", err)
				}
				l += n
				read.IncrN(uint64(n))
			}
		}()
		offset += size
	}
	wg.Wait()
}

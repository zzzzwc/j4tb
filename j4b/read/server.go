package main

import (
	"fmt"
	"net"
	"runtime"
	"sync/atomic"
	. "syscall"
	"time"
)

type (
	statics struct {
		accept  int64
		close   int64
		read    int64
		readMem int64
		keep    int64
		eagain  int64
	}
	worker struct {
		id       int
		errCh    chan error
		listenFd int
		inBuffer []byte
		laddr    Sockaddr
		statics  statics
	}
	Server struct {
		laddr   Sockaddr
		workers []worker
		errCh   chan error
	}
	Option    func(*Server)
	Delimiter interface {
		Write([]byte) ([]byte, bool)
	}
)

func WithSyncHandle() {

}

func Serve(address string) (err error) {
	var s Server
	var ncpu = runtime.NumCPU()
	var tcpAddr *net.TCPAddr
	tcpAddr, err = net.ResolveTCPAddr("tcp", address)
	laddr := &SockaddrInet4{Port: tcpAddr.Port, Addr: [4]byte{tcpAddr.IP[0], tcpAddr.IP[1], tcpAddr.IP[2], tcpAddr.IP[3]}}
	if err != nil {
		return
	}
	s.errCh = make(chan error, ncpu)
	s.workers = make([]worker, ncpu)
	for i := 0; i < ncpu; i++ {
		s.workers[i] = worker{id: i, errCh: s.errCh, inBuffer: make([]byte, 512*1024)}
		go (&s.workers[i]).serve(laddr)
	}
	go func() {
		ticker := time.NewTicker(3 * time.Second)
		for {
			<-ticker.C
			for _, w := range s.workers {
				fmt.Printf("%d: %+v\n", w.id, w.statics)
			}
			fmt.Println()
		}
	}()
	return <-s.errCh
}

func (w *worker) serve(addr Sockaddr) {
	runtime.LockOSThread()
	var err error
	w.listenFd, err = Socket(AF_INET, SOCK_STREAM, IPPROTO_IP)
	if err != nil {
		w.errCh <- err
		return
	}
	err = Reuseport(w.listenFd)
	if err != nil {
		w.errCh <- err
		return
	}
	err = SetNonblock(w.listenFd, true)
	if err != nil {
		w.errCh <- err
		return
	}
	err = Bind(w.listenFd, addr)
	if err != nil {
		w.errCh <- err
		return
	}
	err = Listen(w.listenFd, SOMAXCONN)
	if err != nil {
		w.errCh <- err
		return
	}
	var poll int
	var ready int
	var events = AllocEvents(64)
	poll, err = CreatePoll()
	if err != nil {
		panic(err)
	}
	err = AddFd(poll, w.listenFd)
	for {
		ready, err = Poll(poll, events)
		// fmt.Printf("%d %+v\n", ready, events.Peek(ready))
		if err != nil && err != EINTR {
			panic(err)
			w.errCh <- err
			return
		}
		for i := 0; i < ready; i++ {
			var efd = events.GetFd(i)
			if efd == w.listenFd { // only happens on reuseport enable
				err = w.accept(poll, efd)
				if err != nil && err != EAGAIN {
					panic(fmt.Errorf("%d:%d:%w", w.id, efd, err))
				}
				continue
			}
			if events.GetDataSize(i) == 0 {
				err = w.close(poll, efd)
				if err != nil {
					panic(fmt.Errorf("%d:%d:%w", w.id, efd, err))
				}
				continue
			}
			_, err = w.read(poll, efd)
			if err != nil {
				panic(fmt.Errorf("%d:%d:%w", w.id, efd, err))
			}
		}
	}
}

func (w *worker) accept(poll, efd int) error {
	var err error
	for {
		var cfd int
		// log.Println(w.id, "try accept", efd)
		cfd, _, err = Accept(efd)
		atomic.AddInt64(&w.statics.accept, 1)
		atomic.AddInt64(&w.statics.keep, 1)
		// log.Println(w.id, "accept", efd)
		if err != nil {
			return err
		}
		err = AddFd(poll, cfd)
		if err != nil {
			return err
		}
		err = SetNonblock(cfd, true)
		if err != nil {
			return err
		}
	}
}

func (w *worker) read(poll, efd int) (int, error) {
	var err error
	var l int
	var total int
	for i := 0; ; i++ {
		// log.Println(w.id, "try read", efd, buff[:l])
		l, err = Read(efd, w.inBuffer)
		atomic.AddInt64(&w.statics.read, 1)
		atomic.AddInt64(&w.statics.readMem, int64(l))
		// log.Println(w.id, "read", efd, buff[:l])

		if err == EAGAIN {
			atomic.AddInt64(&w.statics.eagain, 1)
			err = nil
			break
		}
		if err != nil {
			panic(fmt.Errorf("%d:%d:%w", w.id, efd, err))
			w.errCh <- err
			break
		}
		if l == 0 {
			err = w.close(poll, efd)
			if err != nil {
				panic(fmt.Errorf("%d:%d:%w", w.id, efd, err))
			}
			break
		}
		if l < len(w.inBuffer) {
			break
		}
		total += l
	}
	return total, err
}

func (w *worker) close(poll, efd int) error {
	var err error
	// log.Println(w.id, "close", efd)
	err = RemoveFd(poll, efd)
	if err != nil {
		return err
	}
	err = Close(efd)
	if err != nil {
		return err
	}

	atomic.AddInt64(&w.statics.close, 1)
	atomic.AddInt64(&w.statics.keep, -1)
	return err
}

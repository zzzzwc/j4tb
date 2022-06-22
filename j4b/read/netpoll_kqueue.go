//go:build darwin || dragonfly || freebsd || netbsd || openbsd

package main

import (
	. "syscall"
)

type KEvents []Kevent_t

func (e KEvents) Peek(size int) Events {
	return e[:size]
}

func (e KEvents) Size() int {
	return len(e)
}

func (e KEvents) GetFd(idx int) int {
	return int(e[idx].Ident)
}

func (e KEvents) GetDataSize(idx int) int {
	return int(e[idx].Data)
}

func (e KEvents) EOF(idx int) bool {
	return e[idx].Flags&EV_EOF == EV_EOF
}

func CreatePoll() (int, error) {
	return Kqueue()
}

func AddFd(poll, fd int) (err error) {
	var kev Kevent_t
	SetKevent(&kev, fd, EVFILT_READ, EV_ADD|EV_CLEAR)
	_, err = Kevent(poll, []Kevent_t{kev}, nil, nil)
	return
}

func AllocEvents(size int) Events {
	return KEvents(make([]Kevent_t, size))
}

func RemoveFd(poll, fd int) error {
	return nil
}

func Poll(poll int, events Events) (int, error) {
	return Kevent(poll, nil, events.(KEvents), nil)
}

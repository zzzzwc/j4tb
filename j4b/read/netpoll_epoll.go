//go:build linux

package main

import (
	. "syscall"
)

type EpollEvents []EpollEvent

func (e EpollEvents) Peek(size int) Events {
	return e[:size]
}
func (e EpollEvents) Size() int {
	return len(e)
}

func (e EpollEvents) GetFd(idx int) int {
	return int(e[idx].Fd)
}

func (e EpollEvents) GetDataSize(idx int) int {
	return -1
}

func (e EpollEvents) EOF(idx int) bool {
	return false
}

func CreatePoll() (int, error) {
	return EpollCreate1(EPOLL_CLOEXEC)
}

func AddFd(poll, fd int) error {
	var event = EpollEvent{
		Events: EPOLLIN | EPOLLHUP | -EPOLLET,
		Fd:     int32(fd),
	}

	return EpollCtl(poll, EPOLL_CTL_ADD, fd, &event)
}

func AllocEvents(size int) Events {
	return EpollEvents(make([]EpollEvent, size))
}

func RemoveFd(poll, fd int) error {
	return EpollCtl(poll, EPOLL_CTL_DEL, fd, nil)
}

func Poll(poll int, events Events) (int, error) {
	return EpollWait(poll, events.(EpollEvents), -1)
}

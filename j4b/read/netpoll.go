package main

type Events interface {
	Peek(size int) Events
	Size() int
	GetFd(idx int) int
	GetDataSize(idx int) int
	EOF(idx int) bool
}

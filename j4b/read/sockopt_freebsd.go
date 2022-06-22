//go:build freebsd

package main

import (
	. "syscall"

	"golang.org/x/sys/unix"
)

var ReusePortEnable = true

func Reuseport(fd int) error {
	return SetsockoptInt(fd, SOL_SOCKET, unix.SO_REUSEPORT_LB, 1)
}

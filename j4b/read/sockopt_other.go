//go:build !linux && !freebsd

package main

import (
	"errors"
	"fmt"
)

var ErrorNotSupport = errors.New("not support")
var ReusePortEnable = false

func Reuseport(int) error {
	return fmt.Errorf("reuseport %w", ErrorNotSupport)
}

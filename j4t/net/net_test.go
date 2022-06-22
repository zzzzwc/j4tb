package net_test

import (
	"fmt"
	"net"
	"testing"
)

func TestNet(t *testing.T) {
	c, err := net.ListenIP("ip4:ip", &net.IPAddr{
		IP: net.IP{127, 0, 0, 1},
	})
	if err != nil {
		t.Fatal(err)
	}
	for {
		buff := make([]byte, 4096)
		l, err := c.Read(buff)
		if err != nil {
			t.Fatal(err)
		}
		fmt.Println(buff[:l])
	}
}

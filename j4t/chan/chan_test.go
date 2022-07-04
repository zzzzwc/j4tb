package chan_test

import (
	"fmt"
	"testing"
)

func TestChanRecvWhenClosed(t *testing.T) {
	ch := make(chan int, 1)
	ch <- 1
	close(ch)
	v, ok := <-ch
	fmt.Println(v, ok)
	v, ok = <-ch
	fmt.Println(v, ok)
	v, ok = <-ch
	fmt.Println(v, ok)
}

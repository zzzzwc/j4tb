package local

import (
	"fmt"
	"runtime"
	"testing"
	"time"
)

func TestLocal(t *testing.T) {
	for i := 0; i < runtime.GOMAXPROCS(0); i++ {
		go func() {
			runtime.LockOSThread()
			pid := Pin()
			Unpin()
			fmt.Println(pid)
			time.Sleep(10 * time.Second)
		}()
	}
	time.Sleep(1 * time.Second)
}

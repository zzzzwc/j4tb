package process_test

import (
	"fmt"
	"testing"

	"github.com/shirou/gopsutil/v3/process"
)

func TestProcess(t *testing.T) {

	ps, err := process.Processes()
	if err != nil {
		panic(err)
	}
	for _, p := range ps {
		fmt.Println(p.OpenFiles())
	}
}

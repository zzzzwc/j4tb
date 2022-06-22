package main

import (
	"fmt"
	"os"
	"os/exec"
)

func main() {
	args := []string{
		"/Users/tiddar/go/bin/etcd",
		"--data-dir",
		"/Users/tiddar/etcd/data0",
		"--listen-client-urls",
		"http://127.0.0.1:2379",
		"--listen-peer-urls",
		"http://127.0.0.1:2380",
		"--advertise-client-urls",
		"http://127.0.0.1:2379",
		"--initial-advertise-peer-urls",
		"http://127.0.0.1:2380",
		"--initial-cluster",
		"127.0.0.1=http://127.0.0.1:2380",
		"--initial-cluster-state",
		"new",
		"--name",
		"127.0.0.1",
	}
	cmd := exec.Command(args[0])
	cmd.Args = args
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	fmt.Println(cmd)
	err := cmd.Start()
	if err != nil {
		panic(err)
	}
	err = cmd.Wait()
	if err != nil {
		panic(err)
	}
}

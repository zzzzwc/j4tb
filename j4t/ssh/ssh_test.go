package ssh_test

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"

	"golang.org/x/crypto/ssh"
)

func TestSSH(t *testing.T) {
	go serve()
	time.Sleep(2 * time.Second)
	go send()
	time.Sleep(2 * time.Second)
}

func send() {
	nc, err := net.Dial("tcp", "localhost:8989")
	if err != nil {
		panic(err)
	}
	conn, _, _, err := ssh.NewClientConn(nc, "", &ssh.ClientConfig{HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
		return nil
	}})
	if err != nil {
		panic(err)
	}
	ch, reqs, err := conn.OpenChannel("shit", []byte("shitshit"))
	println("client open")
	if err != nil {
		panic(err)
	}
	_, err = ch.SendRequest("shit", false, []byte("shit"))
	if err != nil {
		panic(err)
	}
	_, err = ch.Write([]byte("fuck"))
	if err != nil {
		panic(err)
	}
	for req := range reqs {
		println("cliet read", string(req.Payload))
		buf, err := ioutil.ReadAll(ch)
		if err != nil {
			panic(err)
		}
		println("client read", string(buf))
	}
}

func serve() {
	l, err := net.Listen("tcp", "localhost:8989")
	if err != nil {
		panic(err)
	}
	config := &ssh.ServerConfig{Config: ssh.Config{Ciphers: []string{"aes128-ctr", "aes192-ctr", "aes256-ctr", "aes128-gcm@openssh.com", "arcfour256", "arcfour128"},
		MACs: []string{"hmac-sha2-256-etm@openssh.com", "hmac-sha2-256", "hmac-sha1"}}, NoClientAuth: true}
	keyPath := filepath.Join("/tmp", fmt.Sprintf("ssh_%d.rsa", 8989))
	if _, eStat := os.Stat(keyPath); eStat != nil {
		sshKeygenPath, _ := exec.LookPath("ssh-keygen")

		bufOut := new(bytes.Buffer)
		bufErr := new(bytes.Buffer)

		cmd := exec.Command(sshKeygenPath, "-f", keyPath, "-t", "rsa", "-m", "PEM", "-N", "")
		cmd.Dir = "/tmp"
		cmd.Stdout = bufOut
		cmd.Stderr = bufErr

		err := cmd.Run()
		if err != nil {
			panic(fmt.Sprintf("Failed to generate private key: %v - %s", err, bufErr))
		}
	}

	privateBytes, err := ioutil.ReadFile(keyPath)
	if err != nil {
		panic(err)
	}
	private, err := ssh.ParsePrivateKey(privateBytes)
	if err != nil {
		panic(err)
	}
	config.AddHostKey(private)
	for {
		println("server l try accept")
		c, err := l.Accept()
		println("server l accept")
		if err != nil {
			panic(err)
		}
		_, sch, reqs, err := ssh.NewServerConn(c, config)
		if err != nil {
			panic(err)
		}
		go ssh.DiscardRequests(reqs)
		go func() {
			for ch := range sch {
				println("server try accept")
				c, reqs, err := ch.Accept()
				println("server accept")
				if err != nil {
					panic(err)
				}
				for req := range reqs {
					println("server read", string(req.Payload))
					buf, err := ioutil.ReadAll(c)
					if err != nil {
						panic(err)
					}
					println("server read", string(buf))
				}
			}
		}()
	}
}

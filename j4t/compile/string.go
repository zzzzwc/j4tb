package compile

import (
	"strconv"
)

func connect(a, b int) {
	s := strconv.Itoa(a) + strconv.Itoa(b)
	_ = s
}
func connectEscape(a, b int) string {
	return strconv.Itoa(a) + strconv.Itoa(b)
}

type test struct {
	i int
	m map[string]string
}

var t test

func routine1() {
	for {
		var m map[string]string
		m["shit"] = "fuck"
		t.m = m
	}
}

func routine2() {

}

func main() {
	go routine1()
	routine2()
}

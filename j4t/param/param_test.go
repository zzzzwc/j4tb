package param_test

import "testing"

type name struct {
	a int
}

func (n *name) add() {
	n.a++
}
func TestParam(t *testing.T) {
	var n = name{}
	n.add()
	println(n.a)
}

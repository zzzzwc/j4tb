package ctx

import (
	"context"
	"testing"
)

func TestContext(t *testing.T) {
	parent, cancel := context.WithCancel(context.Background())
	cancel()
	cancel()
	child, childcancel := context.WithCancel(parent)
	cancel()
	<-child.Done()
	childcancel()
}

package url

import (
	"fmt"
	"net/url"
	"testing"
)

func TestURL(t *testing.T) {
	vs, err := url.ParseQuery("a=b&c=的%2B1")
	if err != nil {
		panic(fmt.Errorf("invalid format of connection options, want a=b[&c=d]"))
	}
	vs.Set("aaaa", "bbbbb")
	fmt.Println(vs.Get("c"))
}

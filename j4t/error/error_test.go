package error_test

import (
	"errors"
	"fmt"
	"io"
	"testing"
)

func TestErrorWrap(t *testing.T) {
	err := fmt.Errorf("ss%v", io.EOF)
	println(errors.Is(err, io.EOF))
	err = fmt.Errorf("ss%w", io.EOF)
	println(errors.Is(err, io.EOF))
}

package time_test

import (
	"fmt"
	"log"
	"testing"
	"time"
)

func TestTimestamp(t *testing.T) {
	tt, err := time.Parse("2006", "7683")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(tt.Unix())
}

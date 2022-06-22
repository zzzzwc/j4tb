package csv

import (
	"encoding/csv"
	"math/rand"
	"os"
	"strconv"
	"testing"
	"time"
)

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
func TestWriteCSV(t *testing.T) {
	f, err := os.OpenFile("test.csv", os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		t.Fatal(err)
	}
	w := csv.NewWriter(f)
	w.Comma = '|'
	for i := 0; i < 100; i++ {
		err := w.Write([]string{strconv.Itoa(i), strconv.Itoa(i * i), RandStringRunes(60), RandStringRunes(60), time.Now().Format(time.RFC3339)})
		if err != nil {
			t.Fatal(err)
		}
	}
	w.Flush()
	if err := w.Error(); err != nil {
		t.Fatal(err)
	}
}

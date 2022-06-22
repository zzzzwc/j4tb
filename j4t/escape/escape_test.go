package escape_test

import (
	"fmt"
	"strconv"
	"testing"
)

const (
	singleQuote = '\''
)

func escape(v string) string {
	out := []rune{}

	for _, c := range v {
		switch c {
		case singleQuote:
			out = append(out, singleQuote, singleQuote)
		default:
			out = append(out, c)
		}
	}

	return string(out)
}

func escapeAndQuote(v string) string {
	return fmt.Sprintf(`'%s'`, escape(v))
}

func TestEscape(t *testing.T) {
	fmt.Println(`"'\t'"`)
	fmt.Println(`"、\t"`)
	fmt.Println(strconv.Unquote(`"'\t'"`))
	fmt.Println(strconv.Unquote(`"\t"`))
	fmt.Println(escapeAndQuote("1'2"))
}

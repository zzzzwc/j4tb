package sliceheader_mapheader_test

import (
	"bytes"
	"testing"
)

var (
	input = []byte(`Host: lgtm.lol
Accept-Encoding: deflate, gzip
authority: lgtm.lol
accept: text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9
accept-language: zh-CN,zh;q=0.9,en-US;q=0.8,en;q=0.7,zh-TW;q=0.6
cache-control: no-cache
cookie: _ga=GA1.2.219657678.1648537100
pragma: no-cache
sec-ch-ua: " Not A;Brand";v="99", "Chromium";v="102", "Google Chrome";v="102"
sec-ch-ua-mobile: ?0
sec-ch-ua-platform: "macOS"
sec-fetch-dest: document
sec-fetch-mode: navigate
sec-fetch-site: none
sec-fetch-user: ?1
upgrade-insecure-requests: 1
user-agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/102.0.0.0 Safari/537.36
`)
	queries = []string{
		"Host", "cache-control", "cookie", "user-agent",
	}
)

func loadToMap(input []byte, m map[string]string) {
	var playload = input
	var ki int
	var endi = -1
	for endi < len(playload) {
		playload = playload[endi+1:]
		ki = bytes.IndexByte(playload, ':')
		if ki < 0 {
			break
		}
		k := playload[:ki]
		endi = bytes.IndexByte(playload, '\n')
		if endi < 0 {
			break
		}
		v := playload[ki+2 : endi]
		m[string(k)] = string(v)
	}
}

func queryFromMap(m map[string]string, q string) string {
	return m[q]
}

func queryFromSlice(src []byte, q string) string {
	idx := bytes.Index(src, []byte(q))
	if idx < 0 {
		return ""
	}
	end := bytes.IndexByte(src[idx+len(q)+2:], '\n')
	return string(src[idx+len(q)+2 : idx+len(q)+2+end])
}

func BenchmarkMap(b *testing.B) {
	m := make(map[string]string)
	for i := 0; i < b.N; i++ {
		loadToMap(input, m)
		queryFromMap(m, "Host")
		queryFromMap(m, "cache-control")
		queryFromMap(m, "cookie")
		queryFromMap(m, "user-agent")
	}
}
func BenchmarkSlice(b *testing.B) {
	for i := 0; i < b.N; i++ {
		queryFromSlice(input, "Host")
		queryFromSlice(input, "cache-control")
		queryFromSlice(input, "cookie")
		queryFromSlice(input, "user-agent")
	}
}

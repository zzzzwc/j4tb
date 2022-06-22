package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/go-echarts/go-echarts/charts"
)

func main() {
	line := charts.NewLine()
	f, err := os.Open("charts.data")
	if err != nil {
		log.Fatal(err)
	}
	var xs []int
	scanner := bufio.NewScanner(f)
	var vals [64][]float32
	for i := 0; i < 10 && scanner.Scan(); {
		val := strings.TrimSpace(scanner.Text())
		if len(val) == 0 {
			xs = append(xs, i)
			i++
			continue
		}
		vs := strings.Split(val, " ")
		var vv [3]int
		vv[0], _ = strconv.Atoi(vs[0])
		vv[1], _ = strconv.Atoi(vs[1])
		vv[2], _ = strconv.Atoi(vs[2])
		vals[vv[0]] = append(vals[vv[0]], float32(vv[2]/1024/1024))
	}
	line.AddXAxis(xs)
	for i, val := range vals {
		line.AddYAxis(fmt.Sprintf("worker_%d", i), val, charts.LabelTextOpts{Show: false})
	}

	f, err = os.Create("index.html")
	if err != nil {
		log.Fatal(err)
	}
	err = line.Render(f)
	if err != nil {
		log.Fatal(err)
	}
}

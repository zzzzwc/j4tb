package ini

import (
	"fmt"
	"testing"

	"github.com/ochinchina/go-ini"
)

func TestIniEscape(t *testing.T) {
	section := `a=\n啥啥啥`
	println(section)
	println("解析后:")
	fmt.Print(ini.Load(section).String())
	fmt.Println([]byte(ini.Load(section).String()))
	println("##################################################")
	section = `a=啥啥啥`
	println(section)
	println("解析后:")
	fmt.Print(ini.Load(section).String())
	fmt.Println([]byte(ini.Load(section).String()))
}

func TestPrintf(t *testing.T) {
	str := "哈哈哈"
	println(len(str))
	fmt.Println([]byte("哈"))
	fmt.Println(fmt.Sprintf("%c", 229))
}

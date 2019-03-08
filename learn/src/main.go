// 如果别的包的函数，这个包的函数名字是小学
// 别的包调用，必须 首字母大写
package main

import (
	"fmt"
	"github.com/im-ai/pushm/learn/src/calc"
)

func main() {
	//test3()
	a := calc.Test4(1, 2)
	fmt.Println(a)
}

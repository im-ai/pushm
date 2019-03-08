// 同一个目录下的 go  包名必须一样
// 同一个目录下的，调用另外一个  无需包名引入
package main

import "fmt"

func test3() {
	fmt.Println("hello")
}

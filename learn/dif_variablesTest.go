package main

import "fmt"

var A byte

func main() {
	var A int

	// 不同作用域  就近原则
	fmt.Printf("%T\n", A)
	{
		var A float32
		fmt.Printf("%T\n", A)
	}

	Test2()
}

func Test2() {
	fmt.Printf("%T\n", A) // 全局最近
}

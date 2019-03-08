package main

import "fmt"

func Test() {
	fmt.Println(a)
}

var a int

func main() {
	a = 10

	fmt.Println(a)

	Test()
}

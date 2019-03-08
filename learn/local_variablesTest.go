package main

import (
	"fmt"
)

func test() {
	a := 10
	fmt.Println("a = ", a)
}
func main() {

	//a =111
	{
		i := 10
		fmt.Println("i=", i)
	}

	//i = 111

	if flag := 3; flag == 3 {
		fmt.Println(flag)
	}

	//flag = 111

}

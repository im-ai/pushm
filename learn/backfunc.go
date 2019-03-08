package main

import "fmt"

// 回调函数  函数有一个参数是函数类型

type FuncType func(int, int) int

// 计算器 可以进行四则运算
// 多态，多种形态,调用同一个接口，不同表现，可以实现不同表现  加减乘除
// 先有想法，后实现功能
func Cacle(a, b int, fTest FuncType) (resut int) {

	resut = fTest(a, b)
	fmt.Println("Cacle", resut)
	return
}
func add(i int, i2 int) int {
	return i + i2
}
func minus(i int, i2 int) int {
	return i - i2
}
func mul(i int, i2 int) int {
	return i * i2
}
func main() {
	Cacle(1, 2, add)
	Cacle(1, 2, minus)
	Cacle(1, 2, mul)
}

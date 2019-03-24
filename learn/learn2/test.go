package main

import "fmt"

func main() {
	// 分配一个底层数组，大小为10，并返回
	//长度为0且容量为10的切片
	v := make([]int, 0, 10)
	v = append(v, 1)
	fmt.Println(v)

	// 初始化一个类型的对象 ： slice map chan 三种
	x := make([]int, 10, 10)
	fmt.Println(&x[0]) // 地址：0xc000082050
	x = append(x, 1)   //  11  元素附加到切片的末尾
	fmt.Println(&x[0]) // 地址：0xc000088000  重新分配 内存
	fmt.Println(x)
}

package main

import (
	"fmt"
)

var (
	message = make(chan string, 3) //带有长度chan
)

/**
FIFO first in first out
*/
func sample() {
	message <- "hello goroutine1" //先进先出， 最先和 I m goroutine
	message <- "hello goroutine2"
	message <- "hello goroutine3"
	message <- "hello goroutine4" //
}
func sample1() {
	str := <-message
	str = str + "  I'm goroutine"
	message <- str
	close(message) // 下面 for循环 range 停止消费
}
func main() {
	go sample()
	go sample1()

	//fmt.Println( <- message)
	//fmt.Println( <- message)
	//fmt.Println( <- message)
	//fmt.Println( <- message)

	for str := range message { // 上面 close(message) 才能让 for 停止消费
		fmt.Println(str)
	}

	fmt.Println("hello world!")

}

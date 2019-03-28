package main

import (
	"fmt"
	"time"
)

var (
	messag1e = make(chan string, 3) //带有长度chan
)

/**
FIFO first in first out
*/
func sampl1e() {
	messag1e <- "hello goroutine1" //先进先出， 最先和 I m goroutine
	messag1e <- "hello goroutine2"
	messag1e <- "hello goroutine3"
	messag1e <- "hello goroutine4" //
}
func sampl1e1() {
	time.Sleep(2 * time.Second)
	str := <-messag1e
	str = str + "  I'm goroutine"
	messag1e <- str
}
func main() {
	go sampl1e()
	go sampl1e1()
	time.Sleep(3 * time.Second)

	fmt.Println(<-messag1e)
	fmt.Println(<-messag1e)
	fmt.Println(<-messag1e)
	fmt.Println(<-messag1e)

	fmt.Println("hello world!")

}

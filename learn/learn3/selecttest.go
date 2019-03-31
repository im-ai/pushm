package main

import (
	"fmt"
	"strconv"
	"time"
)

func sample1cn(cn chan string) {
	for i := 0; i < 5; i++ {
		cn <- "I'm sample1 num:" + strconv.Itoa(i)
		time.Sleep(3 * time.Second)
	}
}

func sample1cn2(cn chan int) {
	for i := 0; i < 5; i++ {
		cn <- i
		time.Sleep(60 * time.Second)
	}
	//close(cn)
}
func main() {
	ch1 := make(chan string, 3)
	ch2 := make(chan int, 5)
	for i := 0; i < 10; i++ {
		go sample1cn(ch1)
		go sample1cn2(ch2)
	}
	fmt.Println("hello word!")
	/*
	   chan select 组合使用 ， 选择
	   利用 select  chan 阻塞，
	*/
	for { //这里直接死循环，那个chan 来了，那个就处理
		select {
		case str, ok1 := <-ch1:
			if !ok1 {
				fmt.Println("ch1 failed")
			} else {
				fmt.Println(str)
			}

		case p, ch2 := <-ch2:
			if !ch2 {
				fmt.Println("ch2 failed")
			} else {
				fmt.Println(p)
			}
		}
	}

}

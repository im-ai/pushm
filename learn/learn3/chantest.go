package main

import (
	"fmt"
	"strconv"
	"time"
)

func samplecn(cn chan string) {
	for i := 0; i < 5; i++ {
		cn <- "I'm sample1 num:" + strconv.Itoa(i)
		time.Sleep(1 * time.Second)
	}
}

func samplecn2(cn chan int) {
	for i := 0; i < 5; i++ {
		cn <- i
		time.Sleep(2 * time.Second)
	}
	//close(cn)
}
func main() {
	ch1 := make(chan string, 3)
	ch2 := make(chan int, 5)
	for i := 0; i < 10; i++ {
		go samplecn(ch1)
		go samplecn2(ch2)
	}
	fmt.Println("hello word!")
	/*
	   chan select 组合使用 ， 选择
	*/
	for i := 0; i < 1000; i++ {
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

	time.Sleep(60 * time.Second)
}

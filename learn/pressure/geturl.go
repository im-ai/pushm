package main

import (
	"container/list"
	"fmt"
	"github.com/panjf2000/ants"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

func main() {
	defer ants.Release()

	bodybyte, ferr := ioutil.ReadFile("bookmarks.html")
	if ferr != nil {
		panic(ferr)
		return
	}
	src := string(bodybyte)

	pool, err := ants.NewPool(10000)
	if err != nil {
		panic(err)
		return
	}

	splitarr := strings.Split(src, "<A HREF=\"")
	resultarr := list.New()

	for i := 0; i < len(splitarr); i++ {
		spstr := splitarr[i]
		spstra := strings.Split(spstr, "\" ADD_DATE=")
		spstr = spstra[0]
		if strings.Index(spstr, "http") == 0 {
			resultarr.PushBack(spstr)
		}

	}

	fmt.Println(resultarr.Len())

	start := time.Now()
	ch := make(chan string)
	urlchan := make(chan string)

	go func() {
		for {
			url := <-urlchan
			err := pool.Submit(func() {
				fetch(url, ch)
			})
			if err != nil {
				fmt.Println(err)
			}
		}
	}()

	go func() {
		flag := 0
		for i := resultarr.Front(); i != nil; i = i.Next() {
			urlchan <- fmt.Sprint(i.Value)

			flag++

			if (flag % 1000) == 0 {
				fmt.Printf("running %d work \n", pool.Running())
				fmt.Printf("free %d work \n", pool.Free())
				fmt.Printf("index %d work \n\n", flag)
			}
		}
	}()

	for i := resultarr.Front(); i != nil; i = i.Next() {
		fmt.Println(<-ch)
	}

	fmt.Printf("%.2fs elapsed\n", time.Since(start).Seconds())
}

func fetch(url string, ch chan<- string) {
	start := time.Now()
	resp, err := http.Get(url)
	if err != nil {
		ch <- fmt.Sprint(err)
		return
	}

	nbytes, err := io.Copy(ioutil.Discard, resp.Body)
	resp.Body.Close()
	if err != nil {
		ch <- fmt.Sprintf("while reading %s: %v", url, err)
		return
	}
	secs := time.Since(start).Seconds()
	ch <- fmt.Sprintf("%.2fs  %7d  %s", secs, nbytes, url)
}

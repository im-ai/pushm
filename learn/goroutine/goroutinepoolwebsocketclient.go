package main

import (
	"flag"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/panjf2000/ants"
	"log"
	"net/url"
	"os"
	"os/signal"
	"sync"
	"time"
)

var addr = flag.String("addr", "192.168.9.142:7776", "http service address")

func main() {
	p, _ := ants.NewPool(200000)
	defer p.Release()

	var wg sync.WaitGroup

	go func() {
		for {
			fmt.Printf("pool, running workers number:%d\n", p.Running())
			fmt.Printf("pool, Cap workers number:%d\n", p.Cap())
			fmt.Printf("pool, Free workers number:%d\n", p.Free())
			time.Sleep(1 * time.Second)
		}
	}()

	for i := 0; i < 100000; i++ {
		wg.Add(1)
		_ = p.Submit(func() {
			wscient()
			wg.Done()
		})
	}
	wg.Wait()

}
func wscient() {
	flag.Parse()
	log.SetFlags(0)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	u := url.URL{Scheme: "ws", Host: *addr, Path: "/ws"}
	log.Printf("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	done := make(chan struct{})

	go func() {
		defer close(done)
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}
			log.Printf("recv: %s", message)
		}
	}()

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-done:
			return
		case t := <-ticker.C:
			err := c.WriteMessage(websocket.TextMessage, []byte(t.String()))
			if err != nil {
				log.Println("write:", err)
				return
			}
		case <-interrupt:
			log.Println("interrupt")

			// Cleanly close the connection by sending a close message and then
			// waiting (with timeout) for the server to close the connection.
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("write close:", err)
				return
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			return
		}
	}
}

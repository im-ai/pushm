package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/im-ai/pushm/impl"
	"github.com/panjf2000/ants"
	"net/http"
	"sync/atomic"
	"time"
)

var (
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	running int32
)

func incRunning() {
	atomic.AddInt32(&running, 1)
}
func decRunning() {
	atomic.AddInt32(&running, -1)
}
func main() {
	p, _ := ants.NewPool(200000)
	defer p.Release()

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		incRunning()
		var (
			wsConn *websocket.Conn
			err    error
			data   []byte
			conn   *impl.Connection
		)

		// 升级 为 websocket 请求
		if wsConn, err = upgrader.Upgrade(w, r, nil); err != nil {
			return
		}

		// 初始化
		if conn, err = impl.InitConnetcion(wsConn); err != nil {
			goto ERR
		}

		err = p.Submit(func() {
			var (
				err error
			)
			for {
				if err = conn.WriterMessage([]byte("hearbeat")); err != nil {
					return
				}
				fmt.Printf("pool, running workers number:%d\n", p.Running())
				fmt.Printf("pool, Cap workers number:%d\n", p.Cap())
				fmt.Printf("pool, Free workers number:%d\n", p.Free())
				fmt.Printf("client, client workers number:%d\n\n", running)
				time.Sleep(1 * time.Second)
			}
		})

		for {
			err := p.Submit(func() {
				if data, err = conn.ReadMessage(); err != nil {
					return
				}
				if err = conn.WriterMessage(data); err != nil {
					return
				}
				return
			})
			if err != nil {
				http.Error(w, "throttle limit error", http.StatusInternalServerError)
				goto ERR
			}
		}
		decRunning()
	ERR:
		conn.Close()

	})
	http.ListenAndServe(":7777", nil)
}

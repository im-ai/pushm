package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/im-ai/pushm/impl"
	"github.com/panjf2000/ants"
	"net/http"
	"time"
)

var (
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

func main() {
	p, _ := ants.NewPool(100000)
	defer p.Release()

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
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
				fmt.Printf("pool, Free workers number:%d\n\n", p.Free())
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
	ERR:
		conn.Close()
	})
	http.ListenAndServe(":7777", nil)
}

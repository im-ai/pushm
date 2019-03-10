package main

import (
	"github.com/gorilla/websocket"
	"github.com/im-ai/pushm/impl"
	"net/http"
	"time"
)

var (
	upgrader = websocket.Upgrader{

		// 允许跨域
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

func wsHandler(w http.ResponseWriter, r *http.Request) {

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

	go func() {
		var (
			err error
		)
		for {
			if err = conn.WriterMessage([]byte("hearbeat")); err != nil {
				return
			}
			time.Sleep(5 * time.Second)
		}

	}()

	for {
		if data, err = conn.ReadMessage(); err != nil {
			goto ERR
		}
		if err = conn.WriterMessage(data); err != nil {
			goto ERR
		}
	}

ERR:
	conn.Close()
}

func main() {
	http.HandleFunc("/ws", wsHandler)
	http.ListenAndServe(":7777", nil)
}

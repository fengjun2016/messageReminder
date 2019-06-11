package main

import (
	"net/http"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"github.com/fengjun2016/messageReminder/server/impl"
	"time"
)

var (
	upgrader = websocket.Upgrader{
		//允许跨域
		CheckOrigin:func(req *http.Request) bool {
			return true
		},
	}
)

func wsHandler(rw http.ResponseWriter, req *http.Request) {
	// rw.Write([]byte("hello"))
	var (
		wsConn *websocket.Conn
		err error
		conn *impl.Connection
		data []byte
	)

	//完成ws协议的握手操作
	//Upgrade;websocket
	if wsConn, err = upgrader.Upgrade(rw, req, nil); err != nil {
		return
	}

	if conn, err = impl.InitConnection(wsConn); err != nil {
		goto ERR
	}

	// 启动协程 不断发消息
	go func() {
		var err error
		for {
			if err = conn.WriteMessage([]byte("heartbeat")); err != nil {
				return
			}
			time.Sleep(1*time.Second)
		}
	}()


	for {
		if data, err = conn.ReadMessage(); err != nil {
			goto ERR
		}

		if err = conn.WriteMessage(data); err != nil {
			goto ERR
		}
	}

	ERR:
		conn.Close()
}
func main() {
	http.HandleFunc("/ws", wsHandler)
	logrus.Info("Starting application on address `0.0.0.0:7777`")
	logrus.Println(http.ListenAndServe("0.0.0.0:7777", nil))
}
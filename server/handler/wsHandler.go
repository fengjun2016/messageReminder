package handler

import(
	"time"
	"strconv"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"

	"github.com/fengjun2016/messageReminder/server/impl"
)

var (
	upgrader = websocket.Upgrader{
		//允许跨域
		CheckOrigin:func(req *http.Request) bool {
			return true
		},
	}
)

func WsHandler(rw http.ResponseWriter, req *http.Request) {
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

	// 启动协程 不断发消息 发送心跳包
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
		logrus.Println("Connection closed")
}

// 使用golang提供的time中的timer中的定时器功能 定时给前端发布未读消息提醒
func UnReadMessageNumHandler(rw http.ResponseWriter, req *http.Request) {
		var (
		wsConn *websocket.Conn
		err error
		conn *impl.Connection
	)

	//完成ws协议的握手操作
	//Upgrade;websocket
	if wsConn, err = upgrader.Upgrade(rw, req, nil); err != nil {
		return
	}

	if conn, err = impl.InitConnection(wsConn); err != nil {
		conn.Close()
		logrus.Println("Connection closed")
	}

	d := time.Duration(time.Second * 2)

	t := time.NewTicker(d)
	defer t.Stop()
	num := 0

	for {
		<- t.C

		num++
		// 查找获取数据库里面的未读消息的记录条数
		logrus.Println("unRead message num:", num)
		numString := strconv.Itoa(num)

		if err = wsConn.WriteMessage(websocket.TextMessage, []byte(numString)); err != nil {
			conn.Close()
			logrus.Println("Connection closed")
		}
	}
}

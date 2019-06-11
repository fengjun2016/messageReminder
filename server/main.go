package main

import (
	"net/http"
	"github.com/sirupsen/logrus"

	"github.com/fengjun2016/messageReminder/server/handler"
	"github.com/fengjun2016/messageReminder/server/config"
)

func main() {
	config.Init()

	// listen and serve
	http.HandleFunc("/ws", handler.WsHandler)
	http.HandleFunc("/unread", handler.UnReadMessageNumHandler)
	
	logrus.Info("Starting application on address `0.0.0.0:7777`")
	logrus.Println(http.ListenAndServe("0.0.0.0:7777", nil))
}
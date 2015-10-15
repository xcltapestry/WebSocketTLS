/*
Websocket单向认证例子
服务端
Auther:Xiong Chuan Liang

例子URL： *.xcl.com:8050/xcl

*/
package main

import (
	"golang.org/x/net/websocket"
	"log"
	"net/http"
)

var (
	CACertFile  string = "ca.crt"
	SrvCertFile string = "server.crt"
	SrvKeyFile  string = "server.key"
	BindUrl     string = ":8050"
)

// Echo the data received on the WebSocket.
func EchoServer(ws *websocket.Conn) {
	log.Println("[EchoServer] begin...")

	b := make([]byte, 128)
	for {
		n, err := ws.Read(b)
		if err != nil {
			log.Println("[EchoServer] Read err:", err)
			break
		}

		log.Println("[EchoServer] Read:", string(b))
		if _, err := ws.Write(b[:n]); err != nil {
			log.Println("Write err:", err)
			break
		}
	}
	log.Println("[EchoServer] end...")
}

func WebsocketTLS() {

	http.Handle("/xcl", websocket.Handler(EchoServer))
	err := http.ListenAndServeTLS(BindUrl, SrvCertFile, SrvKeyFile, nil)
	if err != nil {
		panic("[WebsocketTLS] ListenAndServe: " + err.Error())
	}

}

func main() {
	WebsocketTLS()
}

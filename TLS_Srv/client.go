/*
Websocket单向认证例子
客户端
Auther:Xiong Chuan Liang
*/
package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"golang.org/x/net/websocket"
	"io/ioutil"
	"log"
)

var (
	CACertFile  string = "ca.crt"
	CliCertFile string = "client.crt"
	CliKeyFile  string = "client.key"
	BindUrl     string = "*.xcl.com:8050/xcl"
	BindSrvName string = "*.xcl.com"
)

func WebsocketTLS_Client() {

	caCrt, err := ioutil.ReadFile(CACertFile)
	if err != nil {
		log.Fatal("[WebsocketTLS_Client] ReadFile err:", err)
	}

	pool := x509.NewCertPool()
	pool.AppendCertsFromPEM(caCrt)

	origin := fmt.Sprintf("https://%s", BindUrl)
	url := fmt.Sprintf("wss://%s", BindUrl)
	conf, err := websocket.NewConfig(url, origin)
	if err != nil {
		log.Fatal("[WebsocketTLS_Client] NewConfig err:", err)
	}
	conf.TlsConfig = &tls.Config{
		RootCAs:    pool,
		ServerName: BindSrvName,
	}

	ws, err := websocket.DialConfig(conf)
	if err != nil {
		log.Fatal("[WebsocketTLS_Client] DialConfig err:", err)
	}

	if _, err := ws.Write([]byte("Hello World!")); err != nil {
		log.Fatal("[WebsocketTLS_Client] Write err:", err)
	} else {
		log.Println("[WebsocketTLS_Client] Write succeed!")
	}

	var n int
	for {
		var msg = make([]byte, 512)
		if n, err = ws.Read(msg); err != nil {
			log.Fatal("[WebsocketTLS_Client] Read:", err)
		}
		log.Println("[WebsocketTLS_Client] Received: %s.\n", msg[:n])
	}
}

func main() {
	WebsocketTLS_Client()
}

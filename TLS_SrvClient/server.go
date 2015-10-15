/*
Websocket双向认证例子
服务端
Auther:Xiong Chuan Liang

例子URL： *.xcl.com:8051/xcl
*/
package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"golang.org/x/net/websocket"
	"io/ioutil"
	"log"
	"net"
	"net/http"
)

var (
	CACertFile  string = "ca.crt"
	SrvCertFile string = "server.crt"
	SrvKeyFile  string = "server.key"
	BindUrl     string = ":8051"
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

	caCrt, err := ioutil.ReadFile(CACertFile)
	if err != nil {
		fmt.Println("[WebsocketTLS] CA ReadFile err:", err)
		return
	}

	cliCrt, err := tls.LoadX509KeyPair(SrvCertFile, SrvKeyFile)
	if err != nil {
		fmt.Println("[WebsocketTLS] LoadX509KeyPair err:", err)
		return
	}

	pool := x509.NewCertPool()
	pool.AppendCertsFromPEM(caCrt)

	config := &tls.Config{}
	config.ClientAuth = tls.RequireAndVerifyClientCert
	config.ClientCAs = pool
	config.Certificates = []tls.Certificate{cliCrt}

	httpServeMux := http.NewServeMux()
	httpServeMux.Handle("/xcl", websocket.Handler(EchoServer))

	server := &http.Server{Addr: BindUrl, Handler: httpServeMux}
	server.SetKeepAlivesEnabled(true)

	ln, err := net.Listen("tcp", BindUrl)
	if err != nil {
		fmt.Println("[WebsocketTLS] Listen err:", err)
		return
	}

	tlsListener := tls.NewListener(ln, config)
	if err = server.Serve(tlsListener); err != nil {
		log.Println("[WebsocketTLS] Server BindUrl:", BindUrl, " err:", err)
		return
	}
}

func main() {
	WebsocketTLS()
}

package main

import (
	"io"
	"net"
)

func main() {
	//Listens for client
	listener, err := net.Listen("tcp", "127.0.0.1:1111")
	if err != nil {
		panic(err)
	}
	defer listener.Close()

	for {
		clientConn, err := listener.Accept()
		if err != nil {
			panic(err)
		}
		go handleProxyConn(clientConn)
	}
}

func handleProxyConn(client net.Conn) {
	//Makes handshake to server and copies the data of client to give to server
	serverConn, err := net.Dial("tcp", "127.0.0.1:1234")
	if err != nil {
		panic(err)
	}
	defer client.Close()
	go func() {
		_, err = io.Copy(serverConn, client)
		if err != nil {
			panic(err)
		}
	}()
	_, err = io.Copy(client, serverConn)
	if err != nil {
		panic(err)
	}
}

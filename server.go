package main

import (
	"encoding/binary"
	"fmt"
	"net"
	"strings"
	"time"
)

func main() {
	listener, err := net.Listen("tcp", "127.0.0.1:1234")
	if err != nil {
		panic(err)
	}
	defer listener.Close()

	for {
		clientConn, err := listener.Accept()
		if err != nil {
			panic(err)
		}
		go handleServerConn(clientConn) //Menggunakan goroutine untuk dapat multithreading
	}
}

func handleServerConn(client net.Conn) {
	var size uint32
	err := binary.Read(client, binary.LittleEndian, &size)
	if err != nil {
		panic(err)
	}

	// read deadline 10 detik kedepan
	deadline := time.Now().Add(10 * time.Second)
	err = client.SetReadDeadline(deadline)
	if err != nil {
		panic(err)
	}

	bytMsg := make([]byte, size)
	_, err = client.Read(bytMsg)
	if err != nil {
		panic(err)
	}

	strMsg := string(bytMsg)
	fmt.Printf("Recieved: %s\n", strMsg)

	//Membuat message untuk dibalikkan ke client
	var reply string
	if strings.HasSuffix(strMsg, ".zip") {
		reply = "File has been recieved"
	} else if strings.Contains(strMsg, ".") {
		reply = "Only zip files can be sent"
	} else {
		reply = "Message has been recieved"
	}
	err = binary.Write(client, binary.LittleEndian, uint32(len(reply)))
	if err != nil {
		panic(err)
	}

	// Write deadline 5 detik kedepan
	deadline = time.Now().Add(5 * time.Second)
	err = client.SetWriteDeadline(deadline)
	if err != nil {
		panic(err)
	}

	_, err = client.Write([]byte(reply))
	if err != nil {
		panic(err)
	}
	// Sudah berhasil mengirim balik message!
}

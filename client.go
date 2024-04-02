package main

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"net"
	"os"
	"strings"
	"time"
)

func menu() {
	for {
		sendMessageMenu()
	}
}

func sendMessageMenu() {
	scanner := bufio.NewScanner(os.Stdin)
	var msg string
	for {
		fmt.Print("Please insert your message: ")
		scanner.Scan()
		msg = scanner.Text()
		if len(msg) < 10 {
			fmt.Println("Message cannot be less than 10 character")
		} else if strings.Contains(msg, "kasar") {
			fmt.Println("Message cannot contain bad words")
		} else {
			break
		}
	}
	sendMessagetoServer(msg)
}

func sendMessagetoServer(message string) {
	serverConn, err := net.DialTimeout("tcp", "127.0.0.1:1111", 3*time.Second)
	if err != nil {
		panic(err)
	}
	defer serverConn.Close()
	//Karena kita menggunakan dynamic size untuk payload, kita akan memasukan ukuran lessage kedalam serverConn
	err = binary.Write(serverConn, binary.LittleEndian, uint32(len(message)))
	if err != nil {
		panic(err)
	}

	// deadline hingga 5 detik kedepan
	deadline := time.Now().Add(5 * time.Second)
	err = serverConn.SetWriteDeadline(deadline)
	if err != nil {
		panic(err)
	}

	//Datanya apa yang mau ditulis kedalam serverConn. Hanya menerima byte, oleh karenea itu harus di typecast menjadi byte
	_, err = serverConn.Write([]byte(message))
	if err != nil {
		panic(err)
	}
	//================================== Data sudah berhasil dikirimkan ke server! =========================================

	var size uint32
	err = binary.Read(serverConn, binary.LittleEndian, &size)
	if err != nil {
		panic(err)
	}

	// Read deadline hingga 10 detik kedepan
	deadline = time.Now().Add(10 * time.Second)
	err = serverConn.SetReadDeadline(deadline)
	if err != nil {
		panic(err)
	}

	bytReply := make([]byte, size)
	_, err = serverConn.Read(bytReply)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Reply: %s\n", bytReply)
}

func main() {
	menu()
}

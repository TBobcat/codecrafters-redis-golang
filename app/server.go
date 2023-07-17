package main

import (
	"fmt"
	"net"
	"os"
)

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	l, err := net.Listen("tcp", "127.0.0.1:6379")
	if err != nil {
		fmt.Println("Failed to bind to port 6379")
		os.Exit(1)
	}
	conn, err := l.Accept()
	if err != nil {
		fmt.Println("Error accepting connection: ", err.Error())
		os.Exit(1)
	}

	// // incoming request, make buffer size to 1024 byte
	// buffer := make([]byte, 1024)
	// data, err := conn.Read(buffer)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// this converts strings to bytes as response and sends through socket
	conn.Write([]byte("+PONG\r\n"))

	conn.Close()

}

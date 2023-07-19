package main

import (
	"fmt"
	"io"
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

	// this loops keeps the server up and listening
	for {
		connection, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}

		go handler(connection)
	}

}

func handler(conn net.Conn) {

	buffer := make([]byte, 1024)
	// the loop handles data sent through the same connection opened above
	for {
		// sends response so long as data is sent to socket, buffer contains data
		_, err := conn.Read(buffer)
		if err != nil {
			// had to put this within outter err checking block
			// to handle EOF(no more data sent) for some reason
			if err == io.EOF {
				break
			} else {
				fmt.Println("Error reading data sent through connection", err.Error())
				os.Exit(1)
			}
		}

		// this converts strings to bytes as response and sends through socket
		_, err = conn.Write([]byte("+PONG\r\n"))
		if err != nil {
			fmt.Println("Error sending response", err.Error())
			os.Exit(1)
		}
	}
}

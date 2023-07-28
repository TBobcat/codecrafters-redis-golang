package main

import (
	"fmt"
	"io"
	"net"
	"os"
	"strings"
	//"k8s.io/utils/strings/slices" non standard modules not supported for codecrafters' testing
)

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	mem := make(map[string]string)

	l, err := net.Listen("tcp", "127.0.0.1:6379")
	if err != nil {
		fmt.Println("Failed to bind to port 6379")
		os.Exit(1)

	}

	// this loops keeps the server up and listening to handle requests concurrently
	for {
		connection, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}

		go handler(connection, mem)
	}
}

func handler(conn net.Conn, m map[string]string) {

	// buffer is an array of bytes
	buffer := make([]byte, 1024)
	// the loop handles data sent through the same connection opened above
	for {
		// sends response so long as data is sent to socket, buffer contains data
		n, err := conn.Read(buffer)
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

		// buffer is 1024 size, only take up to n as it's the size buffer is used
		// rest is not filled
		request := string(buffer[:n])
		split := strings.Split(request, "\r\n")
		var resp []byte

		// for debugging
		fmt.Println("Received request:", request)
		fmt.Println("Number of parts:", len(split))

		if strings.Contains(request, "ping") {
			resp = []byte("+PONG\r\n")

		} else if strings.Contains(request, "echo") {

			resp = []byte(fmt.Sprintf("+%s\r\n", split[4]))

		} else if strings.Contains(request, "set") {
			m[split[4]] = split[6]
			resp = []byte("+OK\r\n")

		} else if strings.Contains(request, "get") {
			value := m[split[4]]
			if value != "" {
				resp = []byte(fmt.Sprintf("+%s\r\n", value))
			} else {
				resp = []byte(fmt.Sprintf("+%s\r\n", "key not found"))
			}

		}

		// wrap string with redis protocol format and convert to byte array,
		// and send through connection
		_, err = conn.Write(resp)
		if err != nil {
			fmt.Println("Error sending response", err.Error())
			os.Exit(1)
		}

	}
}

/* decode
 */
// func parse([]byte lst) {
// 	decoded_slst = byte[]

// }

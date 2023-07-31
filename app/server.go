package main

import (
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
	//"k8s.io/utils/strings/slices" non standard modules not supported for codecrafters' testing
)

// interface lets Pair p use fields with type checking like p.value.(int64)
type Pair struct {
	value, exp interface{}
}

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	mem := make(map[string]Pair)

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

func handler(conn net.Conn, m map[string]Pair) {

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
			// sets value and expiration time if px is given, using Pair struct
			if strings.Contains(request, "px") {
				// convert string to int64
				i, err := strconv.ParseInt(split[10], 10, 64)
				if err != nil {
					fmt.Println("error converting input milli seconds to int64")
				}

				// make pair as value to store in mem
				exp_time := time.Now().UnixMilli() + i
				p := Pair{split[6], exp_time}
				m[split[4]] = p

			} else {
				m[split[4]] = Pair{split[6], nil}
			}

			resp = []byte("+OK\r\n")

		} else if strings.Contains(request, "get") {
			p, exists := m[split[4]]

			// only check time if p.exp is not nil
			if exists && p.exp != nil {
				val_expired := time.Now().UnixMilli()-p.exp.(int64) > 0

				if !val_expired {
					resp = []byte(fmt.Sprintf("+%s\r\n", p.value))
				} else {
					resp = []byte("+value for gien key expired\r\n")
				}

			} else if exists {
				resp = []byte(fmt.Sprintf("+%s\r\n", p.value))
			} else {
				resp = []byte("-Error: Key not found\r\n")
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

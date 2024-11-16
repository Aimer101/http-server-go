package main

import (
	"fmt"
	"net"
	"os"

	"github.com/codecrafters-io/http-server-starter-go/helper"
)

// Ensures gofmt doesn't remove the "net" and "os" imports above (feel free to remove this!)
var _ = net.Listen
var _ = os.Exit

func main() {

	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}

	defer l.Close()

	fmt.Println("Listening on port 4221")

	for {
		conn, err := l.Accept()

		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			continue
		}

		go handleConnection(conn)
	}

}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	buffer := make([]byte, 1024)

	nBytes, err := conn.Read(buffer)

	if err != nil {
		fmt.Println("Error reading from connection: ", err.Error())
		return
	}

	request, err := helper.ParseRawBuffer(buffer, nBytes)

	fmt.Println("Parsed request: ", request)

	if err != nil {
		fmt.Println("Error parsing request: ", err.Error())
		return
	}

	var response string

	switch request.Endpoint[0] {

	case "abcdefg":
		response = "HTTP/1.1 404 Not Found\r\n\r\n"

	case "echo":
		response = helper.ReturnHttpOkWithResponseBody(request.Endpoint[1])

	case "user-agent":
		response = helper.ReturnHttpOkWithResponseBody(request.Headers["User-Agent"])
	default:
		response = "HTTP/1.1 200 OK\r\n\r\n"

	}

	fmt.Println("Response is : ", response)

	_, err = conn.Write([]byte(response))
	if err != nil {
		fmt.Println("Error writing to connection: ", err.Error())
		return
	}

}

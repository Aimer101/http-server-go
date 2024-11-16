package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"path/filepath"

	"github.com/codecrafters-io/http-server-starter-go/helper"
)

// Ensures gofmt doesn't remove the "net" and "os" imports above (feel free to remove this!)
var _ = net.Listen
var _ = os.Exit

func main() {
	directory := flag.String("directory", "", "directory for serving files")
	flag.Parse()

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

		go handleConnection(conn, *directory)
	}

}

func handleConnection(conn net.Conn, directory string) {
	defer conn.Close()

	buffer := make([]byte, 1024)

	fmt.Println("directory is : ", directory)

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

	case "":
		response = "HTTP/1.1 200 OK\r\n\r\n"

	case "files":
		fmt.Println("file name is: ", request.Endpoint[1])

		fullPath := filepath.Join(directory, request.Endpoint[1])

		file, err := os.ReadFile(fullPath)

		if err != nil {
			response = helper.ReturnHttpNotFound()
		} else {
			response = helper.ReturnFileHttpOkWithResponseBody(&file)
		}

	case "echo":
		response = helper.ReturnHttpOkWithResponseBody(request.Endpoint[1], "text/plain")

	case "user-agent":
		response = helper.ReturnHttpOkWithResponseBody(request.Headers["User-Agent"], "text/plain")
	default:
		response = helper.ReturnHttpNotFound()

	}

	fmt.Println("Response is : ", response)

	_, err = conn.Write([]byte(response))
	if err != nil {
		fmt.Println("Error writing to connection: ", err.Error())
		return
	}

}

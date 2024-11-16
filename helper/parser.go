package helper

import (
	"fmt"
	"strings"
)

type Request struct {
	Method   string
	Endpoint []string
	Headers  map[string]string // Add a map to store headers
}

func ParseRawBuffer(buffer []byte, nBytes int) (Request, error) {
	res := string(buffer[:nBytes])

	fmt.Println("Raw request:" + res)

	lines := strings.Split(res, "\r\n")

	if len(lines) == 0 || len(lines[0]) == 0 {
		return Request{}, fmt.Errorf("invalid HTTP request format")
	}

	requestLine := strings.Split(lines[0], " ") // 0th is the request method, 1 is endpoint, 2 is http version

	headers := make(map[string]string)
	for _, line := range lines[1:] {
		if line == "" { // End of headers section
			break
		}
		// Split header key and value by the first colon
		parts := strings.SplitN(line, ": ", 2)
		if len(parts) == 2 {
			headers[parts[0]] = parts[1]
		}
	}

	return Request{
		Method:   requestLine[0],
		Endpoint: strings.Split(strings.Trim(requestLine[1], "/"), "/"),
		Headers:  headers,
	}, nil
}

func ReturnHttpOkWithResponseBody(body string) string {
	return fmt.Sprintf(
		"HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s",
		len(body),
		body,
	)
}

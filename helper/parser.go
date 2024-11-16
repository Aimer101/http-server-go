package helper

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"strconv"
	"strings"
)

type Request struct {
	Method   string
	Endpoint []string
	Headers  map[string]string // Add a map to store headers
	Body     *string
}

func ParseRawBuffer(buffer []byte, nBytes int) (Request, error) {
	res := string(buffer[:nBytes])

	fmt.Println("Raw request:" + res)

	lines := strings.Split(res, "\r\n")

	if len(lines) == 0 || len(lines[0]) == 0 {
		return Request{}, fmt.Errorf("invalid HTTP request format")
	}

	for i, line := range lines {
		fmt.Printf("%d: %s\n", i, line)
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

	var body *string

	// see if headers contain Content-Length
	if _, ok := headers["Content-Length"]; ok {
		if strings.ToUpper(requestLine[0]) == "POST" {
			body = &lines[len(lines)-1]
		}
	}

	return Request{
		Method:   requestLine[0],
		Endpoint: strings.Split(strings.Trim(requestLine[1], "/"), "/"),
		Headers:  headers,
		Body:     body,
	}, nil
}

func ReturnHttpOkWithResponseBody(body string, headers map[string]string) string {
	var builder strings.Builder
	builder.WriteString("HTTP/1.1 200 OK\r\n")
	builder.WriteString("Content-Type: text/plain\r\n")

	if _, ok := headers["Accept-Encoding"]; ok {
		compression_types := strings.Split(headers["Accept-Encoding"], ",")

		for _, compression_type := range compression_types {

			if strings.ToLower(strings.TrimSpace(compression_type)) == "gzip" {
				var buffer bytes.Buffer

				w := gzip.NewWriter(&buffer)
				w.Write([]byte(body))
				w.Close()

				body = buffer.String()

				builder.WriteString("Content-Length: " + strconv.Itoa(len(body)) + "\r\n")

				builder.WriteString("Content-Encoding: " + compression_type + "\r\n")
			}
		}
	} else {
		builder.WriteString("Content-Length: " + strconv.Itoa(len(body)) + "\r\n")

	}

	builder.WriteString("\r\n")
	builder.WriteString(body)

	return builder.String()
}

func ReturnFileHttpOkWithResponseBody(data *[]byte, headers map[string]string) string {
	var builder strings.Builder
	builder.WriteString("HTTP/1.1 200 OK\r\n")
	builder.WriteString("Content-Type: application/octet-stream\r\n")
	builder.WriteString("Content-Length: " + strconv.Itoa(len(*data)) + "\r\n")

	if _, ok := headers["Accept-Encoding"]; ok {
		compression_types := strings.Split(headers["Accept-Encoding"], ",")

		for _, compression_type := range compression_types {
			if strings.ToLower(strings.TrimSpace(compression_type)) == "gzip" {
				builder.WriteString("Content-Encoding: " + compression_type + "\r\n")
			}
		}
	}

	builder.WriteString("\r\n")
	builder.WriteString(string(*data))
	return builder.String()
}

func ReturnHttpNotFound() string {
	return "HTTP/1.1 404 Not Found\r\n\r\n"
}

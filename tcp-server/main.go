package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"strings"
	"time"
)

func main() {
	l, err := net.Listen("tcp4", ":1234")
	if err != nil {
		log.Fatal(err)
	}

	defer l.Close()

	for {
		c, err := l.Accept()
		if err != nil {
			fmt.Printf("Error while accepting connection, %v\n", err)
			continue
		}

		deadline := time.Now().Add(10 * time.Second)
		c.SetDeadline(deadline)

		go handleConnection(c)
	}
}

func handleConnection(c net.Conn) {
	packet := make([]byte, 1024)
	finalData := make([]byte, 1024)

	defer close(&c)

	var totalBytes int
	for {
		n, err := c.Read(packet)

		if err != nil {
			if err == io.EOF {
				fmt.Println("EOF")
				break
			}
			fmt.Printf("Error while reading %v\n", err)
		}

		finalData = append(finalData, packet[:n]...)

		num, err := c.Write(packet[:n])
		if err != nil {
			fmt.Printf("Error while writing %v\n", err)
			break
		}

		fmt.Printf("Wrote back %d bytes, the payload is %s\n", num, string(finalData))

		if strings.Contains(string(finalData), "\r\n\r\n") {
			idx := strings.Index(string(finalData), "\r\n\r\n")
			totalBytes = len([]byte(finalData[idx+4:]))

			fmt.Printf("totalBytes: %v\n", totalBytes)

			body, headerMap, error := parseHttpRequest(string(finalData))
			if error != nil {
				fmt.Println("Unable to parse request")
				break
			}

			contentLength, exists := headerMap["Content-Length"]
			if !exists {
				fmt.Println("Content length does not exists")
				break
			}

			lengthStr, ok := contentLength.(string)
			if ok {
				data, err := strconv.Atoi(lengthStr)
				if err != nil {
					fmt.Println("Unable to convert Content-Length to int:", err)
					return
				}

				if totalBytes == data {
					fmt.Printf("Body: %v\n", body)
					break
				}
			}

		}
	}
}

func close(conn *net.Conn) {
	fmt.Println("Closing TCP connection")
	(*conn).Close()
}

func parseHttpRequest(rawRequest string) (string, map[string]any, error) {
	parts := strings.SplitN(rawRequest, "\r\n\r\n", 2)
	if len(parts) < 2 {
		return "", nil, fmt.Errorf("invalid request")
	}

	headers := parts[0]
	body := parts[1]

	parsedHeaders := strings.Split(headers, "\r\n")
	fmt.Printf("parsedHeaders: %v\n", parsedHeaders[1:])

	headerMap := make(map[string]any)
	for _, v := range parsedHeaders[1:] {
		_, exists := headerMap[strings.TrimSpace(strings.Split(v, ":")[0])]

		if !exists {
			headerMap[strings.TrimSpace(strings.Split(v, ":")[0])] = strings.TrimSpace(strings.Split(v, ":")[1])
		}
	}

	fmt.Printf("headerMap: %v\n", headerMap)

	return body, headerMap, nil
}

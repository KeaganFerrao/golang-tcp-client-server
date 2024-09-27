package main

import (
	"fmt"
	"io"
	"log"
	"net"
)

func main() {
	tcpAddr, err := net.ResolveTCPAddr("tcp", "127.0.0.1:1234")
	if err != nil {
		log.Fatal(err)
	}

	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		log.Fatal(err)
	}

	defer close(conn)

	fmt.Println("Enter data to send")

	request := "POST /submit HTTP/1.1\r\n" +
		"Host: example.com\r\n" +
		"Content-Length: 21\r\n" +
		"\r\n" +
		"Hello, world12345454!"

	_, err = conn.Write([]byte(request))
	if err != nil {
		log.Fatal(err)
	}

	reply := make([]byte, 2048)

	for {
		n, err := conn.Read(reply)
		if err != nil {
			if err == io.EOF {
				fmt.Println("---EOF---")
				break
			}
			log.Fatal(err)
		}

		fmt.Printf("Reply from server %v\n", string(reply[:n]))
	}
}

func close(conn *net.TCPConn) {
	fmt.Println("Closing TCP connection")
	conn.Close()
}

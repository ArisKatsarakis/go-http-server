package main

import (
	"fmt"
	"log"
	"net"

	"boot.ariskatsarakis.gr/internal/request"
)

func main() {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal("error", "error", err)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal("error", "error", err)
		}

		rl, err := request.RequestFromReader(conn)

		if err != nil {
			log.Fatal("error", "error", err)
		}

		fmt.Printf("Method: %s \n Target: %s \n HttpVersion: %s \n", rl.RequestLine.Method, rl.RequestLine.RequestTarget, rl.RequestLine.HttpVersion)

	}

}

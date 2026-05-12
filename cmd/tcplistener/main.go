package main

import (
	"fmt"
	"log"
	"net"

	"github.com/ibrahimtahadivrik/httpfromtcp/internal/request"
)

func main() {
	listener, err := net.Listen("tcp", ":42069")
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()

	conn, err := listener.Accept()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("New connection from", conn.RemoteAddr())

	request, err := request.RequestFromReader(conn)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Request line:")
	fmt.Printf("- Method: %s\n", request.RequestLine.Method)
	fmt.Printf("- Target: %s\n", request.RequestLine.RequestTarget)
	fmt.Printf("- Version: %s\n", request.RequestLine.HttpVersion)

	log.Println("Closing connection from", conn.RemoteAddr())

}

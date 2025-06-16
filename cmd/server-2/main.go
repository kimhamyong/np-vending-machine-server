package main

import (
	"fmt"
	"net"
)

func main() {
	listener, err := net.Listen("tcp", ":9102")
	if err != nil {
		panic(err)
	}
	fmt.Println("Server-2 started!")

	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		go handleConn(conn)
	}
}

func handleConn(conn net.Conn) {
	defer conn.Close()

	buf := make([]byte, 1024)
	n, _ := conn.Read(buf)
	fmt.Println("Server-2 received:", string(buf[:n]))

	conn.Write([]byte("pong from server-2\n"))
}

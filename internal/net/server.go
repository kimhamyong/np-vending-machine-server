package net

import (
	"fmt"
	"net"
)

func StartServer(name string, port int) error {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return fmt.Errorf("[%s] 리스닝 실패: %w", name, err)
	}
	fmt.Printf("[%s] started! Listening on :%d\n", name, port)

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("[%s] 연결 실패: %v\n", name, err)
			continue
		}
		go handleConn(conn, name)
	}
}

func handleConn(conn net.Conn, name string) {
	defer conn.Close()

	buf := make([]byte, 1024)
	n, _ := conn.Read(buf)
	fmt.Printf("[%s] received: %s\n", name, string(buf[:n]))

	conn.Write([]byte(fmt.Sprintf("pong from %s\n", name)))
}

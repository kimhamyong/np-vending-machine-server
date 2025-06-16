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
	fmt.Printf("[%s] 수신: %s\n", name, string(buf[:n]))

	conn.Write([]byte(fmt.Sprintf("pong from %s\n", name)))

	// 여기서 SYNC 요청 발생시킴
	go func() {
		sender, err := InitRawSocketSender(oppositeIP(name)) // server-1 → server-2 or vice versa
		if err != nil {
			fmt.Println("raw socket 연결 실패:", err)
			return
		}
		err = sender.SendCommand(SyncCommand{Type: 0x01})
		if err != nil {
			fmt.Println("SYNC 요청 실패:", err)
		} else {
			fmt.Println("[SYNC] 요청 전송 완료")
		}
	}()
}

func oppositeIP(name string) string {
	if name == "server-1" {
		return "127.0.0.2" // server-2 주소
	}
	return "127.0.0.1" // server-1 주소
}

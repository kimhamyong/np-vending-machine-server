package main

import (
	"fmt"
	"net"
	"time"
	"sync"

	mynet "vending-system/internal/net"
)

var (
	statusMap map[string]*mynet.ServerStatus
	statusMu  sync.Mutex
	activeMap map[string]bool // 백업 리스닝 여부
)

func main() {
	fmt.Println("Backup Server started!")

	servers := []string{"server-1:9101", "server-2:9102"}
	statusMap = mynet.HealthCheck(servers, 5*time.Second)
	activeMap = make(map[string]bool)

	for {
		for addr, s := range statusMap {
			statusMu.Lock()
			if !s.Alive && !activeMap[addr] {
				fmt.Printf("[Backup] 감지: %s 다운됨, 백업 리스닝 시작\n", addr)
				go startBackupServer(addr)
				activeMap[addr] = true
			}
			statusMu.Unlock()
		}
		time.Sleep(5 * time.Second)
	}
}

// 백업 서버가 다운된 포트에 리스닝 시작
func startBackupServer(addr string) {
	listener, err := net.Listen("tcp", extractPort(addr))
	if err != nil {
		fmt.Printf("[Backup] 포트 리스닝 실패: %v\n", err)
		return
	}
	defer listener.Close()

	fmt.Printf("[Backup] 포트 %s 리스닝 시작\n", addr)

	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		go handleClient(conn, addr)
	}
}

func handleClient(conn net.Conn, from string) {
	defer conn.Close()

	buf := make([]byte, 1024)
	n, _ := conn.Read(buf)
	fmt.Printf("[Backup] %s → 메시지 수신: %s\n", from, string(buf[:n]))

	conn.Write([]byte(fmt.Sprintf("pong from backup (replacing %s)\n", from)))
}

func extractPort(addr string) string {
	_, port, _ := net.SplitHostPort(addr)
	return ":" + port
}

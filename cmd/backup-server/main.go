package main

import (
	"fmt"
	"log"
	"net"
	"sync"
	"time"

	"vending-system/internal/network"

	"vending-system/internal/storage"
)

func main() {
	db := storage.InitDB("/app/data/db.sqlite3", "/app/schema.sql")
    defer db.Close()

	fmt.Println("Backup Server started!")

	servers := []string{"server-1:9101", "server-2:9102"}
	statusMap := network.HealthCheck(servers, 10*time.Second)

	var mu sync.Mutex
	activeListeners := make(map[string]net.Listener)

	for {
		for addr, status := range statusMap {
			if !status.Alive {
				port := targetPort(addr)

				mu.Lock()
				if _, ok := activeListeners[port]; !ok {
					go startBackupListener(port, addr, activeListeners, &mu)
				}
				mu.Unlock()
			}
		}
		time.Sleep(10 * time.Second)
	}
}

func startBackupListener(port string, originAddr string, listeners map[string]net.Listener, mu *sync.Mutex) {
	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Printf("[Backup] 포트 %s 리스닝 실패: %v", port, err)
		return
	}

	mu.Lock()
	listeners[port] = listener
	mu.Unlock()

	fmt.Printf("[Backup] %s 다운됨 → 포트 %s 리스닝 시작\n", originAddr, port)

	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		go handleBackupConn(conn, originAddr)
	}
}

func handleBackupConn(conn net.Conn, origin string) {
	defer conn.Close()

	buf := make([]byte, 1024)
	n, _ := conn.Read(buf)
	fmt.Printf("[Backup] (%s 대신) 수신: %s\n", origin, string(buf[:n]))

	conn.Write([]byte(fmt.Sprintf("pong from backup (replacing %s)\n", origin)))
}

func targetPort(addr string) string {
	_, port, _ := net.SplitHostPort(addr)
	return port
}

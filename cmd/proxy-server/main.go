package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"sync"
	"time"

	mynet "vending-system/internal/net"
)

var (
	servers   = []string{"server-1:9101", "server-2:9102"}
	statusMap map[string]*mynet.ServerStatus
	statusMu  sync.Mutex
	rrIndex   int // 라운드로빈 인덱스
)

func main() {
	statusMap = mynet.HealthCheck(servers, 5*time.Second)

	listener, err := net.Listen("tcp", ":9000")
	if err != nil {
		log.Fatalf("프록시 서버 리스닝 실패: %v", err)
	}
	defer listener.Close()

	fmt.Println("Proxy Server started! Listening on :9000")

	for {
		clientConn, err := listener.Accept()
		if err != nil {
			continue
		}
		go handleClient(clientConn)
	}
}

func handleClient(clientConn net.Conn) {
	defer clientConn.Close()

	target := selectHealthyServer()
	if target == "" {
		fmt.Println("[Proxy] No healthy server available")
		return
	}

	serverConn, err := net.Dial("tcp", target)
	if err != nil {
		fmt.Printf("[Proxy] 서버 연결 실패: %v\n", err)
		return
	}
	defer serverConn.Close()

	fmt.Printf("[Proxy] 연결: 클라이언트 → %s\n", target)
	go io.Copy(serverConn, clientConn)
	io.Copy(clientConn, serverConn)
}

func selectHealthyServer() string {
	statusMu.Lock()
	defer statusMu.Unlock()

	for i := 0; i < len(servers); i++ {
		idx := (rrIndex + i) % len(servers)
		server := servers[idx]
		if statusMap[server].Alive {
			rrIndex = (idx + 1) % len(servers)
			return server
		}
	}
	return ""
}

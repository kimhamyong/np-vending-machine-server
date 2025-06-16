package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"sync"

	mynet "vending-system/internal/net"
)

var (
	// 헬스체크로부터 상태를 가져올 수 있도록 공유 map
	statusMap map[string]*mynet.ServerStatus
	statusMu  sync.Mutex
)

func main() {
	// 서버 주소 목록 (HealthCheck 대상)
	servers := []string{"server-1:9101", "server-2:9102"}

	// 헬스체크 시작 (주기: 5초)
	statusMap = mynet.HealthCheck(servers, 5_000_000_000)

	// 프록시 TCP 리스너 시작
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

// 클라이언트 요청을 처리하여 정상 서버로 릴레이
func handleClient(clientConn net.Conn) {
	defer clientConn.Close()

	// 정상 서버 선택
	target := selectHealthyServer()
	if target == "" {
		fmt.Println("[Proxy] No healthy backend server available")
		return
	}

	// 서버 연결
	serverConn, err := net.Dial("tcp", target)
	if err != nil {
		fmt.Printf("[Proxy] 서버 연결 실패: %v\n", err)
		return
	}
	defer serverConn.Close()

	fmt.Printf("[Proxy] 연결: 클라이언트 → %s\n", target)

	// 양방향 릴레이 (client <-> server)
	go io.Copy(serverConn, clientConn)
	io.Copy(clientConn, serverConn)
}

// 헬스체크 상태 기반으로 서버 선택 (첫 번째 alive 서버)
func selectHealthyServer() string {
	statusMu.Lock()
	defer statusMu.Unlock()

	for addr, s := range statusMap {
		if s.Alive {
			return addr
		}
	}
	return ""
}

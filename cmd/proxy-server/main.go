package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
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
	statusMap = mynet.HealthCheck(servers, 10*time.Second)

	// TCP (9000)
	go startTCPProxy()

	// HTTP (8000)
	startHTTPProxy()
}

func startTCPProxy() {
	listener, err := net.Listen("tcp", ":9000")
	if err != nil {
		log.Fatalf("프록시 TCP 리스닝 실패: %v", err)
	}
	defer listener.Close()

	fmt.Println("Proxy TCP Server started! Listening on :9000")

	for {
		clientConn, err := listener.Accept()
		if err != nil {
			continue
		}
		go forwardTCP(clientConn)
	}
}

func forwardTCP(clientConn net.Conn) {
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

	go io.Copy(serverConn, clientConn)
	io.Copy(clientConn, serverConn)
}

func startHTTPProxy() {
	http.HandleFunc("/api/user_login", httpToTCPHandler("user_login"))
	http.HandleFunc("/api/user_signup", httpToTCPHandler("user_signup"))
	http.HandleFunc("/api/user_change_password", httpToTCPHandler("user_change_password"))
	http.HandleFunc("/api/user_delete_account", httpToTCPHandler("user_delete_account"))

	fmt.Println("Proxy HTTP Server started! Listening on :8000")
	log.Fatal(http.ListenAndServe(":8000", nil))
}

func httpToTCPHandler(action string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "요청 본문 읽기 실패", http.StatusBadRequest)
			return
		}

		var req map[string]interface{}
		if err := json.Unmarshal(body, &req); err != nil {
			http.Error(w, "JSON 파싱 실패", http.StatusBadRequest)
			return
		}
		req["action"] = action

		target := selectHealthyServer()
		if target == "" {
			http.Error(w, "사용 가능한 서버 없음", http.StatusServiceUnavailable)
			return
		}

		conn, err := net.Dial("tcp", target)
		if err != nil {
			http.Error(w, "서버 연결 실패", http.StatusInternalServerError)
			return
		}
		defer conn.Close()

		reqBytes, _ := json.Marshal(req)
		conn.Write(reqBytes)

		respBuf := make([]byte, 1024)
		n, _ := conn.Read(respBuf)

		w.Header().Set("Content-Type", "application/json")
		w.Write(respBuf[:n])
	}
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

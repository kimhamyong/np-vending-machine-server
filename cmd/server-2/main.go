package main

import (
	"fmt"
	"log"
	"net"
	"vending-system/internal/storage"
	"encoding/json"
	"vending-system/internal/handler/user"
	"vending-system/internal/network"
	"database/sql"
)

var syncRequested = false // 동기화 요청을 보냈는지 확인하는 플래그

func main() {
	// DB 초기화
	db := storage.InitDB("/app/data/db.sqlite3", "/app/schema.sql")
	defer db.Close()

	// 1. 서버 실행 시 한 번만 동기화 요청
	err := sendSyncRequest("server-1") // 최초 동기화 요청
	if err != nil {
		log.Fatal("SYNC 요청 실패:", err)
	}

	// 2. 서버 실행
	err = startServerWithDB("server-2", 9102, db)
	if err != nil {
		log.Fatal(err)
	}
}

// 서버 실행 시 한 번만 동기화 요청을 보내는 함수
func sendSyncRequest(name string) error {
	// 이미 동기화 요청을 보냈다면, 다시 보내지 않도록 함
	if syncRequested {
		return nil
	}

	// 동기화 요청 전송
	sender, err := network.InitRawSocketSender("127.0.0.1") // server-1 주소
	if err != nil {
		return fmt.Errorf("raw socket 연결 실패: %v", err)
	}

	// 동기화 요청을 한 번만 보내기
	err = sender.SendCommand(network.SyncCommand{Type: 0x01})
	if err != nil {
		return fmt.Errorf("SYNC 요청 전송 실패: %v", err)
	}

	log.Println("[SYNC] 요청 전송 완료")
	syncRequested = true // 동기화 요청을 보낸 상태로 설정
	return nil
}

// 서버 실행 함수
func startServerWithDB(name string, port int, db *sql.DB) error {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return fmt.Errorf("[%s] 리스닝 실패: %w", name, err)
	}
	fmt.Printf("[%s] started! Listening on :%d\n", name, port)

	// 클라이언트 연결 대기
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("[%s] 연결 실패: %v\n", name, err)
			continue
		}
		go handleConnWithDB(conn, name, db)
	}
}

// 클라이언트 연결을 처리하는 함수
func handleConnWithDB(conn net.Conn, name string, db *sql.DB) {
	defer conn.Close()

	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		log.Printf("[%s] 수신 오류: %v\n", name, err)
		return
	}

	var req map[string]interface{}
	if err := json.Unmarshal(buf[:n], &req); err != nil {
		conn.Write([]byte(`{"success": false, "error": "Invalid JSON"}`))
		return
	}

	action, ok := req["action"].(string)
	if !ok {
		conn.Write([]byte(`{"success": false, "error": "Missing action"}`))
		return
	}

	var response []byte

	// 데이터베이스 변경 작업 수행
	switch action {
	case "user_signup":
		response = user.HandleSignup(req, db)  // db만 전달
	case "user_login":
		response = user.HandleLogin(req, db)  // db만 전달
	case "user_change_password":
		response = user.HandleChangePassword(req, db)  // db만 전달
	case "user_delete_account":
		response = user.HandleDeleteAccount(req, db)  // db만 전달
	default:
		response = []byte(`{"success": false, "error": "Unknown action"}`)
	}

	conn.Write(response)
}

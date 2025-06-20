package main

import (
	"log"
	"time"
	"net"
	"fmt"

	"vending-system/internal/storage"
	"encoding/json"
	"vending-system/internal/handler/user"
	"vending-system/internal/network"
	"database/sql"
)

func main() {
	err := network.StartRawListener(func(cmd network.SyncCommand) {
		log.Printf("[SYNC] 수신된 명령: %+v", cmd)
	})
	if err != nil {
		log.Fatal("raw listener 시작 실패:", err)
	}

	// 서버 실행 시 한 번만 동기화 요청
	err = sendSyncRequest("server-1") // 최초 동기화 요청
	if err != nil {
		log.Fatal("SYNC 요청 실패:", err)
	}

	// DB 초기화
	db := storage.InitDB("/app/data/db.sqlite3", "/app/schema.sql")
	defer db.Close()

	err = startServerWithDB("server-2", 9102, db)
	if err != nil {
		log.Fatal(err)
	}
}

// 동기화 요청을 보내는 함수
func sendSyncRequest(name string) error {

	// 동기화 요청 전송
	time.Sleep(1 * time.Second)
	sender, err := network.InitRawSocketSender(oppositeIP(name)) // server-1 → server-2 or vice versa
	if err != nil {
		fmt.Println("raw socket 연결 실패:", err)
		return err
	}
	err = sender.SendCommand(network.SyncCommand{Type: 0x01})
	if err != nil {
		fmt.Println("SYNC 요청 실패:", err)
		} else {
			fmt.Println("[SYNC] 요청 전송 완료")
		}

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
		if string(response) == `{"success":true}` {
    		fmt.Println("Sending sync request to server-1") // 동기화 요청 전송 로그 확인
    		go sendSyncRequest("server-1") // 서버1에 동기화 요청
		}
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

func oppositeIP(name string) string {
	if name == "server-1" {
		return "127.0.0.2" // server-2 주소
	}
	return "127.0.0.1" // server-1 주소
}
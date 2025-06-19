package net

import (
	"encoding/json"
	"fmt"
	"net"
	"vending-system/internal/handler/user"
	"database/sql"
)

func StartServerWithDB(name string, port int, db *sql.DB) error {
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
		go handleConnWithDB(conn, name, db)
	}
}

func handleConnWithDB(conn net.Conn, name string, db *sql.DB) {
	defer conn.Close()

	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		fmt.Printf("[%s] 수신 오류: %v\n", name, err)
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

	switch action {
	case "user_signup":
		response = user.HandleSignup(req, db) 
	case "user_login":
		response = user.HandleLogin(req, db)  
	case "user_change_password":
		response = user.HandleChangePassword(req, db) 
	case "user_delete_account":
		response = user.HandleDeleteAccount(req, db)
	default:
		response = []byte(`{"success": false, "error": "Unknown action"}`)
	}

	conn.Write(response)

	// SYNC 요청은 여전히 수행
	go func() {
		sender, err := InitRawSocketSender(oppositeIP(name))
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
		return "127.0.0.2"
	}
	return "127.0.0.1"
}

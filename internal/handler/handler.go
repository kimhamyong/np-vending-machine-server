package handler

import (
	"database/sql"
	"encoding/json"
	"log"
	"net"

	"vending-system/internal/network"
	"vending-system/internal/handler/user"
)

func HandleConn(conn net.Conn, name string, db *sql.DB) {
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

	switch action {
	case "user_signup":
		response = user.HandleSignup(req, db)
		if string(response) == `{"success":true}` {
			log.Println("Sending sync request to server-2")
			go network.SendSyncRequest("server-2")
		}
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
}

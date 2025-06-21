package main

import (
	"database/sql"
	"fmt"
	"log"
	"net"

	"vending-system/internal/handler"
	"vending-system/internal/network"
	"vending-system/internal/storage"
)

func main() {
	// raw listener 시작
	err := network.StartRawListener(func(cmd network.SyncCommand) {
		log.Printf("[SYNC] 수신된 명령: %+v", cmd)
	})
	if err != nil {
		log.Fatal("raw listener 시작 실패:", err)
	}

	// 최초 동기화 요청
	err = network.SendSyncRequest("server-2")
	if err != nil {
		log.Fatal("SYNC 요청 실패:", err)
	}

	// DB 초기화
	db := storage.InitDB(
		"data/shared-database/db.sqlite3",
		"internal/storage/schema.sql",
	)
	defer db.Close()

	// TCP 서버 시작
	err = startServerWithDB("server-1", 9101, db)
	if err != nil {
		log.Fatal(err)
	}
}

// TCP 서버 실행 및 클라이언트 연결 수신
func startServerWithDB(name string, port int, db *sql.DB) error {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return fmt.Errorf("[%s] 리스닝 실패: %w", name, err)
	}
	log.Printf("[%s] started! Listening on :%d\n", name, port)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("[%s] 연결 실패: %v\n", name, err)
			continue
		}
		go handler.HandleConn(conn, name, db)
	}
}

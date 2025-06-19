package main

import (
	"log"
	"time"
	"vending-system/internal/net"
    "vending-system/internal/storage"
)

func main() {
    db := storage.InitDB("/app/data/db.sqlite3", "/app/schema.sql")
    defer db.Close()

	// 1. raw listener 먼저 시작 (항상 백그라운드 수신)
	err := net.StartRawListener(func(cmd net.SyncCommand) {
		log.Printf("[SYNC] 수신된 명령: %+v", cmd)
		// TODO: cmd.Type == 0x01 (SYNC 요청) 처리 → store 응답 전송
		// TODO: cmd.Type == 0x02/0x03 → store.Apply(cmd)
	})
	if err != nil {
		log.Fatal("raw listener 시작 실패:", err)
	}

	// 2. peer 연결 및 초기 SYNC 요청 전송
	go func() {
		time.Sleep(1 * time.Second) // listener가 먼저 열릴 수 있도록 딜레이
		sender, err := net.InitRawSocketSender("127.0.0.2") // server-2 주소
		if err != nil {
			log.Fatal("raw socket 연결 실패:", err)
		}

		err = sender.SendCommand(net.SyncCommand{
			Type: 0x01, // SYNC 요청
		})
		if err != nil {
			log.Println("SYNC 요청 실패:", err)
		}
	}()

	// 3. TCP 클라이언트 처리용 서버 시작
	if err := net.StartServer("server-1", 9101); err != nil {
		log.Fatal(err)
	}
}

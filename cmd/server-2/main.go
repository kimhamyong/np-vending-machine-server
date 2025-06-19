package main

import (
	"log"
	"time"

	"vending-system/internal/net"
	"vending-system/internal/storage"
)

func main() {
	// DB 초기화
	db := storage.InitDB("/app/data/db.sqlite3", "/app/schema.sql")
	defer db.Close()

	// 1. raw listener: server-2와의 동기화 수신
	err := net.StartRawListener(func(cmd net.SyncCommand) {
		log.Printf("[SYNC] 수신된 명령: %+v", cmd)
		// TODO: 향후 store.Apply(cmd) 처리 필요
	})
	if err != nil {
		log.Fatal("raw listener 시작 실패:", err)
	}

	// 2. server-2로 SYNC 요청 전송
	go func() {
		time.Sleep(1 * time.Second)
		sender, err := net.InitRawSocketSender("127.0.0.1") // server-1 주소
		if err != nil {
			log.Fatal("raw socket 연결 실패:", err)
		}
		err = sender.SendCommand(net.SyncCommand{Type: 0x01})
		if err != nil {
			log.Println("SYNC 요청 실패:", err)
		}
	}()

	// 3. TCP 서버 시작 (핸들러에 db 넘겨줌)
	if err := net.StartServerWithDB("server-2", 9102, db); err != nil {
		log.Fatal(err)
	}
}

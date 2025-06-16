package main

import (
	"log"
	"time"

	"vending-system/internal/net"
)

func main() {
	err := net.StartRawListener(func(cmd net.SyncCommand) {
		log.Printf("[SYNC] 수신된 명령: %+v", cmd)
	})
	if err != nil {
		log.Fatal("raw listener 시작 실패:", err)
	}

	go func() {
		time.Sleep(1 * time.Second)
		sender, err := net.InitRawSocketSender("127.0.0.1") // server-1 주소
		if err != nil {
			log.Fatal("raw socket 연결 실패:", err)
		}

		err = sender.SendCommand(net.SyncCommand{
			Type: 0x01,
		})
		if err != nil {
			log.Println("SYNC 요청 실패:", err)
		}
	}()

	if err := net.StartServer("server-2", 9102); err != nil {
		log.Fatal(err)
	}
}

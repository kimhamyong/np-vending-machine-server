package network

import (
	"fmt"
)

// 동기화 명령 전송
func SendSyncRequest(name string) error {
	sender, err := InitRawSocketSender(oppositeIP(name))
	if err != nil {
		fmt.Println("raw socket 연결 실패:", err)
		return err
	}
	err = sender.SendCommand(SyncCommand{Type: 0x01})
	if err != nil {
		fmt.Println("SYNC 요청 실패:", err)
	} else {
		fmt.Println("[SYNC] 요청 전송 완료")
	}
	return nil
}

func oppositeIP(name string) string {
	if name == "server-1" {
		return "127.0.0.2"
	}
	return "127.0.0.1"
}

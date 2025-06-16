package main

import (
	"fmt"
	"log"
	"time"

	mynet "vending-system/internal/net"
)

func main() {
	fmt.Println("Backup Server started!")

	// 감시 대상 서버 리스트
	servers := []string{"server-1:9101", "server-2:9102"}

	// 헬스체크 주기 5초로 상태 맵 생성
	statusMap := mynet.HealthCheck(servers, 5*time.Second)

	for {
		for addr, status := range statusMap {
			if !status.Alive {
				log.Printf("[Backup] 감지: %s 다운됨\n", addr)

				// ➤ 향후 동기화/대체 처리
				err := mynet.HandleServerFailure(addr)
				if err != nil {
					log.Printf("[Backup] 처리 실패 (%s): %v\n", addr, err)
				}
			}
		}
		time.Sleep(5 * time.Second)
	}
}

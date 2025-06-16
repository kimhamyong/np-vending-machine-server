package net

import (
	"fmt"
	"net"
	"sync"
	"time"
)

type ServerStatus struct {
	Address   string
	Alive     bool
	LastError error
	Mutex     sync.RWMutex
}

// 서버들의 상태를 주기적으로 확인하고 상태맵을 반환
func HealthCheck(servers []string, interval time.Duration) map[string]*ServerStatus {
	statusMap := make(map[string]*ServerStatus)

	for _, addr := range servers {
		statusMap[addr] = &ServerStatus{Address: addr}
	}

	go func() {
		for {
			for _, addr := range servers {
				status := statusMap[addr]
				conn, err := net.DialTimeout("tcp", addr, 1*time.Second)

				status.Mutex.Lock()
				if err != nil {
					status.Alive = false
					status.LastError = err
					fmt.Printf("[HealthCheck] %s is DOWN\n", addr)
				} else {
					status.Alive = true
					status.LastError = nil
					fmt.Printf("[HealthCheck] %s is UP\n", addr)
					conn.Close()
				}
				status.Mutex.Unlock()
			}
			time.Sleep(interval)
		}
	}()

	return statusMap
}

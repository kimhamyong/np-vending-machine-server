package main

import (
	"log"
	"vending-system/internal/net"
)

func main() {
	err := net.StartServer("server-1", 9101) // server-2는 9102, backup은 9103
	if err != nil {
		log.Fatal(err)
	}
}
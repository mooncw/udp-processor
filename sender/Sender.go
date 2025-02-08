package main

import (
	"fmt"
	"net"
	"sync"
)

func sendMessage(conn net.Conn, message string, wg *sync.WaitGroup) {
	defer wg.Done() // 고루틴 완료 후 WaitGroup의 카운트를 감소시킨다.

	_, err := conn.Write([]byte(message))
	if err != nil {
		fmt.Println("Error sending:", err)
		return
	}
	fmt.Println("Sent:", message)
}

func main() {
	serverAddr := "127.0.0.1:9000" // 수신 서버 주소
	conn, err := net.Dial("udp", serverAddr)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	var wg sync.WaitGroup // 고루틴 동기화를 위한 WaitGroup

	for i := 0; i < 200; i++ {
		message := fmt.Sprintf("Message %d", i)
		wg.Add(1)                          // 고루틴 추가
		go sendMessage(conn, message, &wg) // 고루틴으로 메시지 전송
	}

	// 모든 고루틴이 끝날 때까지 기다림
	wg.Wait()
}

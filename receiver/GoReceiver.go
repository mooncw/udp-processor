package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

var messageCount int
var mu sync.Mutex

func handleRequest(conn *net.UDPConn, remoteAddr *net.UDPAddr, buffer []byte, wg *sync.WaitGroup) {
	defer wg.Done()

	// 데이터 처리
	fmt.Printf("Received from %s: %s\n", remoteAddr, string(buffer))

	// 메시지 카운트 증가
	mu.Lock()
	messageCount++
	mu.Unlock()

	// 예시: 클라이언트로 응답 전송
	_, err := conn.WriteToUDP([]byte("Acknowledged"), remoteAddr)
	if err != nil {
		log.Printf("Error sending acknowledgment to %s: %v\n", remoteAddr, err)
	}
}

func main() {
	addr, err := net.ResolveUDPAddr("udp", ":9000")
	if err != nil {
		log.Fatalf("Error resolving address: %v\n", err)
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		log.Fatalf("Error listening on UDP: %v\n", err)
	}
	defer conn.Close()

	buffer := make([]byte, 1024) // 1KB 버퍼
	fmt.Println("Listening on UDP port 9000...")

	var wg sync.WaitGroup

	// 종료 시그널을 처리하는 채널
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// 고루틴에서 수신 대기
	go func() {
		for {
			n, remoteAddr, err := conn.ReadFromUDP(buffer)
			if err != nil {
				log.Printf("Error receiving data: %v\n", err)
				continue
			}

			// 새로운 슬라이스를 생성하여 데이터 복사
			data := make([]byte, n)
			copy(data, buffer[:n])

			// 요청을 처리하는 별도 고루틴 호출
			wg.Add(1)
			go handleRequest(conn, remoteAddr, data, &wg)
		}
	}()

	// 종료 시그널을 기다림
	<-sigChan
	fmt.Println("Server is shutting down...")

	// 모든 고루틴이 끝날 때까지 기다리기
	wg.Wait()

	// 최종 메시지 수 출력
	mu.Lock()
	fmt.Printf("Total messages received: %d\n", messageCount)
	mu.Unlock()
}

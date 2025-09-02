package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

func main() {
	// TCP 소켓을 8080 포트에서 열기
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal("포트 8080에서 리스닝 실패:", err)
	}
	defer listener.Close()

	fmt.Println("서버가 8080 포트에서 시작되었습니다...")
	fmt.Println("브라우저에서 http://localhost:8080 으로 접속해보세요!")

	for {
		// 클라이언트 연결을 기다림
		conn, err := listener.Accept()
		if err != nil {
			log.Println("연결 수락 실패:", err)
			continue
		}

		fmt.Println("\n=== 새로운 클라이언트 연결:", conn.RemoteAddr(), "===")

		// HTTP 요청 읽고 응답하기
		handleHTTPRequest(conn)

		// 연결 종료
		conn.Close()
		fmt.Println("=== 연결 종료 ===")
	}
}

func handleHTTPRequest(conn net.Conn) {
	// 1. HTTP 요청 읽기
	scanner := bufio.NewScanner(conn)

	fmt.Println("📨 받은 HTTP 요청:")

	lineCount := 0

	for scanner.Scan() {
		line := scanner.Text()

		// 첫 번째 줄 출력 (GET / HTTP/1.1)
		if lineCount == 0 {
			fmt.Printf("🎯 %s\n", line)
		}

		// 빈 줄이 나오면 헤더 끝
		if len(strings.TrimSpace(line)) == 0 {
			break
		}

		lineCount++
	}

	// 2. HTML 파일 읽기
	fmt.Println("📖 HTML 파일 읽는 중...")

	htmlContent, err := os.ReadFile("static/index.html")
	if err != nil {
		// 파일을 못 읽으면 기본 에러 페이지
		errorHTML := `<!DOCTYPE html>
<html><body>
<h1>500 Internal Server Error</h1>
<p>static/index.html 파일을 찾을 수 없습니다.</p>
<p>static 폴더를 만들고 index.html 파일을 넣어주세요.</p>
</body></html>`
		htmlContent = []byte(errorHTML)
		fmt.Println("❌ 파일 읽기 실패:", err)
	} else {
		fmt.Println("✅ HTML 파일 읽기 성공!")
	}

	// 3. HTTP 응답 만들기
	fmt.Println("📤 HTTP 응답 보내는 중...")

	// HTTP 응답 헤더 작성
	response := fmt.Sprintf("HTTP/1.1 200 OK\r\n"+
		"Content-Type: text/html; charset=utf-8\r\n"+
		"Content-Length: %d\r\n"+
		"Connection: close\r\n"+
		"\r\n"+
		"%s", len(htmlContent), string(htmlContent))

	// 응답 보내기
	_, err = conn.Write([]byte(response))
	if err != nil {
		fmt.Println("응답 전송 실패:", err)
		return
	}

	fmt.Println("✅ HTTP 응답 전송 완료!")
}

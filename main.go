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
	listener, err := net.Listen("tcp", ":80")
	if err != nil {
		log.Fatal("포트 80에서 리스닝 실패:", err)
	}
	defer listener.Close()

	fmt.Println("서버가 80 포트에서 시작되었습니다...")
	fmt.Println("브라우저에서 http://localhost 으로 접속해보세요!")

	// for {
	// 	// 클라이언트 연결을 기다림
	// 	conn, err := listener.Accept()
	// 	if err != nil {
	// 		log.Println("연결 수락 실패:", err)
	// 		continue
	// 	}

	// 	fmt.Println("\n=== 새로운 클라이언트 연결:", conn.RemoteAddr(), "===")

	// 	// HTTP 요청 읽고 응답하기
	// 	handleHTTPRequest(conn)

	// 	// 연결 종료
	// 	conn.Close()
	// 	fmt.Println("=== 연결 종료 ===")
	// }
	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}

		// 새로운 고루틴에서 처리 (비동기)
		go func(c net.Conn) {
			fmt.Println("\n=== 새로운 클라이언트 연결:", conn.RemoteAddr(), "===")
			handleHTTPRequest(c)
			c.Close()
			fmt.Println("=== 연결 종료 ===")
		}(conn)
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

			// 새로 추가: 경로 추출
			parts := strings.Split(line, " ")
			if len(parts) >= 2 {
				requestPath := parts[1]
				fmt.Printf("📁 요청 경로: %s\n", requestPath)

				if requestPath == "/" {
					// 기존 HTML 처리 (그대로 유지)
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

				} else if requestPath == "/css/style.css" {
					// 새로 추가: CSS 처리
					fmt.Println("🎨 CSS 파일 읽는 중...")
					cssContent, err := os.ReadFile("static/css/style.css")
					if err != nil {
						fmt.Println("❌ CSS 파일 읽기 실패:", err)
						// 404 에러 응답 보내기
						return
					}

					fmt.Println("✅ CSS 파일 읽기 성공!")
					fmt.Println("📤 CSS 응답 보내는 중...")

					// CSS 응답 생성
					response := fmt.Sprintf("HTTP/1.1 200 OK\r\n"+
						"Content-Type: text/css\r\n"+
						"Content-Length: %d\r\n"+
						"Connection: close\r\n"+
						"\r\n"+
						"%s", len(cssContent), string(cssContent))

					conn.Write([]byte(response))
					fmt.Println("✅ CSS 응답 전송 완료!")
				} else if requestPath == "/js/script.js" {
					// 새로 추가: CSS 처리
					fmt.Println("🎨 JS 파일 읽는 중...")
					cssContent, err := os.ReadFile("static/js/script.js")
					if err != nil {
						fmt.Println("❌ JS 파일 읽기 실패:", err)
						// 404 에러 응답 보내기
						return
					}

					fmt.Println("✅ JS 파일 읽기 성공!")
					fmt.Println("📤 JS 응답 보내는 중...")

					// CSS 응답 생성
					response := fmt.Sprintf("HTTP/1.1 200 OK\r\n"+
						"Content-Type: application/javascript\r\n"+
						"Content-Length: %d\r\n"+
						"Connection: close\r\n"+
						"\r\n"+
						"%s", len(cssContent), string(cssContent))

					conn.Write([]byte(response))
					fmt.Println("✅ JS 응답 전송 완료!")
				}
			}
		}

		// 빈 줄이 나오면 헤더 끝
		if len(strings.TrimSpace(line)) == 0 {
			break
		}

		lineCount++
	}

	// 2. HTML 파일 읽기
	// fmt.Println("📖 HTML 파일 읽는 중...")

	// htmlContent, err := os.ReadFile("static/index.html")

	// 	if err != nil {
	// 		// 파일을 못 읽으면 기본 에러 페이지
	// 		errorHTML := `<!DOCTYPE html>
	// <html><body>
	// <h1>500 Internal Server Error</h1>
	// <p>static/index.html 파일을 찾을 수 없습니다.</p>
	// <p>static 폴더를 만들고 index.html 파일을 넣어주세요.</p>
	// </body></html>`
	// 		htmlContent = []byte(errorHTML)
	// 		fmt.Println("❌ 파일 읽기 실패:", err)
	// 	} else {
	// 		fmt.Println("✅ HTML 파일 읽기 성공!")
	// 	}

	// 	// 3. HTTP 응답 만들기
	// 	fmt.Println("📤 HTTP 응답 보내는 중...")

	// 	// HTTP 응답 헤더 작성
	// 	response := fmt.Sprintf("HTTP/1.1 200 OK\r\n"+
	// 		"Content-Type: text/html; charset=utf-8\r\n"+
	// 		"Content-Length: %d\r\n"+
	// 		"Connection: close\r\n"+
	// 		"\r\n"+
	// 		"%s", len(htmlContent), string(htmlContent))

	// 	// 응답 보내기
	// 	_, err = conn.Write([]byte(response))
	// 	if err != nil {
	// 		fmt.Println("응답 전송 실패:", err)
	// 		return
	// 	}

	fmt.Println("✅ HTTP 응답 전송 완료!")
}

# 🏠 Go 홈서버 프로젝트 - 완성 기록

## 📋 프로젝트 개요

**핵심 목표**: HTML 문서를 읽어서 80포트로 소켓을 열고, HTTP GET 요청에 응답하는 웹서버 구축  
**학습 목표**: 웹개발자로서 네트워크의 전체적인 흐름을 파악을 해보고자 Go로 웹서버 구현

---

## 🎯 Part 1: 로컬에서의 동작(25.09.02)

### ✅ 1단계: 기본 웹서버 구현

#### 핵심 구현 사항
- **TCP 소켓 서버**: 8080포트에서 클라이언트 연결 수락
- **HTTP 요청 파싱**: GET 요청 읽기 및 로깅
- **HTML 파일 서빙**: `static/index.html` 파일 읽어서 응답
- **HTTP 응답 생성**: 적절한 헤더와 상태코드(200 OK) 전송

#### 프로젝트 구조 (1단계)
```
├── main.go              # 메인 서버 코드 (8080포트)
├── static/
│   └── index.html       # 서빙할 HTML 파일
└── go.mod
```

#### 실행 및 테스트
```bash
go run main.go
# 브라우저에서 http://localhost:8080 접속
```

### 🔍 주요 학습 내용

#### HTTP 프로토콜 완전 이해
**요청 구조:**
```
GET / HTTP/1.1
Host: localhost:8080
User-Agent: Mozilla/5.0...
```

**응답 구조:**
```
HTTP/1.1 200 OK
Content-Type: text/html; charset=utf-8
Content-Length: 123
Connection: close

<html>...</html>
```

#### Express.js vs 직접 구현 비교 (1단계)
**Express.js (3줄):**
```javascript
app.get('/', (req, res) => {
  res.sendFile('index.html');
});
```

**Go 직접 구현 (50+ 줄):**
- TCP 소켓 수동 관리
- HTTP 요청 수동 파싱  
- HTTP 응답 헤더 수동 생성
- 바이트 레벨 데이터 전송

#### 와이어샤크를 통한 패킷 분석
- **TCP 3-Way Handshake**: SYN → SYN-ACK → ACK
- **HTTP 메시지 교환**: 요청 → 파싱 → 응답 생성 → 전송
- **연결 종료**: FIN-ACK 교환을 통한 정중한 연결 해제
![image](./images/wireshark.png)

### 💡 1단계 후기

> **"이 프로젝트를 시작한 계기가 된 영상에서 왜 웹은 정말 편리한 도구가 많다고 하는지 조금 체감할 수 있었다. 그나마 좀 익숙한 Express.js를 생각하며 작업을 해보니 더더욱 알 수 있었다.
그리고 HTTP가 얼마나 정교한 프로토콜인지도 깨닫게 되었다."**

---

## 🌐 Part 2: 80포트를 열어서 외부 접속 가능

### ✅ 2단계: 외부 접속 가능한 웹서버

#### 포트 변경 및 네트워크 설정
- **포트 변경**: 8080 → 80 (HTTP 표준 포트, 외부 접속을 위해서.)
- **공인 IP 확인**: 모뎀 직접 연결로 공인 IP 할당 확인(ipconfig로 확인해보니 192.168로 시작되지 않는 것을 확인. -> 포트포워딩 작업을 굳이 할 필요 없었다.)
- **방화벽 설정**: Windows 방화벽 80포트 인바운드 허용

#### 왜 하필 "80"번 포트일까?
- **1990년**: RFC 1060에서 포트 80번이 비어있었음
- **1991년**: 팀 버너스 리가 HTTP 0.9에서 80번을 HTTP 기본 포트로 정의
- **결론**: 딱히 큰 이유는 없고, 단지 "비어 있어서"..

#### 브라우저 동작 변화
**포트 8080 시절:**
```
http://localhost:8080  ← 포트 번호 필수
```

**포트 80 변경 후:**
```
http://localhost       ← 포트 번호 생략 가능
```

### 🔧 동시 접속 처리 구현

#### 문제 발견
- **증상**: 노트북과 폰 동시 접속 불가
- **원인**: 기존 코드에서 순차 처리 방식으로 인한 블로킹

#### Go 고루틴으로 해결
```go
for {
    conn, err := listener.Accept()
    if err != nil {
        continue
    }
    
    // 새로운 고루틴에서 처리 (비동기)
    go func(c net.Conn) {
        handleHTTPRequest(c)
        c.Close()
    }(conn)
}
```

#### Go 고루틴의 특별함
- **매우 가벼움**: 2KB 메모리만 사용 (스레드 대비 1000배 효율)
- **사용법 간단**: `go` 키워드 하나로 동시 처리
- **자동 관리**: Go 런타임이 스케줄링 자동 처리

### 🌍 인터넷 환경 체험

#### 자동화된 봇 트래픽 관찰
**현실적 발견**: 공인 IP 노출 즉시 자동화된 접속 시작
- **5분 이내**: 첫 번째 봇 방문
- **1시간 이내**: 3-5개 다른 IP에서 접속
- **대표 사례**: 65.49.x.x (AWS) 대역에서의 정기적 접속

#### "인터넷 노이즈"의 정체
**봇 트래픽 종류:**
- 검색엔진 크롤러 (Google, Bing)
- 보안 스캐너 (취약점 점검)
- 자동화된 포트 스캔
- 연구기관 모니터링

**깨진 HTTP 요청의 의미:**
- 구버전 HTTP/1.0 사용
- 불완전한 요청 패킷
- 바이너리 데이터 전송 시도

### 💡 2단계 후기
> **"해당 작업은 데스크탑에서 진행을 했는데, 당연히 인터넷이 공유기에 물린줄 알고 뻘짓을 좀 했다. 알고보니 모뎀에 직접 물린거라 포트포워딩 작업이 필요가 없었다. 공인IP와 사설IP에 대해 다시 리마인드 할 수 있는 좋은 경험 이었고, 와이어샤크로 외부에서 들어오는 네트워크의 패킷을 까보다가 이상한 IP와 깨진 HTTP 요청을 발견해서 좀 쫄았다. 그래서 해당 내용을 AI에게 물어보니 수많은 봇들이 24시간 돌아다니며 서버들을 탐색하고 있고, 이건 인터넷의 일상이니 걱정 말라고 했다. 이런걸 인터넷 노이즈 라고 부른다고 한다."**

---

## 📁 Part 3: 정적 파일 서빙 확장

### ✅ 3단계: 정적 파일 서빙 구현

#### 파일 분리 작업
**인라인 → 분리된 파일 구조:**
```
static/
├── index.html          # HTML만 (링크 태그 포함)
├── css/
│   └── style.css      # 분리된 CSS 파일
└── js/
    └── script.js      # 분리된 JavaScript 파일
```

#### 브라우저 자동 요청 메커니즘 발견
**핵심 원리**: HTML 파싱 중 특정 태그를 만나면 자동으로 추가 HTTP 요청
```html
<link rel="stylesheet" href="/css/style.css">  → GET /css/style.css 자동 요청
<script src="/js/script.js"></script>          → GET /js/script.js 자동 요청
```

**실제 관찰된 요청 순서:**
```
1번째: GET / → HTML 응답
2번째: GET /css/style.css → CSS 응답 (HTML 파싱 중 발견)
3번째: GET /js/script.js → JS 응답 (HTML 파싱 중 발견)
```

### 🔧 Go 서버 확장 구현

#### URL 경로 파싱
```go
parts := strings.Split(line, " ")
requestPath := parts[1]  // "/css/style.css" 추출
fmt.Printf("📁 요청 경로: %s\n", requestPath)
```

#### 경로별 분기 처리
```go
if requestPath == "/" {
    // HTML 파일 처리
    response = fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/html\r\n...")
} else if requestPath == "/css/style.css" {
    // CSS 파일 처리
    response = fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/css\r\n...")
} else if requestPath == "/js/script.js" {
    // JS 파일 처리
    response = fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: application/javascript\r\n...")
}
```

#### MIME 타입의 중요성
- **HTML**: `Content-Type: text/html`
- **CSS**: `Content-Type: text/css`
- **JavaScript**: `Content-Type: application/javascript`

### 📚 웹 개발 관례의 진짜 이유

#### CSS/JS 로딩 순서 이해
**CSS를 `<head>`에 두는 이유:**
- 렌더링 최적화 (FOUC 방지)
- HTML 파싱과 동시에 스타일 적용
- 깔끔한 화면 표시

**JavaScript를 `</body>` 직전에 두는 이유:**
- DOM 요소들이 먼저 생성되어야 함
- `document.getElementById()` 등이 정상 동작
- 스크립트 실행 시점에 모든 DOM 요소 존재 보장

#### 성능 최적화 발견 사항
**현재 비효율:**
- 매 요청마다 디스크에서 파일 읽기
- 파일 캐싱 미적용

**개선 방향:**
- 서버 시작 시 메모리에 파일 로드
- 브라우저 캐시 헤더 추가
- 압축 전송 (Gzip) 고려

### 💡 3단계 후기

> **"그렇게 공부했던 브라우저 렌더링이 한 번에 이해가 되었다. script 태그를 맨 마지막에 두라는 이유, CSS는 head에 두라는 이유 등 웹 개발 관례들이 모두 브라우저의 HTML 파싱 순서와 DOM 생성 과정에서 비롯된 합리적인 이유가 있었다. 이 프로젝트는 정말 잘한거 같다."**

---

## 🎉 전체 프로젝트 완성 소감

### 달성한 목표들
✅ **기본 목표**: HTML 문서를 80포트로 서빙하는 웹서버 완성  
✅ **확장 목표**: 외부 접속 가능 + 정적 파일 서빙까지 구현  
✅ **학습 목표**: 네트워크 전체 흐름 완전 이해  

### Before vs After
**이전 (프레임워크 의존):**
- "Express.js 쓰면 되는데?"
- "정적 파일? app.use(express.static()) 하면 되잖아"
- "웹서버? 그냥 작동하는 거 아니야?"

**이후 (원리 이해):**
- "브라우저가 HTML 파싱하면서 자동으로 CSS/JS 요청하는구나"
- "HTTP 상태코드, MIME 타입, 헤더가 이렇게 중요했구나"
- "인터넷에는 이렇게 많은 자동화된 트래픽이 돌아다니는구나"

### 최종 기술 스택
- **언어**: Go (시스템 프로그래밍에 적합한 중급 수준 언어)
- **네트워크**: TCP 소켓, HTTP/1.1 프로토콜
- **동시성**: 고루틴을 활용한 비동기 처리
- **분석 도구**: Wireshark (패킷 캡처 및 분석)
- **인프라**: 공인 IP 직접 할당, Windows 방화벽

---

*2025년 9월 - Go 홈서버 구축 프로젝트 완성*  
*"유튜브에서 우연히 보게된 한 동영상."웹개발자로서 네트워크의 전체적인 흐름을 파악하기 위한 학습 프로젝트를 로우언어로 만들어봐라. "학부시절 해본적이 있던것 같지만 가물가물하기도 하고, 이번 기회에 네트워크 관련 복기를 하기위한 프로젝트였다."*

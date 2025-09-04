// 실시간 시계
function updateClock() {
    const now = new Date();
    const timeString = now.toLocaleTimeString('ko-KR');
    document.getElementById('clock').textContent = `⏰ ${timeString}`;
}

setInterval(updateClock, 1000);
updateClock();

// 방문자 수 시뮬레이션 (실제로는 서버에서 관리해야 함)
let visitorCount = Math.floor(Math.random() * 100) + 1;
document.getElementById('visitors').textContent = visitorCount;

// 연결 상태 확인
function checkConnection() {
    const startTime = performance.now();
    
    // 현재 페이지를 다시 요청해서 응답 시간 측정
    fetch(window.location.href)
        .then(response => {
            const endTime = performance.now();
            const responseTime = Math.round(endTime - startTime);
            
            if (response.ok) {
                alert(`✅ 연결 상태: 정상\n🚀 응답 시간: ${responseTime}ms\n📡 상태 코드: ${response.status}`);
            } else {
                alert(`❌ 연결 상태: 오류\n📡 상태 코드: ${response.status}`);
            }
        })
        .catch(error => {
            alert(`❌ 연결 실패: ${error.message}`);
        });
}

// 서버 정보 표시
function showInfo() {
    const info = `
🖥️ 서버 정보:
• 언어: Go
• 포트: 80 (HTTP 표준)
• 프로토콜: HTTP/1.1
• 소켓: TCP
• 구현: 순수 Go (프레임워크 없음)

🌐 네트워크:
• 공인 IP 직접 할당
• 방화벽: Windows (80포트 허용)
• 외부 접속: 전 세계 가능

📚 학습 목표:
• TCP 소켓 직접 구현
• HTTP 프로토콜 이해
• 네트워크 패킷 분석
• Express.js vs 순수 구현 비교
    `;
    alert(info);
}

// 페이지 로드 애니메이션
window.addEventListener('load', function() {
    document.querySelector('.container').style.opacity = '0';
    document.querySelector('.container').style.transform = 'translateY(50px)';
    
    setTimeout(() => {
        document.querySelector('.container').style.transition = 'all 0.8s ease';
        document.querySelector('.container').style.opacity = '1';
        document.querySelector('.container').style.transform = 'translateY(0)';
    }, 100);
});

// 키보드 이스터에그
let konami = [];
const konamiCode = ['ArrowUp', 'ArrowUp', 'ArrowDown', 'ArrowDown', 'ArrowLeft', 'ArrowRight', 'ArrowLeft', 'ArrowRight'];

document.addEventListener('keydown', function(e) {
    konami.push(e.code);
    if (konami.length > konamiCode.length) {
        konami.shift();
    }
    
    if (JSON.stringify(konami) === JSON.stringify(konamiCode)) {
        document.body.style.background = 'linear-gradient(45deg, #ff0000, #ff7f00, #ffff00, #00ff00, #0000ff, #4b0082, #9400d3)';
        document.body.style.backgroundSize = '400% 400%';
        document.body.style.animation = 'rainbow 2s ease infinite';
        
        const style = document.createElement('style');
        style.textContent = `
            @keyframes rainbow {
                0% { background-position: 0% 50%; }
                50% { background-position: 100% 50%; }
                100% { background-position: 0% 50%; }
            }
        `;
        document.head.appendChild(style);
        
        alert('🎉 숨겨진 기능 발견! 무지개 모드 활성화!');
    }
});
# watch app (PoC)

이 디렉토리는 Apple Watch용 생존 체크 코어 로직을 담는다.

## scope
- 1시간 알림 루프
- 스누즈 제어: 1h / 2h / 3h
- 수면/충전 상태에서는 알림 스킵
- 알림 해제 시 `alive` 이벤트 기록

## current status
- Swift Package 기반 도메인 구현
- `swift test`로 코어 로직 테스트 가능
- 다음 단계: SwiftUI 화면 + UserNotifications + HealthKit 어댑터 연결

# death-project

## 소개
이 프로젝트는 맨날 전화, 슬랙, TG도 안보고 객사하는 대표가 살았는지 죽었는지 알기 위해 시작된 프로젝트입니다.

## 1차 목표 (MVP)
- 1시간 간격으로 대표에게 큰 소리 알림 전송
- 대표가 알림을 끄면 `생존 확인`으로 판정
- 대표 화면에서 `1시간`, `2시간`, `3시간` 일시정지 버튼 제공
- 애플워치 기반 수면 상태면 알림 미전송
- 워치 충전 상태면 알림 미전송

## 기술 스택
- FE: watchOS (Apple Watch), Swift/SwiftUI
- BE: Go
- 데이터 저장소: PostgreSQL (추후 확정)

## 개발 언어 및 버전
- FE 언어: Swift `6.1.2`
- 워치 OS 타깃: watchOS `10+` (Xcode 프로젝트에서 최종 고정)
- BE 언어: Go `1.26.0`
- DB: PostgreSQL `18.3`

## 문서
- [제품 요구사항](./SPEC.md)
- [아키텍처 초안](./docs/architecture.md)

## 현재 코드 구조
- `apps/watch`: 워치 생존 체크 도메인 로직 (Swift Package)
- `apps/backend`: 이벤트 수집/조회 API 스켈레톤 (Go)
- `packages`: 공유 패키지 확장용 디렉토리

## 빠른 시작
```bash
bun install
bun run lint
bun run test
bun run build
```

백엔드 실행(Go 설치 후):
```bash
cd apps/backend
go run ./cmd/server
```

## 개발 아키텍처 상세 명세
### 1. 시스템 구조
이 프로젝트는 `워치 로컬 판정 우선(Local-first)` 구조를 사용한다. 핵심 목적은 네트워크가 불안정해도 1시간 생존 확인 루프가 끊기지 않게 하는 것이다.

1. Watch App (Swift/SwiftUI, watchOS)
2. Health/Sensor Adapter (수면, 심박, 충전 상태 수집)
3. Alert Scheduler (1시간 주기 및 스누즈 계산)
4. Alert UI (알림 표시, 해제 버튼, 1/2/3시간 버튼)
5. Event Queue (오프라인 로컬 큐)
6. Sync Client (Go API로 재전송)
7. Go API Server (이벤트 수신, 조회 API)
8. PostgreSQL (이벤트/상태 영속화)

### 2. 워치 앱 내부 모듈 책임
1. `scheduler`
   - 다음 알림 시각(`next_alert_at`) 계산
   - `snooze_until` 반영
2. `state-detector`
   - 수면 상태: HealthKit Sleep Analysis 기반 판정
   - 충전 상태: WatchKit 배터리 상태 기반 판정
3. `alert-engine`
   - 고강도 알림(소리/진동) 노출
   - 사용자 해제 이벤트 발생
4. `decision-engine`
   - 알림 송신 여부 최종 판정
   - 스킵 사유 코드 생성(`skip_sleeping`, `skip_charging`, `skip_snoozed`)
5. `event-store`
   - 이벤트 UUID 생성
   - 로컬 저장 후 동기화 대기
6. `sync`
   - 재시도 정책(지수 백오프)
   - 성공 시 큐에서 제거

### 3. 상태 판정 우선순위
알림 시점마다 아래 순서로 평가한다.

1. `snoozed`: 사용자가 1/2/3시간 일시정지 설정했는가
2. `sleeping`: 사용자가 수면 상태인가
3. `charging`: 워치가 충전 중인가
4. `alert_required`: 위 조건이 모두 false면 알림 송신

### 4. 생존 판정 규칙
1. 알림을 사용자가 해제하면 즉시 `alive`로 기록한다.
2. 스누즈/수면/충전 중인 경우는 생존 미확정 상태로 `skip_*` 이벤트만 기록한다.
3. 무응답 위험 판정은 백엔드 정책으로 확장한다(예: N회 연속 무응답).

### 5. 핵심 시퀀스
1. 스케줄러가 1시간 주기로 Tick 실행
2. 상태 감지(스누즈, 수면, 충전)
3. 스킵이면 `skip_*` 이벤트 기록 후 종료
4. 스킵이 아니면 고강도 알림 송신
5. 사용자 해제 시 `alive` 이벤트 기록
6. 이벤트를 로컬 큐에 저장
7. 네트워크 가능 시 Go API로 전송

### 6. API 명세 (MVP)
#### `POST /v1/events`
워치 이벤트 업로드

요청 예시:
```json
{
  "event_id": "2f95e7f5-58a4-4f53-a572-0b59e6fbfa4e",
  "device_id": "watch-001",
  "event_type": "alive",
  "occurred_at": "2026-03-03T14:30:00+09:00",
  "metadata": {
    "battery_level": 72,
    "is_sleeping": false,
    "is_charging": false
  }
}
```

응답 예시:
```json
{
  "accepted": true,
  "server_time": "2026-03-03T14:30:01+09:00"
}
```

#### `POST /v1/snoozes`
스누즈 상태 동기화

요청 예시:
```json
{
  "device_id": "watch-001",
  "duration_hours": 2,
  "start_at": "2026-03-03T15:00:00+09:00",
  "end_at": "2026-03-03T17:00:00+09:00"
}
```

#### `GET /v1/devices/{id}/latest-status`
최신 판정 조회

응답 예시:
```json
{
  "device_id": "watch-001",
  "latest_event_type": "skip_sleeping",
  "latest_event_at": "2026-03-03T16:00:00+09:00",
  "is_sleeping": true,
  "is_charging": false,
  "snooze_until": null
}
```

### 7. DB 스키마 명세 (MVP)
1. `devices`
   - `id` (PK), `owner_name`, `timezone`, `created_at`
2. `health_states`
   - `id` (PK), `device_id` (FK), `is_sleeping`, `is_charging`, `heart_rate`, `captured_at`
3. `snoozes`
   - `id` (PK), `device_id` (FK), `duration_hours` (1|2|3), `start_at`, `end_at`, `created_at`
4. `check_events`
   - `id` (PK), `event_id` (UNIQUE), `device_id` (FK), `event_type`, `occurred_at`, `metadata_json`, `created_at`

인덱스:
1. `check_events(device_id, occurred_at DESC)`
2. `health_states(device_id, captured_at DESC)`
3. `snoozes(device_id, end_at DESC)`

### 8. 장애 대응 정책
1. API 전송 실패 시 로컬 큐 보관 후 재시도
2. 서버 장애 시에도 워치 알림/판정 루프는 독립 실행
3. 중복 전송 방지를 위해 `event_id` 멱등성 보장
4. 시간 드리프트 방지를 위해 서버 응답 시간으로 주기 보정

### 9. 보안/개인정보 최소 수집
1. 저장 데이터는 생존 판정에 필요한 최소 필드만 유지
2. 민감 생체 데이터 원본은 장기 보관하지 않고 요약값만 저장
3. API 통신은 TLS 강제
4. 디바이스 키 기반 인증(토큰 만료/재발급 포함) 적용

### 10. 테스트 전략 (BDD)
1. Given 수면 상태, When 알림 시각 도달, Then `skip_sleeping` 기록
2. Given 충전 상태, When 알림 시각 도달, Then `skip_charging` 기록
3. Given 스누즈 2시간 설정, When 2시간 이내 Tick, Then 알림 미송신
4. Given 알림 노출, When 사용자 해제, Then 5초 이내 `alive` 저장
5. Given 네트워크 단절, When 이벤트 발생, Then 로컬 큐 적재 후 재연결 시 동기화

## 상태 판정 규칙
- 알림 해제 이벤트 수신 시: `alive`
- 수면 상태 감지 시: `skip_sleeping`
- 충전 상태 감지 시: `skip_charging`
- 사용자가 1/2/3시간 일시정지를 설정한 경우: `skip_snoozed`

## 다음 단계
- watchOS 알림/권한 PoC
- HealthKit 수면 데이터 가용성 검증(대상 워치 모델 기준)
- Go 백엔드 API 및 이벤트 스키마 구현

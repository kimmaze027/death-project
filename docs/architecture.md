# Architecture Draft

## 구성 요소
- Watch App (Swift / watchOS)
- Scheduler (워치 로컬 주기 엔진)
- State Detector (수면/충전 감지)
- Event Logger (로컬 큐)
- Sync Client (백엔드 전송)
- Go API Server (이벤트 수신/조회)
- PostgreSQL (이벤트 저장)

## 알림 처리 흐름
1. 스케줄러가 다음 알림 시점을 확인한다.
2. 상태 감지기가 `수면`, `충전`, `스누즈` 상태를 우선 평가한다.
3. 스킵 조건이 없으면 알림을 큰 소리로 노출한다.
4. 사용자가 해제하면 `alive` 이벤트를 기록한다.
5. 이벤트는 로컬 큐에 적재 후 백엔드로 동기화한다.

## 상태 우선순위
1. `snoozed`
2. `sleeping`
3. `charging`
4. `alert_required`

## 데이터 모델 (초안)
- `devices`
  - `id`, `owner_name`, `timezone`, `created_at`
- `health_states`
  - `id`, `device_id`, `is_sleeping`, `is_charging`, `captured_at`
- `check_events`
  - `id`, `device_id`, `event_type` (`alert_sent` | `alive` | `skip_sleeping` | `skip_charging` | `skip_snoozed`), `occurred_at`, `metadata`
- `snoozes`
  - `id`, `device_id`, `duration_hours` (1|2|3), `start_at`, `end_at`

## API (초안)
- `POST /v1/events`
  - 워치 이벤트 업로드
- `POST /v1/snoozes`
  - 스누즈 상태 동기화
- `GET /v1/devices/{id}/timeline?from=&to=`
  - 이벤트 타임라인 조회
- `GET /v1/devices/{id}/latest-status`
  - 최신 판정 상태 조회

## 장애 대응
- 네트워크 실패 시 로컬 큐 재시도(지수 백오프)
- 서버 장애 시 워치 로컬 알림 기능은 독립 유지
- 중복 업로드 방지를 위해 이벤트 UUID 사용

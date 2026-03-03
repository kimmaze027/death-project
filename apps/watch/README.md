# watch app (PoC)

This directory contains the Watch-side alive-check core logic.

## scope
- 1 hour alert cycle
- Snooze controls: 1h / 2h / 3h
- Skip alert when sleeping or charging
- Record `alive` when user dismisses the alert

## current status
- Domain-only Kotlin implementation (framework-agnostic)
- Ready to connect with Wear OS UI, notifications, and Health Services adapters

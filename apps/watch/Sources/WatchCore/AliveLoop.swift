import Foundation

public enum AliveLoopError: Error, Equatable {
    case invalidSnoozeHours
}

public final class AliveLoop {
    private let config: LoopConfig
    private let decisionEngine: DecisionEngine

    private var nextAlertAt: Date?
    private var snoozeUntil: Date?

    public init(
        config: LoopConfig = LoopConfig(),
        decisionEngine: DecisionEngine = DecisionEngine()
    ) {
        self.config = config
        self.decisionEngine = decisionEngine
    }

    public func bootstrap(startAt: Date) {
        nextAlertAt = startAt.addingTimeInterval(config.intervalSeconds)
    }

    public func currentNextAlertAt() -> Date? {
        nextAlertAt
    }

    public func currentSnoozeUntil() -> Date? {
        snoozeUntil
    }

    public func applySnooze(now: Date, hours: Int) throws {
        guard [1, 2, 3].contains(hours) else {
            throw AliveLoopError.invalidSnoozeHours
        }
        snoozeUntil = now.addingTimeInterval(config.intervalSeconds * Double(hours))
    }

    public func tick(now: Date, isSleeping: Bool, isCharging: Bool) -> CheckEvent? {
        guard let scheduledAt = nextAlertAt, now >= scheduledAt else {
            return nil
        }

        let state = WatchState(
            now: now,
            isSleeping: isSleeping,
            isCharging: isCharging,
            snoozeUntil: snoozeUntil
        )

        let decision = decisionEngine.decide(state)
        let event: CheckEvent

        switch decision {
        case .sendAlert:
            event = CheckEvent(type: .alertSent, occurredAt: now)
        case .skip(let reason):
            event = CheckEvent(
                type: mapSkipReasonToEventType(reason),
                occurredAt: now,
                metadata: ["reason": reason.rawValue]
            )
        }

        nextAlertAt = now.addingTimeInterval(config.intervalSeconds)
        if let snoozeUntil, now >= snoozeUntil {
            self.snoozeUntil = nil
        }

        return event
    }

    public func onAlertDismissed(now: Date) -> CheckEvent {
        CheckEvent(type: .alive, occurredAt: now)
    }

    private func mapSkipReasonToEventType(_ reason: SkipReason) -> EventType {
        switch reason {
        case .snoozed: return .skipSnoozed
        case .sleeping: return .skipSleeping
        case .charging: return .skipCharging
        }
    }
}

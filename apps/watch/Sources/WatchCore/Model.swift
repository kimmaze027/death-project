import Foundation

public enum EventType: String, Equatable {
    case alertSent = "alert_sent"
    case alive = "alive"
    case skipSnoozed = "skip_snoozed"
    case skipSleeping = "skip_sleeping"
    case skipCharging = "skip_charging"
}

public enum SkipReason: String, Equatable {
    case snoozed
    case sleeping
    case charging
}

public struct WatchState: Equatable {
    public let now: Date
    public let isSleeping: Bool
    public let isCharging: Bool
    public let snoozeUntil: Date?

    public init(now: Date, isSleeping: Bool, isCharging: Bool, snoozeUntil: Date?) {
        self.now = now
        self.isSleeping = isSleeping
        self.isCharging = isCharging
        self.snoozeUntil = snoozeUntil
    }
}

public struct CheckEvent: Equatable {
    public let id: UUID
    public let type: EventType
    public let occurredAt: Date
    public let metadata: [String: String]

    public init(
        id: UUID = UUID(),
        type: EventType,
        occurredAt: Date,
        metadata: [String: String] = [:]
    ) {
        self.id = id
        self.type = type
        self.occurredAt = occurredAt
        self.metadata = metadata
    }
}

public struct LoopConfig: Equatable {
    public let intervalSeconds: TimeInterval

    public init(intervalSeconds: TimeInterval = 3600) {
        self.intervalSeconds = intervalSeconds
    }
}

public enum Decision: Equatable {
    case sendAlert
    case skip(SkipReason)
}

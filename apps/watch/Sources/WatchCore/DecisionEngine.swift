import Foundation

public final class DecisionEngine {
    public init() {}

    public func decide(_ state: WatchState) -> Decision {
        if let snoozeUntil = state.snoozeUntil, state.now < snoozeUntil {
            return .skip(.snoozed)
        }
        if state.isSleeping {
            return .skip(.sleeping)
        }
        if state.isCharging {
            return .skip(.charging)
        }
        return .sendAlert
    }
}

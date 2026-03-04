package io.deathproject.watch.core

class DecisionEngine {
    fun decide(state: WatchState): Decision {
        if (state.snoozeUntil != null && state.now.isBefore(state.snoozeUntil)) {
            return Decision.Skip(SkipReason.SNOOZED)
        }
        if (state.isSleeping) {
            return Decision.Skip(SkipReason.SLEEPING)
        }
        if (state.isCharging) {
            return Decision.Skip(SkipReason.CHARGING)
        }
        return Decision.SendAlert
    }
}

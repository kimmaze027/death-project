package io.deathproject.watch.core

import java.time.Instant

class AliveLoop(
    private val config: LoopConfig = LoopConfig(),
    private val decisionEngine: DecisionEngine = DecisionEngine(),
) {
    private var nextAlertAt: Instant? = null
    private var snoozeUntil: Instant? = null

    fun bootstrap(startAt: Instant) {
        nextAlertAt = startAt.plus(config.interval)
    }

    fun getNextAlertAt(): Instant? = nextAlertAt

    fun getSnoozeUntil(): Instant? = snoozeUntil

    fun applySnooze(now: Instant, hours: Int) {
        require(hours in setOf(1, 2, 3)) { "hours must be one of 1,2,3" }
        snoozeUntil = now.plus(config.interval.multipliedBy(hours.toLong()))
    }

    fun tick(now: Instant, isSleeping: Boolean, isCharging: Boolean): CheckEvent? {
        val scheduledAt = nextAlertAt ?: return null
        if (now.isBefore(scheduledAt)) {
            return null
        }

        val state = WatchState(
            now = now,
            isSleeping = isSleeping,
            isCharging = isCharging,
            snoozeUntil = snoozeUntil,
        )

        val event = when (val decision = decisionEngine.decide(state)) {
            Decision.SendAlert -> CheckEvent(type = EventType.ALERT_SENT, occurredAt = now)
            is Decision.Skip -> CheckEvent(
                type = skipReasonToEventType(decision.reason),
                occurredAt = now,
                metadata = mapOf("reason" to decision.reason.name.lowercase()),
            )
        }

        nextAlertAt = now.plus(config.interval)
        if (snoozeUntil != null && !now.isBefore(snoozeUntil)) {
            snoozeUntil = null
        }

        return event
    }

    fun onAlertDismissed(now: Instant): CheckEvent {
        return CheckEvent(type = EventType.ALIVE, occurredAt = now)
    }

    private fun skipReasonToEventType(reason: SkipReason): EventType {
        return when (reason) {
            SkipReason.SNOOZED -> EventType.SKIP_SNOOZED
            SkipReason.SLEEPING -> EventType.SKIP_SLEEPING
            SkipReason.CHARGING -> EventType.SKIP_CHARGING
        }
    }
}

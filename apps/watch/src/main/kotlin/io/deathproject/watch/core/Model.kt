package io.deathproject.watch.core

import java.time.Duration
import java.time.Instant
import java.util.UUID

enum class EventType {
    ALERT_SENT,
    ALIVE,
    SKIP_SNOOZED,
    SKIP_SLEEPING,
    SKIP_CHARGING,
}

enum class SkipReason {
    SNOOZED,
    SLEEPING,
    CHARGING,
}

data class WatchState(
    val now: Instant,
    val isSleeping: Boolean,
    val isCharging: Boolean,
    val snoozeUntil: Instant?,
)

data class CheckEvent(
    val id: UUID = UUID.randomUUID(),
    val type: EventType,
    val occurredAt: Instant,
    val metadata: Map<String, String> = emptyMap(),
)

data class LoopConfig(
    val interval: Duration = Duration.ofHours(1),
)

sealed interface Decision {
    data object SendAlert : Decision
    data class Skip(val reason: SkipReason) : Decision
}

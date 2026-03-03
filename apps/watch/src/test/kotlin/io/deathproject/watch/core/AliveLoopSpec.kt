package io.deathproject.watch.core

import java.time.Instant

/**
 * BDD test scenarios for AliveLoop.
 *
 * Given the local environment does not include Kotlin/JDK toolchain yet,
 * this file captures executable test intent to run once Gradle/JDK is wired.
 */
class AliveLoopSpec {
    // Given not sleeping/charging and no snooze
    // When tick reaches schedule
    // Then ALERT_SENT should be generated

    // Given sleeping
    // When tick reaches schedule
    // Then SKIP_SLEEPING should be generated

    // Given charging
    // When tick reaches schedule
    // Then SKIP_CHARGING should be generated

    // Given snooze 2h at T0
    // When tick occurs before T0+2h
    // Then SKIP_SNOOZED should be generated

    // Given an alert was shown
    // When user dismisses
    // Then ALIVE should be generated

    @Suppress("unused")
    private fun sampleInvocation() {
        val loop = AliveLoop()
        val t0 = Instant.parse("2026-03-03T00:00:00Z")
        loop.bootstrap(t0)
        loop.tick(now = t0.plusSeconds(3600), isSleeping = false, isCharging = false)
        loop.onAlertDismissed(t0.plusSeconds(3602))
    }
}

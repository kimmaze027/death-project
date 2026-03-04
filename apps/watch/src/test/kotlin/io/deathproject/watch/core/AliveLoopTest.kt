package io.deathproject.watch.core

import java.time.Instant
import kotlin.test.Test
import kotlin.test.assertEquals
import kotlin.test.assertNotNull
import kotlin.test.assertNull
import kotlin.test.assertTrue
import kotlin.test.assertFailsWith

class AliveLoopTest {
    @Test
    fun `Given start time When bootstrap Then next alert is one hour later`() {
        val loop = AliveLoop()
        val start = Instant.parse("2026-03-03T00:00:00Z")

        loop.bootstrap(start)

        assertEquals(Instant.parse("2026-03-03T01:00:00Z"), loop.getNextAlertAt())
    }

    @Test
    fun `Given scheduled loop When tick before schedule Then no event`() {
        val loop = AliveLoop()
        val start = Instant.parse("2026-03-03T00:00:00Z")
        loop.bootstrap(start)

        val event = loop.tick(
            now = Instant.parse("2026-03-03T00:59:59Z"),
            isSleeping = false,
            isCharging = false,
        )

        assertNull(event)
    }

    @Test
    fun `Given awake and not charging When tick on schedule Then alert event is emitted`() {
        val loop = AliveLoop()
        val start = Instant.parse("2026-03-03T00:00:00Z")
        loop.bootstrap(start)

        val event = loop.tick(
            now = Instant.parse("2026-03-03T01:00:00Z"),
            isSleeping = false,
            isCharging = false,
        )

        assertNotNull(event)
        assertEquals(EventType.ALERT_SENT, event.type)
        assertEquals(Instant.parse("2026-03-03T02:00:00Z"), loop.getNextAlertAt())
    }

    @Test
    fun `Given sleeping state When tick on schedule Then skip sleeping event is emitted`() {
        val loop = AliveLoop()
        val start = Instant.parse("2026-03-03T00:00:00Z")
        loop.bootstrap(start)

        val event = loop.tick(
            now = Instant.parse("2026-03-03T01:00:00Z"),
            isSleeping = true,
            isCharging = false,
        )

        assertNotNull(event)
        assertEquals(EventType.SKIP_SLEEPING, event.type)
        assertEquals("sleeping", event.metadata["reason"])
    }

    @Test
    fun `Given charging state When tick on schedule Then skip charging event is emitted`() {
        val loop = AliveLoop()
        val start = Instant.parse("2026-03-03T00:00:00Z")
        loop.bootstrap(start)

        val event = loop.tick(
            now = Instant.parse("2026-03-03T01:00:00Z"),
            isSleeping = false,
            isCharging = true,
        )

        assertNotNull(event)
        assertEquals(EventType.SKIP_CHARGING, event.type)
        assertEquals("charging", event.metadata["reason"])
    }

    @Test
    fun `Given snoozed and sleeping and charging When tick on schedule Then snoozed has priority`() {
        val loop = AliveLoop()
        val start = Instant.parse("2026-03-03T00:00:00Z")
        loop.bootstrap(start)
        loop.applySnooze(now = start, hours = 2)

        val event = loop.tick(
            now = Instant.parse("2026-03-03T01:00:00Z"),
            isSleeping = true,
            isCharging = true,
        )

        assertNotNull(event)
        assertEquals(EventType.SKIP_SNOOZED, event.type)
        assertEquals("snoozed", event.metadata["reason"])
    }

    @Test
    fun `Given expired snooze When tick on boundary Then alert is emitted and snooze is cleared`() {
        val loop = AliveLoop()
        val start = Instant.parse("2026-03-03T00:00:00Z")
        loop.bootstrap(start)
        loop.applySnooze(now = start, hours = 1)

        val event = loop.tick(
            now = Instant.parse("2026-03-03T01:00:00Z"),
            isSleeping = false,
            isCharging = false,
        )

        assertNotNull(event)
        assertEquals(EventType.ALERT_SENT, event.type)
        assertNull(loop.getSnoozeUntil())
    }

    @Test
    fun `Given unsupported snooze hour When applySnooze Then exception is thrown`() {
        val loop = AliveLoop()
        val start = Instant.parse("2026-03-03T00:00:00Z")

        val ex = assertFailsWith<IllegalArgumentException> {
            loop.applySnooze(now = start, hours = 4)
        }

        assertTrue(ex.message?.contains("1,2,3") == true)
    }

    @Test
    fun `Given alert dismissed When onAlertDismissed Then alive event is emitted`() {
        val loop = AliveLoop()
        val now = Instant.parse("2026-03-03T01:00:02Z")

        val event = loop.onAlertDismissed(now)

        assertEquals(EventType.ALIVE, event.type)
        assertEquals(now, event.occurredAt)
    }
}

import Foundation
import Testing
@testable import WatchCore

private func isoDate(_ value: String) -> Date {
    let formatter = ISO8601DateFormatter()
    formatter.formatOptions = [.withInternetDateTime]
    guard let date = formatter.date(from: value) else {
        fatalError("invalid iso date: \(value)")
    }
    return date
}

@Suite("AliveLoop")
struct AliveLoopTests {
    @Test("Given start time When bootstrap Then next alert is one hour later")
    func bootstrapSetsNextAlert() {
        let loop = AliveLoop()
        let start = isoDate("2026-03-03T00:00:00Z")

        loop.bootstrap(startAt: start)

        #expect(loop.currentNextAlertAt() == isoDate("2026-03-03T01:00:00Z"))
    }

    @Test("Given schedule not reached When tick Then no event")
    func tickBeforeScheduleReturnsNil() {
        let loop = AliveLoop()
        let start = isoDate("2026-03-03T00:00:00Z")
        loop.bootstrap(startAt: start)

        let event = loop.tick(
            now: isoDate("2026-03-03T00:59:59Z"),
            isSleeping: false,
            isCharging: false
        )

        #expect(event == nil)
    }

    @Test("Given awake and not charging When tick on schedule Then alert_sent")
    func tickOnScheduleSendsAlert() {
        let loop = AliveLoop()
        let start = isoDate("2026-03-03T00:00:00Z")
        loop.bootstrap(startAt: start)

        let event = loop.tick(
            now: isoDate("2026-03-03T01:00:00Z"),
            isSleeping: false,
            isCharging: false
        )

        #expect(event?.type == .alertSent)
        #expect(loop.currentNextAlertAt() == isoDate("2026-03-03T02:00:00Z"))
    }

    @Test("Given sleeping When tick Then skip_sleeping")
    func sleepingSkipsAlert() {
        let loop = AliveLoop()
        let start = isoDate("2026-03-03T00:00:00Z")
        loop.bootstrap(startAt: start)

        let event = loop.tick(
            now: isoDate("2026-03-03T01:00:00Z"),
            isSleeping: true,
            isCharging: false
        )

        #expect(event?.type == .skipSleeping)
        #expect(event?.metadata["reason"] == "sleeping")
    }

    @Test("Given charging When tick Then skip_charging")
    func chargingSkipsAlert() {
        let loop = AliveLoop()
        let start = isoDate("2026-03-03T00:00:00Z")
        loop.bootstrap(startAt: start)

        let event = loop.tick(
            now: isoDate("2026-03-03T01:00:00Z"),
            isSleeping: false,
            isCharging: true
        )

        #expect(event?.type == .skipCharging)
        #expect(event?.metadata["reason"] == "charging")
    }

    @Test("Given snoozed+sleeping+charging When tick Then snoozed has priority")
    func snoozeHasPriority() throws {
        let loop = AliveLoop()
        let start = isoDate("2026-03-03T00:00:00Z")
        loop.bootstrap(startAt: start)
        try loop.applySnooze(now: start, hours: 2)

        let event = loop.tick(
            now: isoDate("2026-03-03T01:00:00Z"),
            isSleeping: true,
            isCharging: true
        )

        #expect(event?.type == .skipSnoozed)
        #expect(event?.metadata["reason"] == "snoozed")
    }

    @Test("Given snooze expired at boundary Then alert_sent and snooze cleared")
    func expiredSnoozeClearsAndSendsAlert() throws {
        let loop = AliveLoop()
        let start = isoDate("2026-03-03T00:00:00Z")
        loop.bootstrap(startAt: start)
        try loop.applySnooze(now: start, hours: 1)

        let event = loop.tick(
            now: isoDate("2026-03-03T01:00:00Z"),
            isSleeping: false,
            isCharging: false
        )

        #expect(event?.type == .alertSent)
        #expect(loop.currentSnoozeUntil() == nil)
    }

    @Test("Given unsupported snooze hour When applySnooze Then throw")
    func invalidSnoozeThrows() {
        let loop = AliveLoop()

        #expect(throws: AliveLoopError.invalidSnoozeHours) {
            try loop.applySnooze(now: isoDate("2026-03-03T00:00:00Z"), hours: 4)
        }
    }

    @Test("Given alert dismissed When onAlertDismissed Then alive")
    func dismissCreatesAliveEvent() {
        let loop = AliveLoop()
        let now = isoDate("2026-03-03T01:00:02Z")

        let event = loop.onAlertDismissed(now: now)

        #expect(event.type == .alive)
        #expect(event.occurredAt == now)
    }
}

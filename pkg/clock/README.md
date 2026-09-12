# Clock

Implemented against [CONTRACT.md](CONTRACT.md), revision 1. System and immutable fixed wall clocks, with exact instant preservation.

No monotonic measurement, sleeping, advancement or scheduler. Pure domains receive timestamp values. IDs are not needed by these APIs.

6 focused tests pass. The combined errors/clock/secret integration gate is
complete. See project `work/handoffs/native-clock.md` for scenario evidence.

# feedback

Presentational state signals read semantic tokens and expose accessible visual states:

- `Badge`, `Status`, `Spinner`, and `Progress` for compact status and progress.
- `Alert` and `Banner` for persistent messages with optional dismissal and actions.
- `EmptyState` and `Skeleton` for no-content and loading states.
- `Toast` and `Toaster` for host-owned transient notifications.

Consumers own event and service semantics. `Toaster` is intentionally controlled:
the host provides the queue and handles dismissal.

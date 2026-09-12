# Worker build order

Source of task status: [manifest.json](manifest.json). The first three leaves and their integration gate are complete; no implementation
task is ready until its next contract is specified. Run `python3 tools/work/verify_manifest.py` before dispatch.

```text
native-errors ─┐
native-clock ──┼→ native-leaves-integrate → preferences-domain → preferences-app
native-secret ┘                              ↑                    ├→ preferences-memory ───┐
preferences-contract ────────────────────────┘                    └→ desktop-wire-contract ─┤
                                                                                           ↓
                                                                             preferences-native-slice
                                                                                           ↓
                                                                             preferences-persistence
```

The leaf integration task can review implemented handoffs together; mark the
leaves complete only after that gate passes. Other dependent tasks require their
prerequisites complete before dispatch. Planned tasks require contracts and task
files even when all dependencies are complete.

The graph schedules the calibration batch before domain implementation; it does
not force preferences to import clocks or secrets. Config and IDs are now complete. Logger/file IO and
source architecture enforcement remain separately planned in the manifest.
Preferences, Workspace and Jobs contracts are specified independently of the
leaf batch. Preferences remains the first implementation slice.

Do not parallelize changes to Cargo module registration, root wiring, dependency
manifests or lockfiles. Keep one coordinator responsible for those files.

## Completed first assignments

- [native-errors](tasks/native-errors.md)
- [native-clock](tasks/native-clock.md)
- [native-secret](tasks/native-secret.md)

Then the coordinator runs [native-leaves-integrate](tasks/native-leaves-integrate.md).
Setup evidence: [worker setup verification](../notes/worker-setup-verification.md).

Config is also complete after errors and secrets. See [native-config](tasks/native-config.md)
and its [handoff](handoffs/native-config.md). No new task becomes ready merely
because these dependencies are complete; planned work still needs its contract.

IDs are complete after errors: [native-id](tasks/native-id.md),
[verification](handoffs/native-id.md). Provenance is now complete; logger is now complete.

Completed: [native-provenance](tasks/native-provenance.md),
[verification](handoffs/native-provenance.md), after errors/IDs/clocks.
Logger uses provenance for its explicit projection contract.

Completed: [native-logger](tasks/native-logger.md),
[verification](handoffs/native-logger.md). Console is composed at native startup.
Preferences, Workspace and Jobs remain planned for implementation; their
revision 1 contracts are complete.

Completed contract tasks: [preferences-contract](tasks/preferences-contract.md),
[workspace-contract](tasks/workspace-contract.md) and
[jobs-contract](tasks/jobs-contract.md), with specification evidence in their
matching handoffs.

Completed: [native-root-bootstrap](tasks/native-root-bootstrap.md),
[verification](handoffs/native-root-bootstrap.md). The root now creates the
bootstrap ID/provenance scope and binds it to lifecycle logging.

# jobs-contract handoff

Revision 1 contract specified 2026-09-12. See `internal/jobs/CONTRACT.md` for
the job state machine, progress and cancellation rules, provenance linkage,
conditional store outcomes, returned event facts and worker boundary.

The domain records operation state only; it does not schedule work, spawn
processes, decide commit proof or publish through a global bus. No runtime code
changed.

Verification: `python3 tools/work/verify_manifest.py` passed after the contract
task metadata was updated.

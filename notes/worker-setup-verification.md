# Worker setup verification · 2026-09-11

Added module-contract/worker-task templates, workflow guidance, a 16-task dependency
manifest and revision-1 native errors/clock/secret contracts with bounded task
files. All three leaves are ready for assignment; their implementation is absent.
Later tasks are explicitly planned or waiting, not implicitly ready.

Checks from the project root:

- `python3 tools/work/verify_manifest.py`: passes; 16 tasks, 3 ready.
- `python3 -m unittest discover -s tools/work -p 'test_*.py'`: 10 tests pass.

Negative metadata controls cover missing task files, revision drift, unknown
prerequisites, cycles, unfinished prerequisites, path escape, missing completion
handoffs, duplicate IDs and an unspecified task marked ready.

This verifies coordination tooling only. Source import enforcement, native domain
behavior, worker handoff review and the reference preference slice remain future
work. No runtime code or dependency versions changed in this setup step.

# jobs-contract: specify the revision 1 boundary

Record the Jobs state machine, cancellation/progress rules, provenance linkage,
conditional store outcomes, returned event facts and worker/adapter boundaries.
Keep scheduling, retries, process execution and delivery infrastructure outside
the first contract.

The contract is documentation-only. Do not implement jobs domain, application,
memory, worker, desktop or persistence code in this task.

Checks:

```sh
python3 tools/work/verify_manifest.py
```

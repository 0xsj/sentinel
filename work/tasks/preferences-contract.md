# preferences-contract: specify the revision 1 boundary

Record the Preferences bounded-context vocabulary, validation rules, conditional
store outcomes, returned event facts and adapter/transport boundaries. Keep
workspace identity opaque and keep secrets outside ordinary preference values.

The contract is documentation-only. Do not implement domain, application,
memory, desktop or persistence code in this task.

Checks:

```sh
python3 tools/work/verify_manifest.py
```

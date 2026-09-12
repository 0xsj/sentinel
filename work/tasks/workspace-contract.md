# workspace-contract: specify the revision 1 boundary

Record the Workspace registry vocabulary, canonical-location rules, lifecycle
transitions, conditional outcomes, returned event facts and path/store adapter
boundaries. Keep file operations and peer contexts outside the domain.

The contract is documentation-only. Do not implement workspace domain,
application, memory, path, desktop or persistence code in this task.

Checks:

```sh
python3 tools/work/verify_manifest.py
```

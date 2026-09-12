# Worker setup verification

From repository root:

```sh
python3 tools/work/verify_manifest.py
python3 -m unittest discover -s tools/work -p 'test_*.py'
```

Manifest checks cover references, states, revisions, paths and dependency cycles.
Tests include invalid metadata controls. They do not verify implementation claims,
source imports, scenario coverage or actual modification scope. Coordinator review
is still required. Checks use Python's standard library.

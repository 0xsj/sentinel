# preferences-contract handoff

Revision 1 contract specified 2026-09-12. See
`internal/preferences/CONTRACT.md` for scope/key/value validation, compare and
replace outcomes, revision rules, returned event facts, narrow store port and
explicit exclusions.

The contract deliberately keeps workspace IDs opaque, routes secret settings to
the existing Secret capability, and leaves event publication to the application
after a successful store transaction. No runtime code changed.

Verification: `python3 tools/work/verify_manifest.py` passed after the contract
task metadata was updated.

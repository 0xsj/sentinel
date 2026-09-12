# workspace-contract handoff

Revision 1 contract specified 2026-09-12. See
`internal/workspace/CONTRACT.md` for metadata vocabulary, canonical location
ownership, active/archived/forgotten lifecycle, conditional outcomes, returned
event facts and the application-owned path capability.

The domain owns registry metadata only; it never reads the filesystem, moves
files or imports peer contexts. No runtime code changed.

Verification: `python3 tools/work/verify_manifest.py` passed after the contract
task metadata was updated.

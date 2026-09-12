# Errors preserve conditions and separate public data

Accepted direction from the user, 2026-09-11. Implement n2f's responsibility split
in both Atelier backends. Go uses condition-preserving immutable derivation;
Rust uses explicit classification and typed context. Both provide a safe,
transport-neutral public snapshot separate from private details and causes.

Revision 2 contracts replace the minimal revision 1 exclusions. This changes
constructor message semantics: non-Internal messages and field problems must be
public-safe. Internal/unknown failures project a fixed fallback. There are no
production callers of the initial leaf, so no application migration was needed.

Retain Atelier's exact diagnostic display and original kind parsing. Public
projection does not finalize desktop encoding, frontend Failure mappings, retries,
logging policy or authentication. Source implementations are locally adapted from
n2f, not imported. Ownership and behavioral evidence live in each errors contract
and work/handoffs/native-errors.md.

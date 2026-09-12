# IDs

Implemented: UUID values, UUIDv7 production generation, and finite deterministic
sequences. Adapted locally from n2f; no sibling imports. See [contract](CONTRACT.md)
for parsing, rollback, exhaustion, error and concurrency guarantees.

The composition root supplies the existing clock capability. It now creates the
production generator for the bootstrap provenance scope; later consumers can own
smaller generation ports. Domain-specific wrappers belong to their domains.

Production uses cryptographic OS entropy; tests inject controlled entropy and
private mutable clocks. No mutable clock API is added to the shared clock module.

# Secret

Implemented against [CONTRACT.md](CONTRACT.md), revision 1. An explicit string wrapper with redacted supported diagnostic formats.

No encryption, memory wiping, secret store or reflection/debugger protection. Revealing/exposing contents is explicit; only the contract-listed formatting forms are guaranteed.

4 focused tests pass. The combined errors/clock/secret integration gate is
complete. See project `work/handoffs/native-secret.md` for scenario evidence.

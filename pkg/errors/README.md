# Native errors

Implemented against [CONTRACT.md](CONTRACT.md), revision 2. This follows n2f's
separation of classification, public information, private diagnostics and native
condition preservation. The Go and Rust implementations share meaning, not
identical APIs. No sibling imports or external dependencies.

```text
kind.go       classification vocabulary
errors.go     failures, sentinel identity, immutable builders
inspect.go    single-chain classification and private inspection
public.go     disclosure policy and owned public snapshots
doc.go        package responsibility and usage
```

Constructor messages and field problems for non-Internal kinds must be public-safe.
Use details/causes for private diagnostics. Internal and unknown projections hide
message/type/fields. Error/Display and diagnostic inspection are not UI APIs.
Projection selects one outer frame and never inherits missing data from a cause.

```go
var ErrNameTaken = failure.New(failure.Conflict, "name already used").
    WithType("workspace.name_taken")

err := ErrNameTaken.WithCause(databaseError).WithDetail("operation", "rename")
// errors.Is(err, ErrNameTaken) remains true.
view, present := failure.Public(err) // pass view to the future wire encoder
```

Go distinguishes literal nil from unknown failures and retains sentinel identity
across derivations. Rust uses ordinary Result and typed Context<E>/Classified;
public_info(None) means an actual unknown failure, not success. Foreign Rust errors
need explicit adapter classification. Neither language automatically retries.

Nineteen tests pass per implementation. See project `work/handoffs/native-errors.md`
for scenario mapping and native integration evidence. The combined leaf-batch gate is complete. JSON/desktop bridge encoding remains separate.

# Native errors: first contract

Revision: 1 · State: specified · Task: native-errors

## Purpose and dependency boundary

Provide native failure classification and diagnostic cause preservation for future
application/adapters. The first scheduled consumer is preference operation error
mapping. This is a pure standard-library-only leaf; no native framework, sibling
module, serialization, logging or IO imports.

## Vocabulary and semantics

Ten canonical names: unauthenticated, forbidden, rate_limited, unavailable,
timeout, canceled, internal, not_found, invalid, conflict. Names are exact lowercase
ASCII; parsing neither trims nor folds case. Native classification is not an
assertion that every operation can produce every kind. Domain-owned error types
may remain more specific and be mapped explicitly by application/transport.

A failure owns its kind and diagnostic message. Attaching a cause preserves those
outer values; it does not infer or replace classification from the cause's text
or kind. Different occurrences with equal contents are not the same occurrence.
Diagnostics are not approved public UI messages. No serialization/public projection,
fields/retry metadata, generic retry predicate, aggregate errors, cause traversal
or transport policy is part of this slice. Those require their own contract.

### Go API and language-specific cases

Package `errors` under `pkg/errors` (alias it `failure` at ambiguous call sites).

```go
type Kind uint8
const (
    Internal Kind = iota
    Unauthenticated
    Forbidden
    RateLimited
    Unavailable
    Timeout
    Canceled
    NotFound
    Invalid
    Conflict
)
func (k Kind) String() string
func ParseKind(text string) (Kind, bool)
func Kinds() []Kind

type Failure struct { /* private fields */ }
func New(kind Kind, diagnostic string) *Failure
func (f *Failure) Kind() Kind
func (f *Failure) Diagnostic() string
func (f *Failure) Error() string
func (f *Failure) Unwrap() error
func (f *Failure) WithCause(cause error) *Failure
```

Kind zero is Internal. Invalid numeric Kind.String returns "unknown"; New normalizes
invalid numeric kinds to Internal. ParseKind's failure value is Internal,false.
Kinds returns caller-owned storage. New always returns non-nil; the zero Failure
value is Internal with empty diagnostic and nil cause. Nil *Failure receivers are
outside this API contract; literal nil error represents no failure.

WithCause returns a new occurrence without changing its receiver and replaces any
previous cause. Literal nil removes the cause; a typed-nil error is retained as
supplied and never called/inspected by the builder. Standard errors.Is reaches a
supplied sentinel; errors.As reaches a supplied typed cause. Do not implement
content-based Is. Foreign causes are references, not deep copies; no claim of
immutability of foreign error objects is made.

E08: mutate Kinds output and verify later enumeration unchanged. E09: test zero
Failure, invalid numeric kind construction, and nil-cause removal. E10: verify
WithCause leaves the receiver unchanged, replacement, typed-nil retention and
standard Is/As. Source review verifies no mutable fields are exported.


## Acceptance scenarios

| ID | Given / when | Required observable outcome |
| --- | --- | --- |
| E01 | Enumerate kinds and parse each canonical name | Exactly the ten names, no duplicates, round trip for each; enumeration order not promised |
| E02 | Parse empty, uppercase, whitespace-padded or unknown text | Rejected; never silently classified as internal |
| E03 | Construct every kind with an ordinary or empty diagnostic | Kind and exact diagnostic preserved; display uses the diagnostic verbatim |
| E04 | Attach a distinct classified cause | Outer kind/diagnostic unchanged; native cause inspection reaches the supplied cause |
| E05 | Attach a foreign error whose message says timeout | No automatic classification; outer classification preserved |
| E06 | Construct two equal-looking occurrences | No content-based native error identity/equality policy is introduced |
| E07 | Read the API | No mutation path to stored fields; no serialization or retry/public-message API |

E07 includes source/API review, not just runtime assertions. Add language-specific
checks below. An empty diagnostic is allowed and must not panic or acquire copy.

## Exclusions and completion

No hidden timestamps/IDs, logging, IO, dependencies or host wiring. Workers may not
edit this contract. Map E01–E07 and language-specific cases to tests/review in the
handoff. Error message redaction happens at a future explicit bridge projection;
never expose this diagnostic Error/Display directly to the frontend.

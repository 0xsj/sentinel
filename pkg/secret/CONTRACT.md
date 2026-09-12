# Redacted secret: first contract

Revision: 1 · State: specified · Task: native-secret

## Purpose and boundary

An explicit string wrapper prevents routine diagnostic formatting from exposing a
secret. Its scheduled consumer is process configuration/adapter credentials. It is
not a domain requirement, encryption, a keychain, access control or secure memory.
Only standard-library dependencies; no IO, runtime, config or logger imports.

### Go API

```go
// package secret
type Secret struct { /* private string */ }
func New(value string) Secret
func (s Secret) Reveal() string
func (s Secret) String() string
func (s Secret) GoString() string
```

String and GoString return exactly [REDACTED]. The zero value reveals an empty
string and still formats as [REDACTED]. S02 must cover fmt.Sprint, %s, %v, %+v and
%#v. These are the supported formatting forms; do not claim coverage for arbitrary
verbs or reflection. Add S05 for zero-value behavior. No MarshalJSON/MarshalText,
unmarshal methods, byte getters, Set methods or third-party secret container.


## Acceptance scenarios

| ID | Given / when | Required observable outcome |
| --- | --- | --- |
| S01 | Construct empty, Unicode, newline or ordinary string | Explicit reveal returns the exact input; no trimming/validation |
| S02 | Format through supported diagnostic forms | Exactly the redaction marker [REDACTED], independent of contents |
| S03 | Store two different values | Reveals remain independent; diagnostic output does not distinguish contents |
| S04 | Inspect production surface | No implicit string conversion exposing contents, serialization, setters or equality API |

No memory wiping promise: language strings and explicit reveals can retain copies.
No reflection/debugger protection. Debug formats not specified by this contract
are not a security boundary. Revealing a value transfers responsibility to the
caller. Public fields holding raw data are forbidden. This module does not
validate required credentials or decide where they may be transmitted.

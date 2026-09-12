# Process configuration

Implemented against [CONTRACT.md](CONTRACT.md), revision 1, in both Go and Rust.
This is the captured process/environment reader, separate from persisted user
preferences. Root owns concrete settings, defaults and cross-field validation.
No application-specific startup settings are installed yet.

```text
lookup       owned map snapshot for injection/tests
os           one-time native environment capture
parse        strict key/integer helpers
reader       typed reads, accumulated problems, redacted sorted manifest
```

The reader uses the local errors and secret packages. Missing and present-empty
values remain distinct; defaults apply only to absence. Required strings and
secrets reject empty values. Int/bool/enum parsing is exact. Invalid definitions,
duplicate keys and malformed input fail explicitly, without returning raw values
in diagnostics. String/Required/Enum are for caller-declared nonsecret values.

```go
lookup, err := config.OS()
if err != nil { return err }
r := config.New(lookup)
// Illustrative definitions only; not installed Atelier settings.
workers := r.Int("EXAMPLE_WORKERS", 2, 1, 16)
token := r.Secret("EXAMPLE_TOKEN")
if err := r.Err(); err != nil { return err }
// Only now may root construct resources using workers and token.
```

Reads return placeholders on failure to allow collecting all problems. Never use
those values before Err/check succeeds. Manifest also refuses when any problem
exists; successful manifests contain only declared settings and redact secrets.
They are not dumps of the environment. Reader is sequential; no global instance,
live reload, dotenv loading, remote resolution or ID generation.

Fifteen config tests pass per backend. Scenario mapping and integration evidence
are in project `work/handoffs/native-config.md`. Implementation and initial V01–V06
tests adapt n2f locally, with additional strictness, ownership and source tests;
there is no sibling dependency. Rust uses Cargo tests because config consumes
other leaves; the standalone std-only leaf runner does not apply.

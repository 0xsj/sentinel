# Process configuration reader

Revision: 1 · State: specified · Task: native-config

## Values and behavior

Lookup returns text plus explicit presence. Map snapshots own a copy of caller data.
OS capture snapshots the environment once; invalid native text is a classified
invalid / config.source failure with no raw key/value. It never resolves files or
remote secrets. Lookup callbacks are trusted synchronous adapters.

Reader exposes String(key, fallback), Required(key), Int(key, fallback, min, max),
Bool(key, fallback), Enum(key, fallback, allowed), Secret(key), Err/result and
Manifest. Native method casing follows the language. Values returned while the
reader has errors are not usable configuration.

- String uses a fallback only when absent. Present empty is an actual string.
- Required and Secret refuse absent or empty. Secret has no fallback.
- Int accepts only `-?(0|[1-9][0-9]*)`, with value in the inclusive supplied
  bounds and the common safe integer range ±9007199254740991. No whitespace,
  exponent, plus sign, leading zero, fractional or trailing text is accepted.
- Bool accepts exactly true/false text; empty is invalid. Enum accepts exactly one
  supplied choice; neither parser trims or case-folds.
- Defaults and bounds/choices are validated before a lookup. An invalid default,
  bounds or empty/duplicate choice list is invalid_definition, not operator input.
- Keys use `[A-Z][A-Z0-9_]*`. A key can be read once; a duplicate produces
  duplicate_key. Definitions and lookups are not lazy.
- All failures accumulate. Err returns invalid / config.invalid, message
  "invalid configuration", and copied key→reason field problems. Reasons are fixed
  required, invalid_integer, out_of_range, invalid_boolean, invalid_choice,
  invalid_definition, invalid_key, duplicate_key, invalid_source. Never echo values.
  Invalid key problems use a fixed `<key>` placeholder instead of the rejected key.
- Manifest refuses whenever Err is present. Otherwise it returns an owned,
  key-sorted list of {key,value,source,secret}, source environment/default.
  Secret values are always `[REDACTED]`. Integer/boolean entries use canonical text.
  String/enum entries reflect explicitly nonsecret settings selected by the root.
  This is not a dump of every environment variable.

No logger, root, provider SDK, global reader or registration dependency. Reader
uses secret and errors; leaf lookup/parser helpers need no project dependency.
A root must check all reads before constructing resources or printing a boot summary.
Retries, live reload, dotenv precedence and remote secret resolution are later
owners; this slice reads the supplied environment only.

## Comparable scenarios

| ID | Requirement |
| --- | --- |
| V01 | Absent, empty and set are distinct for String, Required and Secret. |
| V02 | Integers/booleans/enums accept exact supported text and refuse malformed/boundary input. |
| V03 | Invalid definitions refuse before lookup; fallback validation cannot be bypassed by a supplied value. |
| V04 | Collect independent problems without raw-value disclosure; any problem prevents a usable manifest. |
| V05 | Owned sorted manifest preserves public values and redacts secrets; input/output mutation cannot rewrite it. |
| V06 | Duplicate keys refuse; OS capture and map fixtures preserve snapshot semantics. |

These requirements precede tests/implementation. Compiler failures, executed
assertions and mutation evidence are recorded separately in module notes.

## Additional boundary rules

Reader construction is explicit and reads are sequential; it is not a concurrent
mutable settings store. Map/OS adapters capture once; custom callbacks are trusted
and may have their own behavior. There is no global reader, dotenv search, JSON
loader, reload watcher or process mutation. Root captures once, defines fields,
checks Err/check, then constructs resources; returned placeholders on failed reads
must never be used to start the app. Manifest success is a safe finalization gate.
No application settings are chosen by this module.

Problems are keyed by setting: independent keys accumulate, the latest refusal
for a repeated key replaces that key's reason. Errors remain sticky; a later valid
read does not clear an earlier invalid definition. Invalid definitions are checked
before source access. Once a syntactically valid read is attempted, a second read
of the key refuses without invoking lookup again. Invalid definitions/keys never
invoke lookup. A valid configuration with no reads has an empty successful manifest.

Bounds must be within the common safe range. Out-of-safe-range/overflow source
integers are invalid_integer; syntactically valid safe integers outside the field's
bounds are out_of_range. -0 is accepted and manifests as 0. Required/Secret reject
only absent or empty, not whitespace. Public String/Required/Enum values are
explicitly nonsecret by caller choice; the library cannot infer sensitivity.

An invalid Go key is replaced by <key>; a nil Go Lookup yields <source> with
invalid_source. OS capture fails closed on invalid UTF-8 native keys/values using
config.source and no raw key/value/cause. Custom Go string lookups are trusted;
Rust String inputs are already valid UTF-8. Map owns its input snapshot.

| ID | Further acceptance evidence |
| --- | --- |
| V07 | Invalid key, duplicate and invalid definition do not invoke lookup; errors stay sticky |
| V08 | Invalid native text capture is classified and redacted; source capture is a snapshot |
| V09 | Exact integer limits, -0 canonicalization and out-of-range reasons match Go/Rust |
| V10 | Empty manifest succeeds; returned problem maps cannot mutate future errors; nil Go lookup fails closed |

Tests may adapt n2f's scenarios locally, supplemented with boundary cases. No
sibling runtime imports or added external dependencies. Config consumes errors
revision 2 and secret revision 1; clocks and IDs are not dependencies. Its scheduled
consumer is native root configuration when real process settings are specified.

## Go API

Package config, under pkg/config:

```go
type Lookup func(string) (string, bool)
func Map(values map[string]string) Lookup
func OS() (Lookup, error)
type Var struct { Key, Value, Source string; Secret bool }
type Reader struct { /* private fields */ }
func New(lookup Lookup) *Reader
func (r *Reader) String(key, fallback string) string
func (r *Reader) Required(key string) string
func (r *Reader) Int(key string, fallback, min, max int64) int64
func (r *Reader) Bool(key string, fallback bool) bool
func (r *Reader) Enum(key, fallback string, allowed []string) string
func (r *Reader) Secret(key string) secret.Secret
func (r *Reader) Err() error
func (r *Reader) Manifest() ([]Var, error)
```

Construct Reader via New; its zero value is not usable. Failures use local
errors.Failure. Zero-value read placeholders are empty string/secret, 0 or false.
Only Err()==nil or successful Manifest authorizes using the collected values.

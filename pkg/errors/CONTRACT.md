# Native errors: classification, ownership and public projection

Revision: 2 · State: specified · Task: native-errors

Supersedes [revision 1](CONTRACT.v1.md). The user selected n2f's responsibility
split and condition-preservation semantics, mirrored idiomatically in Go/Rust.
Coordinator revision before implementation; no native bridge encoding is selected.

## Changes from revision 1

`New` / `Failure::new` now receives an explicitly public-safe message for
non-Internal kinds. Internal text remains diagnostic-only. Private context belongs
in details and causes. The existing Diagnostic/diagnostic accessor and Error/Display
retain the exact constructor message, including empty text, and do not append a
source. They are diagnostic APIs, never the frontend projection. Existing v1
call sites must be audited for sensitive text before public projection is used;
currently only leaf tests call these constructors.

Go derivation now preserves condition identity (not content equality):
`errors.Is(template.WithCause(cause), template)` is true. Rust preserves local
cases through typed Context<E>, without introducing pointer identity or requiring
Failure/source equality or cloning. Revision 1's prohibition on metadata and public
projection is replaced by the APIs below. Retry and transport policy remain excluded.

## Shared promises

Ten exact lowercase kinds remain unchanged. Parsing never trims or folds case.
A condition identifier is an optional module-owned string, e.g.
`workspace.name_taken`. It is not a new kind or a Rust concrete type name. Empty
identifier means absent. Public field problems and private details are independent
string maps; merges retain unrelated keys and incoming values replace collisions.
No global error registry or cross-domain imports.

The OUTERMOST explicit classification owns kind, message, type and fields together.
Do not inherit missing outer metadata from a cause. Context annotation preserves
classification; explicit translation replaces it and retains diagnostic source.
Unclassified failures stay unclassified until a boundary explicitly translates
them. Internal and unknown public projections are exactly `internal error`, with
Internal kind, no identifier and no fields. Empty messages on other classified
kinds project `request failed`. Private details and sources never appear in public
snapshots. Public-safe content is a caller contract, not an automatic sanitizer.

Public snapshots are independent owned values. Builders cannot mutate their
inputs or earlier Go errors; Rust consuming builders transfer ownership. Foreign
cause objects are opaque, not deeply immutable. Replacing a source preserves all
outer classification and metadata, not the previous source. Public projection is
transport-neutral: no JSON envelope, HTTP status, retry-after or bridge registration.

## Go public API

Keep Failure, New, Kind, ParseKind, Kinds, Kind(), Diagnostic(), Error(), Unwrap()
and WithCause() from revision 1, including numeric-kind normalization and zero
Failure behavior. Add:

```go
func (f *Failure) Is(target error) bool
func (f *Failure) WithType(value string) *Failure
func (f *Failure) WithField(key, value string) *Failure
func (f *Failure) WithFields(values map[string]string) *Failure
func (f *Failure) WithDetail(key, value string) *Failure
func (f *Failure) WithDetails(values map[string]string) *Failure
func Wrap(err error, kind Kind, message string) error
func KindOf(err error) (Kind, bool)
func IsKind(err error, kind Kind) bool
func DiagnosticTypeOf(err error) string
func DetailsOf(err error) map[string]string

type PublicInfo struct {
    Kind Kind
    Message string
    Type string // empty means absent
    Fields map[string]string
}
func Public(err error) (PublicInfo, bool)
func Message(err error) string
func TypeOf(err error) string
func FieldsOf(err error) map[string]string
```

With* methods clone metadata and preserve the original condition identity.
Separate New calls never match by contents. WithCause(nil) constructs a derived
failure with no cause; Wrap(nil, ...) returns a literal nil error. Typed-nil causes
are retained without invoking methods. Nil *Failure receivers remain outside the
instance API contract; inspectors must recognize typed nil safely as unknown.

Inspection follows single Unwrap() error edges and selects the first *Failure.
Stop before Unwrap() []error aggregates unless an explicit outer Failure was
already selected. Recognize standard context.Canceled and context.DeadlineExceeded
as Canceled and Timeout. Do not infer from strings. At most 128 frames are
examined: a cycle or deeper chain without classification returns unclassified.
Arbitrary foreign Unwrap panics are not recovered by this API.

KindOf reports classification; nil and unknown both return Internal,false.
Public reports failure presence: only literal nil gives false; unknown/typed-nil
errors give the Internal fallback,true. For standard context errors, public message
is request failed and metadata is empty. DetailsOf and DiagnosticTypeOf inspect
only the selected private frame (including Internal); they copy maps. Public
accessors share Public's disclosure policy and return zero values on success.


## Acceptance scenarios

Retain v1 E01–E10 coverage except the explicitly revised API/identity exclusions.
New cross-language scenarios:

| ID | Observable promise |
| --- | --- |
| P01 | Non-Internal public view contains the selected kind/message/identifier/fields; details/source excluded |
| P02 | Internal and unclassified errors expose only the fallback, even with private message/type/fields |
| P03 | Empty identifier is absent; empty non-Internal message uses request failed |
| P04 | Input maps, derived maps and independent public snapshots cannot mutate stored metadata |
| P05 | Metadata merges retain unrelated keys, incoming collision wins; public fields and private details stay separate |
| P06 | Context preserves classification, condition/local variant and its source |
| P07 | Explicit outer translation retains the source but never inherits its public metadata or private details |
| P08 | Replacing a source preserves all outer metadata and replaces only the source |
| P09 | Actual unknown error differs from success/absence; aggregate without explicit summary is unclassified regardless of item order |
| P10 | Cancellation/timeout survive context; text alone never classifies a foreign error |

G01: derivations preserve sentinel identity while independent New calls do not.
G02: nil/typed nil, context sentinels, aggregate ordering and cyclic unwrap are
handled as specified. G03: all public/private map-returning accessors are isolated.


No new external dependencies, logging, time, ID generation, implicit retry or
native SDK imports. Domains may use the shared vocabulary, but own their cases.
Return errors through intermediate layers; observation belongs to the handling
boundary. Tests are implementation-visible. Worker handoff maps old and new
scenarios, records exact commands, API review and limitations.

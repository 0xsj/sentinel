// Package errors provides native classification, condition identity and safe
// public projections. Domains own their condition sentinels; the shared Kind is
// the cross-layer vocabulary. Use fmt.Errorf with %w to annotate, and Wrap to
// deliberately translate at an adapter boundary. Builders preserve sentinel
// identity and copy metadata. Private details and causes never enter PublicInfo.
// Constructor messages and field problems for non-Internal kinds must be public-safe.
// Error and Diagnostic are diagnostic APIs; transports must use Public instead.
// No JSON, host SDK, retry or logging policy belongs in this leaf.
package errors

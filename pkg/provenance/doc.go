// Package provenance describes logical work and its individual executions.
// Values own their snapshots; factories consume explicit clock and ID ports.
// Attribution is descriptive, never authority or proof of a committed result.
// Factory calls require caller-owned serialization, including shared dependencies.
// See CONTRACT.md for transitions, validation, effects and limits.
package provenance

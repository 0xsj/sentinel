// Package logger provides synchronous structured diagnostics with swappable sinks.
package logger

import (
	faults "github.com/0xsj/atelier-wails/pkg/errors"
	"github.com/0xsj/atelier-wails/pkg/secret"
	"sync"
	"time"
)

type Level uint8

const (
	Debug Level = iota
	Info
	Warn
	Error
)

func (l Level) String() string {
	if l > Error {
		return "unknown"
	}
	return [...]string{"debug", "info", "warn", "error"}[l]
}
func ParseLevel(s string) (Level, error) {
	for l := Debug; l <= Error; l++ {
		if l.String() == s {
			return l, nil
		}
	}
	return 0, invalid()
}

// Value admits scalars only; arbitrary object formatting is deliberately absent.
type Value struct {
	kind   uint8
	text   string
	number int64
	flag   bool
}

func Text(s string) Value          { return Value{text: s} }
func Int(n int64) Value            { return Value{kind: 1, number: n} }
func Bool(b bool) Value            { return Value{kind: 2, flag: b} }
func Secret(_ secret.Secret) Value { return Text("[REDACTED]") }
func (v Value) Scalar() any {
	switch v.kind {
	case 1:
		return v.number
	case 2:
		return v.flag
	default:
		return v.text
	}
}

type Fields map[string]Value

func copyFields(f Fields) Fields {
	out := Fields{}
	for k, v := range f {
		out[k] = v
	}
	return out
}
func merge(a, b Fields) Fields {
	out := copyFields(a)
	for k, v := range b {
		out[k] = v
	}
	return out
}

// Record owns each map. A sink may retain it; it shares no mutable logger state.
type Record struct {
	Time                              time.Time
	Level                             Level
	Message, Service                  string
	Fields, Scope, Error, ErrorFields Fields
}

func (r Record) clone() Record {
	r.Fields = copyFields(r.Fields)
	r.Scope = copyFields(r.Scope)
	r.Error = copyFields(r.Error)
	r.ErrorFields = copyFields(r.ErrorFields)
	return r
}

type Clock interface{ Now() time.Time }

// Sink implementations must support concurrent calls. Enabled must be effect-free.
// Write and Flush return delivery errors; callers own sink lifetime.
type Sink interface {
	Enabled(Level) bool
	Write(Record) error
	Flush() error
}
type state struct {
	mu    sync.Mutex
	clock Clock
	sink  Sink
}
type Logger struct {
	state                               *state
	floor                               Level
	service                             string
	fields, scope, failure, errorFields Fields
}

func invalid() error {
	return faults.New(faults.Invalid, "invalid logger configuration").WithType("logger.invalid_configuration")
}
func New(c Clock, s Sink, service string, floor Level) (*Logger, error) {
	if c == nil || s == nil || service == "" || floor > Error {
		return nil, invalid()
	}
	return &Logger{state: &state{clock: c, sink: s}, floor: floor, service: service}, nil
}
func (l *Logger) With(f Fields) *Logger { n := *l; n.fields = merge(l.fields, f); return &n }
func (l *Logger) Enabled(level Level) bool {
	return level <= Error && level >= l.floor && l.state.sink.Enabled(level)
}
func delivery(e error) error {
	if e == nil {
		return nil
	}
	return faults.New(faults.Unavailable, "logger output unavailable").WithType("logger.sink").WithCause(e)
}
func (l *Logger) Log(level Level, message string, fields Fields) error {
	if level > Error {
		return invalid()
	}
	if !l.Enabled(level) {
		return nil
	}
	l.state.mu.Lock()
	defer l.state.mu.Unlock()
	t := l.state.clock.Now()
	if t.Before(time.UnixMilli(0)) || !t.Before(time.UnixMilli(1<<48)) {
		return faults.New(faults.Invalid, "invalid logger time").WithType("logger.invalid_time")
	}
	r := Record{Time: time.UnixMilli(t.UnixMilli()).UTC(), Level: level, Message: message, Service: l.service, Fields: merge(l.fields, fields), Scope: copyFields(l.scope), Error: copyFields(l.failure), ErrorFields: copyFields(l.errorFields)}
	return delivery(l.state.sink.Write(r))
}
func (l *Logger) Flush() error {
	l.state.mu.Lock()
	defer l.state.mu.Unlock()
	return delivery(l.state.sink.Flush())
}

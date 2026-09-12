package logger

import (
	"context"
	"io"
	"log/slog"
	"sort"
	"sync"
)

func nativeLevel(l Level) slog.Level {
	return [...]slog.Level{slog.LevelDebug, slog.LevelInfo, slog.LevelWarn, slog.LevelError}[l]
}

type slogSink struct {
	mu      sync.Mutex
	handler slog.Handler
	flush   func() error
}

// NewSlog adapts an injected handler. nil flush means no buffered output is owned.
func NewSlog(h slog.Handler, flush func() error) (Sink, error) {
	if h == nil {
		return nil, invalid()
	}
	return &slogSink{handler: h, flush: flush}, nil
}
func (s *slogSink) Enabled(l Level) bool {
	return l <= Error && s.handler.Enabled(context.Background(), nativeLevel(l))
}
func attrs(f Fields) []slog.Attr {
	keys := make([]string, 0, len(f))
	for k := range f {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	out := make([]slog.Attr, 0, len(keys))
	for _, k := range keys {
		out = append(out, slog.Any(k, f[k].Scalar()))
	}
	return out
}
func (s *slogSink) Write(r Record) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	record := slog.NewRecord(r.Time, nativeLevel(r.Level), r.Message, 0)
	record.AddAttrs(slog.String("service", r.Service), slog.GroupAttrs("fields", attrs(r.Fields)...), slog.GroupAttrs("scope", attrs(r.Scope)...), slog.GroupAttrs("error", attrs(r.Error)...), slog.GroupAttrs("error_fields", attrs(r.ErrorFields)...))
	return s.handler.Handle(context.Background(), record)
}
func (s *slogSink) Flush() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.flush != nil {
		return s.flush()
	}
	return nil
}

// NewConsole borrows the writer; it never closes stderr or caller-owned files.
func NewConsole(w io.Writer) (Sink, error) {
	if w == nil {
		return nil, invalid()
	}
	var flush func() error
	if f, ok := w.(interface{ Flush() error }); ok {
		flush = f.Flush
	}
	return NewSlog(slog.NewTextHandler(w, &slog.HandlerOptions{Level: slog.LevelDebug}), flush)
}

type Discard struct{}

func (Discard) Enabled(Level) bool { return false }
func (Discard) Write(Record) error { return nil }
func (Discard) Flush() error       { return nil }

type Memory struct {
	mu      sync.Mutex
	records []Record
}

func (*Memory) Enabled(Level) bool { return true }
func (m *Memory) Write(r Record) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.records = append(m.records, r.clone())
	return nil
}
func (*Memory) Flush() error { return nil }
func (m *Memory) Records() []Record {
	m.mu.Lock()
	defer m.mu.Unlock()
	out := make([]Record, len(m.records))
	for i, r := range m.records {
		out[i] = r.clone()
	}
	return out
}

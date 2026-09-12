package id

import (
	"crypto/rand"
	"io"
	"sync"
	"time"

	faults "github.com/0xsj/atelier-wails/pkg/errors"
)

// WallClock is owned by this consumer; elapsed time is not required.
type WallClock interface{ Now() time.Time }

// V7 serializes its own dependencies and state. Construct it with NewV7 or
// NewV7WithEntropy; its zero value is unusable. Do not copy it after first use.
// Dependencies must not reenter this generator and must be valid non-nil values.
type V7 struct {
	mu          sync.Mutex
	clock       WallClock
	entropy     io.Reader
	initialized bool
	last        int64
	counter     uint16
}

func NewV7(clock WallClock) (*V7, error) { return NewV7WithEntropy(clock, rand.Reader) }

// NewV7WithEntropy requires a cryptographic reader in production. It checks nil
// interfaces; a typed-nil implementation remains a caller configuration mistake.
func NewV7WithEntropy(clock WallClock, entropy io.Reader) (*V7, error) {
	if clock == nil || entropy == nil {
		return nil, faults.New(faults.Invalid, "invalid ID generator configuration").WithType("id.configuration")
	}
	return &V7{clock: clock, entropy: entropy}, nil
}

func (g *V7) NewID() (ID, error) {
	g.mu.Lock()
	defer g.mu.Unlock()
	wall := g.clock.Now()
	if wall.Before(time.Unix(0, 0)) || !wall.Before(time.UnixMilli(1<<48)) {
		return ID{}, faults.New(faults.Invalid, "ID timestamp out of range").WithType("id.time_range")
	}
	ms := wall.UnixMilli()
	fresh := !g.initialized || ms > g.last
	var counter uint16
	if !fresh {
		ms = g.last
		if g.counter == 4095 {
			return ID{}, faults.New(faults.Unavailable, "ID counter exhausted").WithType("id.exhausted")
		}
		counter = g.counter + 1
	}
	var random [10]byte
	if _, e := io.ReadFull(g.entropy, random[:]); e != nil {
		return ID{}, faults.New(faults.Unavailable, "ID entropy unavailable").WithType("id.entropy").WithCause(e)
	}
	if fresh {
		counter = (uint16(random[0])<<8 | uint16(random[1])) & 0x07ff
	}
	v := encodeV7(ms, counter, random[2:])
	g.last = ms
	g.counter = counter
	g.initialized = true
	return v, nil
}

// Package clock supplies explicit wall-clock reads. Pure domain transitions
// receive timestamps as values; they do not read a clock themselves.
package clock

import "time"

type Clock interface{ Now() time.Time }

type System struct{}

func (System) Now() time.Time { return time.Now().UTC() }

// Fixed holds one immutable instant. Its zero value returns Go's zero time.
type Fixed struct{ at time.Time }

func NewFixed(at time.Time) Fixed { return Fixed{at: at.UTC()} }
func (f Fixed) Now() time.Time    { return f.at }

package id

import (
	"sync"

	faults "github.com/0xsj/atelier-wails/pkg/errors"
)

// Sequence is a finite deterministic fixture generator. Duplicates are allowed.
// Do not copy a used Sequence; calls are serialized, but caller return order is not.
type Sequence struct {
	mu   sync.Mutex
	ids  []ID
	next int
}

func NewSequence(ids ...ID) (*Sequence, error) {
	for _, v := range ids {
		if v.IsZero() {
			return nil, invalidID()
		}
	}
	return &Sequence{ids: append([]ID(nil), ids...)}, nil
}
func (s *Sequence) NewID() (ID, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.next >= len(s.ids) {
		return ID{}, faults.New(faults.Unavailable, "ID sequence exhausted").WithType("id.sequence_exhausted")
	}
	v := s.ids[s.next]
	s.next++
	return v, nil
}

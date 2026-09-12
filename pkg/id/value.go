package id

import (
	"encoding/hex"
	"time"

	faults "github.com/0xsj/atelier-wails/pkg/errors"
)

// ID is a comparable UUID value. Only its zero (absence) value is unvalidated.
type ID struct{ bytes [16]byte }

func (v ID) IsZero() bool { return v == ID{} }
func (v ID) Version() int { return int(v.bytes[6] >> 4) }

// Time extracts UUIDv7's approximate generator timestamp, not event occurrence.
func (v ID) Time() (time.Time, bool) {
	if v.Version() != 7 {
		return time.Time{}, false
	}
	var ms int64
	for _, b := range v.bytes[:6] {
		ms = ms<<8 | int64(b)
	}
	return time.UnixMilli(ms).UTC(), true
}
func (v ID) String() string {
	var out [36]byte
	hex.Encode(out[0:8], v.bytes[0:4])
	out[8] = '-'
	hex.Encode(out[9:13], v.bytes[4:6])
	out[13] = '-'
	hex.Encode(out[14:18], v.bytes[6:8])
	out[18] = '-'
	hex.Encode(out[19:23], v.bytes[8:10])
	out[23] = '-'
	hex.Encode(out[24:36], v.bytes[10:16])
	return string(out[:])
}
func invalidID() error { return faults.New(faults.Invalid, "invalid ID").WithType("id.invalid") }

// Parse accepts only canonical-width hex-and-dash syntax and variant 10.
func Parse(s string) (ID, error) {
	if len(s) != 36 || s[8] != '-' || s[13] != '-' || s[18] != '-' || s[23] != '-' {
		return ID{}, invalidID()
	}
	var compact [32]byte
	n := 0
	for i := range len(s) {
		if i == 8 || i == 13 || i == 18 || i == 23 {
			continue
		}
		compact[n] = s[i]
		n++
	}
	var v ID
	if _, e := hex.Decode(v.bytes[:], compact[:]); e != nil || v.bytes[8]&0xc0 != 0x80 {
		return ID{}, invalidID()
	}
	return v, nil
}
func (v ID) MarshalText() ([]byte, error) {
	if v.IsZero() {
		return nil, invalidID()
	}
	return []byte(v.String()), nil
}
func (v *ID) UnmarshalText(b []byte) error {
	parsed, e := Parse(string(b))
	if e != nil {
		return e
	}
	*v = parsed
	return nil
}

func encodeV7(ms int64, counter uint16, suffix []byte) ID {
	var v ID
	for i := 5; i >= 0; i-- {
		v.bytes[i] = byte(ms)
		ms >>= 8
	}
	v.bytes[6] = 0x70 | byte(counter>>8)
	v.bytes[7] = byte(counter)
	copy(v.bytes[8:], suffix)
	v.bytes[8] = (v.bytes[8] & 0x3f) | 0x80
	return v
}

package errors

// Kind is a native failure category. Its zero value is Internal.
type Kind uint8

const (
	Internal Kind = iota
	Unauthenticated
	Forbidden
	RateLimited
	Unavailable
	Timeout
	Canceled
	NotFound
	Invalid
	Conflict
)

func (k Kind) String() string {
	switch k {
	case Internal:
		return "internal"
	case Unauthenticated:
		return "unauthenticated"
	case Forbidden:
		return "forbidden"
	case RateLimited:
		return "rate_limited"
	case Unavailable:
		return "unavailable"
	case Timeout:
		return "timeout"
	case Canceled:
		return "canceled"
	case NotFound:
		return "not_found"
	case Invalid:
		return "invalid"
	case Conflict:
		return "conflict"
	default:
		return "unknown"
	}
}

// ParseKind accepts exact canonical names only.
func ParseKind(text string) (Kind, bool) {
	for _, kind := range Kinds() {
		if kind.String() == text {
			return kind, true
		}
	}
	return Internal, false
}

// Kinds returns caller-owned storage; enumeration order is not a contract.
func Kinds() []Kind {
	return []Kind{Internal, Unauthenticated, Forbidden, RateLimited, Unavailable, Timeout, Canceled, NotFound, Invalid, Conflict}
}

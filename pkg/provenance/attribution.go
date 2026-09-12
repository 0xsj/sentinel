package provenance

type AttributionSpec struct {
	Initiator, OnBehalfOf *Actor
	Tenant                *string
}
type Attribution struct{ data AttributionSpec }

func NewAttribution(s AttributionSpec) (Attribution, error) {
	if s.Initiator != nil && !s.Initiator.valid() {
		return Attribution{}, invalid("invalid_attribution")
	}
	if s.OnBehalfOf != nil && (s.Initiator == nil || !s.Initiator.named() || !s.OnBehalfOf.named() || *s.Initiator == *s.OnBehalfOf) {
		return Attribution{}, invalid("invalid_attribution")
	}
	if s.Tenant != nil && !identitySyntax.MatchString(*s.Tenant) {
		return Attribution{}, invalid("invalid_attribution")
	}
	return Attribution{AttributionSpec{clone(s.Initiator), clone(s.OnBehalfOf), clone(s.Tenant)}}, nil
}
func (a Attribution) Snapshot() AttributionSpec {
	return AttributionSpec{clone(a.data.Initiator), clone(a.data.OnBehalfOf), clone(a.data.Tenant)}
}

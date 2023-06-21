package normalizer

// Prepend creates a normalizer that strip the normalized string inplace.
type Prepend struct {
	Prepend string `json:"prepend"`
}

func NewPrepend(prepend string) *Prepend {
	return &Prepend{prepend}
}

// Implement Normalizer for Prepend
func (p *Prepend) Normalize(normalized *NormalizedString) (*NormalizedString, error) {
	if normalized.IsEmpty() {
		return nil, nil
	}

	return normalized.Prepend(p.Prepend), nil
}

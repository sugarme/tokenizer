package normalizer

type Strip struct {
	stripLeft  bool
	stripRight bool
}

func NewStrip(stripLeft, stripRight bool) (retVal Strip) {
	return Strip{
		stripLeft:  stripLeft,
		stripRight: stripRight,
	}
}

// Implement Normalizer interface for Strip:
// =========================================

func (s Strip) Normalize(normalized NormalizedString) (retVal NormalizedString, err error) {

	if s.stripLeft && s.stripRight {
		return normalized.Strip(), nil
	} else {
		if s.stripLeft {
			return normalized.LStrip(), nil
		}

		if s.stripRight {
			return normalized.RStrip(), nil
		}
	}

	return normalized, nil
}

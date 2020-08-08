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

func (s Strip) Normalize(normalized Normalized) (retVal Normalized, err error) {

	if s.stripLeft && s.stripRight {
		normalized.Strip()
	} else {
		if s.stripLeft {
			normalized.LStrip()
		}

		if s.stripRight {
			normalized.RStrip()
		}
	}

	return normalized, nil
}

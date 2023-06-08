package normalizer

import (
	"strings"

	"github.com/sugarme/tokenizer/spm"

	"github.com/rivo/uniseg"
)

func replace(transformations []ChangeMap, oldPart, newPart string) []ChangeMap {
	oldCount := len(strings.Split(oldPart, ""))
	newCount := len(strings.Split(newPart, ""))
	diff := newCount - oldCount

	// If just replacing characters, all changes should be == 0
	for _, c := range strings.Split(newPart, "") {
		t := ChangeMap{
			RuneVal: c,
			Changes: 0,
		}
		transformations = append(transformations, t)
	}

	// log.Printf("transformations: %+v - diff: %v\n", transformations, diff)

	n := len(transformations)
	switch {
	case diff > 0:
		// If adding some characters, the last diff characters should be == 1
		for i := n - 1; i >= diff || i >= 0; i-- {
			transformations[i].Changes = 1
		}
	case diff < 0:
		// If removing some characters, the last one should include the diff
		n := len(transformations)
		transformations[n-1].Changes += diff
	}

	return transformations
}

type Precompiled struct {
	*spm.Precompiled
}

// Implement Normalizer for spm.Precompiled
func (m *Precompiled) Normalize(normalized *NormalizedString) (*NormalizedString, error) {
	original := normalized.GetNormalized()
	var transformations []ChangeMap
	for i := 0; i < len(original); i++ {
		transformations = append(transformations, ChangeMap{})
	}

	graphemes := uniseg.NewGraphemes(original)

	var modified bool

	for graphemes.Next() {
		grapheme := graphemes.Str()

		if len(grapheme) < 6 {
			norm := m.Transform(grapheme)
			if len(norm) > 0 {
				modified = true
				transformations = replace(transformations, grapheme, norm)

				continue
			}
		}

		// TT. This is a hacky way to turn non-spacing marks into hexa string
		// and pass the unit tests.
		grapheme = spm.NormalizeMn(grapheme)
		var charIdx int = 0
		for _, r := range grapheme {
			part := string(grapheme)[charIdx : charIdx+len(string(r))]
			norm := m.Transform(part)
			norm = spm.NormalizeMn(norm)
			if len(norm) > 0 {
				modified = true
				transformations = replace(transformations, part, norm)
			} else {
				t := ChangeMap{
					RuneVal: string(r),
					Changes: 0,
				}
				transformations = append(transformations, t)
			}
			// charIdx += len(string(r))
		}
	}

	if modified {
		normalized = normalized.Transform(transformations, 0)
	}

	return normalized, nil
}

package decoder

import (
	"strings"
)

type Strip struct {
	*DecoderBase

	Content string
	Start   int
	Stop    int
}

func NewStrip(content string, start, stop int) *Strip {
	base := new(DecoderBase)

	return &Strip{
		DecoderBase: base,
		Content:     content,
		Start:       start,
		Stop:        stop,
	}
}

func (d *Strip) DecodeChain(tokens []string) []string {
	var toks []string

	for _, token := range tokens {
		chars := strings.Split(token, "")

		startCut := 0
		for i := 0; i < d.Start; i++ {
			c := chars[i]
			if c == d.Content {
				startCut = i + 1
				continue
			} else {
				break
			}
		}

		stopCut := len(chars)
		for i := 0; i < d.Stop; i++ {
			index := len(chars) - i - 1
			if chars[index] == d.Content {
				stopCut = index
				continue
			} else {
				break
			}
		}

		newToken := strings.Join(chars[startCut:stopCut], "")
		toks = append(toks, newToken)
	}

	return toks
}

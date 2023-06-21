package decoder

import (
	"encoding/hex"
	"strings"
	"unicode/utf8"

	"github.com/sugarme/tokenizer"
)

type ByteFallback struct {
	*DecoderBase

	typ string
}

func NewByteFallback() *ByteFallback {
	base := new(DecoderBase)

	d := &ByteFallback{
		DecoderBase: base,
		typ:         "ByteFallback",
	}

	d.DecoderBase.Decoder = interface{}(d).(tokenizer.Decoder)

	return d
}

var _ tokenizer.Decoder = new(ByteFallback)

func (d *ByteFallback) DecodeChain(tokens []string) []string {
	var (
		newTokens          []string
		previousByteTokens []byte
	)

	for _, token := range tokens {
		var bytes []byte
		var err error

		if len(token) == 6 && strings.HasPrefix(token, "<0x") && strings.HasSuffix(token, ">") {
			// convert hex string to bytes
			bytes, err = hex.DecodeString(token[3:5])
			if err != nil {
				panic(err)
			}
		}

		if len(bytes) > 0 {
			previousByteTokens = append(previousByteTokens, bytes...)
		} else {
			if len(previousByteTokens) > 0 {
				if utf8.Valid(previousByteTokens) {
					tok := string(previousByteTokens)
					newTokens = append(newTokens, tok)
				} else {
					for i := 0; i < len(previousByteTokens); i++ {
						newTokens = append(newTokens, "�")
					}
				}
				previousByteTokens = []byte{}
			}
			newTokens = append(newTokens, token)
		}
	}

	// last one
	if len(previousByteTokens) > 0 {
		if utf8.Valid(previousByteTokens) {
			tok := string(previousByteTokens)
			newTokens = append(newTokens, tok)
		} else {
			for i := 0; i < len(previousByteTokens); i++ {
				newTokens = append(newTokens, "�")
			}
		}
	}

	return newTokens
}

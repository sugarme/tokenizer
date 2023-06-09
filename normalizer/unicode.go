package normalizer

import (
	"golang.org/x/text/unicode/norm"
)

// Basic Unicode normal form composing and decomposing - NFC, NFD, NFKC, NFKD
// Ref. https://blog.golang.org/normalization

type UnicodeNormalizer struct {
	Form norm.Form
}

func NewUnicodeNormalizer(form norm.Form) *UnicodeNormalizer {
	return &UnicodeNormalizer{
		Form: form,
	}
}

func (un *UnicodeNormalizer) Normalize(n *NormalizedString) (*NormalizedString, error) {
	switch un.Form {
	case norm.NFC:
		return n.NFC(), nil
	case norm.NFD:
		return n.NFD(), nil
	case norm.NFKC:
		return n.NFKC(), nil
	case norm.NFKD:
		return n.NFKD(), nil
	}

	return n, nil
}

type NFC struct{}

func NewNFC() *NFC {
	return new(NFC)
}

func (n *NFC) Normalize(norm *NormalizedString) (*NormalizedString, error) {
	return norm.NFC(), nil
}

type NFKC struct{}

func NewNFKC() *NFKC {
	return new(NFKC)
}

func (n *NFKC) Normalize(norm *NormalizedString) (*NormalizedString, error) {
	return norm.NFKC(), nil
}

type NFD struct{}

func NewNFD() *NFD {
	return new(NFD)
}

func (n *NFD) Normalize(norm *NormalizedString) (*NormalizedString, error) {
	return norm.NFD(), nil
}

type NFKD struct{}

func NewNFKD() *NFKD {
	return new(NFKD)
}

func (n *NFKD) Normalize(norm *NormalizedString) (*NormalizedString, error) {
	return norm.NFKD(), nil
}

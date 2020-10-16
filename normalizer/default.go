// Basic text preprocessing tasks are:
// 1. Remove HTML tags
// 2. Remove extra whitespaces
// 3. Convert accented characters to ASCII characters
// 4. Expand contractions
// 5. Remove special characters
// 6. Lowercase all texts
// 7. Convert number words to numeric form
// 8. Remove numbers
// 9. Remove stopwords
// 10. Lemmatization
package normalizer

import (
	"encoding/csv"
	"encoding/json"
	"log"
	"os"
)

type DefaultNormalizer struct {
	lower bool // to lowercase
	strip bool // trim leading and trailing whitespaces
	// ExtraWhitespace bool // remove extra-whitespaces
	// Contraction     bool // expand contraction
}

type DefaultOption func(*DefaultNormalizer)

func WithLowercase(lowercase bool) DefaultOption {
	return func(o *DefaultNormalizer) {
		o.lower = lowercase
	}
}

func WithStrip(strip bool) DefaultOption {
	return func(o *DefaultNormalizer) {
		o.strip = strip
	}
}

/*
 * func WithContractionExpansion() DefaultOption {
 *   return func(o *DefaultNormalizer) {
 *     o.Contraction = true
 *   }
 * }
 *  */

func (dn *DefaultNormalizer) Normalize(n *NormalizedString) (*NormalizedString, error) {

	var normalized *NormalizedString = n

	if dn.lower {
		normalized = normalized.Lowercase()
	}

	if dn.strip {
		normalized = normalized.Strip()
	}

	return normalized, nil
}

func NewDefaultNormalizer(opts ...DefaultOption) *DefaultNormalizer {

	dn := DefaultNormalizer{
		lower: true,
		strip: true,
		// Contraction:     false,
	}

	for _, o := range opts {
		o(&dn)
	}

	return &dn

}

// TODO: move this func to `normalized` file
func expandContraction(txt string) string {
	var cMap map[string]string
	cMap, err := loadContractionMap()
	if err != nil {
		log.Fatal(err)
	}

	if k, ok := cMap[txt]; ok {
		return k
	}

	return txt
}

func loadContractionMap() (map[string]string, error) {
	const file = "contraction.csv"
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}

	defer f.Close()

	lines, err := csv.NewReader(f).ReadAll()
	if err != nil {
		return nil, err
	}

	type Contract struct {
		Contraction string
		Expansion   string
	}

	var cList []Contract

	for _, line := range lines {
		cList = append(cList, Contract{
			Contraction: line[0],
			Expansion:   line[1],
		})
	}

	cMap := map[string]string{}
	inrec, _ := json.Marshal(cList)
	json.Unmarshal(inrec, &cMap)

	return cMap, nil

}

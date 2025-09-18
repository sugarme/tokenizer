package pretokenizer

import (
	"log"
	"unicode"

	"github.com/sugarme/tokenizer"
	"github.com/sugarme/tokenizer/normalizer"
)

// GetScript returns key to script in `unicode.Scripts`.
func GetScript(r rune) string {
	for n, rt := range unicode.Scripts {
		if unicode.In(r, rt) {
			return n
		}
	}

	return ""
}

type UnicodeScript struct{}

func NewUnicodeScript() *UnicodeScript {
	return new(UnicodeScript)
}

func DefaultUnicodeScript() *UnicodeScript {
	return new(UnicodeScript)
}

func FixedScript(c rune) string {
	rawScript := GetScript(c)

	if uint32(c) == 0x30FC {
		return "Han"
	} else if c == ' ' {
		return "Any"
	} else {
		if rawScript == "Hiragana" || rawScript == "Katakana" {
			return "Han"
		} else {
			return rawScript
		}
	}
}

// Implemente tokenizer.PreTokenizer

var _ tokenizer.PreTokenizer = new(UnicodeScript)

func (us *UnicodeScript) PreTokenize(pretokenized *tokenizer.PreTokenizedString) (*tokenizer.PreTokenizedString, error) {
	pretok := pretokenized.Split(func(noop int, normalized *normalizer.NormalizedString) []tokenizer.SplitIdx {
		lastScript := ""
		offset := 0
		var ranges []int
		for _, c := range normalized.GetNormalized() {
			script := FixedScript(c)
			var result int
			if script != "Any" && lastScript != "Any" && lastScript != script {
				result = offset
			} else {
				result = -1
			}

			offset += len(string(c))
			if script != "Any" {
				lastScript = script
			}

			if result >= 0 {
				ranges = append(ranges, result)
			}
		}

		ranges = append(ranges, len(normalized.GetNormalized()))

		log.Printf("ranges: %v\n", ranges)

		var splits []normalizer.NormalizedString
		for i := 0; i < len(ranges)-1; i++ { // windows(2)
			var start, end int
			if i == 0 {
				start = 0
				end = 1
			} else {
				start = i
				end = i + 1
			}
			inputRange := normalizer.NewRange(ranges[start], ranges[end], normalizer.NormalizedTarget)
			log.Printf("inputRange: %+v\n", inputRange)
			split := normalized.Slice(inputRange)
			splits = append(splits, *split)
		}

		var splitIdxs []tokenizer.SplitIdx
		for _, s := range splits {
			normalized := s
			splitIdx := tokenizer.SplitIdx{Normalized: &normalized, Tokens: nil}
			splitIdxs = append(splitIdxs, splitIdx)
		}

		return splitIdxs
	})

	return pretok, nil
}

/*
            Ok(ranges
                .windows(2)
                .map(|item| {
                    normalized
                        .slice(Range::Normalized(item[0]..item[1]))
                        .expect("NormalizedString bad split")
                })
                .collect::<Vec<_>>())
impl PreTokenizer for UnicodeScripts {
    fn pre_tokenize(&self, pretokenized: &mut PreTokenizedString) -> Result<()> {
        pretokenized.split(|_, normalized| {
            let mut last_script = None;
            let mut offset = 0;
            let mut ranges: Vec<_> = normalized
                .get()
                .chars()
                .filter_map(|c| {
                    let script = Some(fixed_script(c));
                    let result = if script != Some(Script::Any)
                        && last_script != Some(Script::Any)
                        && last_script != script
                    {
                        Some(offset)
                    } else {
                        None
                    };
                    offset += c.len_utf8();
                    if script != Some(Script::Any) {
                        last_script = script;
                    }

                    result
                })
                .collect();
            ranges.push(normalized.get().len());
            Ok(ranges
                .windows(2)
                .map(|item| {
                    normalized
                        .slice(Range::Normalized(item[0]..item[1]))
                        .expect("NormalizedString bad split")
                })
                .collect::<Vec<_>>())
        })
    }
}
*/

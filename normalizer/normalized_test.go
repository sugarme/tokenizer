package normalizer_test

import (
	"fmt"
	"reflect"

	// "strings"
	"testing"
	"unicode"

	// "golang.org/x/text/transform"
	// "golang.org/x/text/unicode/norm"

	"github.com/sugarme/tokenizer/normalizer"
	// "github.com/sugarme/tokenizer/util"
)

func TestNormalized_NFDAddsNewChars(t *testing.T) {
	n := normalizer.NewNormalizedFrom("élégant").NFD()

	wantN := [][]int{{0, 2}, {0, 2}, {0, 2}, {2, 3}, {3, 5}, {3, 5}, {3, 5}, {5, 6}, {6, 7}, {7, 8}, {8, 9}}
	gotN := n.Alignments()

	wantO := [][]int{{0, 3}, {0, 3}, {3, 4}, {4, 7}, {4, 7}, {7, 8}, {8, 9}, {9, 10}, {10, 11}}
	gotO := n.AlignmentsOriginal()

	if !reflect.DeepEqual(wantN, gotN) {
		t.Errorf("Want normalized: %v\n", wantN)
		t.Errorf("Got normalized: %v\n", gotN)
	}

	if !reflect.DeepEqual(wantO, gotO) {
		t.Errorf("Want original: %v\n", wantO)
		t.Errorf("Got original: %v\n", gotO)
	}
}

func TestNormalized_RemoveCharsAddedByNFD(t *testing.T) {
	n := normalizer.NewNormalizedFrom("élégant").NFD()
	/*
	 *   n = n.Filter(func(r rune) bool {
	 *     return unicode.Is(unicode.Mn, r)
	 *   })
	 *  */
	n = n.RemoveAccents()
	wantN := [][]int{{0, 2}, {2, 3}, {3, 5}, {5, 6}, {6, 7}, {7, 8}, {8, 9}}
	gotN := n.Alignments()
	wantO := [][]int{{0, 1}, {0, 1}, {1, 2}, {2, 3}, {2, 3}, {3, 4}, {4, 5}, {5, 6}, {6, 7}}
	gotO := n.AlignmentsOriginal()

	if !reflect.DeepEqual(wantN, gotN) {
		t.Errorf("Want normalized: %v\n", wantN)
		t.Errorf("Got normalized: %v\n", gotN)
	}

	if !reflect.DeepEqual(wantO, gotO) {
		t.Errorf("Want original: %v\n", wantO)
		t.Errorf("Got original: %v\n", gotO)
	}
}

func TestNormalized_RemoveChars(t *testing.T) {
	n := normalizer.NewNormalizedFrom("élégant")

	n = n.Filter(func(r rune) bool {
		return r != 'n'
	})

	wantN := [][]int{{0, 2}, {0, 2}, {2, 3}, {3, 5}, {3, 5}, {5, 6}, {6, 7}, {8, 9}}
	gotN := n.Alignments()
	wantO := [][]int{{0, 2}, {0, 2}, {2, 3}, {3, 5}, {3, 5}, {5, 6}, {6, 7}, {7, 7}, {7, 8}}
	gotO := n.AlignmentsOriginal()

	if !reflect.DeepEqual(wantN, gotN) {
		t.Errorf("Want normalized: %v\n", wantN)
		t.Errorf("Got normalized: %v\n", gotN)
	}

	if !reflect.DeepEqual(wantO, gotO) {
		t.Errorf("Want original: %v\n", wantO)
		t.Errorf("Got original: %v\n", gotO)
	}
}

func TestNormalized_MixedAdditionRemoval(t *testing.T) {
	n := normalizer.NewNormalizedFrom("élégant").NFD()

	n = n.Filter(func(r rune) bool {
		return r != 'n' && !unicode.Is(unicode.Mn, r) // Mark non-spacing
	})

	wantN := [][]int{{0, 2}, {2, 3}, {3, 5}, {5, 6}, {6, 7}, {8, 9}}
	gotN := n.Alignments()
	wantO := [][]int{{0, 1}, {0, 1}, {1, 2}, {2, 3}, {2, 3}, {3, 4}, {4, 5}, {5, 5}, {5, 6}}
	gotO := n.AlignmentsOriginal()

	if !reflect.DeepEqual(wantN, gotN) {
		t.Errorf("Want normalized: %v\n", wantN)
		t.Errorf("Got normalized: %v\n", gotN)
	}

	if !reflect.DeepEqual(wantO, gotO) {
		t.Errorf("Want original: %v\n", wantO)
		t.Errorf("Got original: %v\n", gotO)
	}
}

func TestNormalized_RangeConversion(t *testing.T) {
	n := normalizer.NewNormalizedFrom("    __Hello__   ")
	n = n.Filter(func(r rune) bool {
		return r != ' '
	}).Lowercase()

	fmt.Printf("n: %+v\n", n)

	helloRange := n.ConvertOffset(normalizer.NewRange(6, 11, normalizer.OriginalTarget))
	gotRange := helloRange.Values()
	wantRange := []int{2, 7}

	gotN := n.Range(helloRange)
	wantN := "hello"
	gotO := n.RangeOriginal(helloRange)
	wantO := "Hello"

	if !reflect.DeepEqual(wantRange, gotRange) {
		t.Errorf("Want range: %v\n", wantRange)
		t.Errorf("Got range: %v\n", gotRange)
	}

	if !reflect.DeepEqual(wantN, gotN) {
		t.Errorf("Want normalized: %v\n", wantN)
		t.Errorf("Got normalized: %v\n", gotN)
	}

	if !reflect.DeepEqual(wantO, gotO) {
		t.Errorf("Want original: %v\n", wantO)
		t.Errorf("Got original: %v\n", gotO)
	}

	// Make sure we get None only in specific cases
	fmt.Printf("Len original: %v\n", n.LenOriginal())
	testRange(t, n, []int{0, 0}, []int{0, 0}, normalizer.OriginalTarget)
	testRange(t, n, []int{3, 3}, []int{3, 3}, normalizer.OriginalTarget)
	testRange(t, n, []int{15, n.LenOriginal()}, []int{9, 9}, normalizer.OriginalTarget)
	testRange(t, n, []int{16, n.LenOriginal() + 1}, []int{16, 16}, normalizer.OriginalTarget)
	testRange(t, n, []int{17, n.LenOriginal() + 1}, nil, normalizer.OriginalTarget)
	testRange(t, n, []int{0, 0}, []int{0, 0}, normalizer.NormalizedTarget)
	testRange(t, n, []int{3, 3}, []int{3, 3}, normalizer.NormalizedTarget)
	testRange(t, n, []int{9, n.Len() + 1}, []int{9, 9}, normalizer.NormalizedTarget)
	testRange(t, n, []int{10, n.Len() + 1}, nil, normalizer.NormalizedTarget)
}

func testRange(t *testing.T, n *normalizer.NormalizedString, input, wantR []int, indexOn normalizer.IndexOn) {
	gotR := n.ConvertOffset(normalizer.NewRange(input[0], input[1], indexOn)).Values()
	if !reflect.DeepEqual(wantR, gotR) {
		t.Errorf("Want range: %v\n", wantR)
		t.Errorf("Got range: %v\n", gotR)
	}
}

func TestNormalized_AddedAroundEdge(t *testing.T) {
	n := normalizer.NewNormalizedFrom("Hello")

	changeMap := []normalizer.ChangeMap{
		{" ", 1},
		{"H", 0},
		{"e", 0},
		{"l", 0},
		{"l", 0},
		{"o", 0},
		{" ", 1},
	}

	n.Transform(changeMap, 0)

	want := " Hello "
	got := n.GetNormalized()
	if !reflect.DeepEqual(want, got) {
		t.Errorf("Want: %v\n", want)
		t.Errorf("Got: %v\n", got)
	}
}

func TestNormalized_AddedCharactersAlignment(t *testing.T) {
	n := normalizer.NewNormalizedFrom("野口 No")

	chars := []rune(n.GetNormalized())
	var changeMap []normalizer.ChangeMap
	for _, r := range chars {
		if r > 0x4E00 {
			changeMap = append(changeMap, []normalizer.ChangeMap{
				{" ", 0},
				{string(r), 1},
				{" ", 1},
			}...)
		} else {
			changeMap = append(changeMap, normalizer.ChangeMap{string(r), 0})
		}
	}

	n.Transform(changeMap, 0)

	original := "野口 No"
	normalized := " 野  口  No"
	alignments := [][]int{{0, 3}, {0, 3}, {0, 3}, {0, 3}, {0, 3}, {3, 6}, {3, 6}, {3, 6}, {3, 6}, {3, 6}, {6, 7}, {7, 8}, {8, 9}}
	alignmentsOriginal := [][]int{{0, 5}, {0, 5}, {0, 5}, {5, 10}, {5, 10}, {5, 10}, {10, 11}, {11, 12}, {12, 13}}
	originalShift := 0

	want := normalizer.NewNormalizedString(original, normalized, alignments, alignmentsOriginal, originalShift)
	got := n
	if !reflect.DeepEqual(want, got) {
		t.Errorf("Want: %v\n", want)
		t.Errorf("Got: %v\n", got)
	}
}

func TestNormalized_RemoveAtBeginning(t *testing.T) {
	n := normalizer.NewNormalizedFrom("     Hello")
	n.Filter(func(r rune) bool {
		return r != ' '
	})

	got1 := n.RangeOriginal(normalizer.NewRange(1, len("Hello"), normalizer.NormalizedTarget))
	want1 := "ello"

	got2 := n.RangeOriginal(normalizer.NewRange(0, len(n.GetNormalized()), normalizer.NormalizedTarget))
	want2 := "Hello"

	if !reflect.DeepEqual(want1, got1) {
		t.Errorf("Want: %v\n", want1)
		t.Errorf("Got: %v\n", got1)
	}

	if !reflect.DeepEqual(want2, got2) {
		t.Errorf("Want: %v\n", want2)
		t.Errorf("Got: %v\n", got2)
	}
}

func TestNormalized_RemoveAtEnd(t *testing.T) {
	n := normalizer.NewNormalizedFrom("Hello    ")
	n.Filter(func(r rune) bool {
		return r != ' '
	})

	got1 := n.RangeOriginal(normalizer.NewRange(0, 4, normalizer.NormalizedTarget))
	want1 := "Hell"

	got2 := n.RangeOriginal(normalizer.NewRange(0, len(n.GetNormalized()), normalizer.NormalizedTarget))
	want2 := "Hello"

	if !reflect.DeepEqual(want1, got1) {
		t.Errorf("Want: %v\n", want1)
		t.Errorf("Got: %v\n", got1)
	}

	if !reflect.DeepEqual(want2, got2) {
		t.Errorf("Want: %v\n", want2)
		t.Errorf("Got: %v\n", got2)
	}
}

func TestNormalized_AroundBothEdges(t *testing.T) {
	n := normalizer.NewNormalizedFrom("  Hello  ")
	n.Filter(func(r rune) bool {
		return r != ' '
	})

	got0 := n.GetNormalized()
	want0 := "Hello"

	got1 := n.RangeOriginal(normalizer.NewRange(0, len("Hello"), normalizer.NormalizedTarget))
	want1 := "Hello"

	got2 := n.RangeOriginal(normalizer.NewRange(1, len("Hell"), normalizer.NormalizedTarget))
	want2 := "ell"

	if !reflect.DeepEqual(want0, got0) {
		t.Errorf("Want: %v\n", want0)
		t.Errorf("Got: %v\n", got0)
	}

	if !reflect.DeepEqual(want1, got1) {
		t.Errorf("Want: %v\n", want1)
		t.Errorf("Got: %v\n", got1)
	}

	if !reflect.DeepEqual(want2, got2) {
		t.Errorf("Want: %v\n", want2)
		t.Errorf("Got: %v\n", got2)
	}
}

func TestNormalized_LStrip(t *testing.T) {
	n := normalizer.NewNormalizedFrom("  This is an example  ")
	n.LStrip()

	got0 := n.GetNormalized()
	want0 := "This is an example  "

	got1 := n.RangeOriginal(normalizer.NewRange(0, len(n.GetNormalized()), normalizer.NormalizedTarget))
	want1 := "This is an example  "

	if !reflect.DeepEqual(want0, got0) {
		t.Errorf("Want: %v\n", want0)
		t.Errorf("Got: %v\n", got0)
	}

	if !reflect.DeepEqual(want1, got1) {
		t.Errorf("Want: %v\n", want1)
		t.Errorf("Got: %v\n", got1)
	}
}

func TestNormalized_RStrip(t *testing.T) {
	n := normalizer.NewNormalizedFrom("  This is an example  ")
	n.RStrip()

	got0 := n.GetNormalized()
	want0 := "  This is an example"

	got1 := n.RangeOriginal(normalizer.NewRange(0, len(n.GetNormalized()), normalizer.NormalizedTarget))
	want1 := "  This is an example"

	if !reflect.DeepEqual(want0, got0) {
		t.Errorf("Want: %v\n", want0)
		t.Errorf("Got: %v\n", got0)
	}

	if !reflect.DeepEqual(want1, got1) {
		t.Errorf("Want: %v\n", want1)
		t.Errorf("Got: %v\n", got1)
	}
}

func TestNormalized_Strip(t *testing.T) {
	n := normalizer.NewNormalizedFrom("  This is an example  ")
	n.Strip()

	got0 := n.GetNormalized()
	want0 := "This is an example"

	got1 := n.RangeOriginal(normalizer.NewRange(0, len(n.GetNormalized()), normalizer.NormalizedTarget))
	want1 := "This is an example"

	if !reflect.DeepEqual(want0, got0) {
		t.Errorf("Want: %v\n", want0)
		t.Errorf("Got: %v\n", got0)
	}

	if !reflect.DeepEqual(want1, got1) {
		t.Errorf("Want: %v\n", want1)
		t.Errorf("Got: %v\n", got1)
	}
}

func TestNormalized_Prepend(t *testing.T) {
	n := normalizer.NewNormalizedFrom("there")
	n.Prepend("Hey ")

	got0 := n.Alignments()
	want0 := [][]int{{0, 1}, {0, 1}, {0, 1}, {0, 1}, {0, 1}, {1, 2}, {2, 3}, {3, 4}, {4, 5}}

	got1 := n.ConvertOffset(normalizer.NewRange(0, 4, normalizer.NormalizedTarget)).Values()
	want1 := []int{0, 1}

	got2 := n.GetNormalized()
	want2 := "Hey there"

	if !reflect.DeepEqual(want0, got0) {
		t.Errorf("Want: %v\n", want0)
		t.Errorf("Got: %v\n", got0)
	}

	if !reflect.DeepEqual(want1, got1) {
		t.Errorf("Want: %v\n", want1)
		t.Errorf("Got: %v\n", got1)
	}

	if !reflect.DeepEqual(want2, got2) {
		t.Errorf("Want: %v\n", want2)
		t.Errorf("Got: %v\n", got2)
	}
}

func TestNormalized_Append(t *testing.T) {
	n := normalizer.NewNormalizedFrom("Hey")
	n.Append(" there")

	got0 := n.Alignments()
	want0 := [][]int{{0, 1}, {1, 2}, {2, 3}, {2, 3}, {2, 3}, {2, 3}, {2, 3}, {2, 3}, {2, 3}}

	got1 := n.ConvertOffset(normalizer.NewRange(3, len(" there"), normalizer.NormalizedTarget)).Values()
	want1 := []int{2, 3}

	got2 := n.GetNormalized()
	want2 := "Hey there"

	if !reflect.DeepEqual(want0, got0) {
		t.Errorf("Want: %v\n", want0)
		t.Errorf("Got: %v\n", got0)
	}

	if !reflect.DeepEqual(want1, got1) {
		t.Errorf("Want: %v\n", want1)
		t.Errorf("Got: %v\n", got1)
	}

	if !reflect.DeepEqual(want2, got2) {
		t.Errorf("Want: %v\n", want2)
		t.Errorf("Got: %v\n", got2)
	}
}

func TestNormalized_GetRange(t *testing.T) {
	s := "Hello my name is John 👋"
	start := 0
	end := len(s)

	got0 := normalizer.RangeOf(s, []int{start, end})
	want0 := s

	got1 := normalizer.RangeOf(s, []int{17, end})
	want1 := "John 👋"

	if !reflect.DeepEqual(want0, got0) {
		t.Errorf("Want: %v\n", want0)
		t.Errorf("Got: %v\n", got0)
	}

	if !reflect.DeepEqual(want1, got1) {
		t.Errorf("Want: %v\n", want1)
		t.Errorf("Got: %v\n", got1)
	}
}

func TestNormalized_Slice(t *testing.T) {
	n := normalizer.NewNormalizedFrom("𝔾𝕠𝕠𝕕 𝕞𝕠𝕣𝕟𝕚𝕟𝕘")
	n = n.NFKC()

	oSlice := n.Slice(normalizer.NewRange(0, 4, normalizer.OriginalTarget))
	got1 := oSlice.GetNormalized()
	want1 := "G"
	got2 := oSlice.GetOriginal()
	want2 := "𝔾"
	testSlice(t, want1, got1)
	testSlice(t, want2, got2)

	nSlice := n.Slice(normalizer.NewRange(0, 4, normalizer.NormalizedTarget))
	got3 := nSlice.GetNormalized()
	want3 := "Good"
	got4 := nSlice.GetOriginal()
	want4 := "𝔾𝕠𝕠𝕕"
	testSlice(t, want3, got3)
	testSlice(t, want4, got4)

	// Make sure the sliced NormalizedString is still aligned as expected
	n1 := normalizer.NewNormalizedFrom("   Good Morning!   ")
	n1 = n1.Strip()

	// 1. If we keep the whole slice
	sliceO := n1.Slice(normalizer.NewRange(0, len(n1.GetOriginal()), normalizer.OriginalTarget))
	got5 := sliceO.RangeOriginal(normalizer.NewRange(0, 4, normalizer.NormalizedTarget))
	want5 := "Good"
	testSlice(t, want5, got5)

	sliceN := n1.Slice(normalizer.NewRange(0, len(n1.GetOriginal()), normalizer.NormalizedTarget))
	got6 := sliceN.RangeOriginal(normalizer.NewRange(0, 4, normalizer.OriginalTarget))
	want6 := "Good"
	testSlice(t, want6, got6)

	// 2. If we keep after the modified piece
	sliceAM := n1.Slice(normalizer.NewRange(4, 15, normalizer.OriginalTarget))
	got7 := sliceAM.RangeOriginal(normalizer.NewRange(0, 3, normalizer.NormalizedTarget))
	want7 := "ood"
	testSlice(t, want7, got7)

	// 3. If we keep only the modified piece
	sliceM := n1.Slice(normalizer.NewRange(3, 16, normalizer.OriginalTarget))
	got8 := sliceM.RangeOriginal(normalizer.NewRange(0, 4, normalizer.NormalizedTarget))
	want8 := "Good"
	testSlice(t, want8, got8)
}

func testSlice(t *testing.T, want, got interface{}) {
	if !reflect.DeepEqual(want, got) {
		t.Errorf("Want: %v\n", want)
		t.Errorf("Got: %v\n", got)
	}
}

func TestNormalized_Replace(t *testing.T) {
	// Simple
	n1 := normalizer.NewNormalizedFrom(" Hello   friend ")
	n1 = n1.Replace(normalizer.NewRunePattern(' '), "_")
	want1 := "_Hello___friend_"
	got1 := n1.GetNormalized()
	if !reflect.DeepEqual(want1, got1) {
		t.Errorf("Want: %v\n", want1)
		t.Errorf("Got: %v\n", got1)
	}

	n2 := normalizer.NewNormalizedFrom("aaaab")
	n2 = n2.Replace(normalizer.NewRunePattern('a'), "b")
	want2 := "bbbbb"
	got2 := n2.GetNormalized()
	if !reflect.DeepEqual(want2, got2) {
		t.Errorf("Want: %v\n", want2)
		t.Errorf("Got: %v\n", got2)
	}

	// overlapping
	n3 := normalizer.NewNormalizedFrom("aaaab")
	n3 = n3.Replace(normalizer.NewStringPattern("aaa"), "b")
	want3 := "bab"
	got3 := n3.GetNormalized()
	if !reflect.DeepEqual(want3, got3) {
		t.Errorf("Want: %v\n", want3)
		t.Errorf("Got: %v\n", got3)
	}

	// Regexp
	n4 := normalizer.NewNormalizedFrom(" Hello   friend ")
	n4 = n4.Replace(normalizer.NewRegexpPattern(`\s+`), "_")
	want4 := "_Hello_friend_"
	got4 := n4.GetNormalized()
	if !reflect.DeepEqual(want4, got4) {
		t.Errorf("Want: %v\n", want4)
		t.Errorf("Got: %v\n", got4)
	}
}

func TestNormalized_Split(t *testing.T) {
	n := normalizer.NewNormalizedFrom("The-final--countdown")

	want1 := []string{"The", "final", "countdown"}
	testSplit(t, normalizer.RemovedBehavior, n, want1)

	want2 := []string{"The", "-", "final", "-", "-", "countdown"}
	testSplit(t, normalizer.IsolatediBehavior, n, want2)

	want3 := []string{"The-", "final-", "-", "countdown"}
	testSplit(t, normalizer.MergedWithPreviousBehavior, n, want3)

	want4 := []string{"The", "-final", "-", "-countdown"}
	testSplit(t, normalizer.MergedWithNextBehavior, n, want4)
}

func testSplit(t *testing.T, behavior normalizer.SplitDelimiterBehavior, n *normalizer.NormalizedString, want []string) {
	pattern := normalizer.NewStringPattern("-")

	splits := n.Split(pattern, behavior)

	var got []string
	for _, split := range splits {
		got = append(got, split.GetNormalized())
	}

	if !reflect.DeepEqual(want, got) {
		t.Errorf("Want: %v\n", want)
		t.Errorf("Got: %v\n", got)
	}
}

func TestNormalized_TransformRange_SingleBytes(t *testing.T) {

	n1 := normalizer.NewNormalizedFrom("Hello friend")

	// Removing at the beginning
	changeMap1 := []normalizer.ChangeMap{
		{"Y", 0},
	}
	got1 := n1.TransformRange(normalizer.NewRange(0, 4, normalizer.OriginalTarget), changeMap1, 3)
	original := "Hello friend"
	normalized := "Yo friend"
	alignments := [][]int{{3, 4}, {4, 5}, {5, 6}, {6, 7}, {7, 8}, {8, 9}, {9, 10}, {10, 11}, {11, 12}}
	alignmentsOriginal := [][]int{{0, 0}, {0, 0}, {0, 0}, {0, 1}, {1, 2}, {2, 3}, {3, 4}, {4, 5}, {5, 6}, {6, 7}, {7, 8}, {8, 9}}
	originalShift := 0
	want1 := normalizer.NewNormalizedString(original, normalized, alignments, alignmentsOriginal, originalShift)
	test(t, want1, got1)

	// Removing in the middle
	n2 := normalizer.NewNormalizedFrom("Hello friend")
	changeMap2 := []normalizer.ChangeMap{
		{"_", 0},
		{"F", 0},
		{"R", -2},
	}
	got2 := n2.TransformRange(normalizer.NewRange(3, 10, normalizer.OriginalTarget), changeMap2, 2)
	original = "Hello friend"
	normalized = "Hel_FRnd"
	alignments = [][]int{{0, 1}, {1, 2}, {2, 3}, {5, 6}, {6, 7}, {7, 8}, {10, 11}, {11, 12}}
	alignmentsOriginal = [][]int{{0, 1}, {1, 2}, {2, 3}, {3, 3}, {3, 3}, {3, 4}, {4, 5}, {5, 6}, {6, 6}, {6, 6}, {6, 7}, {7, 8}}
	originalShift = 0
	want2 := normalizer.NewNormalizedString(original, normalized, alignments, alignmentsOriginal, originalShift)
	test(t, want2, got2)

	// Removing at the end
	n3 := normalizer.NewNormalizedFrom("Hello friend")
	changeMap3 := []normalizer.ChangeMap{
		{"_", 0},
		{"F", -5},
	}
	got3 := n3.TransformRange(normalizer.NewRange(5, len([]byte(n3.GetOriginal())), normalizer.OriginalTarget), changeMap3, 0)
	original = "Hello friend"
	normalized = "Hello_F"
	alignments = [][]int{{0, 1}, {1, 2}, {2, 3}, {3, 4}, {4, 5}, {5, 6}, {6, 7}}
	alignmentsOriginal = [][]int{{0, 1}, {1, 2}, {2, 3}, {3, 4}, {4, 5}, {5, 6}, {6, 7}, {7, 7}, {7, 7}, {7, 7}, {7, 7}, {7, 7}}
	originalShift = 0
	want3 := normalizer.NewNormalizedString(original, normalized, alignments, alignmentsOriginal, originalShift)
	test(t, want3, got3)

	// Adding at the beginning
	n4 := normalizer.NewNormalizedFrom("Hello friend")
	changeMap4 := []normalizer.ChangeMap{
		{"H", 1},
		{"H", 0},
	}
	got4 := n4.TransformRange(normalizer.NewRange(0, 1, normalizer.OriginalTarget), changeMap4, 0)
	original = "Hello friend"
	normalized = "HHello friend"
	alignments = [][]int{{0, 0}, {0, 1}, {1, 2}, {2, 3}, {3, 4}, {4, 5}, {5, 6}, {6, 7}, {7, 8}, {8, 9}, {9, 10}, {10, 11}, {11, 12}}
	alignmentsOriginal = [][]int{{1, 2}, {2, 3}, {3, 4}, {4, 5}, {5, 6}, {6, 7}, {7, 8}, {8, 9}, {9, 10}, {10, 11}, {11, 12}, {12, 13}}
	originalShift = 0
	want4 := normalizer.NewNormalizedString(original, normalized, alignments, alignmentsOriginal, originalShift)
	test(t, want4, got4)

	// Equivalent to the previous one
	n5 := normalizer.NewNormalizedFrom("Hello friend")
	changeMap5 := []normalizer.ChangeMap{
		{"H", 1},
	}
	got5 := n5.TransformRange(normalizer.NewRange(0, 0, normalizer.OriginalTarget), changeMap5, 0)
	original = "Hello friend"
	normalized = "HHello friend"
	alignments = [][]int{{0, 0}, {0, 1}, {1, 2}, {2, 3}, {3, 4}, {4, 5}, {5, 6}, {6, 7}, {7, 8}, {8, 9}, {9, 10}, {10, 11}, {11, 12}}
	alignmentsOriginal = [][]int{{1, 2}, {2, 3}, {3, 4}, {4, 5}, {5, 6}, {6, 7}, {7, 8}, {8, 9}, {9, 10}, {10, 11}, {11, 12}, {12, 13}}
	originalShift = 0
	want5 := normalizer.NewNormalizedString(original, normalized, alignments, alignmentsOriginal, originalShift)
	test(t, want5, got5)

	// Adding as part of the first character
	n6 := normalizer.NewNormalizedFrom("Hello friend")
	changeMap6 := []normalizer.ChangeMap{
		{"H", 0},
		{"H", 1},
	}
	got6 := n6.TransformRange(normalizer.NewRange(0, 1, normalizer.OriginalTarget), changeMap6, 0)
	original = "Hello friend"
	normalized = "HHello friend"
	alignments = [][]int{{0, 1}, {0, 1}, {1, 2}, {2, 3}, {3, 4}, {4, 5}, {5, 6}, {6, 7}, {7, 8}, {8, 9}, {9, 10}, {10, 11}, {11, 12}}
	alignmentsOriginal = [][]int{{0, 2}, {2, 3}, {3, 4}, {4, 5}, {5, 6}, {6, 7}, {7, 8}, {8, 9}, {9, 10}, {10, 11}, {11, 12}, {12, 13}}
	originalShift = 0
	want6 := normalizer.NewNormalizedString(original, normalized, alignments, alignmentsOriginal, originalShift)
	test(t, want6, got6)

	// Adding in the middle
	n7 := normalizer.NewNormalizedFrom("Hello friend")
	changeMap7 := []normalizer.ChangeMap{
		{"_", 0},
		{"m", 1},
		{"y", 1},
		{"_", 1},
	}
	got7 := n7.TransformRange(normalizer.NewRange(5, 6, normalizer.OriginalTarget), changeMap7, 0)
	original = "Hello friend"
	normalized = "Hello_my_friend"
	alignments = [][]int{{0, 1}, {1, 2}, {2, 3}, {3, 4}, {4, 5}, {5, 6}, {5, 6}, {5, 6}, {5, 6}, {6, 7}, {7, 8}, {8, 9}, {9, 10}, {10, 11}, {11, 12}}
	alignmentsOriginal = [][]int{{0, 1}, {1, 2}, {2, 3}, {3, 4}, {4, 5}, {5, 9}, {9, 10}, {10, 11}, {11, 12}, {12, 13}, {13, 14}, {14, 15}}
	originalShift = 0
	want7 := normalizer.NewNormalizedString(original, normalized, alignments, alignmentsOriginal, originalShift)
	test(t, want7, got7)

	// Adding at the end
	n8 := normalizer.NewNormalizedFrom("Hello friend")
	changeMap8 := []normalizer.ChangeMap{
		{"d", 0},
		{"_", 1},
		{"!", 1},
	}
	got8 := n8.TransformRange(normalizer.NewRange(11, len([]byte(n8.GetNormalized())), normalizer.OriginalTarget), changeMap8, 0)
	original = "Hello friend"
	normalized = "Hello friend_!"
	alignments = [][]int{{0, 1}, {1, 2}, {2, 3}, {3, 4}, {4, 5}, {5, 6}, {6, 7}, {7, 8}, {8, 9}, {9, 10}, {10, 11}, {11, 12}, {11, 12}, {11, 12}}
	alignmentsOriginal = [][]int{{0, 1}, {1, 2}, {2, 3}, {3, 4}, {4, 5}, {5, 6}, {6, 7}, {7, 8}, {8, 9}, {9, 10}, {10, 11}, {11, 14}}
	originalShift = 0
	want8 := normalizer.NewNormalizedString(original, normalized, alignments, alignmentsOriginal, originalShift)
	test(t, want8, got8)
}

func TestNormalized_TransformRange_MultipleBytes(t *testing.T) {

	n1 := normalizer.NewNormalizedFrom("𝔾𝕠𝕠𝕕")

	// Removing at the beginning
	changeMap1 := []normalizer.ChangeMap{
		{"G", -1},
	}
	got1 := n1.TransformRange(normalizer.NewRange(0, 8, normalizer.OriginalTarget), changeMap1, 0)
	original := "𝔾𝕠𝕠𝕕"
	normalized := "G𝕠𝕕"
	alignments := [][]int{{0, 4}, {8, 12}, {8, 12}, {8, 12}, {8, 12}, {12, 16}, {12, 16}, {12, 16}, {12, 16}}
	alignmentsOriginal := [][]int{{0, 1}, {0, 1}, {0, 1}, {0, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 5}, {1, 5}, {1, 5}, {1, 5}, {5, 9}, {5, 9}, {5, 9}, {5, 9}}
	originalShift := 0
	want1 := normalizer.NewNormalizedString(original, normalized, alignments, alignmentsOriginal, originalShift)
	test(t, want1, got1)

	got2 := n1.Range(normalizer.NewRange(0, 8, normalizer.OriginalTarget))
	want2 := "G"
	test(t, want2, got2)

	got3 := n1.Range(normalizer.NewRange(0, 4, normalizer.OriginalTarget))
	want3 := "G"
	test(t, want3, got3)

	got4 := n1.RangeOriginal(normalizer.NewRange(0, 4, normalizer.OriginalTarget))
	want4 := "𝔾"
	test(t, want4, got4)

	got5 := n1.RangeOriginal(normalizer.NewRange(0, 8, normalizer.OriginalTarget))
	want5 := "𝔾𝕠"
	test(t, want5, got5)

	// Removing in the middle
	n2 := normalizer.NewNormalizedFrom("𝔾𝕠𝕠𝕕")
	changeMap2 := []normalizer.ChangeMap{
		{"o", -1},
	}
	got6 := n2.TransformRange(normalizer.NewRange(4, 12, normalizer.OriginalTarget), changeMap2, 0)
	original = "𝔾𝕠𝕠𝕕"
	normalized = "𝔾o𝕕"
	alignments = [][]int{{0, 4}, {0, 4}, {0, 4}, {0, 4}, {4, 8}, {12, 16}, {12, 16}, {12, 16}, {12, 16}}
	alignmentsOriginal = [][]int{{0, 4}, {0, 4}, {0, 4}, {0, 4}, {4, 5}, {4, 5}, {4, 5}, {4, 5}, {5, 5}, {5, 5}, {5, 5}, {5, 5}, {5, 9}, {5, 9}, {5, 9}, {5, 9}}
	originalShift = 0
	want6 := normalizer.NewNormalizedString(original, normalized, alignments, alignmentsOriginal, originalShift)
	test(t, want6, got6)

	// Removing at the end
	n3 := normalizer.NewNormalizedFrom("𝔾𝕠𝕠𝕕")
	changeMap3 := []normalizer.ChangeMap{
		{"d", 0},
		{"!", 1},
	}
	got7 := n3.TransformRange(normalizer.NewRange(12, len([]byte(n3.GetNormalized())), normalizer.OriginalTarget), changeMap3, 0)
	original = "𝔾𝕠𝕠𝕕"
	normalized = "𝔾𝕠𝕠d!"
	alignments = [][]int{{0, 4}, {0, 4}, {0, 4}, {0, 4}, {4, 8}, {4, 8}, {4, 8}, {4, 8}, {8, 12}, {8, 12}, {8, 12}, {8, 12}, {12, 16}, {12, 16}}
	alignmentsOriginal = [][]int{{0, 4}, {0, 4}, {0, 4}, {0, 4}, {4, 8}, {4, 8}, {4, 8}, {4, 8}, {8, 12}, {8, 12}, {8, 12}, {8, 12}, {12, 14}, {12, 14}, {12, 14}, {12, 14}}
	originalShift = 0
	want7 := normalizer.NewNormalizedString(original, normalized, alignments, alignmentsOriginal, originalShift)
	test(t, want7, got7)

	// Adding at the beginning
	n4 := normalizer.NewNormalizedFrom("𝔾𝕠𝕠𝕕")
	changeMap4 := []normalizer.ChangeMap{
		{"_", 1},
		{"𝔾", 0},
	}
	got8 := n4.TransformRange(normalizer.NewRange(0, 1, normalizer.OriginalTarget), changeMap4, 0)
	original = "𝔾𝕠𝕠𝕕"
	normalized = "_𝔾𝕠𝕠𝕕"
	alignments = [][]int{{0, 0}, {0, 4}, {0, 4}, {0, 4}, {0, 4}, {4, 8}, {4, 8}, {4, 8}, {4, 8}, {8, 12}, {8, 12}, {8, 12}, {8, 12}, {12, 16}, {12, 16}, {12, 16}, {12, 16}}
	alignmentsOriginal = [][]int{{1, 5}, {1, 5}, {1, 5}, {1, 5}, {5, 9}, {5, 9}, {5, 9}, {5, 9}, {9, 13}, {9, 13}, {9, 13}, {9, 13}, {13, 17}, {13, 17}, {13, 17}, {13, 17}}
	originalShift = 0
	want8 := normalizer.NewNormalizedString(original, normalized, alignments, alignmentsOriginal, originalShift)
	test(t, want8, got8)

	got9 := n4.Range(normalizer.NewRange(0, 8, normalizer.OriginalTarget))
	want9 := "𝔾𝕠"
	test(t, want9, got9)

	got10 := n4.Range(normalizer.NewRange(0, 4, normalizer.OriginalTarget))
	want10 := "𝔾"
	test(t, want10, got10)

	got11 := n4.RangeOriginal(normalizer.NewRange(0, 4, normalizer.OriginalTarget))
	want11 := "𝔾"
	test(t, want11, got11)

	got12 := n4.RangeOriginal(normalizer.NewRange(0, 8, normalizer.OriginalTarget))
	want12 := "𝔾𝕠"
	test(t, want12, got12)

	// Equivalent to the previous one
	n5 := normalizer.NewNormalizedFrom("𝔾𝕠𝕠𝕕")
	changeMap5 := []normalizer.ChangeMap{
		{"_", 1},
	}
	got13 := n5.TransformRange(normalizer.NewRange(0, 0, normalizer.OriginalTarget), changeMap5, 0)
	original = "𝔾𝕠𝕠𝕕"
	normalized = "_𝔾𝕠𝕠𝕕"
	alignments = [][]int{{0, 0}, {0, 4}, {0, 4}, {0, 4}, {0, 4}, {4, 8}, {4, 8}, {4, 8}, {4, 8}, {8, 12}, {8, 12}, {8, 12}, {8, 12}, {12, 16}, {12, 16}, {12, 16}, {12, 16}}
	alignmentsOriginal = [][]int{{1, 5}, {1, 5}, {1, 5}, {1, 5}, {5, 9}, {5, 9}, {5, 9}, {5, 9}, {9, 13}, {9, 13}, {9, 13}, {9, 13}, {13, 17}, {13, 17}, {13, 17}, {13, 17}}
	originalShift = 0
	want13 := normalizer.NewNormalizedString(original, normalized, alignments, alignmentsOriginal, originalShift)
	test(t, want13, got13)

	got14 := n5.Range(normalizer.NewRange(0, 8, normalizer.OriginalTarget))
	want14 := "𝔾𝕠"
	test(t, want14, got14)

	got15 := n5.Range(normalizer.NewRange(0, 4, normalizer.OriginalTarget))
	want15 := "𝔾"
	test(t, want15, got15)

	got16 := n5.RangeOriginal(normalizer.NewRange(0, 4, normalizer.OriginalTarget))
	want16 := "𝔾"
	test(t, want16, got16)

	got17 := n5.RangeOriginal(normalizer.NewRange(0, 8, normalizer.OriginalTarget))
	want17 := "𝔾𝕠"
	test(t, want17, got17)

	// Adding as part of the first character
	n6 := normalizer.NewNormalizedFrom("𝔾𝕠𝕠𝕕")
	changeMap6 := []normalizer.ChangeMap{
		{"𝔾", 0},
		{"o", 1},
	}
	got18 := n6.TransformRange(normalizer.NewRange(0, 1, normalizer.OriginalTarget), changeMap6, 0)
	original = "𝔾𝕠𝕠𝕕"
	normalized = "𝔾o𝕠𝕠𝕕"
	alignments = [][]int{{0, 4}, {0, 4}, {0, 4}, {0, 4}, {0, 4}, {4, 8}, {4, 8}, {4, 8}, {4, 8}, {8, 12}, {8, 12}, {8, 12}, {8, 12}, {12, 16}, {12, 16}, {12, 16}, {12, 16}}
	alignmentsOriginal = [][]int{{0, 5}, {0, 5}, {0, 5}, {0, 5}, {5, 9}, {5, 9}, {5, 9}, {5, 9}, {9, 13}, {9, 13}, {9, 13}, {9, 13}, {13, 17}, {13, 17}, {13, 17}, {13, 17}}
	originalShift = 0
	want18 := normalizer.NewNormalizedString(original, normalized, alignments, alignmentsOriginal, originalShift)
	test(t, want18, got18)

	got19 := n6.Range(normalizer.NewRange(0, 8, normalizer.OriginalTarget))
	want19 := "𝔾o𝕠"
	test(t, want19, got19)

	got20 := n6.Range(normalizer.NewRange(0, 4, normalizer.OriginalTarget))
	want20 := "𝔾o"
	test(t, want20, got20)

	got21 := n6.RangeOriginal(normalizer.NewRange(0, 4, normalizer.OriginalTarget))
	want21 := "𝔾"
	test(t, want21, got21)

	got22 := n6.RangeOriginal(normalizer.NewRange(0, 8, normalizer.OriginalTarget))
	want22 := "𝔾𝕠"
	test(t, want22, got22)

	// Adding in the middle
	n7 := normalizer.NewNormalizedFrom("𝔾𝕠𝕠𝕕")
	changeMap7 := []normalizer.ChangeMap{
		{"𝕠", 0},
		{"o", 1},
		{"o", 1},
		{"o", 1},
	}
	got23 := n7.TransformRange(normalizer.NewRange(4, 8, normalizer.OriginalTarget), changeMap7, 0)
	original = "𝔾𝕠𝕠𝕕"
	normalized = "𝔾𝕠ooo𝕠𝕕"
	alignments = [][]int{{0, 4}, {0, 4}, {0, 4}, {0, 4}, {4, 8}, {4, 8}, {4, 8}, {4, 8}, {4, 8}, {4, 8}, {4, 8}, {8, 12}, {8, 12}, {8, 12}, {8, 12}, {12, 16}, {12, 16}, {12, 16}, {12, 16}}
	alignmentsOriginal = [][]int{{0, 4}, {0, 4}, {0, 4}, {0, 4}, {4, 11}, {4, 11}, {4, 11}, {4, 11}, {11, 15}, {11, 15}, {11, 15}, {11, 15}, {15, 19}, {15, 19}, {15, 19}, {15, 19}}
	originalShift = 0
	want23 := normalizer.NewNormalizedString(original, normalized, alignments, alignmentsOriginal, originalShift)
	test(t, want23, got23)

	// Adding at the end
	n8 := normalizer.NewNormalizedFrom("𝔾𝕠𝕠𝕕")
	changeMap8 := []normalizer.ChangeMap{
		{"!", 1},
	}
	got24 := n8.TransformRange(normalizer.NewRange(16, len([]byte(n8.GetNormalized())), normalizer.OriginalTarget), changeMap8, 0)
	original = "𝔾𝕠𝕠𝕕"
	normalized = "𝔾𝕠𝕠𝕕!"
	alignments = [][]int{{0, 4}, {0, 4}, {0, 4}, {0, 4}, {4, 8}, {4, 8}, {4, 8}, {4, 8}, {8, 12}, {8, 12}, {8, 12}, {8, 12}, {12, 16}, {12, 16}, {12, 16}, {12, 16}, {12, 16}}
	alignmentsOriginal = [][]int{{0, 4}, {0, 4}, {0, 4}, {0, 4}, {4, 8}, {4, 8}, {4, 8}, {4, 8}, {8, 12}, {8, 12}, {8, 12}, {8, 12}, {12, 17}, {12, 17}, {12, 17}, {12, 17}}
	originalShift = 0
	want24 := normalizer.NewNormalizedString(original, normalized, alignments, alignmentsOriginal, originalShift)
	test(t, want24, got24)
}

func test(t *testing.T, want, got interface{}) {
	if !reflect.DeepEqual(want, got) {
		t.Errorf("Want: %v\n", want)
		t.Errorf("Got: %v\n", got)
	}
}

package normalizer

import (
	"reflect"
	"testing"
)

func TestPrecompiled_Normalize(t *testing.T) {
	var transformations []ChangeMap = make([]ChangeMap, 10)

	n := NewNormalizedFrom("™\x1eg")
	transformations = replace(transformations, "™", "TM")

	transformations = replace(transformations, "\x1e", "")

	transformations = append(transformations, ChangeMap{
		RuneVal: "g",
		Changes: 0,
	})

	n = n.Transform(transformations, 0)

	got := n.GetNormalized()
	want := "TMg"

	if !reflect.DeepEqual(want, got) {
		t.Errorf("want %s, got %s\n", want, got)
	}
}

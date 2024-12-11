package normalizer

import (
	"reflect"
	"testing"
)

func Test_validateRange(t *testing.T) {

	tests := []struct {
		name       string
		s          string
		inputRange *Range
		want       *Range
	}{
		{
			name:       "Valid range 0,4",
			s:          "Lion Löwe 老虎 Léopard",
			inputRange: NewRange(0, 4, OriginalTarget),
			want:       NewRange(0, 4, OriginalTarget),
		},
		{
			name:       "Valid range 0,3",
			s:          "老",
			inputRange: NewRange(0, 3, OriginalTarget),
			want:       NewRange(0, 3, OriginalTarget),
		},
		{
			name:       "Valid range 10,13",
			s:          "Lion Löwe 老虎 Léopard",
			inputRange: NewRange(11, 14, OriginalTarget),
			want:       NewRange(11, 14, OriginalTarget),
		},
		// character is 3 bytes, so 1-2 is invalid
		{
			name:       "Invalid range",
			s:          `老`,
			inputRange: NewRange(1, 2, OriginalTarget),
			want:       nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ns := NewNormalizedFrom(tt.s)
			got := ns.validateRange(tt.inputRange)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Expected range %v, got %v", got, tt.want)
			}
		})
	}
}

func Benchmark_validateRange(b *testing.B) {
	s := "Lion Löwe 老虎 Léopard"
	ns := NewNormalizedFrom(s)

	var r *Range
	for i := 0; i < b.N; i++ {
		r = ns.validateRange(NewRange(0, len(s), NormalizedTarget))
	}

	_ = r
}
